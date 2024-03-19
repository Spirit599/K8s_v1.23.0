package balancedallocation

import (
	"context"
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	schedutil "k8s.io/kubernetes/pkg/scheduler/util"
)

var curDiskIO = map[string]float64{}
var maxLimitDiskIO = map[string]float64{
	"master": 0,
	"node1":  30,
	"node2":  30,
	"node3":  30,
	"node4":  30,
	"node5":  30,
	"node6":  30,
	"node7":  30,
	"node8":  30,
	"node9":  30,
	"node10": 30,
}

var curNetIO = map[string]float64{}
var maxLimitNetIO = map[string]float64{
	"master": 0,
	"node1":  1,
	"node2":  1,
	"node3":  1,
	"node4":  1,
	"node5":  1,
	"node6":  1,
	"node7":  1,
	"node8":  1,
	"node9":  1,
	"node10": 1,
}

var curCpuRate = map[string]float64{}
var curMemRate = map[string]float64{}

var nodeScore = map[string]float64{}

// var weights []float64 = []float64{1, 1, 1, 1}

type BalancedAllocation struct {
	handle framework.Handle
}

var _ framework.PreScorePlugin = &BalancedAllocation{}
var _ framework.ScorePlugin = &BalancedAllocation{}
var _ framework.PreFilterPlugin = &BalancedAllocation{}

const Name = "BalancedAllocation"

func (pl *BalancedAllocation) Name() string {
	return Name
}

func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &BalancedAllocation{handle: h}, nil
}

func (pl *BalancedAllocation) calculatePodResourceRequest(pod *v1.Pod, resource v1.ResourceName) int64 {
	var podRequest int64
	for i := range pod.Spec.Containers {
		container := &pod.Spec.Containers[i]
		value := schedutil.GetRequestForResource(resource, &container.Resources.Requests, false)
		podRequest += value
	}

	for i := range pod.Spec.InitContainers {
		initContainer := &pod.Spec.InitContainers[i]
		value := schedutil.GetRequestForResource(resource, &initContainer.Resources.Requests, false)
		if podRequest < value {
			podRequest = value
		}
	}

	return podRequest
}

func (pl *BalancedAllocation) PreFilter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod) *framework.Status {
	klog.Infof("PreFilter BalancedAllocation")
	return nil
}

func (pl *BalancedAllocation) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// PreScore
func (pl *BalancedAllocation) PreScore(
	pCtx context.Context,
	cycleState *framework.CycleState,
	pod *v1.Pod,
	nodes []*v1.Node,
) *framework.Status {

	getResourceData()

	reqDiskIO, _ := strconv.ParseFloat(pod.Annotations["DiskIO"], 64)
	reqNetIO, _ := strconv.ParseFloat(pod.Annotations["NetIO"], 64)
	reqNetIO = reqNetIO / 1000
	klog.Infof("reqDiskIO:%f reqNetIO:%f", reqDiskIO, reqNetIO)

	var RequestedMatrix [][]float64
	var NeedMatrix [][]float64
	var idxToName []string

	for _, node := range nodes {
		name := node.Name
		nodeCpuRate := curCpuRate[name]
		nodeMenRate := curMemRate[name]
		nodeDiskRate := curDiskIO[name] / maxLimitDiskIO[name]
		nodeNetRate := curNetIO[name] / maxLimitNetIO[name]
		klog.Infof("name %s cpu:%f men:%f diskio:%f netio:%f", name, nodeCpuRate, nodeMenRate, nodeDiskRate, nodeNetRate)
		RequestedMatrix = append(RequestedMatrix, []float64{nodeCpuRate, nodeMenRate, nodeDiskRate, nodeNetRate})

		nodeInfo, err := pl.handle.SnapshotSharedLister().NodeInfos().Get(name)
		if err != nil {
			return framework.AsStatus(fmt.Errorf("getting node %q from Snapshot: %w", name, err))
		}

		podCpu := (float64)(pl.calculatePodResourceRequest(pod, "cpu"))
		podMen := (float64)(pl.calculatePodResourceRequest(pod, "memory"))

		podCpuRate := (podCpu) / float64(nodeInfo.Allocatable.MilliCPU)
		podMenRate := (podMen) / float64(nodeInfo.Allocatable.Memory)
		podDiskRate := reqDiskIO / maxLimitDiskIO[name]
		podNetRate := reqNetIO / maxLimitNetIO[name]
		klog.Infof("name %s cpu:%f men:%f diskio:%f netio:%f", pod.Name, podCpuRate, podMenRate, podDiskRate, podNetRate)
		NeedMatrix = append(NeedMatrix, []float64{podCpuRate, podMenRate, podDiskRate, podNetRate})

		idxToName = append(idxToName, name)
	}

	scores := pearson(RequestedMatrix, NeedMatrix)
	klog.Infof("scores_len:%d", len(scores))
	for i, s := range scores {
		klog.Infof("name:%s score:%f", idxToName[i], (1-s)/2)
		nodeScore[idxToName[i]] = (1 - s) / 2
	}

	return nil
}

// Score invoked at the score extension point.
func (pl *BalancedAllocation) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {

	// nodeInfo, err := pl.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	// if err != nil {
	// 	return 0, framework.AsStatus(fmt.Errorf("getting node %q from Snapshot: %w", nodeName, err))
	// }

	return int64(nodeScore[nodeName]), nil
}

// ScoreExtensions of the Score plugin.
func (pl *BalancedAllocation) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
