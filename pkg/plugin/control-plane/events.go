package control_plane

import (
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func decodeEvents(manifestBytes []byte) ([]*v1.Event, error) {
	listObj, err := runtime.Decode(decoder, manifestBytes)
	if err != nil {
		return nil, err
	}
	listItems := listObj.(*corev1.List).Items
	result := make([]*v1.Event, len(listItems))
	for i, item := range listItems {
		operatorObj, err := runtime.Decode(decoder, item.Raw)
		if err != nil {
			return nil, err
		}
		result[i] = operatorObj.(*v1.Event)
	}
	return result, nil
}

func filterControlPlaneEvents(list *[]*v1.Event) {
	filteredList := []*v1.Event{}
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

func formatControlPlaneEvents(events []*v1.Event) [][]string {
	output := [][]string{}
	sort.Slice(events, func(i, j int) bool {
		return events[i].LastTimestamp.Time.Before(events[j].LastTimestamp.Time)
	})
	for _, event := range events {
		if len(event.Message) > 160 {
			event.Message = event.Message[0:160] + "..."
		}
		eventDetails := []string{
			event.LastTimestamp.String(),
			strings.TrimSuffix(strings.TrimPrefix(event.Namespace, "openshift-"), "-operator"), // keep the first column small
			event.Type,
			event.Message,
		}
		output = append(output, eventDetails)
	}
	return output
}
