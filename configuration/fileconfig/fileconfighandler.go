package fileconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/janmbaco/copier"
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/configuration/events"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
	"github.com/janmbaco/go-infrastructure/logs"
)

const maxTries = 10

type fileConfigHandler struct {
	errorCatcher errors.ErrorCatcher
	oldconfig    interface{}
	dataconfig   interface{}
	newConfig    interface{}
	fromFile     interface{}
	*events.RestoredEventHandler
	*events.ModifiedEventHandler
	eventManager *eventsmanager.EventManager
	stopRefresh  chan bool
	filePath     string
	period       configuration.Period
	mutex        sync.RWMutex
	isFreezed    bool
}

// NewFileConfigHandler returns a ConfigHandler
func NewFileConfigHandler(filePath string, defaults interface{}, errorCatcher errors.ErrorCatcher, eventManager *eventsmanager.EventManager, filechangeNotifier disk.FileChangedNotifier, logger logs.Logger) (configuration.ConfigHandler, error) {
	modifiedSubs := eventsmanager.NewSubscriptions[events.ModifiedEvent]()
	restoredSubs := eventsmanager.NewSubscriptions[events.RestoredEvent]()
	modifiedPub := eventsmanager.NewPublisher(modifiedSubs, logger)
	restoredPub := eventsmanager.NewPublisher(restoredSubs, logger)
	eventsmanager.Register(eventManager, modifiedPub)
	eventsmanager.Register(eventManager, restoredPub)
	fileConfigHandler := &fileConfigHandler{
		filePath:             filePath,
		eventManager:         eventManager,
		ModifiedEventHandler: events.NewModifiedEventHandler(modifiedSubs),
		RestoredEventHandler: events.NewRestoredEventHandler(restoredSubs),
		errorCatcher:         errorCatcher,
		stopRefresh:          make(chan bool, 1),
	}
	fileConfigHandler.dataconfig = reflect.New(reflect.TypeOf(defaults).Elem()).Interface()
	if err := copier.Copy(fileConfigHandler.dataconfig, defaults); err != nil {
		return nil, err
	}
	if !disk.ExistsPath(fileConfigHandler.filePath) {
		if err := fileConfigHandler.writeFile(); err != nil {
			return nil, err
		}
	}
	if err := filechangeNotifier.Subscribe(fileConfigHandler.onModifiedConfigFile); err != nil {
		return nil, err
	}
	if err := fileConfigHandler.readFile(); err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(fileConfigHandler.fromFile, fileConfigHandler.dataconfig) {
		if err := copier.Copy(fileConfigHandler.dataconfig, fileConfigHandler.fromFile); err != nil {
			return nil, err
		}
	}
	return fileConfigHandler, nil
}

// SetRefreshTime sets the period to refresh the config
func (f *fileConfigHandler) SetRefreshTime(period configuration.Period) error {
	if f.stopRefresh == nil {
		return newFileConfigHandlerError(UnexpectedError, "handler not initialized", nil)
	}
	f.stopRefresh <- true
	f.period = period
	go f.refreshLoop()
	return nil
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
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.dataconfig
}

// SetConfig updates the configuration and writes it to file
func (f *fileConfigHandler) SetConfig(newConfig interface{}) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.oldconfig = f.createConfig()
	if err := copier.Copy(f.oldconfig, f.dataconfig); err != nil {
		return f.pipeError(err)
	}

	// Si newConfig es un map (vino de JSON), necesitamos re-marshalearlo y unmarshalearlo
	// para convertirlo al tipo correcto
	configBytes, err := json.Marshal(newConfig)
	if err != nil {
		return f.pipeError(err)
	}

	tempConfig := f.createConfig()
	if err := json.Unmarshal(configBytes, tempConfig); err != nil {
		return f.pipeError(err)
	}

	if err := copier.Copy(f.dataconfig, tempConfig); err != nil {
		return f.pipeError(err)
	}

	if err := f.writeFile(); err != nil {
		return f.pipeError(err)
	}

	eventsmanager.Publish(f.eventManager, events.ModifiedEvent{})
	return nil
}

// ForceRefresh forces a refresh of the configuration
func (f *fileConfigHandler) ForceRefresh() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.newConfig != nil && !reflect.DeepEqual(f.newConfig, f.dataconfig) {
		if err := copier.Copy(f.dataconfig, f.newConfig); err != nil {
			return f.pipeError(err)
		}
	}
	return nil
}

// CanRestore indicates if the config can be restored
func (f *fileConfigHandler) CanRestore() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.oldconfig != nil
}

// Restore restores the configuration to an older version
func (f *fileConfigHandler) Restore() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.oldconfig == nil {
		return f.pipeError(newFileConfigHandlerError(OldConfigNilError, "it is no posible restore to old config because is nil", nil))
	}
	f.newConfig = f.createConfig()
	if err := copier.Copy(f.newConfig, f.dataconfig); err != nil {
		return f.pipeError(err)
	}
	if err := copier.Copy(f.dataconfig, f.oldconfig); err != nil {
		return f.pipeError(err)
	}
	f.oldconfig = nil
	if err := f.writeFile(); err != nil {
		return f.pipeError(err)
	}
	eventsmanager.Publish(f.eventManager, events.RestoredEvent{})
	return nil
}

func (f *fileConfigHandler) readFile() error {
	var content []byte
	var err error
	try := 1
	for len(content) == 0 && try < maxTries {
		content, err = os.ReadFile(f.filePath)
		if err != nil {
			return err
		}
		try++
	}
	ret := reflect.New(reflect.TypeOf(f.dataconfig)).Interface()
	if err := json.Unmarshal(content, ret); err != nil {
		return err
	}
	f.fromFile = f.createConfig()
	return copier.Copy(f.fromFile, ret)
}

func (f *fileConfigHandler) writeFile() error {
	var content []byte
	var err error
	content, err = json.MarshalIndent(f.dataconfig, "", "\t")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(f.filePath), 0o755); err != nil {
		return err
	}
	return disk.CreateFile(f.filePath, content)
}

func (f *fileConfigHandler) onModifiedConfigFile() {
	f.errorCatcher.TryCatchError( //nolint:errcheck // TryCatchError handles errors internally via callback
		func() error {
			f.mutex.Lock()
			defer f.mutex.Unlock()

			if err := f.readFile(); err != nil {
				return err
			}
			if !reflect.DeepEqual(f.fromFile, f.dataconfig) {
				f.newConfig = f.createConfig()
				if err := copier.Copy(f.newConfig, f.fromFile); err != nil {
					return err
				}
				if !f.isFreezed {
					f.oldconfig = f.createConfig()
					if err := copier.Copy(f.oldconfig, f.dataconfig); err != nil {
						return err
					}
					if err := copier.Copy(f.dataconfig, f.newConfig); err != nil {
						return err
					}
					f.newConfig = nil
					eventsmanager.Publish(f.eventManager, events.ModifiedEvent{})
				}
			}
			return nil
		},
		func(err error) {
			_ = f.recoveryFile() //nolint:errcheck // recovery errors are not actionable when already in error state
		})
}

func (f *fileConfigHandler) createConfig() interface{} {
	return reflect.New(reflect.TypeOf(f.dataconfig).Elem()).Interface()
}

func (f *fileConfigHandler) recoveryFile() error {
	if err := disk.Copy(f.filePath, f.filePath+".badconfig"); err != nil {
		return err
	}
	return f.writeFile()
}

func (f *fileConfigHandler) refreshLoop() {
	exit := false
	for {
		select {
		case <-time.After(time.Minute):
			if f.period.IsFinished() {
				_ = f.ForceRefresh() //nolint:errcheck // periodic refresh errors should not crash the system
			}
		case <-f.stopRefresh:
			exit = true
		}
		if exit {
			break
		}
	}
}

func (f *fileConfigHandler) pipeError(err error) error {
	resultError := err

	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); !errType.Implements(reflect.TypeOf((*HandlerError)(nil)).Elem()) {
		resultError = newFileConfigHandlerError(UnexpectedError, err.Error(), err)
	}

	return resultError
}
