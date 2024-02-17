package kube

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

type Pods []PodInfo

type PodsFilter interface {
	Filter(v1.ContainerStatus) bool
}

type PodInfo struct {
	Pod v1.Pod
	// Memory         MemoryInfo
	// ContainerName  string // Name of the container within the pod that was terminated, in the case of multi-container pods.
	TerminatedTime string // When the pod was terminated
	StartTime      string // When the pod was started during the termination period.
	// Internal representation of TerminatedTime, used for operations which require
	// the explicit time.Time type, such as sorting.
	terminatedTime time.Time
	RestartCount   string
}

func (t Pods) SortByTimestamp() {
	sort.Slice(t, func(i, j int) bool {
		return t[i].terminatedTime.Before(t[j].terminatedTime)
	})
}

func GetRestartdPods(configFlags *genericclioptions.ConfigFlags, namesapce, filterType string) (Pods, error) {
	client, _, err := getK8sClientAndConfig(configFlags)
	pods, err := getPodList(client, namesapce)
	if err != nil {
		return nil, err
	}
	var podsFilter PodsFilter
	switch filterType {
	case "oom":
		podsFilter = OOMFilter{}
	case "restart":
		podsFilter = RestartFilter{}
	default:
		podsFilter = RestartFilter{}
	}
	terminatedPods := filteringPods(pods, podsFilter)
	podInfo := convertPodsToPodInfo(terminatedPods)
	return podInfo, nil
}

// TerminatedPodsFilter is used to filter for pods that contain a terminated container, with an exit code of 137 (OOMKilled).
func filteringPods(pods []v1.Pod, podsFilter PodsFilter) []v1.Pod {

	var filterdPods []v1.Pod
	for _, pod := range pods {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			// The terminated state may be nil, i.e. not terminated, we must check this first.
			if podsFilter.Filter(containerStatus) {
				filterdPods = append(filterdPods, pod)
			}
		}
	}
	return filterdPods
}

func getPodList(client *kubernetes.Clientset, namespace string) ([]v1.Pod, error) {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	return pods.Items, nil
}

func convertPodsToPodInfo(pods []v1.Pod) []PodInfo {
	var podInfo []PodInfo
	for _, pod := range pods {
		for _, containerStatus := range pod.Status.ContainerStatuses {

			if terminated := containerStatus.LastTerminationState.Terminated; terminated != nil {
				pi := PodInfo{
					Pod: pod,
					// ContainerName:  containerStatus.Name,
					TerminatedTime: terminated.FinishedAt.String(),
					StartTime:      terminated.StartedAt.String(),
					terminatedTime: terminated.FinishedAt.Time,
					RestartCount:   strconv.FormatInt(int64(containerStatus.RestartCount), 10),
				}
				podInfo = append(podInfo, pi)
			}
		}
	}
	return podInfo
}

type RestartFilter struct{}

func (f RestartFilter) Filter(containerStatus v1.ContainerStatus) bool {
	return containerStatus.RestartCount > 0
}

type OOMFilter struct{}

func (f OOMFilter) Filter(containerStatus v1.ContainerStatus) bool {
	if terminated := containerStatus.LastTerminationState.Terminated; terminated != nil {
		return terminated.ExitCode == 137
	}
	return false
}
