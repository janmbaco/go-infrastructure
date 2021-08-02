package Config_test

import (
	"encoding/json"
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/configuration/events"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	dependencyinjection_test "github.com/janmbaco/go-infrastructure/dependencyinjection.test"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"sync"
	"testing"
	"time"
)

type configFile struct {
	Options string `json:"options"`
}

var filePath = "configuration.json"
var initialOptions = "New Options"
var BadOptions = "Bad Options"
var GoodOtions = "Good Options"
var CancelMessage = "Is Bad Options"

func TestNewFileConfigHandler(t *testing.T) {
	container := dependencyinjection.NewContainer()
	dependencyinjection_test.Registerfacade(container)

	errorCatcher := container.Resolver.Type(new(errors.ErrorCatcher), nil).(errors.ErrorCatcher)
	errorCatcher.TryCatchErrorAndFinally(func() {

		configHandler := container.Resolver.Type(
			new(configuration.ConfigHandler),
			map[string]interface{}{
				"filePath": filePath,
				"defaults": &configFile{Options: initialOptions},
			},
		).(configuration.ConfigHandler)

		wg := sync.WaitGroup{}
		wg.Add(1)

		onModifyingConfig := func(eventArgs *events.ModifyingEventArgs) {
			t.Logf("Modifying: %v to %v", configHandler.GetConfig().(*configFile).Options, eventArgs.Config.(*configFile).Options)
			if eventArgs.Config.(*configFile).Options == BadOptions {
				eventArgs.Cancel = true
				eventArgs.CancelMessage = CancelMessage
			}
		}
		configHandler.ModifyingSubscribe(&onModifyingConfig)

		onModificationCanceled := func(eventArgs *events.ModificationCanceledEventArgs) {
			t.Logf("ModificationCanceled: %v ", eventArgs.CancelMessage)
			if eventArgs.CancelMessage != CancelMessage {
				t.Errorf("The message was not the correct! message: %v, correct: %v", eventArgs.CancelMessage, CancelMessage)
			}
			wg.Done()
		}

		configHandler.ModificationCanceledSubscribe(&onModificationCanceled)

		onModifiedConfig := func() {
			t.Logf("Modified: %v ", configHandler.GetConfig().(*configFile).Options)
			if configHandler.GetConfig().(*configFile).Options != GoodOtions {
				t.Errorf("The message was not the correct! message: %v, correct: %v", configHandler.GetConfig().(*configFile).Options, GoodOtions)
			}
			wg.Done()
		}

		configHandler.ModifiedSubscribe(&onModifiedConfig)

		modifyWith(BadOptions)
		wg.Wait()
		wg.Add(1)

		modifyWith(GoodOtions)
		wg.Wait()

		t.Log("Test finalized!")

	}, func(err error) {
		t.Error(err)
	}, func() {
		<-time.After(5 * time.Millisecond)
		disk.DeleteFile("configuration.json")
		disk.DeleteFile("configuration.json.badconfig")
	})
}

func modifyWith(text string) {
	lcontent, lerr := json.MarshalIndent(&configFile{
		Options: text,
	}, "", "\t")
	errorschecker.TryPanic(lerr)
	disk.CreateFile(filePath, lcontent)
}
