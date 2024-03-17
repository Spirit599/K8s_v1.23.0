package iopriorty

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

var curDiskIO = map[string]float64{}

var curNetIO = map[string]float64{}

var minDisk float64 = 10000000
var maxDisk float64 = 0

var minNet float64 = 10000000
var maxNet float64 = 0

type IOPriorty struct {
	handle framework.Handle
}

var _ framework.PreScorePlugin = &IOPriorty{}
var _ framework.ScorePlugin = &IOPriorty{}
var _ framework.PreFilterPlugin = &IOPriorty{}

const Name = "IOPriorty"

func (pl *IOPriorty) Name() string {
	return Name
}

func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &IOPriorty{handle: h}, nil
}

func (pl *IOPriorty) PreFilter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod) *framework.Status {
	klog.Infof("PreFilter IOPriorty")
	return nil
}

func (pl *IOPriorty) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// PreScore
func (pl *IOPriorty) PreScore(
	pCtx context.Context,
	cycleState *framework.CycleState,
	pod *v1.Pod,
	nodes []*v1.Node,
) *framework.Status {

	getDiskIOData()
	getNetIOData()

	minDisk = 10000000
	maxDisk = 0

	minNet = 10000000
	maxNet = 0

	for key, value := range curDiskIO {
		if key != "master" {
			minDisk = MinFloat64(minDisk, value)
			maxDisk = MaxFloat64(maxDisk, value)
		}
	}

	for key, value := range curNetIO {
		if key != "master" {
			minNet = MinFloat64(minNet, value)
			maxNet = MaxFloat64(maxNet, value)
		}
	}

	return nil
}

// Score invoked at the score extension point.
func (pl *IOPriorty) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {

	// nodeInfo, err := pl.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	// if err != nil {
	// 	return 0, framework.AsStatus(fmt.Errorf("getting node %q from Snapshot: %w", nodeName, err))
	// }

	score1 := 100 * (1 - (curDiskIO[nodeName]-minDisk)/(maxDisk-minDisk))
	score2 := 100 * (1 - (curNetIO[nodeName]-minNet)/(maxNet-minNet))
	score := int64(score1+score2) / 2

	klog.Infof("nodename:%s minDisk:%f curDisk:%f maxDisk:%f", nodeName, minDisk, curDiskIO[nodeName], maxDisk)
	klog.Infof("nodename:%s minNet:%f curNet:%f maxNet:%f", nodeName, minNet, curNetIO[nodeName], maxNet)
	klog.Infof("nodename:%s score:%d", nodeName, score)

	return score, nil
}

// ScoreExtensions of the Score plugin.
func (pl *IOPriorty) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
