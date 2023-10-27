package executor

import (
	"github.com/amidgo/cloud-resources/internal/resourcemanager"
)

type Delta struct {
	// указывает на количество ресурса которое нужно прибавить,
	// отрицательное значение указывает на то что необходимо отнять данный ресурс
	// положительное значение указывает на то что необходимо добавить данные ресурс
	CPU, RAM float32
}

func (d *Delta) MinusResourceReport(report resourcemanager.ResourceReport) {
	d.CPU -= report.CPU
	d.RAM -= report.RAM
}

func (d *Delta) PlusResourceReport(report resourcemanager.ResourceReport) {
	d.CPU += report.CPU
	d.RAM += report.RAM
}
