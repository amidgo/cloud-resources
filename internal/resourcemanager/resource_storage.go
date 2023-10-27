package resourcemanager

import (
	"context"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
)

//go:generate mockgen -source resource_storage.go -destination resource_storage_mocks.go -package resourcemanager

type ResourceStorage interface {
	AddResource(ctx context.Context, addResource resourcemodel.AddResource) (resourcemodel.Resource, error)
	DeleteResource(ctx context.Context, id int) error
	UpdateResource(ctx context.Context, id int, updateResource resourcemodel.UpdateResource) error
	ResourceList(ctx context.Context) ([]*resourcemodel.Resource, error)
}

// враппер который фильтрует список ресурсов по типу
type ResourceStorageTypeWrapper struct {
	Type resourcetype.ResourceType
	ResourceStorage
}

func (r *ResourceStorageTypeWrapper) ResourceList(ctx context.Context) ([]*resourcemodel.Resource, error) {
	resourceList := make([]*resourcemodel.Resource, 0)
	list, err := r.ResourceStorage.ResourceList(ctx)
	for i := range list {
		if list[i].Type != r.Type {
			continue
		}
		resourceList = append(resourceList, list[i])
	}
	return resourceList, err
}
