package control_plane

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	configv1 "github.com/openshift/api/config/v1"
)

var (
	scheme  = runtime.NewScheme()
	decoder runtime.Decoder
)

func init() {
	utilruntime.Must(corev1.AddToScheme(scheme)) // to get v1.List
	utilruntime.Must(configv1.Install(scheme))
	decoder = serializer.NewCodecFactory(scheme).UniversalDecoder(corev1.SchemeGroupVersion, configv1.GroupVersion)
}

func decodeClusterOperators(manifestBytes []byte) ([]*configv1.ClusterOperator, error) {
	listObj, err := runtime.Decode(decoder, manifestBytes)
	if err != nil {
		return nil, err
	}
	listItems := listObj.(*corev1.List).Items
	result := make([]*configv1.ClusterOperator, len(listItems))
	for i, item := range listItems {
		operatorObj, err := runtime.Decode(decoder, item.Raw)
		if err != nil {
			return nil, err
		}
		result[i] = operatorObj.(*configv1.ClusterOperator)
	}
	return result, nil
}

func formatDegradedControlPlaneOperators(clusterOperators []*configv1.ClusterOperator) []string {
	output := []string{}
	for _, operator := range clusterOperators {
		if !controlPlaneOperators.Has(operator.Name) {
			continue
		}
		for _, c := range operator.Status.Conditions {
			switch c.Type {
			case configv1.OperatorDegraded:
				if c.Status != configv1.ConditionFalse {
					output = append(output, fmt.Sprintf("Operator %q is %s since %s because %s: %q", operator.Name, c.Type, c.LastTransitionTime, c.Reason, c.Message))
				}
			case configv1.OperatorAvailable:
				if c.Status != configv1.ConditionTrue {
					output = append(output, fmt.Sprintf("Operator %q is not %s since %s because %s: %q", operator.Name, c.Type, c.LastTransitionTime, c.Reason, c.Message))
				}
			case configv1.OperatorUpgradeable:
				if c.Status != configv1.ConditionTrue {
					output = append(output, fmt.Sprintf("Operator %q is not %s since %s because %s: %q", operator.Name, c.Type, c.LastTransitionTime, c.Reason, c.Message))
				}
			}
		}
	}
	return output
}

func formatControlPlaneClusterOperators(clusterOperators []*configv1.ClusterOperator) [][]string {
	output := [][]string{}
	for _, operator := range clusterOperators {
		if !controlPlaneOperators.Has(operator.Name) {
			continue
		}
		conditions := []string{}
		for _, c := range operator.Status.Conditions {
			conditions = append(conditions, fmt.Sprintf("%s=%s", c.Type, c.Status))
		}
		output = append(output, append([]string{operator.Name}, conditions...))
	}
	return output
}
