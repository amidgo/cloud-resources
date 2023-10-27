package executor

import (
	"context"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/resourcemanager"
)

//go:generate mockgen -source delta_executor.go -destination delta_executor_mocks.go -package executor

type DeltaExecutor interface {
	ExecuteDelta(ctx context.Context, delta Delta) error
}

// структура которая напрямую занимается изменением ресурсов и взаимодействием с resourceManager
type deltaExecutor struct {
	log             Logger
	resourceManager resourcemanager.ResourceManager
}

func NewDeltaExecutor(
	log Logger,
	resourceManager resourcemanager.ResourceManager,
) DeltaExecutor {
	return &deltaExecutor{log: log, resourceManager: resourceManager}
}

// в зависимости от значений структуры данных delta мы добавляем или отнимаем ресурсы
// задача этого метода максимально приблизить значения delta к Delta{CPU: 0, RAM: 0} то есть к идеальной (целевой) нагрузке
// достигается это за счёт того что каждый метод resourceManager возвращает нам информацию о добавленных/удалённых ресурсах (report)
// после вызова каждого метода resourceManager мы учитываем количество отнятых/прибавленных ресурсов при следующем вызове в рамках этой функции
func (d *deltaExecutor) ExecuteDelta(ctx context.Context, delta Delta) (resErr error) {
	if delta.CPU > 0 {
		d.log.Printf("add cpu %f", delta.CPU)
		report, err := d.resourceManager.AddCPU(ctx, delta.CPU)
		// отнимаем прибавленные ресурсы от текущей delta, для того чтобы последующие вызовы были в курсе что ресурсов стало больше
		delta.MinusResourceReport(report)
		if err != nil {
			resErr = fmt.Errorf("failed add cpu, %w", err)
		}
	}
	if delta.RAM > 0 {
		d.log.Printf("add ram %f", delta.RAM)
		report, err := d.resourceManager.AddRAM(ctx, delta.RAM)
		delta.MinusResourceReport(report)
		if err != nil {
			resErr = fmt.Errorf("failed add ram, %w", err)
		}
	}
	if delta.CPU < 0 {
		d.log.Printf("remove cpu %f", delta.CPU)
		report, err := d.resourceManager.RemoveCPU(ctx, delta.CPU, delta.RAM)
		// отнимаем прибавленные ресурсы от текущей delta, для того чтобы последующие вызовы были в курсе что ресурсов стало меньше
		delta.MinusResourceReport(report)
		if err != nil {
			resErr = fmt.Errorf("failed remove cpu, %w", err)
		}
	}
	if delta.RAM < 0 {
		d.log.Printf("remove ram %f", delta.RAM)
		report, err := d.resourceManager.RemoveRAM(ctx, delta.RAM, delta.CPU)
		delta.MinusResourceReport(report)
		if err != nil {
			resErr = fmt.Errorf("failed remove ram, %w", err)
		}
	}
	return
}
