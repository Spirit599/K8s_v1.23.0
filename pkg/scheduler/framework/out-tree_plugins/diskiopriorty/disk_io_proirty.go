package diskiopriorty

import (
	"context"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

var allRate float64
var avgRate float64

var avgDiff float64
var minDiff float64
var maxDiff float64

var curDiskIO = map[string]float64{}

var oldRate = map[string]float64{}

var diff = map[string]float64{}

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

type DiskIOPriorty struct {
	handle framework.Handle
}

var _ framework.PreScorePlugin = &DiskIOPriorty{}
var _ framework.ScorePlugin = &DiskIOPriorty{}
var _ framework.PreFilterPlugin = &DiskIOPriorty{}

const Name = "DiskIOPriorty"

func (pl *DiskIOPriorty) Name() string {
	return Name
}

func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &DiskIOPriorty{handle: h}, nil
}

func (pl *DiskIOPriorty) PreFilter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod) *framework.Status {
	klog.Infof("PreFilter DiskIOPriorty")
	return nil
}

func (pl *DiskIOPriorty) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// PreScore
func (pl *DiskIOPriorty) PreScore(
	pCtx context.Context,
	cycleState *framework.CycleState,
	pod *v1.Pod,
	nodes []*v1.Node,
) *framework.Status {

	getDiskIOData()

	totalNumNodes := float64(len(curDiskIO) - 1)

	allRate = 0
	avgDiff = 0
	maxDiff = 0
	minDiff = 100000000000

	for key, value := range curDiskIO {
		if key != "master" {
			rate := value / maxLimitDiskIO[key]
			klog.Infof("%s, %f, %f, %f", key, value, maxLimitDiskIO[key], rate)
			oldRate[key] = rate
			allRate += rate
		}
	}

	avgRate = allRate / totalNumNodes
	klog.Infof("avgRate %f nodesNum %d", avgRate, len(nodes))

	reqIO, _ := strconv.ParseFloat(pod.Annotations["DiskIO"], 64)
	klog.Infof("reqIO %f", reqIO)

	for key, value := range curDiskIO {
		if key != "master" {
			rate := (value + reqIO) / maxLimitDiskIO[key]
			klog.Infof("%s, %f, %f, %f", key, value+reqIO, maxLimitDiskIO[key], rate)
			newAllRate := allRate + rate - oldRate[key]
			newAvgRate := newAllRate / totalNumNodes

			curDiff := (newAvgRate - rate) * (newAvgRate - rate)

			for key2, _ := range curDiskIO {
				if key2 != "master" && key2 != key {
					curDiff += (newAvgRate - oldRate[key2]) * (newAvgRate - oldRate[key2])
				}
			}
			curDiff = curDiff / totalNumNodes
			diff[key] = curDiff
			avgDiff += curDiff
			maxDiff = MaxFloat64(maxDiff, curDiff)
			minDiff = MinFloat64(minDiff, curDiff)

			klog.Infof("nodename %s diff[key] %f", key, curDiff)
		}
	}

	avgDiff = avgDiff / totalNumNodes

	klog.Infof("maxDiff %f avgDiff %f minDiff %f", maxDiff, avgDiff, minDiff)

	return nil
}

// Score invoked at the score extension point.
func (pl *DiskIOPriorty) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {

	// nodeInfo, err := pl.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	// if err != nil {
	// 	return 0, framework.AsStatus(fmt.Errorf("getting node %q from Snapshot: %w", nodeName, err))
	// }
	score := 100 - 100*(diff[nodeName]-minDiff)/(maxDiff-minDiff)
	klog.Infof("nodename %s score %f", nodeName, score)

	return 0, nil
}

// ScoreExtensions of the Score plugin.
func (pl *DiskIOPriorty) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
