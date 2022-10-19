package demo

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type Demo struct{}

var _ framework.PreFilterPlugin = &Demo{}

const Name = "Demo"

func (pl *Demo) Name() string {
	return Name
}

func New(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &Demo{}, nil
}

func (pl *Demo) PreFilter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod) *framework.Status {
	klog.Infof("PreFilter Demo")
	return nil
}

func (pl *Demo) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}
