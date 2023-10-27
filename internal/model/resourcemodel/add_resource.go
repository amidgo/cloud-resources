package resourcemodel

import "github.com/amidgo/cloud-resources/internal/model/resourcetype"

type AddResource struct {
	CPU  float32                   `json:"cpu"`
	RAM  float32                   `json:"ram"`
	Type resourcetype.ResourceType `json:"type"`
}

func NewAddResource(cpu, ram float32, resourceType resourcetype.ResourceType) AddResource {
	return AddResource{
		CPU:  cpu,
		RAM:  ram,
		Type: resourceType,
	}
}
