package apply

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	clienttesting "k8s.io/client-go/testing"

	workapiv1 "open-cluster-management.io/api/work/v1"

	testingcommon "open-cluster-management.io/ocm/pkg/common/testing"
)

const defaultOwner = "test-owner"

func TestServerSideApply(t *testing.T) {
	cases := []struct {
		name            string
		owner           metav1.OwnerReference
		existing        *unstructured.Unstructured
		required        *unstructured.Unstructured
		gvr             schema.GroupVersionResource
		conflict        bool
		validateActions func(t *testing.T, actions []clienttesting.Action)
	}{
		{
			name:            "server side apply successfully",
			owner:           metav1.OwnerReference{APIVersion: "v1", Name: "test", UID: defaultOwner},
			existing:        nil,
			required:        testingcommon.NewUnstructured("v1", "Namespace", "", "test"),
			gvr:             schema.GroupVersionResource{Version: "v1", Resource: "namespaces"},
			validateActions: func(t *testing.T, actions []clienttesting.Action) {},
		},
		{
			name:     "server side apply successfully conflict",
			owner:    metav1.OwnerReference{APIVersion: "v1", Name: "test", UID: defaultOwner},
			existing: testingcommon.NewUnstructured("v1", "Secret", "ns1", "test"),
			required: testingcommon.NewUnstructured("v1", "Secret", "ns1", "test"),
			gvr:      schema.GroupVersionResource{Version: "v1", Resource: "secrets"},
			conflict: true,
			validateActions: func(t *testing.T, actions []clienttesting.Action) {
				testingcommon.AssertActions(t, actions, "patch")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var objects []runtime.Object
			if c.existing != nil {
				objects = append(objects, c.existing)
			}
			scheme := runtime.NewScheme()
			dynamicClient := fakedynamic.NewSimpleDynamicClient(scheme, objects...)

			// The fake client does not support PatchType ApplyPatchType, add an reactor to mock apply patch
			// see issue: https://github.com/kubernetes/kubernetes/issues/103816
			reactor := &reactor{}
			reactors := []clienttesting.Reactor{reactor}
			dynamicClient.Fake.ReactionChain = append(reactors, dynamicClient.Fake.ReactionChain...)

			applier := NewServerSideApply(dynamicClient)

			syncContext := testingcommon.NewFakeSyncContext(t, "test")
			option := &workapiv1.ManifestConfigOption{
				UpdateStrategy: &workapiv1.UpdateStrategy{
					Type: workapiv1.UpdateStrategyTypeServerSideApply,
				},
			}
			obj, err := applier.Apply(
				context.TODO(), c.gvr, c.required, c.owner, option, syncContext.Recorder())
			c.validateActions(t, dynamicClient.Actions())
			if !c.conflict {
				if err != nil {
					t.Errorf("expect no error, but got %v", err)
				}

				accessor, err := meta.Accessor(obj)
				if err != nil {
					t.Errorf("type %t cannot be accessed: %v", obj, err)
				}
				if accessor.GetNamespace() != c.required.GetNamespace() || accessor.GetName() != c.required.GetName() {
					t.Errorf("Expect resource %s/%s, but %s/%s",
						c.required.GetNamespace(), c.required.GetName(), accessor.GetNamespace(), accessor.GetName())
				}
				return
			}

			var ssaConflict *ServerSideApplyConflictError
			if !errors.As(err, &ssaConflict) {
				t.Errorf("expect serverside apply conflict error, but got %v", err)
			}

		})
	}
}

type reactor struct {
}

func (r *reactor) Handles(action clienttesting.Action) bool {
	switch action := action.(type) {
	case clienttesting.PatchActionImpl:
		if action.GetPatchType() == types.ApplyPatchType {
			return true
		}

	default:
		return false
	}
	return true
}

// React handles the action and returns results.  It may choose to
// delegate by indicated handled=false.
func (r *reactor) React(action clienttesting.Action) (handled bool, ret runtime.Object, err error) {
	switch action.GetResource().Resource {
	case "namespaces":
		return true, testingcommon.NewUnstructured("v1", "Namespace", "", "test"), nil
	case "secrets":
		return true, nil, apierrors.NewApplyConflict([]metav1.StatusCause{
			{
				Type:    metav1.CauseTypeFieldManagerConflict,
				Message: "field managed configl",
				Field:   "metadata.annotations",
			},
		}, "server side apply secret failed")
	}

	return true, nil, fmt.Errorf("PatchType is not supported")
}

func TestRemoveCreationTime(t *testing.T) {
	cases := []struct {
		name         string
		required     *unstructured.Unstructured
		validateFunc func(t *testing.T, obj *unstructured.Unstructured)
	}{
		{
			name:     "remove creationTimestamp from a kube object",
			required: newDeployment(),
			validateFunc: func(t *testing.T, obj *unstructured.Unstructured) {
				_, existing, err := unstructured.NestedFieldCopy(obj.Object, "metadata", "creationTimestamp")
				if err != nil {
					t.Fatal(err)
				}
				if existing {
					t.Errorf("unexpected creationTimestamp in `metadata.creationTimestamp`")
				}
				_, existing, err = unstructured.NestedFieldCopy(obj.Object, "spec", "template", "metadata", "creationTimestamp")
				if err != nil {
					t.Fatal(err)
				}
				if existing {
					t.Errorf("unexpected creationTimestamp in `spec.template.metadata.creationTimestamp`")
				}
			},
		},
		{
			name:     "remove creationTimestamp from a manifestwork",
			required: newManifestWork(),
			validateFunc: func(t *testing.T, obj *unstructured.Unstructured) {
				_, existing, err := unstructured.NestedFieldCopy(obj.Object, "metadata", "creationTimestamp")
				if err != nil {
					t.Fatal(err)
				}
				if existing {
					t.Errorf("unexpected creationTimestamp in `metadata.creationTimestamp`")
				}

				manifests, existing, err := unstructured.NestedSlice(obj.Object, "spec", "workload", "manifests")
				if err != nil {
					t.Fatal(err)
				}
				if !existing {
					t.Fatalf("no manifests")
				}

				_, existing, err = unstructured.NestedFieldCopy(manifests[0].(map[string]interface{}), "metadata", "creationTimestamp")
				if err != nil {
					t.Fatal(err)
				}
				if existing {
					t.Errorf("unexpected creationTimestamp in `spec.workload.manifests[0].metadata.creationTimestamp`")
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			removeCreationTimeFromMetadata(c.required.Object)
			c.validateFunc(t, c.required)
		})
	}
}

func newDeployment() *unstructured.Unstructured {
	var replicas int32 = 3
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"test": "test"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "test",
							Image: "test",
						},
					},
				},
			},
		},
	}

	obj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(deploy)
	return &unstructured.Unstructured{Object: obj}
}

func newManifestWork() *unstructured.Unstructured {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm-test",
			Namespace: "default",
		},
		Data: map[string]string{
			"some": "data",
		},
	}

	cmObj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(cm)
	raw, _ := (&unstructured.Unstructured{Object: cmObj}).MarshalJSON()
	manifest := workapiv1.Manifest{}
	manifest.Raw = raw

	work := &workapiv1.ManifestWork{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "work.open-cluster-management.io/v1",
			Kind:       "ManifestWork",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: workapiv1.ManifestWorkSpec{
			Workload: workapiv1.ManifestsTemplate{
				Manifests: []workapiv1.Manifest{manifest},
			},
		},
	}

	obj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(work)
	return &unstructured.Unstructured{Object: obj}
}
