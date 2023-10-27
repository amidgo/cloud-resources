package resourcemanager

import (
	"github.com/amidgo/cloud-resources/internal/model/pricemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
)

// своего рода состояние (не путать с паттерном State) которое инкапсулирует ряд переменных
// по сути не обязательна, но тогда метод resourceManager.Remove превращается в непонятное месиво
type RemoveState struct {
	// основной и побочный ресурс которые мы не можем превышать даже если основной все ёщё положительный
	/*
		Например:
			Если мы убавляем CPU то основным ресурсом будет cpu а побочным ram
			Допустим что у нас накопился излишек ресурсов (излишек должен быть отрицательным):
				CPU: -10
				RAM: -3
			в таком случае мы можем удалять CPU до тех пор пока RAM или CPU больше 0 см. RemoveState.IsOverflowReport
	*/
	// при увеличении количества переменных есть смысл создать отдельную структуру
	// которая будет инкапсулировать все эти перемнные внутри себя, но пока мы можем так не запариваться, то оставим как есть
	resource, limit Float
	// минимальный максимальный ресурсы
	min, max pricemodel.Price
	// счётчик активных ресурсов
	availabelResourceCounter AvailabelResourceCounter
	report                   ResourceReport

	op func(resource, limit *Float, report ResourceReport)
}

func NewRemoveState(
	resource, limit float32,
	min, max pricemodel.Price,
	op func(resource, limit *Float, report ResourceReport),
) RemoveState {
	return RemoveState{
		resource: Float{Value: resource},
		limit:    Float{Value: limit},
		min:      min,
		max:      max,
		op:       op,
	}
}

func (r *RemoveState) Resource() float32 {
	return r.resource.Value
}

func (r *RemoveState) Limit() float32 {
	return r.limit.Value
}

// обновляем текущее количество активных ресурсов
func (a *RemoveState) SetResources(rlist []*resourcemodel.Resource) {
	var availabelResourceCount int
	for i := range rlist {
		if rlist[i].Failed {
			continue
		}
		availabelResourceCount++
	}
	a.availabelResourceCounter.SetResourceCount(availabelResourceCount)
}

// применяем изменения из report, увеличиваем количество недоступных ресурсов
func (a *RemoveState) Remove(report ResourceReport) {
	a.op(&a.resource, &a.limit, report)
	a.report.Plus(report)
	a.availabelResourceCounter.IncFailed()
}

// возвращаем информацию о том можно ли удалить данный ресурс
/*
	Условия для удаления ресурса:
		1. Он активный
		2. Помимо него у нас должна быть хотя бы 1 активная машина (при чем любая)
		(система устроена так чтобы мы не могли оставить 0 активных машин, хотя бы одна но должна функционировать)
*/
func (a *RemoveState) CanRemoveResource(resource *resourcemodel.Resource) bool {
	return !resource.Failed && a.availabelResourceCounter.CanIncFailed()
}

// возвращает информацию о том можем ли мы применить изменения из отчета report
// нельзя допустить чтобы один из ресурсов перешёл черту, лучше иметь небольшой профицит чем доводить до недостатка ресурсов
func (a *RemoveState) IsOverflowReport(report ResourceReport) bool {
	res, lim := a.resource, a.limit
	a.op(&res, &lim, report)
	return res.Value > 0 || lim.Value > 0
}

func (r *RemoveState) Report() ResourceReport {
	return r.report
}
