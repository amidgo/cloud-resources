package executor_test

import (
	"context"
	"database/sql"
	"github.com/amidgo/cloud-resources/internal/executor"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"github.com/amidgo/cloud-resources/internal/model/statiscticsmodel"
	"github.com/amidgo/cloud-resources/internal/resourcemanager"
	"log"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_ExecuteInitter_Init(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	resourceManager := resourcemanager.NewMockResourceManager(ctrl)
	resourceStorage := executor.NewMockResourceStorage(ctrl)

	initCpu := float32(10)
	initRam := float32(9)
	initter := executor.NewInitter(resourceManager, resourceStorage, initCpu, initRam)

	cases := []struct {
		addCpu struct {
			cpu    float32
			report resourcemanager.ResourceReport
			err    error
		}
		addRam struct {
			ram    float32
			report resourcemanager.ResourceReport
			err    error
		}
		getResources struct {
			resources []*resourcemodel.Resource
			err       error
		}

		err error
	}{
		{
			addCpu: struct {
				cpu    float32
				report resourcemanager.ResourceReport
				err    error
			}{
				cpu: initCpu,
			},
			addRam: struct {
				ram    float32
				report resourcemanager.ResourceReport
				err    error
			}{
				ram: initRam,
			},
			getResources: struct {
				resources []*resourcemodel.Resource
				err       error
			}{
				err: sql.ErrNoRows,
			},
			err: sql.ErrNoRows,
		},
		{
			addCpu: struct {
				cpu    float32
				report resourcemanager.ResourceReport
				err    error
			}{
				cpu: initCpu,
			},
			addRam: struct {
				ram    float32
				report resourcemanager.ResourceReport
				err    error
			}{
				ram: initRam,
			},
			getResources: struct {
				resources []*resourcemodel.Resource
				err       error
			}{
				resources: []*resourcemodel.Resource{
					NewResourceCpuRam(10, 10, false),
					NewResourceCpuRam(10, 10, false),
					NewResourceCpuRam(10, 10, false),
					NewResourceCpuRam(10, 10, false),
					NewResourceCpuRam(10, 10, false),
				},
			},
		},
	}

	for _, cs := range cases {
		resourceManager.EXPECT().AddCPU(gomock.Any(), cs.addCpu.cpu).Return(cs.addCpu.report, cs.addCpu.err).Times(1)
		resourceManager.EXPECT().AddRAM(gomock.Any(), cs.addRam.ram).Return(cs.addRam.report, cs.addRam.err).Times(1)
		resourceStorage.EXPECT().ResourceList(gomock.Any()).Return(cs.getResources.resources, cs.getResources.err)

		resources, err := initter.Init(ctx)
		assert.DeepEqual(t, resources, cs.getResources.resources)
		assert.ErrorIs(t, err, cs.err, "wrong err")
	}
}

func Test_Executor_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	logger := executor.ResourceTypeLoggerWrapper{
		ResourceType: resourcetype.DB,
		Logger:       log.Default(),
	}
	initter := executor.NewMockExecuteInitter(ctrl)
	resourceStorage := executor.NewMockResourceStorage(ctrl)
	statisticsStorage := executor.NewMockStatiscticsStorage(ctrl)
	deltaExecutor := executor.NewMockDeltaExecutor(ctrl)
	deltaCounter := executor.NewMockDeltaCounter(ctrl)

	ex := executor.NewExecutor(logger, initter, resourceStorage, statisticsStorage, deltaExecutor, deltaCounter)

	cases := []struct {
		statisticsCall struct {
			call       bool
			statistics statiscticsmodel.Statistics
			err        error
		}
		resourcesCall struct {
			call      bool
			resources []*resourcemodel.Resource
			err       error
		}
		initterCall struct {
			call      bool
			resources []*resourcemodel.Resource
			err       error
		}
		deltaCounterCall struct {
			delta executor.Delta
			call  bool
		}
		executeDeltaCall struct {
			call bool
			err  error
		}
		err error
	}{
		{
			statisticsCall: struct {
				call       bool
				statistics statiscticsmodel.Statistics
				err        error
			}{
				call: true,
				err:  sql.ErrConnDone,
			},
			err: sql.ErrConnDone,
		},
		{
			statisticsCall: struct {
				call       bool
				statistics statiscticsmodel.Statistics
				err        error
			}{
				call: true,
			},
			resourcesCall: struct {
				call      bool
				resources []*resourcemodel.Resource
				err       error
			}{
				call: true,
				err:  sql.ErrNoRows,
			},
			err: sql.ErrNoRows,
		},

		{
			statisticsCall: struct {
				call       bool
				statistics statiscticsmodel.Statistics
				err        error
			}{
				call: true,
			},
			resourcesCall: struct {
				call      bool
				resources []*resourcemodel.Resource
				err       error
			}{
				call: true,
			},
			initterCall: struct {
				call      bool
				resources []*resourcemodel.Resource
				err       error
			}{
				call: true,
				err:  sql.ErrNoRows,
			},
			err: sql.ErrNoRows,
		},

		{
			statisticsCall: struct {
				call       bool
				statistics statiscticsmodel.Statistics
				err        error
			}{
				call: true,
			},
			resourcesCall: struct {
				call      bool
				resources []*resourcemodel.Resource
				err       error
			}{
				resources: []*resourcemodel.Resource{
					NewResourceCpuRam(10, 10, false),
					NewResourceCpuRam(10, 10, true),
					NewResourceCpuRam(10, 10, false),
				},
				call: true,
			},
			deltaCounterCall: struct {
				delta executor.Delta
				call  bool
			}{
				call: true,
				delta: executor.Delta{
					CPU: 100,
					RAM: 100,
				},
			},
			executeDeltaCall: struct {
				call bool
				err  error
			}{
				call: true,
				err:  http.ErrBodyNotAllowed,
			},
			err: http.ErrBodyNotAllowed,
		},

		{
			statisticsCall: struct {
				call       bool
				statistics statiscticsmodel.Statistics
				err        error
			}{
				call: true,
			},
			resourcesCall: struct {
				call      bool
				resources []*resourcemodel.Resource
				err       error
			}{
				call: true,
				resources: []*resourcemodel.Resource{
					NewResourceCpuRam(10, 10, false),
					NewResourceCpuRam(10, 10, false),
					NewResourceCpuRam(10, 10, true),
					NewResourceCpuRam(10, 10, false),
				},
			},
			deltaCounterCall: struct {
				delta executor.Delta
				call  bool
			}{
				call: true,
				delta: executor.Delta{
					CPU: 100,
					RAM: 100,
				},
			},
			executeDeltaCall: struct {
				call bool
				err  error
			}{
				call: true,
			},
		},
	}

	for _, cs := range cases {
		if cs.statisticsCall.call {
			statisticsStorage.EXPECT().Statisctics(gomock.Any()).Return(cs.statisticsCall.statistics, cs.statisticsCall.err).Times(1)
		}
		if cs.resourcesCall.call {
			resourceStorage.EXPECT().ResourceList(gomock.Any()).Return(cs.resourcesCall.resources, cs.resourcesCall.err).Times(1)
		}
		if cs.initterCall.call {
			initter.EXPECT().Init(gomock.Any()).Return(cs.initterCall.resources, cs.initterCall.err).Times(1)
		}
		if cs.deltaCounterCall.call {
			deltaCounter.EXPECT().Delta(cs.statisticsCall.statistics, cs.resourcesCall.resources).Return(cs.deltaCounterCall.delta).Times(1)
		}
		if cs.executeDeltaCall.call {
			deltaExecutor.EXPECT().ExecuteDelta(gomock.Any(), cs.deltaCounterCall.delta).Return(cs.executeDeltaCall.err).Times(1)
		}
		err := ex.Execute(ctx)
		assert.ErrorIs(t, err, cs.err, "wrong err")
	}

}
