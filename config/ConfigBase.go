package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errorhandler"
)

const (
	configFileChangedEvent = "ConfigFileChangedEvent"
)

type ConfigBase struct {
	filePath                   string
	configChangedSubscriptions []func()
	onModifiedConfigFile       func()
	watcherActive              bool
	defaults                   *ConfigBase
}

func NewConfigBase(defaults *ConfigBase) *ConfigBase {
	return &ConfigBase{defaults: defaults}
}

func (this *ConfigBase) Load(filePath string) {
	this.filePath = filePath
	if !disk.ExistsPath(this.filePath) {
		this.writeFile(true)
		disk.NewFileChangedNotifier(this.filePath).Subscribe(this.onModifiedConfigFile)
	}
	errorhandler.TryPanic(this.readFile())
}

func (this *ConfigBase) readFile() error {
	//it's possible that it still modifying
	time.Sleep(100)
	content, err := ioutil.ReadFile(this.filePath)
	if err != nil {
		err = json.Unmarshal(content, this)
	}
	return err
}

func (this *ConfigBase) modifiedConfigFie() {

	time.Sleep(100)
	errorhandler.TryCatchError(
		func() {
			errorhandler.TryPanic(this.readFile())
			this.configFileChangedPublish()
		},
		func(err error) {
			this.writeFile(false)
		})
}

func (this *ConfigBase) writeFile(byDefaults bool) {
	var content []byte
	var err error
	if byDefaults {
		content, err = json.Marshal(this.defaults)
	} else {
		content, err = json.Marshal(this)
	}
	errorhandler.TryPanic(err)
	_ = os.Mkdir(filepath.Dir(this.filePath), 0666)
	errorhandler.TryPanic(disk.CreateFile(this.filePath, content))
}

func (this *ConfigBase) Subscribe(subscribeFunc func(config *ConfigBase)) {
	onModifiedConfigFile := func() {
		subscribeFunc(this)
	}
	this.configChangedSubscriptions = append(this.configChangedSubscriptions, onModifiedConfigFile)
}

func (this *ConfigBase) configFileChangedPublish() {
	for _, f := range this.configChangedSubscriptions {
		errorhandler.OnErrorContinue(f)
	}
}
