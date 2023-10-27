package resourcemanager

import (
	"context"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/model/pricemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/internal/pricemanager"
	"log"
	"sort"
)

//go:generate mockgen -source resource_manager.go -destination resource_manager_mocks.go -package resourcemanager

type ResourceReport struct {
	CPU, RAM float32
}

func (r *ResourceReport) Plus(report ResourceReport) {
	r.CPU += report.CPU
	r.RAM += report.RAM
}

func (r *ResourceReport) PlusPrice(price pricemodel.Price) {
	r.CPU += price.CPU
	r.RAM += price.RAM
}

func (r *ResourceReport) MinusPrice(price pricemodel.Price) {
	r.CPU -= price.CPU
	r.RAM -= price.RAM
}

func (r *ResourceReport) MinusResource(resource resourcemodel.Resource) {
	r.CPU -= resource.CPU
	r.RAM -= resource.RAM
}

// resource manager expect positive on add and negative on remove
type ResourceManager interface {
	AddCPU(ctx context.Context, cpu float32) (report ResourceReport, err error)
	RemoveCPU(ctx context.Context, cpu, ramLimit float32) (report ResourceReport, err error)

	AddRAM(ctx context.Context, ram float32) (report ResourceReport, err error)
	RemoveRAM(ctx context.Context, ram, cpuLimit float32) (report ResourceReport, err error)
}

// является ядром нашего сервиса, именно он решает по какой стратегии мы будем добавлять/удалять ресурсы
type resourceManager struct {
	maxMachineCount int

	resourceStorage ResourceStorage
	upgrader        ResourceAddUpgrader
	remover         ResourceRemover
	priceManager    pricemanager.PriceManager
}

func New(
	resourceStorage ResourceStorage,
	upgrader ResourceAddUpgrader,
	remover ResourceRemover,
	priceManager pricemanager.PriceManager,
) *resourceManager {
	return &resourceManager{
		resourceStorage: resourceStorage,
		upgrader:        upgrader,
		remover:         remover,
		priceManager:    priceManager,
	}
}

func NewResourceManager(
	resourceStorage ResourceStorage,
	priceManager pricemanager.PriceManager,
	maxMachineCount int,
) ResourceManager {
	return &resourceManager{
		resourceStorage: resourceStorage,
		upgrader:        NewUpgrader(resourceStorage),
		remover:         NewRemover(resourceStorage),
		priceManager:    priceManager,

		maxMachineCount: maxMachineCount,
	}
}

// ссылочный тип float32, на подобие Float в Java
type Float struct {
	Value float32
}

// базовый метод для добавляния ресурса
func (r *resourceManager) Add(
	ctx context.Context,
	sortF func([]*resourcemodel.Resource),
	state AddState,
) (ResourceReport, error) {
	// если у ресурса нет дефицита или он достиг состояния дзена то прекращаем выполнение функции как будто её и не вызывали
	if state.Resource() <= 0 {
		return ResourceReport{}, nil
	}
	// получаем список ресурсов
	resourceList, err := r.resourceStorage.ResourceList(ctx)
	if err != nil {
		return ResourceReport{}, fmt.Errorf("failed get resource list, %w", err)
	}
	// сортируем список ресурсов
	sortF(resourceList)
	// обновляем ресурсы у state
	state.SetResources(resourceList)
	// обновляем количество максимальных машин у state
	state.SetMaxMachineCount(r.maxMachineCount)

	// Шаг 1: Добавление маленьких ресурсов до максимального предела количества машин
	// код будет выполняться пока мы не превысим максимальное количество машин
	// иначе говоря пока текущее количество машин меньше максимального мы будем выполнять нижеописанные инструкции
	for !state.MaxMachineCountOverflow() {
		// получаем информацию о ресурсах которые будут добавлены при добавлении минимального ресурса
		addMinResourceReport := r.upgrader.AddResourceReport(state.Min())
		// если мы не превышаем предел то добавляем ресурс
		if !state.IsOverflowReport(addMinResourceReport) {
			err := r.upgrader.AddResource(ctx, state.Min())
			// если добавление прошло неудачно то прекращаем добавлять ресурсы и переходим к следующему шагу
			if err != nil {
				log.Printf("failed add min resource, err: %s", err)
				break
			}
			// если добавления прошло успешно то применяем изменения к нашему состоянию
			state.Add(addMinResourceReport)
			continue
		}
		break
	}

	// Шаг 2: Обновление минимальных ресурсов до масимальных
	// Если после первого шага мы добавили достаточно ресурсов то мы выйдем из цикла в начале первой итерации
	// Если у нас нет ресурсов которые мы могли бы проапргейдить, то Шаг 2 не внесёт изменений в систему
	for _, resource := range resourceList {
		// если уже усё то усё
		if state.Overflow() {
			break
		}
		// если ресурс не прошёл фейс контроль то шлём его подальше и кричим: СЛЕДУЮЩИЙ!!
		if !state.CanUpgradeResource(resource) {
			continue
		}
		// получаем информацию об изменениях при апдейте ресурса
		report := r.upgrader.UpgradeResourceReport(*resource, state.Min(), state.Max())
		//	апдейтим ресурс
		err := r.upgrader.UpgradeResource(ctx, *resource, state.Max())
		if err != nil {
			return state.Report(), fmt.Errorf("failed upgrade resource %d, err: %s", resource.ID, err)
		}
		// применяем изменения для текущего сотояния
		state.Upgrade(report)
	}
	// возращаем отчёт об изменённых ресурсах
	return state.Report(), nil
}

// базовый метод для удаления ресурса
func (r *resourceManager) Remove(
	ctx context.Context,
	sortF func([]*resourcemodel.Resource),
	removeState RemoveState,
) (ResourceReport, error) {
	if removeState.Resource() >= 0 || removeState.Limit() >= 0 {
		return ResourceReport{}, nil
	}
	resourceList, err := r.resourceStorage.ResourceList(ctx)
	if err != nil {
		return ResourceReport{}, fmt.Errorf("failed get resource list, %w", err)
	}
	sortF(resourceList)
	removeState.SetResources(resourceList)
	for _, resource := range resourceList {
		if !removeState.CanRemoveResource(resource) {
			continue
		}
		downReport := r.remover.RemoveReport(ctx, *resource)
		if removeState.IsOverflowReport(downReport) {
			continue
		}
		err = r.remover.RemoveResource(ctx, *resource)
		if err != nil {
			err = fmt.Errorf("failed remove resource, %w", err)
			log.Printf("%s", err)
			continue
		}
		removeState.Remove(downReport)
	}
	return removeState.Report(), nil
}

func (r *resourceManager) AddCPU(ctx context.Context, cpu float32) (report ResourceReport, err error) {
	min, max := r.priceManager.MinCPU(), r.priceManager.MaxCPU()
	state := NewAddState(
		cpu,
		min,
		max,
		func(f Float) bool { return -f.Value >= min.CPU },
		func(f *Float, rr ResourceReport) { f.Value -= rr.CPU },
	)
	return r.Add(
		ctx,
		func(r []*resourcemodel.Resource) { sort.Sort(SortResourcesByCPUAsc(r)) },
		state,
	)
}

func (r *resourceManager) AddRAM(ctx context.Context, ram float32) (report ResourceReport, err error) {
	min, max := r.priceManager.MinRAM(), r.priceManager.MaxRAM()
	state := NewAddState(
		ram,
		min,
		max,
		func(f Float) bool { return -f.Value >= min.RAM },
		func(f *Float, rr ResourceReport) { f.Value -= rr.RAM },
	)
	return r.Add(
		ctx,
		func(r []*resourcemodel.Resource) { sort.Sort(SortResourcesByRAMAsc(r)) },
		state,
	)
}

func (r *resourceManager) RemoveCPU(ctx context.Context, cpu, ramLimit float32) (reporrt ResourceReport, err error) {
	min, max := r.priceManager.MinCPU(), r.priceManager.MaxCPU()
	removeState := NewRemoveState(
		cpu,
		ramLimit,
		min,
		max,
		func(resource, limit *Float, report ResourceReport) {
			resource.Value -= report.CPU
			limit.Value -= report.RAM
		},
	)
	return r.Remove(
		ctx,
		func(r []*resourcemodel.Resource) { sort.Sort(SortResourcesByCPUDesc(r)) },
		removeState,
	)
}

func (r *resourceManager) RemoveRAM(ctx context.Context, ram, cpuLimit float32) (report ResourceReport, err error) {
	min, max := r.priceManager.MinRAM(), r.priceManager.MaxRAM()
	removeState := NewRemoveState(
		ram,
		cpuLimit,
		min,
		max,
		func(resource, limit *Float, report ResourceReport) {
			resource.Value -= report.RAM
			limit.Value -= report.CPU
		},
	)
	return r.Remove(
		ctx,
		func(r []*resourcemodel.Resource) { sort.Sort(SortResourcesByRAMDesc(r)) },
		removeState,
	)
}
