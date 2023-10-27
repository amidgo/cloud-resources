package executor_test

import (
	"github.com/amidgo/cloud-resources/internal/executor"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"github.com/amidgo/cloud-resources/internal/model/statiscticsmodel"
	"log"
	"math"
	"testing"

	"gotest.tools/v3/assert"
)

func float32Equal(x, y float32) bool {
	return math.Abs(float64(x-y)) < 0.01
}

func NewResourceCpuRam(cpu float32, ram float32, failed bool) *resourcemodel.Resource {
	return &resourcemodel.Resource{
		CPU:    cpu,
		RAM:    ram,
		Failed: failed,
	}
}

func Test_DeltaCalculator_ReduceResources(t *testing.T) {
	cases := []struct {
		resources []*resourcemodel.Resource

		cpuPotential, cpuCurrent, ramPotential, ramCurrent float32
	}{
		{
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(12, 10, false),
				NewResourceCpuRam(12, 10, true),
				NewResourceCpuRam(1, 2.2, false),
				NewResourceCpuRam(2.2, 2.2, true),
			},
			cpuPotential: 27.2,
			cpuCurrent:   13,
			ramPotential: 24.4,
			ramCurrent:   12.2,
		},
	}

	for _, cs := range cases {
		calc := executor.DeltaCalculator{}
		cpu, ram := calc.ReduceResourcesCpuRam(cs.resources)

		assert.Equal(t, true, float32Equal(cpu.Potential, cs.cpuPotential), "p cpu not equal")
		assert.Equal(t, true, float32Equal(ram.Potential, cs.ramPotential), "p ram not equal")
		assert.Equal(t, true, float32Equal(cpu.Current, cs.cpuCurrent), "c cpu not equal")
		assert.Equal(t, true, float32Equal(ram.Current, cs.ramCurrent), "c ram not equal")
	}

}

func Test_ResourceStatisticsCounter(t *testing.T) {
	cases := []struct {
		load, current, potential float32
		offload                  float32
		k                        float32
		minHealthLoadDelta       float32
		maxHealthLoadDelta       float32
		required                 float32
		delta                    float32
	}{
		{
			load:               159,
			offload:            1,
			current:            29,
			potential:          90,
			minHealthLoadDelta: 5,
			maxHealthLoadDelta: 100,

			k: 0.76,

			required: 60.67,
			delta:    -29.3289,
		},
		{
			load:               18.9,
			current:            189,
			offload:            1,
			potential:          200,
			k:                  0.89,
			minHealthLoadDelta: 5,
			maxHealthLoadDelta: 100,

			required: 40.721,
			delta:    -159.279,
		},

		{
			load:               80,
			current:            10,
			potential:          10,
			offload:            1,
			k:                  0.8,
			minHealthLoadDelta: 5,
			maxHealthLoadDelta: 100,

			required: 13,
			delta:    3,
		},
		{
			load:               80,
			current:            60,
			potential:          60,
			offload:            1,
			k:                  0.8,
			minHealthLoadDelta: 5,
			maxHealthLoadDelta: 100,
			required:           60,
			delta:              0,
		},
		{
			load:               58.1,
			current:            60,
			potential:          60,
			offload:            1,
			k:                  0.8,
			minHealthLoadDelta: 5,
			maxHealthLoadDelta: 100,
			required:           43.575,
			delta:              -16.425,
		},

		{
			load:               90.91,
			current:            1000,
			potential:          1000,
			offload:            0.95,
			k:                  0.8,
			minHealthLoadDelta: 1,
			maxHealthLoadDelta: 70,
			required:           933.645,
			delta:              -66.355,
		},
	}

	for _, cs := range cases {
		counter := executor.NewResourceStatisticCounter(cs.load, cs.current, cs.potential)

		deltaLoad := executor.ResourceDeltaLoad{
			Load:     cs.k,
			Offload:  cs.offload,
			MaxDelta: cs.maxHealthLoadDelta,
			MinDelta: cs.minHealthLoadDelta,
		}
		required := counter.Required(deltaLoad)
		delta := counter.Delta(deltaLoad)

		log.Println(delta, cs.delta)
		log.Println(required, cs.required)
		assert.Equal(t, true, float32Equal(required, cs.required), "required not equal")
		assert.Equal(t, true, float32Equal(delta, cs.delta), "delta not equal")
	}
}

func DBStatistics(cpuLoad, cpu, ramLoad, ram float32) statiscticsmodel.Statistics {
	return statiscticsmodel.Statistics{
		DBCPU:     cpu,
		DBCPULoad: cpuLoad,
		DBRAM:     ram,
		DBRAMLoad: ramLoad,
	}
}

func VMStatistics(cpuLoad, cpu, ramLoad, ram float32) statiscticsmodel.Statistics {
	return statiscticsmodel.Statistics{
		VMCPU:     cpu,
		VMCPULoad: cpuLoad,
		VMRAM:     ram,
		VMRAMLoad: ramLoad,
	}
}

func Test_DeltaCounterFabric_Counter_Delta(t *testing.T) {
	cases := []struct {
		statistics   statiscticsmodel.Statistics
		resources    []*resourcemodel.Resource
		resourceType resourcetype.ResourceType

		vm, db executor.HealthLoad

		delta executor.Delta
	}{
		{
			statistics: DBStatistics(159, 29, 34, 89), // cpu load: 159, current: 29, potential: 90; ram load: 34, current: 89, potential 108
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(10, 14, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(11, 1, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(10, 1, true),

				NewResourceCpuRam(10, 40, false),
				NewResourceCpuRam(10, 40, false),
				NewResourceCpuRam(9, 9, false),
			},
			db: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
			},
			resourceType: resourcetype.DB,

			delta: executor.Delta{
				CPU: -29.3289,
				RAM: -68.1842,
			},
		},
		{
			statistics: DBStatistics(0, 29, 0, 89), // cpu load: 159, current: 29, potential: 90; ram load: 34, current: 89, potential 108
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(10, 14, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(11, 1, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(10, 1, true),

				NewResourceCpuRam(10, 40, false),
				NewResourceCpuRam(10, 40, false),
				NewResourceCpuRam(9, 9, false),
			},
			db: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
			},
			resourceType: resourcetype.DB,

			delta: executor.Delta{
				CPU: 0,
				RAM: 0,
			},
		},

		{
			statistics: DBStatistics(80, 10, 17.51, 10),
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(10, 10, false),
			},
			db: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.8,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.8,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
			},
			resourceType: resourcetype.DB,

			delta: executor.Delta{
				CPU: 3,
				RAM: -3.249,
			},
		},

		{
			statistics: DBStatistics(80, 10, 17.51, 10),
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(10, 10, false),
			},

			resourceType: resourcetype.DB,
			db: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.9,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.3,
					MinDelta: 6,
					MaxDelta: 70,
					Offload:  1,
				},
			},

			delta: executor.Delta{
				CPU: 3,
				RAM: -2.249,
			},
		},
		{
			statistics: VMStatistics(159, 29, 34, 89), // cpu load: 159, current: 29, potential: 90; ram load: 34, current: 89, potential 108
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(10, 14, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(11, 1, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(10, 1, true),

				NewResourceCpuRam(10, 40, false),
				NewResourceCpuRam(10, 40, false),
				NewResourceCpuRam(9, 9, false),
			},
			resourceType: resourcetype.VM,
			vm: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
			},

			delta: executor.Delta{
				CPU: -29.3289,
				RAM: -68.1842,
			},
		},
		{
			statistics: VMStatistics(0, 29, 0, 89), // cpu load: 159, current: 29, potential: 90; ram load: 34, current: 89, potential 108
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(10, 14, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(11, 1, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(10, 1, true),
				NewResourceCpuRam(10, 1, true),

				NewResourceCpuRam(10, 40, false),
				NewResourceCpuRam(10, 40, false),
				NewResourceCpuRam(9, 9, false),
			},
			resourceType: resourcetype.VM,
			vm: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
			},

			delta: executor.Delta{
				CPU: 0,
				RAM: 0,
			},
		},

		{
			statistics: VMStatistics(80, 10, 17.51, 10),
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(10, 10, false),
			},

			resourceType: resourcetype.VM,
			vm: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.76,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
			},

			delta: executor.Delta{
				CPU: 3,
				RAM: -3.249,
			},
		},

		{
			statistics: VMStatistics(80, 10, 17.51, 10),
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(10, 10, false),
			},

			resourceType: resourcetype.VM,
			vm: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.9,
					MinDelta: 5,
					MaxDelta: 70,
					Offload:  1,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.3,
					MinDelta: 6,
					MaxDelta: 70,
					Offload:  1,
				},
			},

			delta: executor.Delta{
				CPU: 3,
				RAM: -2.249,
			},
		},
		/*
			{
				load:               90.91,
				current:            1000,
				potential:          1000,
				offload:            0.95,
				k:                  0.8,
				minHealthLoadDelta: 1,
				maxHealthLoadDelta: 70,
				required:           933.645,
				delta:              -66.355,
			},
		*/
		{
			statistics: DBStatistics(90.91, 1000, 90.91, 1000),
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(1000, 1000, false),
			},
			resourceType: resourcetype.DB,

			db: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.8,
					Offload:  0.95,
					MinDelta: 1,
					MaxDelta: 70,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.9,
					Offload:  0.93,
					MinDelta: 1,
					MaxDelta: 65,
				},
			},
			delta: executor.Delta{
				CPU: -66.355,
				RAM: -89.537,
			},
		},

		{
			statistics: VMStatistics(90.91, 1000, 90.91, 1000),
			resources: []*resourcemodel.Resource{
				NewResourceCpuRam(1000, 1000, false),
			},
			resourceType: resourcetype.VM,

			vm: executor.HealthLoad{
				CPU: executor.ResourceDeltaLoad{
					Load:     0.8,
					Offload:  0.95,
					MinDelta: 1,
					MaxDelta: 70,
				},
				RAM: executor.ResourceDeltaLoad{
					Load:     0.9,
					Offload:  0.93,
					MinDelta: 1,
					MaxDelta: 65,
				},
			},
			delta: executor.Delta{
				CPU: -66.355,
				RAM: -89.537,
			},
		},
	}

	for _, cs := range cases {
		fabric := executor.NewDeltaCounterFabric(cs.vm, cs.db)
		counter := fabric.DeltaCounter(cs.resourceType)
		delta := counter.Delta(cs.statistics, cs.resources)
		t.Log(delta)
		assert.Equal(t, true, float32Equal(delta.CPU, cs.delta.CPU), "delta cpu not equal, %s", cs.resourceType)
		assert.Equal(t, true, float32Equal(delta.RAM, cs.delta.RAM), "delta ram not equal, %s", cs.resourceType)
	}
}
