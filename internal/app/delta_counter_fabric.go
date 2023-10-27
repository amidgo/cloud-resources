package app

import (
	"github.com/amidgo/cloud-resources/config"
	"github.com/amidgo/cloud-resources/internal/executor"
)

// парсит конфиг и возвращает executor.DeltaCounterFabric
func DeltaCounterFabric(cnf *config.AppConfig) executor.DeltaCounterFabric {
	vmHealthLoadConfig := cnf.VM.HealthLoad
	vmHealthLoad := executor.HealthLoad{
		CPU: executor.ResourceDeltaLoad{
			Offload:  vmHealthLoadConfig.CPU.Offload,
			Load:     vmHealthLoadConfig.CPU.Load,
			MinDelta: vmHealthLoadConfig.CPU.MinDelta,
			MaxDelta: vmHealthLoadConfig.CPU.MaxDelta,
		},
		RAM: executor.ResourceDeltaLoad{
			Offload:  vmHealthLoadConfig.RAM.Offload,
			Load:     vmHealthLoadConfig.RAM.Load,
			MinDelta: vmHealthLoadConfig.RAM.MinDelta,
			MaxDelta: vmHealthLoadConfig.RAM.MaxDelta,
		},
	}
	dbHealthLoadConfig := cnf.DB.HealthLoad
	dbHealthLoad := executor.HealthLoad{
		CPU: executor.ResourceDeltaLoad{
			Offload:  dbHealthLoadConfig.CPU.Offload,
			Load:     dbHealthLoadConfig.CPU.Load,
			MinDelta: dbHealthLoadConfig.CPU.MinDelta,
			MaxDelta: dbHealthLoadConfig.CPU.MaxDelta,
		},
		RAM: executor.ResourceDeltaLoad{
			Offload:  dbHealthLoadConfig.RAM.Offload,
			Load:     dbHealthLoadConfig.RAM.Load,
			MinDelta: dbHealthLoadConfig.RAM.MinDelta,
			MaxDelta: dbHealthLoadConfig.RAM.MaxDelta,
		},
	}
	fabric := executor.NewDeltaCounterFabric(vmHealthLoad, dbHealthLoad)
	return fabric
}
