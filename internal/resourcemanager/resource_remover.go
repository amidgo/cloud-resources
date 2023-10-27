package resourcemanager

import (
	"context"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
)

//go:generate mockgen -source resource_remover.go -destination resource_remover_mocks.go -package resourcemanager

type ResourceRemover interface {
	RemoveResource(ctx context.Context, resource resourcemodel.Resource) (err error)
	RemoveReport(ctx context.Context, resource resourcemodel.Resource) (r ResourceReport)
}

// структура ответсвенная за удаления ресурса
// имеет два метода: отчёт об изменениях при удалении ресурса, и само удаление ресурса
type resourceRemover struct {
	storage ResourceStorage
}

func NewRemover(resourceStorage ResourceStorage) ResourceRemover {
	return &resourceRemover{
		storage: resourceStorage,
	}
}

func (r *resourceRemover) RemoveResource(ctx context.Context, resource resourcemodel.Resource) (err error) {
	err = r.storage.DeleteResource(ctx, resource.ID)
	if err != nil {
		err = fmt.Errorf("failed delete resource, %w", err)
		return
	}
	return
}

func (r *resourceRemover) RemoveReport(ctx context.Context, resource resourcemodel.Resource) (report ResourceReport) {
	report.MinusResource(resource)
	return
}
