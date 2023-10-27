package executor

import (
	"context"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/internal/model/statiscticsmodel"
	"github.com/amidgo/cloud-resources/internal/resourcemanager"
)

//go:generate mockgen -source executor.go -destination executor_mocks.go -package executor

// интерфейс для инициальзации ресурсов
type ExecuteInitter interface {
	Init(ctx context.Context) ([]*resourcemodel.Resource, error)
}

type initter struct {
	resourceManager  resourcemanager.ResourceManager
	resourceStorage  ResourceStorage
	InitCPU, InitRAM float32
}

// вызов этого метода добавляет начальные ресурсы по cpu и ram
// возвращает обновлённые ресурсы
func (e *initter) Init(ctx context.Context) ([]*resourcemodel.Resource, error) {
	e.resourceManager.AddCPU(ctx, e.InitCPU)
	e.resourceManager.AddRAM(ctx, e.InitRAM)
	resources, err := e.resourceStorage.ResourceList(ctx)
	if err != nil {
		err = fmt.Errorf("failed get resource list, %w", err)
		return nil, err
	}
	return resources, nil
}

func NewInitter(
	resourceManager resourcemanager.ResourceManager,
	resourceStorage ResourceStorage,
	initCpu, initRam float32,
) ExecuteInitter {
	return &initter{
		resourceManager: resourceManager,
		resourceStorage: resourceStorage,
		InitCPU:         initCpu,
		InitRAM:         initRam,
	}
}

type Executor interface {
	Execute(ctx context.Context) error
}

type StatiscticsStorage interface {
	Statisctics(ctx context.Context) (statiscticsmodel.Statistics, error)
}
type ResourceStorage interface {
	ResourceList(ctx context.Context) ([]*resourcemodel.Resource, error)
}

// основная структура которая имплементирует интерфейс Executor
// вызов метода Execute является точкой входа в плане изменения ресурсов
type executor struct {
	log     Logger
	initter ExecuteInitter

	resourceStorage   ResourceStorage
	statisticsStorage StatiscticsStorage
	deltaExecutor     DeltaExecutor
	deltaCounter      DeltaCounter
}

func NewExecutor(
	log Logger,
	initter ExecuteInitter,

	resourceStorage ResourceStorage,
	statisticsStorage StatiscticsStorage,
	deltaExecutor DeltaExecutor,
	deltaCounter DeltaCounter,
) Executor {
	return &executor{
		log:               log,
		initter:           initter,
		resourceStorage:   resourceStorage,
		statisticsStorage: statisticsStorage,
		deltaExecutor:     deltaExecutor,
		deltaCounter:      deltaCounter,
	}
}

func (e *executor) Execute(ctx context.Context) error {
	// получаем текущую статистику
	statistics, err := e.statisticsStorage.Statisctics(ctx)
	if err != nil {
		return fmt.Errorf("failed get statistics %w", err)
	}
	// получаем текущие ресурсы
	resources, err := e.resourceStorage.ResourceList(ctx)
	if err != nil {
		return fmt.Errorf("failed get resources, %w", err)
	}
	e.log.Printf("machine count %d", len(resources))
	// если ресурсов нет то вызываем метод Init у нашего ExecuteInitter
	if len(resources) == 0 {
		e.log.Printf("zero vm machine count")
		resources, err = e.initter.Init(ctx)
		if err != nil {
			return fmt.Errorf("failed init resources, %w", err)
		}
	}
	// считаем разницу ресурсов
	delta := e.deltaCounter.Delta(statistics, resources)
	e.log.Printf("delta, cpu: %f ram: %f", delta.CPU, delta.RAM)
	// добавляем/убавляем ресурсы исходя из значений delta
	err = e.deltaExecutor.ExecuteDelta(ctx, delta)
	if err != nil {
		return fmt.Errorf("failed execute delta, %w", err)
	}
	return nil
}
