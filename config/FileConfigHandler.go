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
	"github.com/janmbaco/go-infrastructure/events"
)

const maxTries = 10

// FileConfigHandler  is the object responsible for loading the configuration from a file and observing its changes
type FileConfigHandler struct {
	*configSubscriber
	filePath             string
	dataconfig           interface{}
	onModifiedConfigFile func()
	fileChangeNotifier   disk.FileChangedNotifier
	isSubscribed         bool
}

// NewFileConfigHandler returns a FielConfigHandler
func NewFileConfigHandler(filePath string) *FileConfigHandler {
	return &FileConfigHandler{filePath: filePath, configSubscriber: &configSubscriber{eventPublisher: events.NewPublisher()}, fileChangeNotifier: disk.NewFileChangedNotifier(filePath)}
}

// Load loads the default configuration if the file not exits.
func (fileConfigHandler *FileConfigHandler) Load(defaults interface{}) {
	fileConfigHandler.dataconfig = defaults
	if !disk.ExistsPath(fileConfigHandler.filePath) {
		fileConfigHandler.writeFile()
	}
	fileConfigHandler.onModifiedConfigFile = func() {
		errorhandler.TryCatchError(
			func() {
				fileConfigHandler.readFile()
				if fileConfigHandler.onModifyingConfigPublish() {
					panic(fileConfigHandler.cancelMessage)
				}
				fileConfigHandler.onModifiedConfigPublish()
			},
			func(err error) {
				_ = disk.Copy(fileConfigHandler.filePath, fileConfigHandler.filePath+".badconfig")
				fileConfigHandler.writeFile()
			})
	}
	fileConfigHandler.fileChangeNotifier.Subscribe(&fileConfigHandler.onModifiedConfigFile)
	fileConfigHandler.isSubscribed = true
	fileConfigHandler.readFile()
}

func (fileConfigHandler *FileConfigHandler) readFile() {
	var content []byte
	var err error
	try := 1
	for len(content) == 0 && try < maxTries {
		content, err = ioutil.ReadFile(fileConfigHandler.filePath)
		errorhandler.TryPanic(err)
		try++
	}
	ret := reflect.New(reflect.TypeOf(fileConfigHandler.dataconfig)).Interface()
	errorhandler.TryPanic(json.Unmarshal(content, ret))
	errorhandler.TryPanic(copier.Copy(fileConfigHandler.dataconfig, ret))
}

func (fileConfigHandler *FileConfigHandler) writeFile() {
	var content []byte
	var err error
	content, err = json.Marshal(fileConfigHandler.dataconfig)
	errorhandler.TryPanic(err)
	_ = os.Mkdir(filepath.Dir(fileConfigHandler.filePath), 0666)
	errorhandler.TryPanic(disk.CreateFile(fileConfigHandler.filePath, content))
}
