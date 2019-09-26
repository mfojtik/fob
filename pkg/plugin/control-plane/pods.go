package control_plane

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func filterControlPlanePods(list *[]*v1.Pod) {
	filteredList := []*v1.Pod{}
	for _, pod := range *list {
		found := false
		if controlPlaneOperators.Has(pod.Namespace) {
			found = true
		}
		if !found && !controlPlaneOperators.Has(strings.TrimSuffix(strings.TrimPrefix(pod.Namespace, "openshift-"), "-operator")) {
			continue
		}
		filteredList = append(filteredList, pod)
	}
	*list = filteredList
}

func formatControlPlanePods(pods []*v1.Pod) [][]string {
	output := [][]string{}
	for _, pod := range pods {
		restartsCount := int32(0)
		for _, c := range pod.Status.ContainerStatuses {
			restartsCount += c.RestartCount
		}
		for _, c := range pod.Status.InitContainerStatuses {
			restartsCount += c.RestartCount
		}
		podDetails := []string{
			pod.Namespace + "/" + pod.Name,
			string(pod.Status.Phase),
			fmt.Sprintf("%d restarts", restartsCount),
		}
		output = append(output, podDetails)
	}
	return output
}

func decodePods(manifestBytes []byte) ([]*v1.Pod, error) {
	listObj, err := runtime.Decode(decoder, manifestBytes)
	if err != nil {
		return nil, err
	}
	listItems := listObj.(*corev1.List).Items
	result := make([]*v1.Pod, len(listItems))
	for i, item := range listItems {
		operatorObj, err := runtime.Decode(decoder, item.Raw)
		if err != nil {
			return nil, err
		}
		result[i] = operatorObj.(*v1.Pod)
	}
	return result, nil
}
