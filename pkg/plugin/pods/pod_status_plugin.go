package pods

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mfojtik/fob/pkg/artifact"
	"github.com/mfojtik/fob/pkg/plugin"
)

type PodStatusPlugin struct{}

func (s *PodStatusPlugin) Register(manager plugin.PluginManager) error {
	manager.Add(s)
	return nil
}

func (s *PodStatusPlugin) Execute(ctx context.Context, o plugin.PluginOptions) error {
	podsBytes, err := artifact.Get(ctx, o.JobUrl, "artifacts/e2e-aws/pods.json")
	if err != nil {
		return err
	}

	podList := podList{}
	if err := json.Unmarshal(podsBytes, &podList); err != nil {
		return err
	}

	for _, pod := range podList.Items {
		restartsMsgs := []string{}
		for _, s := range pod.Status.ContainerStatuses {
			if s.RestartCount > 0 {
				restartsMsgs = append(restartsMsgs, fmt.Sprintf("container %q restarted %d times (%d)", s.Name, s.RestartCount, s.LastState.Terminated.ExitCode))
			}
		}
		for _, s := range pod.Status.InitContainerStatuses {
			if s.RestartCount > 0 {
				restartsMsgs = append(restartsMsgs, fmt.Sprintf("init container %q restarted %d times (%d)", s.Name, s.RestartCount, s.LastState.Terminated.ExitCode))
			}
		}
		restarstMsg := strings.Join(restartsMsgs, ",")
		if len(restartsMsgs) > 0 {
			restarstMsg = ", " + restarstMsg
			conditions := []string{}
			for _, c := range pod.Status.Conditions {
				conditions = append(conditions, fmt.Sprintf("%s=%s(%s)", c.Type, c.Reason, c.Status))
			}
			o.Printf("Pod %q is %q%s %s", pod.Namespace+"/"+pod.Name, pod.Status.Phase, restarstMsg, strings.Join(conditions, ","))
			for _, s := range pod.Status.ContainerStatuses {
				if len(s.LastState.Terminated.Message) > 0 {
					o.Printf("container %q:\n%s", s.Name, s.LastState.Terminated.Message)
				}
			}
		}
	}

	return nil
}

type podList struct {
	Items []pod `json:"items"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type pod struct {
	Metadata `json:"metadata"`
	Status   PodStatus `json:"status"`
}

type PodCondition struct {
	Reason string `json:"reason"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

type PodStatus struct {
	Conditions            []PodCondition    `json:"conditions"`
	ContainerStatuses     []ContainerStatus `json:"containerStatuses"`
	InitContainerStatuses []ContainerStatus `json:"initContainerStatuses"`
	Phase                 string            `json:"phase"`
	StartTime             string            `json:"startTime"`
}

type ContainerStatus struct {
	ContainerID  string              `json:"containerID"`
	Image        string              `json:"image"`
	ImageID      string              `json:"imageID"`
	LastState    ContainerLastStatus `json:"lastState"`
	Name         string              `json:"name"`
	Ready        bool                `json:"ready"`
	RestartCount int64               `json:"restartCount"`
	State        ContainerState      `json:"state"`
}

type ContainerLastStatus struct {
	Terminated PodList_sub54 `json:"terminated"`
}

type RunningState struct {
	StartedAt string `json:"startedAt"`
}

type PodList_sub54 struct {
	ContainerID string `json:"containerID"`
	ExitCode    int64  `json:"exitCode"`
	FinishedAt  string `json:"finishedAt"`
	Message     string `json:"message"`
	Reason      string `json:"reason"`
	StartedAt   string `json:"startedAt"`
}

type TerminatedState struct {
	ContainerID string `json:"containerID"`
	ExitCode    int64  `json:"exitCode"`
	FinishedAt  string `json:"finishedAt"`
	Reason      string `json:"reason"`
	StartedAt   string `json:"startedAt"`
}

type ContainerState struct {
	Running    RunningState    `json:"running"`
	Terminated TerminatedState `json:"terminated"`
}
