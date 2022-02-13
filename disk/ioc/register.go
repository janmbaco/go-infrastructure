package ioc

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/disk"
)

func init(){
	static.Container.Register().AsType(new(disk.FileChangedNotifier), disk.NewFileChangedNotifier, map[uint]string{0: "filePath"})
}