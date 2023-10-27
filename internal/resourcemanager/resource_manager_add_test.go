package resourcemanager_test

import (
	"context"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"github.com/amidgo/cloud-resources/internal/resourcemanager"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func Test_AddCPU(t *testing.T) {
	ctx := context.Background()
	tp := resourcetype.VM
	priceManager := PriceManager{
		resourceType: tp,
	}

	cases := []struct {
		resources             []*resourcemodel.Resource
		resourceCall          bool
		addCount, updateCount int
		cpu                   float32

		maxMachineCount int
		report          resourcemanager.ResourceReport
	}{
		{
			addCount: 10,
			cpu:      10,
			report: resourcemanager.ResourceReport{
				CPU: 10,
				RAM: 10,
			},
			maxMachineCount: 10,
			resourceCall:    true,
		},
		{
			cpu: 0,
		},
		{
			cpu: -1,
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MaxResource(tp, false),
			},
			maxMachineCount: 10,
			cpu:             28,
			addCount:        2,
			updateCount:     3,
			report: resourcemanager.ResourceReport{
				CPU: 29,
				RAM: 29,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, false),
				MaxResource(tp, false),
			},
			maxMachineCount: 10,
			cpu:             28.75,
			updateCount:     1,
			addCount:        2,
			report: resourcemanager.ResourceReport{
				CPU: 11,
				RAM: 11,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, false),
				MaxResource(tp, true),
			},
			maxMachineCount: 10,
			cpu:             28.75,
			addCount:        2,
			report: resourcemanager.ResourceReport{
				CPU: 2,
				RAM: 2,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			cpu:             28.75,
			updateCount:     1,
			addCount:        2,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 11,
				RAM: 11,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
			},
			cpu:             9,
			addCount:        9,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 9,
				RAM: 9,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
			},
			cpu:             18,
			addCount:        9,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 9,
				RAM: 9,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
			},
			cpu:             19,
			addCount:        9,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 9,
				RAM: 9,
			},
		},

		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			cpu:             19,
			addCount:        7,
			updateCount:     2,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 25,
				RAM: 25,
			},
		},
	}

	for _, cs := range cases {
		storage := MockResourceStorage(t, cs.addCount, 0, cs.updateCount)
		if cs.resourceCall {
			storage.EXPECT().ResourceList(gomock.Any()).Return(cs.resources, nil).Times(1)
		}
		resourceManager := resourcemanager.NewResourceManager(storage, &priceManager, cs.maxMachineCount)
		rep, err := resourceManager.AddCPU(ctx, cs.cpu)

		assert.NilError(t, err, "failed add cpu")
		assert.Equal(t, rep, cs.report)
	}
}

func Test_AddRAM(t *testing.T) {
	ctx := context.Background()
	tp := resourcetype.VM
	priceManager := PriceManager{
		resourceType: tp,
	}

	cases := []struct {
		resources             []*resourcemodel.Resource
		resourceCall          bool
		addCount, updateCount int
		ram                   float32

		maxMachineCount int
		report          resourcemanager.ResourceReport
	}{
		{
			addCount: 10,
			ram:      10,
			report: resourcemanager.ResourceReport{
				CPU: 10,
				RAM: 10,
			},
			maxMachineCount: 10,
			resourceCall:    true,
		},
		{
			ram: 0,
		},
		{
			ram: -1,
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MaxResource(tp, false),
			},
			maxMachineCount: 10,
			ram:             28,
			addCount:        2,
			updateCount:     3,
			report: resourcemanager.ResourceReport{
				CPU: 29,
				RAM: 29,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, false),
				MaxResource(tp, false),
			},
			maxMachineCount: 10,
			ram:             28.75,
			updateCount:     1,
			addCount:        2,
			report: resourcemanager.ResourceReport{
				CPU: 11,
				RAM: 11,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, false),
				MaxResource(tp, true),
			},
			maxMachineCount: 10,
			ram:             28.75,
			addCount:        2,
			report: resourcemanager.ResourceReport{
				CPU: 2,
				RAM: 2,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			ram:             28.75,
			updateCount:     1,
			addCount:        2,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 11,
				RAM: 11,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
			},
			ram:             9,
			addCount:        9,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 9,
				RAM: 9,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
			},
			ram:             18,
			addCount:        9,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 9,
				RAM: 9,
			},
		},
		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
			},
			ram:             19,
			addCount:        9,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 9,
				RAM: 9,
			},
		},

		{
			resourceCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			ram:             19,
			addCount:        7,
			updateCount:     2,
			maxMachineCount: 10,
			report: resourcemanager.ResourceReport{
				CPU: 25,
				RAM: 25,
			},
		},
	}

	for _, cs := range cases {
		storage := MockResourceStorage(t, cs.addCount, 0, cs.updateCount)
		if cs.resourceCall {
			storage.EXPECT().ResourceList(gomock.Any()).Return(cs.resources, nil).Times(1)
		}
		resourceManager := resourcemanager.NewResourceManager(storage, &priceManager, cs.maxMachineCount)
		rep, err := resourceManager.AddRAM(ctx, cs.ram)

		assert.NilError(t, err, "failed add cpu")
		assert.Equal(t, rep, cs.report)
	}
}
