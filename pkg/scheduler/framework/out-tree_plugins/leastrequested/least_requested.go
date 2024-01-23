package leastrequested

import (
	"context"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

var curDiskIO = map[string]float64{}
var maxLimitDiskIO = map[string]float64{
	"master": 0,
	"node1":  500,
	"node2":  500,
	"node3":  500,
	"node4":  500,
	"node5":  500,
	"node6":  500,
	"node7":  500,
	"node8":  500,
	"node9":  500,
	"node10": 500,
}

var curNetIO = map[string]float64{}
var maxLimitNetIO = map[string]float64{
	"master": 0,
	"node1":  500,
	"node2":  500,
	"node3":  500,
	"node4":  500,
	"node5":  500,
	"node6":  500,
	"node7":  500,
	"node8":  500,
	"node9":  500,
	"node10": 500,
}

var curCpuRate = map[string]float64{}
var curMemRate = map[string]float64{}

var nodeScore = map[string]float64{}

var weights []float64 = []float64{1, 1, 1, 1}

type LeastRequested struct {
	handle framework.Handle
}

var _ framework.PreScorePlugin = &LeastRequested{}
var _ framework.ScorePlugin = &LeastRequested{}
var _ framework.PreFilterPlugin = &LeastRequested{}

const Name = "LeastRequested"

func (pl *LeastRequested) Name() string {
	return Name
}

func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &LeastRequested{handle: h}, nil
}

func (pl *LeastRequested) PreFilter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod) *framework.Status {
	klog.Infof("PreFilter LeastRequested")
	return nil
}

func (pl *LeastRequested) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// PreScore
func (pl *LeastRequested) PreScore(
	pCtx context.Context,
	cycleState *framework.CycleState,
	pod *v1.Pod,
	nodes []*v1.Node,
) *framework.Status {

	getResourceData()

	reqDiskIO, _ := strconv.ParseFloat(pod.Annotations["DiskIO"], 64)
	reqNetIO, _ := strconv.ParseFloat(pod.Annotations["NetIO"], 64)
	klog.Infof("reqDiskIO:%f reqNetIO:%f", reqDiskIO, reqNetIO)

	var matrix [][]float64
	var idxToName []string

	for _, node := range nodes {
		name := node.Name
		nodeCpuRate := curCpuRate[name]
		nodeMenRate := curMemRate[name]
		nodeDiskRate := (curDiskIO[name] + reqDiskIO) / maxLimitDiskIO[name]
		nodeNetRate := (curNetIO[name] + reqNetIO) / maxLimitNetIO[name]
		klog.Infof("name %s cpu:%f men:%f diskio:%f netio:%f", name, nodeCpuRate, nodeMenRate, nodeDiskRate, nodeNetRate)
		matrix = append(matrix, []float64{nodeCpuRate, nodeMenRate, nodeDiskRate, nodeNetRate})
		idxToName = append(idxToName, name)
	}

	score := topsis(matrix, weights)
	for i, s := range score {
		klog.Infof("name:%s score:%f", idxToName[i], s)
		nodeScore[idxToName[i]] = s
	}

	return nil
}

// Score invoked at the score extension point.
func (pl *LeastRequested) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {

	// nodeInfo, err := pl.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	// if err != nil {
	// 	return 0, framework.AsStatus(fmt.Errorf("getting node %q from Snapshot: %w", nodeName, err))
	// }

	return 0, nil
}

// ScoreExtensions of the Score plugin.
func (pl *LeastRequested) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
