package control_plane

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/mfojtik/fob/pkg/artifact"
	"github.com/mfojtik/fob/pkg/plugin"
)

var controlPlaneOperators = sets.NewString(
	"kube-apiserver",
	"kube-controller-manager",
	"kube-scheduler",
	"openshift-apiserver",
	"openshift-controller-manager",
	"authentication",
)

type ControlPlanePlugin struct{}

func (c *ControlPlanePlugin) Register(m plugin.PluginManager) error {
	m.Add(c)
	return nil
}

func (c *ControlPlanePlugin) Execute(ctx context.Context, o plugin.PluginOptions) error {
	clusterOperatorBytes, err := artifact.Get(ctx, o.JobUrl, "/artifacts/e2e-aws/clusteroperators.json")
	if err != nil {
		return err
	}

	clusterOperators, err := decodeClusterOperators(clusterOperatorBytes)
	if err != nil {
		return err
	}
	o.PrintTable("ClusterOperator Status", formatControlPlaneClusterOperators(clusterOperators))

	o.Printf("Degraded Operators:\n\n%s\n", strings.Join(formatDegradedControlPlaneOperators(clusterOperators), "\n"))

	podsBytes, err := artifact.Get(ctx, o.JobUrl, "/artifacts/e2e-aws/pods.json")
	if err != nil {
		return err
	}

	pods, err := decodePods(podsBytes)
	if err != nil {
		return err
	}
	filterControlPlanePods(&pods)

	o.PrintTable("Control Plane Pods Status", formatControlPlanePods(pods))

	return nil
}
