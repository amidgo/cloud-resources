package resourcemodel

import "github.com/amidgo/cloud-resources/internal/model/resourcetype"

type Resource struct {
	ID          int                       `json:"id"`
	Cost        float32                   `json:"cost"`
	CPU         float32                   `json:"cpu"`
	CPULoad     float32                   `json:"cpu_load"`
	RAM         float32                   `json:"ram"`
	RAMLoad     float32                   `json:"ram_load"`
	Failed      bool                      `json:"failed"`
	FailedUntil string                    `json:"failed_until"`
	Type        resourcetype.ResourceType `json:"type"`
}

func NewResource(id int, cost, cpu, ram, cpuLoad, ramLoad float32, failed bool, failedUntil string, resourceType resourcetype.ResourceType) Resource {
	return Resource{
		ID:          id,
		Cost:        cost,
		CPU:         cpu,
		CPULoad:     cpuLoad,
		RAM:         ram,
		RAMLoad:     ramLoad,
		Failed:      failed,
		FailedUntil: failedUntil,
		Type:        resourceType,
	}
}
