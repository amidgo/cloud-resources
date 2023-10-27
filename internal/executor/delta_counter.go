package executor

import (
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"github.com/amidgo/cloud-resources/internal/model/statiscticsmodel"
	"log"
)

//go:generate mockgen -source delta_counter.go -destination delta_counter_mocks.go -package executor

type DeltaCounter interface {
	Delta(statistics statiscticsmodel.Statistics, resources []*resourcemodel.Resource) (delta Delta)
}

type DeltaCounterFabric interface {
	DeltaCounter(resourcetype.ResourceType) DeltaCounter
}

// структура которая создержит конфигурационные данные для подсчёта необходимого количества ресурсов
type ResourceDeltaLoad struct {
	Load, Offload, MinDelta, MaxDelta float32
}

// структура которая содержит конфигурацинные данные по cpu/ram
type HealthLoad struct {
	CPU, RAM ResourceDeltaLoad
}

// фабрика которая возвращает нам нужный DeltaCounter в зависимости от типа (vm или db)
type deltaCounterFabric struct {
	vm, db HealthLoad
}

func NewDeltaCounterFabric(vm, db HealthLoad) DeltaCounterFabric {
	return deltaCounterFabric{
		vm: vm,
		db: db,
	}
}

func (d deltaCounterFabric) DeltaCounter(resourceType resourcetype.ResourceType) (counter DeltaCounter) {
	switch resourceType {
	case resourcetype.DB:
		counter = DBDeltaCounter{
			DeltaCalculator: DeltaCalculator{
				HealthLoad: d.db,
			},
		}
	case resourcetype.VM:
		counter = VMDeltaCounter{
			DeltaCalculator: DeltaCalculator{
				HealthLoad: d.vm,
			},
		}
	}
	return
}

// структура которая считает необходимое нам количество ресурсов
type ResourceStatisticCounter struct {
	Load      float32
	Current   float32
	Potential float32
}

func NewResourceStatisticCounter(load, current, potential float32) ResourceStatisticCounter {
	return ResourceStatisticCounter{
		Load:      load,
		Current:   current,
		Potential: potential,
	}
}

/*
Подсчёт необходимых ресурсов
*/
func (r ResourceStatisticCounter) Required(deltaLoad ResourceDeltaLoad) (delta float32) {
	// подсчёт текущей нагрузки при 100 процентах
	required100percent := r.Current * (r.Load / 100)
	// подсчёт нагрузки с учетом коэффицента ресурса при котором машина отключается
	required100percent = required100percent * deltaLoad.Offload
	// считаем целевое количество ресурсов которое мы хотим поддерживать, тобишь если нам нужна 80% нагрузка то deltaLoad.Load должен быть равен 0.8
	requiredK := required100percent / deltaLoad.Load
	// считаем целевое количество ресурсов с учетом минимального запаса, ниже которого мы не хотим опускаться
	// например при необходимых ресурсах 10 и целевой нагрузке 80% в нашем резерве будет 2
	// однако мы хотели бы иметь минимальный запас ресурсов для предотвращения скачков на нижнем сегменте
	requiredMinDelta := required100percent + deltaLoad.MinDelta
	// считаем целевое количество ресурсов с учетом максимального запаса
	// мы не можем превысить это значение даже если целевая нагрузка больше обозначеного процента
	requiredMaxDelta := required100percent + deltaLoad.MaxDelta
	// считаем максимум из целевой нагрузки и нагрузки с учетом минимального запаса
	// если целевая нагрузка меньше то мы используем нагрузку минимального запаса ресурсов
	delta = max(requiredK, requiredMinDelta)
	// считаем минимум из полученной ранее необходимой нам нагрузки и нагрузки с максимально возможным запасом
	// если ранее полученная нагрузка больше чем максимальная, используем нагрузку requiredMaxDelta
	delta = min(delta, requiredMaxDelta)
	return
}

func (r ResourceStatisticCounter) Delta(deltaLoad ResourceDeltaLoad) (delta float32) {
	// апишка показывает неверные данные поэтому используем такую затычку, чтобы не шалило
	if r.Load == 0 {
		return 0
	}
	required := r.Required(deltaLoad)
	// возвращаем разницу которую нам необходимо прибавить/убавить
	return required - r.Potential
}

type DeltaCalculator struct {
	HealthLoad HealthLoad
}

type ResourceStatistics struct {
	Potential, Current float32
}

// считаем текущие/потенциальные cpu/ram из списка ресурсов
/*
поскольку текущие ресурсы из запроса /resource не всегда соответсвуют из текущих ресурсов статистики то
было принято решения брать текущие ресурсы только из статистики, но мне лень было менять что-то за несколько дней
до чемпионата поэтому так, не осерчайте как говорится
*/
func (d DeltaCalculator) ReduceResourcesCpuRam(resources []*resourcemodel.Resource) (cpu, ram ResourceStatistics) {
	for _, r := range resources {
		cpu.Potential += r.CPU
		ram.Potential += r.RAM
		if !r.Failed {
			cpu.Current += r.CPU
			ram.Current += r.RAM
		}
	}
	return
}

// считаем текущую разницу ресурсов
func (d DeltaCalculator) CalculateDelta(cpu, ram ResourceStatisticCounter) (delta Delta) {
	delta = Delta{
		CPU: cpu.Delta(d.HealthLoad.CPU),
		RAM: ram.Delta(d.HealthLoad.RAM),
	}
	return
}

type DBDeltaCounter struct {
	DeltaCalculator
}

// текущие ресурсы только из статистики, апишка говна
func (c DBDeltaCounter) Delta(statistics statiscticsmodel.Statistics, resources []*resourcemodel.Resource) (delta Delta) {
	cpu, ram := c.ReduceResourcesCpuRam(resources)
	log.Printf("potential db resources, cpu: %f, ram: %f", cpu.Potential, ram.Potential)
	log.Printf("current db resources, cpu: %f, ram: %f", cpu.Current, ram.Current)
	log.Printf("db load, cpu: %f, ram: %f", statistics.DBCPULoad, statistics.DBRAMLoad)
	cpuStat := NewResourceStatisticCounter(statistics.DBCPULoad, statistics.DBCPU, cpu.Potential)
	ramStat := NewResourceStatisticCounter(statistics.DBRAMLoad, statistics.DBRAM, ram.Potential)
	return c.CalculateDelta(cpuStat, ramStat)
}

type VMDeltaCounter struct {
	DeltaCalculator
}

// текущие ресурсы только из статистики, апишка говна
func (c VMDeltaCounter) Delta(statistics statiscticsmodel.Statistics, resources []*resourcemodel.Resource) (delta Delta) {
	cpu, ram := c.ReduceResourcesCpuRam(resources)
	log.Printf("potential vm resources, cpu: %f, ram: %f", cpu.Potential, ram.Potential)
	log.Printf("current vm resources, cpu: %f, ram: %f", cpu.Current, ram.Current)
	log.Printf("vm load, cpu: %f, ram: %f", statistics.VMCPULoad, statistics.VMRAMLoad)
	cpuStat := NewResourceStatisticCounter(statistics.VMCPULoad, statistics.VMCPU, cpu.Potential)
	ramStat := NewResourceStatisticCounter(statistics.VMRAMLoad, statistics.VMRAM, ram.Potential)
	return c.CalculateDelta(cpuStat, ramStat)
}
