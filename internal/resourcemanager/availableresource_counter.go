package resourcemanager

type AvailabelResourceCounter struct {
	availableResourceCount int
}

func (r *AvailabelResourceCounter) CanIncFailed() bool {
	return r.availableResourceCount > 1
}

func (r *AvailabelResourceCounter) IncFailed() {
	r.availableResourceCount--
}

func (r *AvailabelResourceCounter) SetResourceCount(count int) {
	r.availableResourceCount = count
}
