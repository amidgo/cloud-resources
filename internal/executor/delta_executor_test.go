package executor_test

import (
	"context"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/executor"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"github.com/amidgo/cloud-resources/internal/resourcemanager"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

type FloatMatcher struct {
	Value float32
}

func (f FloatMatcher) Matches(x interface{}) bool {
	value, ok := x.(float32)
	if !ok {
		return false
	}
	return float32Equal(value, f.Value)
}

func (f FloatMatcher) String() string {
	return fmt.Sprintf("is equal %f", f.Value)
}

type ResourceManagerCall struct {
	Value, Limit float32
	Report       resourcemanager.ResourceReport
	Err          error
}

func (r ResourceManagerCall) Call() bool {
	return r != ResourceManagerCall{}
}

func Test_DeltaExecutor_ExecuteDelta(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	logger := executor.ResourceTypeLoggerWrapper{
		Logger:       log.Default(),
		ResourceType: resourcetype.DB,
	}
	resourceManager := resourcemanager.NewMockResourceManager(ctrl)
	cases := []struct {
		delta executor.Delta

		addCpu, removeCpu, addRam, removeRam ResourceManagerCall

		err error
	}{
		{
			delta: executor.Delta{
				CPU: 100,
				RAM: 45,
			},
			addCpu: ResourceManagerCall{
				Value:  100,
				Report: resourcemanager.ResourceReport{CPU: 100, RAM: 100},
			},
			removeRam: ResourceManagerCall{
				Value:  -55,
				Limit:  0,
				Report: resourcemanager.ResourceReport{CPU: 0, RAM: 0},
			},
		},
		{
			delta: executor.Delta{
				CPU: -100,
				RAM: -19,
			},
			removeCpu: ResourceManagerCall{
				Value:  -100,
				Limit:  -19,
				Report: resourcemanager.ResourceReport{CPU: -19, RAM: -19},
			},
		},
		{
			delta: executor.Delta{
				CPU: -19,
				RAM: -100,
			},
			removeCpu: ResourceManagerCall{
				Value:  -19,
				Limit:  -100,
				Report: resourcemanager.ResourceReport{CPU: -19, RAM: -19},
			},
			removeRam: ResourceManagerCall{
				Value:  -81,
				Limit:  0,
				Report: resourcemanager.ResourceReport{CPU: 0, RAM: 0},
			},
		},
	}

	deltaExecutor := executor.NewDeltaExecutor(logger, resourceManager)

	for _, cs := range cases {
		if cs.addCpu.Call() {
			resourceManager.EXPECT().AddCPU(gomock.Any(), cs.addCpu.Value).Return(cs.addCpu.Report, cs.addCpu.Err).Times(1)
		}
		if cs.addRam.Call() {
			resourceManager.EXPECT().AddRAM(gomock.Any(), cs.addRam.Value).Return(cs.addRam.Report, cs.addRam.Err).Times(1)
		}
		if cs.removeCpu.Call() {
			resourceManager.EXPECT().RemoveCPU(gomock.Any(), cs.removeCpu.Value, cs.removeCpu.Limit).Return(cs.removeCpu.Report, cs.removeCpu.Err).Times(1)
		}
		if cs.removeRam.Call() {
			resourceManager.EXPECT().RemoveRAM(gomock.Any(), cs.removeRam.Value, cs.removeRam.Limit).Return(cs.removeRam.Report, cs.removeRam.Err).Times(1)
		}
		err := deltaExecutor.ExecuteDelta(ctx, cs.delta)
		assert.ErrorIs(t, err, cs.err, "wrong err of execute delta")
	}
}
