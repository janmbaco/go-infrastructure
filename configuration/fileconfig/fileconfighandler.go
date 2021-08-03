package fileconfig

import (
	"encoding/json"
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/configuration/events"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/janmbaco/copier"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

const maxTries = 10

type fileConfigHandler struct {
	*events.ModifiedEventHandler
	*events.ModifyingEventHandler
	*events.ModificationCanceledEventHandler
	*events.RestoredEventHandler
	filePath               string
	oldconfig              interface{}
	dataconfig             interface{}
	newConfig              interface{}
	isFreezed              bool
	fromFile               interface{}
	onModifiedConfigFileFn func()
	publisher              eventsmanager.Publisher
	period                 configuration.Period
	errorCatcher           errors.ErrorCatcher
	errorDefer             errors.ErrorDefer
	stopRefresh            chan bool
}

// NewFileConfigHandler returns a ConfigHandler
func NewFileConfigHandler(filePath string, defaults interface{}, errorCatcher errors.ErrorCatcher, errorThrower errors.ErrorThrower) configuration.ConfigHandler {
	errorschecker.CheckNilParameter(map[string]interface{}{"defaults": defaults, "errorCatcher": errorCatcher, "errorThrower": errorThrower})

	subscriptions := eventsmanager.NewSubscriptions(errorThrower)
	errorHandler := errors.NewErrorDefer(errorThrower, &fileConfigHandleErrorPipe{})
	fileConfigHandler := &fileConfigHandler{
		filePath:                         filePath,
		publisher:                        eventsmanager.NewPublisher(subscriptions, errorCatcher),
		ModifiedEventHandler:             events.NewModifiedEventHandler(subscriptions),
		ModifyingEventHandler:            events.NewModifyingEventHandler(subscriptions),
		RestoredEventHandler:             events.NewRestoredEventHandler(subscriptions),
		ModificationCanceledEventHandler: events.NewModificationCanceledEventHandler(subscriptions),
		errorCatcher:                     errorCatcher,
		stopRefresh:                      make(chan bool, 1),
		errorDefer:                       errorHandler,
	}
	fileConfigHandler.dataconfig = reflect.New(reflect.TypeOf(defaults).Elem()).Interface()
	errorschecker.TryPanic(copier.Copy(fileConfigHandler.dataconfig, defaults))
	if !disk.ExistsPath(fileConfigHandler.filePath) {
		fileConfigHandler.writeFile()
	}
	disk.NewFileChangedNotifier(fileConfigHandler.filePath, errorCatcher, errorThrower).Subscribe(fileConfigHandler.onModifiedConfigFile)
	fileConfigHandler.readFile()
	return fileConfigHandler
}

// SetRefreshTime sets the period to refresh the config
func (f *fileConfigHandler) SetRefreshTime(period configuration.Period) {
	errorschecker.CheckNilParameter(map[string]interface{}{"period": period})
	defer f.errorDefer.TryThrowError()
	f.stopRefresh <- true
	f.period = period
	go f.refreshLoop()
}

// Freeze causes configuration changes to not be made until the end of the specified period
func (f *fileConfigHandler) Freeze() {
	f.isFreezed = true
}

// Unfreeze causes configuration changes to be made when they occur
func (f *fileConfigHandler) Unfreeze() {
	f.isFreezed = false
}

// GetConfig get the current config applied
func (f *fileConfigHandler) GetConfig() interface{} {
	return f.dataconfig
}

// ForceRefresh forces a refresh of the configuration
func (f *fileConfigHandler) ForceRefresh() {
	defer f.errorDefer.TryThrowError()
	if f.newConfig != nil && !reflect.DeepEqual(f.newConfig, f.dataconfig) {
		errorschecker.TryPanic(copier.Copy(f.dataconfig, f.newConfig))
	}
}

// CanRestore indicates if the config can be restored
func (f *fileConfigHandler) CanRestore() bool {
	return f.oldconfig != nil
}

// Restore restores the configuration to an older version
func (f *fileConfigHandler) Restore() {
	defer f.errorDefer.TryThrowError()
	if !f.CanRestore() {
		panic(newFileConfigHandlerError(OldConfigNilError, "it is no posible restore to old config because is nil"))
	}
	f.newConfig = f.createConfig()
	errorschecker.TryPanic(copier.Copy(f.newConfig, f.dataconfig))
	errorschecker.TryPanic(copier.Copy(f.dataconfig, f.oldconfig))
	f.oldconfig = nil
	f.publisher.Publish(&events.RestoredEvent{})
}

func (f *fileConfigHandler) readFile() {
	var content []byte
	var err error
	try := 1
	for len(content) == 0 && try < maxTries {
		content, err = ioutil.ReadFile(f.filePath)
		errorschecker.TryPanic(err)
		try++
	}
	ret := reflect.New(reflect.TypeOf(f.dataconfig)).Interface()
	errorschecker.TryPanic(json.Unmarshal(content, ret))
	f.fromFile = f.createConfig()
	errorschecker.TryPanic(copier.Copy(f.fromFile, ret))
}

func (f *fileConfigHandler) writeFile() {
	var content []byte
	var err error
	content, err = json.MarshalIndent(f.dataconfig, "", "\t")
	errorschecker.TryPanic(err)
	_ = os.Mkdir(filepath.Dir(f.filePath), 0666)
	disk.CreateFile(f.filePath, content)
}

func (f *fileConfigHandler) onModifiedConfigFile() {
	f.errorCatcher.TryCatchError(
		func() {
			f.readFile()
			if !reflect.DeepEqual(f.fromFile, f.dataconfig) {
				eventArgs := &events.ModifyingEventArgs{Config: f.fromFile}
				f.publisher.Publish(&events.ModifyingEvent{
					EventArgs: eventArgs,
				})
				if eventArgs.Cancel {
					f.recoveryFile()
					f.publisher.Publish(&events.ModificationCanceledEvent{EventArgs: &events.ModificationCanceledEventArgs{CancelMessage: eventArgs.CancelMessage}})
				} else {

					f.newConfig = f.createConfig()
					errorschecker.TryPanic(copier.Copy(f.newConfig, f.fromFile))
					if !f.isFreezed {
						f.oldconfig = f.createConfig()
						errorschecker.TryPanic(copier.Copy(f.oldconfig, f.dataconfig))
						errorschecker.TryPanic(copier.Copy(f.dataconfig, f.newConfig))
						f.newConfig = nil
						f.publisher.Publish(&events.ModifiedEvent{})
					}
				}
			}
		},
		func(err error) {
			f.recoveryFile()
		})
}

func (f *fileConfigHandler) createConfig() interface{} {
	return reflect.New(reflect.TypeOf(f.dataconfig).Elem()).Interface()
}

func (f *fileConfigHandler) recoveryFile() {
	disk.Copy(f.filePath, f.filePath+".badconfig")
	f.writeFile()
}

func (f *fileConfigHandler) refreshLoop() {
	for {
		select {
		case <-time.After(time.Minute):
			if f.period.IsFinished() {
				f.ForceRefresh()
			}
		case <-f.stopRefresh:
			break
		}
	}
}
