package resourcemanager

import (
	"github.com/amidgo/cloud-resources/internal/model/pricemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
)

func NewAddState(resource float32, min, max pricemodel.Price, overflow func(f Float) bool, op func(f *Float, rr ResourceReport)) AddState {
	return AddState{
		resource: Float{Value: resource},
		min:      min,
		max:      max,
		overflow: overflow,
		op:       op,
	}
}

// своего рода состояние (не путать с паттерном State) которое инкапсулирует ряд переменных
// по сути не обязательна, но тогда метод resourceManager.Add превращается в непонятное месиво
type AddState struct {
	// текущее значение ресурса
	resource Float
	// минимальный и максимальный ресурс для конкретного типа
	/*
		Напимер:
			для VM из CPU мы выделяем два ресурса: минимальный и максимальный
		Работу мы ведём исключительно с этими двумя экземплярами ресурсов
	*/
	min, max pricemodel.Price
	// счётчик текущий ресурсов со значением Failed = false
	availableResourceCounter AvailabelResourceCounter
	// максимальное количество машин
	maxMachineCount int
	// текущее количество машин
	machineCount int

	/*
		чисто технически можно было бы сделать для для следующий двух полей отдельный структуры
		которые бы имплементировали интерфейс с методом но:

			1. Мне было уже лень
			2. Как будто бы ещё чуть чуть и уже оверинжиром запахнет
		Возможно однажды я так и сделаю, но из-за олимпиды я целых 3 дня не занимался своими делами
		так что мне лень что-то здесь концептуально менять, + скоро апишка перестанет работать
		и для прогонки надо будет писать свой mock service, а мне лень
	*/

	// функция которая сообщает нам что: Братан, чот ты дофига всего создал, думаю с этого вызова функции достаточно
	overflow func(f Float) bool
	// функция которое изменяет значение переменной
	op func(f *Float, rr ResourceReport)

	// итоговый отчёт в котором будет указано сколько ресурсов cpu и ram мы добавили
	report ResourceReport
}

func (a *AddState) Resource() float32 {
	return a.resource.Value
}

func (a *AddState) Min() pricemodel.Price {
	return a.min
}

func (a *AddState) Max() pricemodel.Price {
	return a.max
}

func (a *AddState) SetMaxMachineCount(maxMachineCount int) {
	a.maxMachineCount = maxMachineCount
}

// обновляет количество машин, а также количество ресурсов которое в строю (Failed = false)
func (a *AddState) SetResources(rlist []*resourcemodel.Resource) {
	a.machineCount = len(rlist)
	var availabelResourceCount int
	for i := range rlist {
		if rlist[i].Failed {
			continue
		}
		availabelResourceCount++
	}
	a.availableResourceCounter.SetResourceCount(availabelResourceCount)
}

// сообщает о том можем ли мы проапгрейдить данный ресурс
/*
	Для этого необходимо выполнение следующих условий:
		1. Данный ресурс должен быть доступен
		2. Данный ресурс должен быть минимальным (нет смысла грейдить уже максимальный ресурс)
		3. Помимо него у нас должна быть хотя бы 1 активная машина (при чем любая)
		(система устроена так чтобы мы не могли оставить 0 активных машин, хотя бы одна но должна функционировать)
*/
func (a *AddState) CanUpgradeResource(resource *resourcemodel.Resource) bool {
	isMinResource := resource.CPU == a.min.CPU && resource.RAM == a.min.RAM
	return !resource.Failed && isMinResource && a.availableResourceCounter.CanIncFailed()
}

// возвращает информацию о том приведёт ли данный отчет к переполнению ресурса или нет
// состояние AddState при этом не меняется
func (a *AddState) IsOverflowReport(report ResourceReport) bool {
	res := a.resource
	a.op(&res, report)
	return a.overflow(res)
}

// возращает информацию достигли ли мы предела по количеству машин
func (a *AddState) MaxMachineCountOverflow() bool {
	return a.machineCount >= a.maxMachineCount
}

func (a *AddState) Overflow() bool {
	return a.overflow(a.resource)
}

// применяет изменения upgradeReport к текущему состоянию, увеличивает счётчик неактивных машин
func (a *AddState) Upgrade(upgradeReport ResourceReport) {
	a.availableResourceCounter.IncFailed()
	a.op(&a.resource, upgradeReport)
	a.report.Plus(upgradeReport)
}

// применяет изменения addReport к текущему состоянию, увеличивает счётчик текущего количества машин
func (a *AddState) Add(addReport ResourceReport) {
	a.op(&a.resource, addReport)
	a.report.Plus(addReport)
	a.machineCount++
}

func (a *AddState) Report() ResourceReport {
	return a.report
}
