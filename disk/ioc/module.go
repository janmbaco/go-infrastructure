package ioc

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/disk"
)

// DiskModule implements Module for disk services
type DiskModule struct{}

// NewDiskModule creates a new disk module
func NewDiskModule() *DiskModule {
	return &DiskModule{}
}

// RegisterServices registers all disk services
func (m *DiskModule) RegisterServices(register dependencyinjection.Register) error {
	dependencyinjection.RegisterTypeWithParams[disk.FileChangedNotifier](
		register,
		disk.NewFileChangedNotifier,
		map[int]string{0: "filePath"},
	)

	return nil
}
