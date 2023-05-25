package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_AppliedClusterResourceQuota = map[string]string{
	"":       "AppliedClusterResourceQuota mirrors ClusterResourceQuota at a project scope, for projection into a project.  It allows a project-admin to know which ClusterResourceQuotas are applied to his project and their associated usage.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"spec":   "Spec defines the desired quota",
	"status": "Status defines the actual enforced quota and its current usage",
}

func (AppliedClusterResourceQuota) SwaggerDoc() map[string]string {
	return map_AppliedClusterResourceQuota
}

var map_AppliedClusterResourceQuotaList = map[string]string{
	"":      "AppliedClusterResourceQuotaList is a collection of AppliedClusterResourceQuotas\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"items": "Items is a list of AppliedClusterResourceQuota",
}

func (AppliedClusterResourceQuotaList) SwaggerDoc() map[string]string {
	return map_AppliedClusterResourceQuotaList
}

var map_ClusterResourceQuota = map[string]string{
	"":       "ClusterResourceQuota mirrors ResourceQuota at a cluster scope.  This object is easily convertible to synthetic ResourceQuota object to allow quota evaluation re-use.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"spec":   "Spec defines the desired quota",
	"status": "Status defines the actual enforced quota and its current usage",
}

func (ClusterResourceQuota) SwaggerDoc() map[string]string {
	return map_ClusterResourceQuota
}

var map_ClusterResourceQuotaList = map[string]string{
	"":      "ClusterResourceQuotaList is a collection of ClusterResourceQuotas\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"items": "Items is a list of ClusterResourceQuotas",
}

func (ClusterResourceQuotaList) SwaggerDoc() map[string]string {
	return map_ClusterResourceQuotaList
}

var map_ClusterResourceQuotaSelector = map[string]string{
	"":            "ClusterResourceQuotaSelector is used to select projects.  At least one of LabelSelector or AnnotationSelector must present.  If only one is present, it is the only selection criteria.  If both are specified, the project must match both restrictions.",
	"labels":      "LabelSelector is used to select projects by label.",
	"annotations": "AnnotationSelector is used to select projects by annotation.",
}

func (ClusterResourceQuotaSelector) SwaggerDoc() map[string]string {
	return map_ClusterResourceQuotaSelector
}

var map_ClusterResourceQuotaSpec = map[string]string{
	"":         "ClusterResourceQuotaSpec defines the desired quota restrictions",
	"selector": "Selector is the selector used to match projects. It should only select active projects on the scale of dozens (though it can select many more less active projects).  These projects will contend on object creation through this resource.",
	"quota":    "Quota defines the desired quota",
}

func (ClusterResourceQuotaSpec) SwaggerDoc() map[string]string {
	return map_ClusterResourceQuotaSpec
}

var map_ClusterResourceQuotaStatus = map[string]string{
	"":           "ClusterResourceQuotaStatus defines the actual enforced quota and its current usage",
	"total":      "Total defines the actual enforced quota and its current usage across all projects",
	"namespaces": "Namespaces slices the usage by project.  This division allows for quick resolution of deletion reconciliation inside of a single project without requiring a recalculation across all projects.  This can be used to pull the deltas for a given project.",
}

func (ClusterResourceQuotaStatus) SwaggerDoc() map[string]string {
	return map_ClusterResourceQuotaStatus
}

var map_ResourceQuotaStatusByNamespace = map[string]string{
	"":          "ResourceQuotaStatusByNamespace gives status for a particular project",
	"namespace": "Namespace the project this status applies to",
	"status":    "Status indicates how many resources have been consumed by this project",
}

func (ResourceQuotaStatusByNamespace) SwaggerDoc() map[string]string {
	return map_ResourceQuotaStatusByNamespace
}

// AUTO-GENERATED FUNCTIONS END HERE
