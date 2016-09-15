package chart

import (
	"log"

	"github.com/kubernetes/dashboard/src/app/backend/client"
	"k8s.io/helm/pkg/helm"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

// AppDeploymentFromChartSpec is a specification for a chart deployment.
type AppDeploymentFromChartSpec struct {
	// Name of the chart.
	ChartName string `json:"chartName"`

	// Name of the release.
	ReleaseName string `json:"releaseName"`

	// Namespace for release.
	Namespace string `json:"namespace"`
}

// AppDeploymentFromChartResponse is a specification for a chart deployment.
type AppDeploymentFromChartResponse struct {
	// Name of the chart.
	ChartName string `json:"chartName"`

	// Name of the release.
	ReleaseName string `json:"releaseName"`

	// Namespace for release.
	Namespace string `json:"namespace"`

	// Error after deploying chart
	Error string `json:"error"`
}

// TODO: relocate to chart resource pakcage
// DeployChart deploys an chart based on the given configuration.
func DeployChart(spec *AppDeploymentFromChartSpec, c unversioned.Interface) error {
	log.Printf("Deploying chart %s with release name %s", spec.ChartName, spec.ReleaseName)

	// TODO: pre-init tiller client and provide as param to this func
	tc, err := client.CreateTillerClient()
	if err != nil {
		log.Printf("Error creating tiller client: %s", err)
		return err
	}

	// if res, err := tc.ListReleases(); err != nil {
	// 	log.Printf("Error listing releases: %s", err)
	// 	return err
	// }
	// log.Printf("helm releases: %s", res.Releases)

	chartPath, err := locateChartPath(spec.ChartName)
	if err != nil {
		log.Printf("Failed to find chart: %s", err)
		return err
	}

	res, err := tc.InstallRelease(
		chartPath,
		spec.Namespace,
		helm.ValueOverrides(nil),
		helm.ReleaseName(spec.ReleaseName),
	)
	if err != nil {
		log.Printf("Error installing release: %s", err)
		return err
	}
	log.Printf("Release response: %s", res)
	return nil
}
