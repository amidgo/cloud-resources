package pricemodel

import "github.com/amidgo/cloud-resources/internal/model/resourcetype"

type Price struct {
	ID   int                       `json:"id"`
	Cost float32                   `json:"cost"`
	CPU  float32                   `json:"cpu"`
	RAM  float32                   `json:"ram"`
	Name string                    `json:"name"`
	Type resourcetype.ResourceType `json:"type"`
}

func NewPrice(
	id int,
	cost, cpu, ram float32,
	name string,
	resourceType resourcetype.ResourceType,
) Price {
	return Price{
		ID:   id,
		Cost: cost,
		CPU:  cpu,
		RAM:  ram,
		Name: name,
		Type: resourceType,
	}
}

func (p Price) IsZero() bool {
	return p == Price{}
}
