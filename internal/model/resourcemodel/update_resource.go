package resourcemodel

import "github.com/amidgo/cloud-resources/internal/model/resourcetype"

type UpdateResource struct {
	CPU  float32                   `json:"cpu"`
	RAM  float32                   `json:"ram"`
	Type resourcetype.ResourceType `json:"type"`
}

func NewUpdateResource(cpu, ram float32, resourceType resourcetype.ResourceType) UpdateResource {
	return UpdateResource{
		CPU:  cpu,
		RAM:  ram,
		Type: resourceType,
	}
}
