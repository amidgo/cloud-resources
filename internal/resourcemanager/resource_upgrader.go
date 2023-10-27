package resourcemanager

import (
	"context"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/model/pricemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
)

//go:generate mockgen -source resource_upgrader.go -destination resource_upgrader_mocks.go -package resourcemanager

type ResourceAddUpgrader interface {
	UpgradeResource(ctx context.Context, resource resourcemodel.Resource, max pricemodel.Price) (err error)
	UpgradeResourceReport(resource resourcemodel.Resource, min, max pricemodel.Price) (report ResourceReport)
	AddResource(ctx context.Context, price pricemodel.Price) (err error)
	AddResourceReport(price pricemodel.Price) (report ResourceReport)
}

// структура ответсвенная за добавление ресурсов с помощью создания/апгрейда ресурсов
// имеет методы для получения отчета об измении и самом изменении
type resourceAddUpgrader struct {
	resourceStorage ResourceStorage
}

func NewUpgrader(resourceStorage ResourceStorage) ResourceAddUpgrader {
	return &resourceAddUpgrader{
		resourceStorage: resourceStorage,
	}
}

// когда мы апрейдим минимальный ресурс то изменения равняется разнице между макисмальным и минимальным ресурсом
func (r *resourceAddUpgrader) UpgradeResourceReport(resource resourcemodel.Resource, min, max pricemodel.Price) (report ResourceReport) {
	report.PlusPrice(max)
	report.MinusPrice(min)
	return
}

func (r *resourceAddUpgrader) UpgradeResource(
	ctx context.Context,
	resource resourcemodel.Resource,
	max pricemodel.Price,
) (err error) {
	err = r.UpgradeByPrice(ctx, resource.ID, max)
	if err != nil {
		err = fmt.Errorf("failed upgrade resource %d to max, %w", resource.ID, err)
		return
	}
	return
}

func (r *resourceAddUpgrader) AddResource(ctx context.Context, price pricemodel.Price) (err error) {
	err = r.AddByPrice(ctx, price)
	if err != nil {
		err = fmt.Errorf("failed add max cpu, %w", err)
		return
	}
	return
}
func (r *resourceAddUpgrader) AddResourceReport(price pricemodel.Price) (report ResourceReport) {
	report.PlusPrice(price)
	return
}

func (r *resourceAddUpgrader) UpgradeByPrice(ctx context.Context, resourceId int, price pricemodel.Price) (err error) {
	err = r.resourceStorage.UpdateResource(ctx, resourceId, resourcemodel.NewUpdateResource(price.CPU, price.RAM, price.Type))
	if err == nil {
		return
	}
	return
}

func (r *resourceAddUpgrader) AddByPrice(ctx context.Context, price pricemodel.Price) (err error) {
	_, err = r.resourceStorage.AddResource(ctx, resourcemodel.NewAddResource(price.CPU, price.RAM, price.Type))
	if err != nil {
		return
	}
	return
}
