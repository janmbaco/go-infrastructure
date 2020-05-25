package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/janmbaco/copier"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/event"
)

const maxTries = 10

type FileConfigHandler struct {
	*ConfigSubscriber
	filePath   string
	dataconfig interface{}
}

func NewFileConfigHandler(filePath string) *FileConfigHandler {
	return &FileConfigHandler{filePath: filePath, ConfigSubscriber: &ConfigSubscriber{eventPublisher: event.NewEventPublisher()}}
}

func (this *FileConfigHandler) Load(defaults interface{}) {
	this.dataconfig = defaults
	if !disk.ExistsPath(this.filePath) {
		this.writeFile()
	}
	disk.NewFileChangedNotifier(this.filePath).Subscribe(this.onModifiedConfigFile)
	this.readFile()
}

func (this *FileConfigHandler) readFile() {
	var content []byte
	var err error
	try := 1
	for len(content) == 0 || try == maxTries {
		content, err = ioutil.ReadFile(this.filePath)
		errorhandler.TryPanic(err)
		try++
	}
	ret := reflect.New(reflect.TypeOf(this.dataconfig)).Interface()
	errorhandler.TryPanic(json.Unmarshal(content, ret))
	errorhandler.TryPanic(copier.Copy(this.dataconfig, ret))
}

func (this *FileConfigHandler) writeFile() {
	var content []byte
	var err error
	content, err = json.Marshal(this.dataconfig)
	errorhandler.TryPanic(err)
	_ = os.Mkdir(filepath.Dir(this.filePath), 0666)
	errorhandler.TryPanic(disk.CreateFile(this.filePath, content))
}

func (this *FileConfigHandler) onModifiedConfigFile() {
	errorhandler.TryCatchError(
		func() {
			this.readFile()
			this.onModifiedConfigPublish()
		},
		func(err error) {
			this.writeFile()
		})
}
