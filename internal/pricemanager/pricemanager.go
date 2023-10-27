package pricemanager

import (
	"context"
	"github.com/amidgo/cloud-resources/internal/model/pricemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"log"
	"time"
)

type PriceStorage interface {
	PriceList(ctx context.Context) ([]*pricemodel.Price, error)
}

// чисто технически нам не обязательно возвращать САМЫЙ минимальны и САМЫЙ максимальный ресурс
// достаточно соблюдать правило что значение всех ресурсов MinCPU должно быть меньше значений ресурсов MaxCPU, тоже самое с RAM
// можно просто выбрать оптимальный набор ресурсов который нам будет эффективнее использовать
type PriceManager interface {
	MinCPU() pricemodel.Price
	MinRAM() pricemodel.Price
	MaxCPU() pricemodel.Price
	MaxRAM() pricemodel.Price
}

// набор из минимальных и максимальных цен по CPU/RAM
type PriceSet struct {
	MinRAM, MinCPU, MaxRAM, MaxCPU pricemodel.Price
}

// как вы могли заметить мне было лень тестить эту залупу, поэтому на соревах я просто смотрел как он парсит ресы и вроде было норм
// и если я буду как-то менять структуру этого проекта, возможно я заменю функции setMin...|setMax на что-то более ООПэшное
// но в Москве время было ограничено
type useCase struct {
	// тип по которому мы смотрим ресурсы
	resourceType resourcetype.ResourceType
	priceSet     PriceSet
	priceStorage PriceStorage
}

func New(priceStorage PriceStorage, resourceType resourcetype.ResourceType) PriceManager {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	u := &useCase{
		resourceType: resourceType,
		priceStorage: priceStorage,
	}
	// обновляем список цен при создании нового PriceManager
	u.updatePriceSet(ctx)
	return u
}

// итерируемся по списку цен, находим максимальные и минимальные цены
func (u *useCase) updatePriceSet(ctx context.Context) {
	priceList, err := u.priceStorage.PriceList(ctx)
	if err != nil {
		log.Printf("%s", err)
	}
	for _, price := range priceList {
		if price.Type != u.resourceType {
			continue
		}
		u.setMinCPU(*price)
		u.setMaxCPU(*price)
		u.setMinRAM(*price)
		u.setMaxRAM(*price)
	}
}

func (u *useCase) setMinCPU(price pricemodel.Price) {
	switch {
	case u.MinCPU().IsZero():
		u.priceSet.MinCPU = price
	case u.MinCPU().CPU > price.CPU:
		u.priceSet.MinCPU = price
	case u.MinCPU().CPU == price.CPU:
		isCostLower := price.Cost < u.MinCPU().Cost
		if isCostLower {
			u.priceSet.MinCPU = price
		}
		isCostEqual := u.MinCPU().Cost == price.Cost
		isRamBigger := u.MinCPU().RAM < price.RAM
		if isCostEqual && isRamBigger {
			u.priceSet.MinCPU = price
		}
	}

}

func (u *useCase) setMinRAM(price pricemodel.Price) {
	switch {
	case u.MinRAM().IsZero():
		u.priceSet.MinRAM = price
	case u.MinRAM().RAM > price.RAM:
		u.priceSet.MinRAM = price
	case u.MinCPU().RAM == price.RAM:
		isCostLower := price.Cost < u.MinRAM().Cost
		if isCostLower {
			u.priceSet.MinRAM = price
		}
		isCostEqual := u.MinRAM().Cost == price.Cost
		isCpuBigger := u.MinRAM().CPU < price.CPU
		if isCostEqual && isCpuBigger {
			u.priceSet.MinRAM = price
		}
	}
}

func (u *useCase) setMaxCPU(price pricemodel.Price) {
	switch {
	case u.MaxCPU().CPU < price.CPU:
		u.priceSet.MaxCPU = price
	case u.MaxCPU().CPU == price.CPU:
		isCostLower := price.Cost < u.MaxCPU().Cost
		if isCostLower {
			u.priceSet.MaxCPU = price
		}
		isCostEqual := u.MaxCPU().Cost == price.Cost
		isRamBigger := u.MaxCPU().RAM < price.RAM
		if isCostEqual && isRamBigger {
			u.priceSet.MaxCPU = price
		}
	}
}

func (u *useCase) setMaxRAM(price pricemodel.Price) {
	switch {
	case u.MaxRAM().RAM < price.RAM:
		u.priceSet.MaxRAM = price
	case u.MaxRAM().RAM == price.RAM:
		isCostLower := price.Cost < u.MaxRAM().Cost
		if isCostLower {
			u.priceSet.MaxRAM = price
		}
		isCostEqual := u.MaxRAM().Cost == price.Cost
		isCpuBigger := u.MaxRAM().CPU < price.CPU
		if isCostEqual && isCpuBigger {
			u.priceSet.MaxRAM = price
		}
	}
}

func (u *useCase) MinCPU() pricemodel.Price {
	return u.priceSet.MinCPU
}

func (u *useCase) MinRAM() pricemodel.Price {
	return u.priceSet.MinRAM
}

func (u *useCase) MaxCPU() pricemodel.Price {
	return u.priceSet.MaxCPU
}

func (u *useCase) MaxRAM() pricemodel.Price {
	return u.priceSet.MaxRAM
}
