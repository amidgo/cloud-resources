package executor

import (
	"fmt"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
)

//go:generate mockgen -source logger.go -destination logger_mocks.go -package executor

type Logger interface {
	Printf(string, ...any)
}

type ResourceTypeLoggerWrapper struct {
	ResourceType resourcetype.ResourceType
	Logger
}

func (l ResourceTypeLoggerWrapper) Printf(format string, args ...any) {
	typePrefix := fmt.Sprintf("resource type %s, ", l.ResourceType)
	l.Logger.Printf(typePrefix+format, args...)
}
