package resourcetype

// возможно его следовало бы назвать MachineType но уже поздо...
type ResourceType string

const (
	DB ResourceType = "db"
	VM ResourceType = "vm"
)
