package resourcemanager_test

import (
	"context"
	"github.com/amidgo/cloud-resources/internal/model/pricemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"github.com/amidgo/cloud-resources/internal/resourcemanager"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

type PriceManager struct {
	resourceType resourcetype.ResourceType
}

func (p *PriceManager) MinCPU() pricemodel.Price {
	return pricemodel.NewPrice(1, 1, 1, 1, "", p.resourceType)
}
func (p *PriceManager) MinRAM() pricemodel.Price {
	return pricemodel.NewPrice(1, 1, 1, 1, "", p.resourceType)
}
func (p *PriceManager) MaxCPU() pricemodel.Price {
	return pricemodel.NewPrice(2, 8, 10, 10, "", p.resourceType)
}
func (p *PriceManager) MaxRAM() pricemodel.Price {
	return pricemodel.NewPrice(2, 8, 10, 10, "", p.resourceType)
}

func MockResourceStorage(t *testing.T, addCount, deleteCount, updateCount int) *resourcemanager.MockResourceStorage {
	ctrl := gomock.NewController(t)
	storage := resourcemanager.NewMockResourceStorage(ctrl)

	storage.EXPECT().AddResource(gomock.Any(), gomock.Any()).Return(resourcemodel.Resource{}, nil).Times(addCount)
	storage.EXPECT().DeleteResource(gomock.Any(), gomock.Any()).Return(nil).Times(deleteCount)
	storage.EXPECT().UpdateResource(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(updateCount)
	return storage
}

func MinResource(resourceType resourcetype.ResourceType, failed bool) *resourcemodel.Resource {
	return &resourcemodel.Resource{
		ID:     rand.Int(),
		CPU:    1,
		RAM:    1,
		Type:   resourceType,
		Failed: failed,
	}
}

func MaxResource(resourceType resourcetype.ResourceType, failed bool) *resourcemodel.Resource {
	return &resourcemodel.Resource{
		ID:     rand.Int(),
		CPU:    10,
		RAM:    10,
		Type:   resourceType,
		Failed: failed,
	}
}

func Test_RemoveCPU(t *testing.T) {
	ctx := context.Background()
	tp := resourcetype.VM
	priceManager := PriceManager{
		resourceType: tp,
	}
	cases := []struct {
		resourcesCall bool
		resources     []*resourcemodel.Resource

		cpu         float32
		ramLimit    float32
		deleteCount int

		report resourcemanager.ResourceReport
	}{
		{
			cpu: -1,
		},
		{
			ramLimit: -1,
		},
		{
			cpu:      0,
			ramLimit: 0,
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			cpu:         -5,
			ramLimit:    -10,
			deleteCount: 5,
			report: resourcemanager.ResourceReport{
				CPU: -5,
				RAM: -5,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			cpu:         -5,
			ramLimit:    -10,
			deleteCount: 4,
			report: resourcemanager.ResourceReport{
				CPU: -4,
				RAM: -4,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			cpu:         -5,
			ramLimit:    -2,
			deleteCount: 2,
			report: resourcemanager.ResourceReport{
				CPU: -2,
				RAM: -2,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, true),
				MinResource(tp, true),
			},
			cpu:         -5,
			ramLimit:    -50,
			deleteCount: 2,
			report: resourcemanager.ResourceReport{
				CPU: -2,
				RAM: -2,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MaxResource(tp, false),
				MinResource(tp, true),
				MinResource(tp, true),
				MaxResource(tp, false),
			},
			cpu:         -21,
			ramLimit:    -50,
			deleteCount: 3,
			report: resourcemanager.ResourceReport{
				CPU: -21,
				RAM: -21,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MaxResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, false),
			},
			cpu:         -6,
			ramLimit:    -50,
			deleteCount: 3,
			report: resourcemanager.ResourceReport{
				CPU: -3,
				RAM: -3,
			},
		},
	}

	for _, cs := range cases {
		storage := MockResourceStorage(t, 0, cs.deleteCount, 0)
		if cs.resourcesCall {
			storage.EXPECT().ResourceList(gomock.Any()).Return(cs.resources, nil).Times(1)
		}
		resourceManager := resourcemanager.NewResourceManager(storage, &priceManager, 0)
		rep, err := resourceManager.RemoveCPU(ctx, cs.cpu, cs.ramLimit)
		assert.NilError(t, err, "failed add cpu")
		assert.Equal(t, rep, cs.report)
	}
}

func Test_RemoveRAM(t *testing.T) {
	ctx := context.Background()
	tp := resourcetype.VM
	priceManager := PriceManager{
		resourceType: tp,
	}
	cases := []struct {
		resourcesCall bool
		resources     []*resourcemodel.Resource

		ram         float32
		cpuLimit    float32
		deleteCount int

		report resourcemanager.ResourceReport
	}{
		{
			ram: -1,
		},
		{
			cpuLimit: -1,
		},
		{
			ram:      0,
			cpuLimit: 0,
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			ram:         -5,
			cpuLimit:    -10,
			deleteCount: 5,
			report: resourcemanager.ResourceReport{
				CPU: -5,
				RAM: -5,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			ram:         -5,
			cpuLimit:    -10,
			deleteCount: 4,
			report: resourcemanager.ResourceReport{
				CPU: -4,
				RAM: -4,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
			},
			ram:         -5,
			cpuLimit:    -2,
			deleteCount: 2,
			report: resourcemanager.ResourceReport{
				CPU: -2,
				RAM: -2,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, true),
				MinResource(tp, true),
			},
			ram:         -5,
			cpuLimit:    -50,
			deleteCount: 2,
			report: resourcemanager.ResourceReport{
				CPU: -2,
				RAM: -2,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MaxResource(tp, false),
				MinResource(tp, true),
				MinResource(tp, true),
				MaxResource(tp, false),
			},
			ram:         -21,
			cpuLimit:    -50,
			deleteCount: 3,
			report: resourcemanager.ResourceReport{
				CPU: -21,
				RAM: -21,
			},
		},
		{
			resourcesCall: true,
			resources: []*resourcemodel.Resource{
				MinResource(tp, false),
				MinResource(tp, false),
				MinResource(tp, false),
				MaxResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, true),
				MinResource(tp, false),
			},
			ram:         -6,
			cpuLimit:    -50,
			deleteCount: 3,
			report: resourcemanager.ResourceReport{
				CPU: -3,
				RAM: -3,
			},
		},
	}

	for _, cs := range cases {
		storage := MockResourceStorage(t, 0, cs.deleteCount, 0)
		if cs.resourcesCall {
			storage.EXPECT().ResourceList(gomock.Any()).Return(cs.resources, nil).Times(1)
		}
		resourceManager := resourcemanager.NewResourceManager(storage, &priceManager, 0)
		rep, err := resourceManager.RemoveRAM(ctx, cs.ram, cs.cpuLimit)
		assert.NilError(t, err, "failed add cpu")
		assert.Equal(t, rep, cs.report)
	}
}
