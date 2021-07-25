package Config_test

import (
	"encoding/json"
	config2 "github.com/janmbaco/go-infrastructure/config"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"sync"
	"testing"
	"time"
)

func TestNewFileConfigHandler(t *testing.T) {
	_ = disk.DeleteFile("config.json")
	_ = disk.DeleteFile("config.json.badconfig")
	type config struct {
		Options string `json:"options"`
	}

	myConfig := &config{
		Options: "New options",
	}

	configHandler := config2.NewFileConfigHandler("config.json")
	configHandler.Load(myConfig)
	wg := sync.WaitGroup{}
	wg.Add(1)
	onModifiedConfigFunc := func() {
		t.Log(myConfig.Options)
		if myConfig.Options != "other options" {
			t.Log("Options does not changed")
		} else {
			t.Log("Opción correcta.")
		}
		wg.Done()
	}
	configHandler.OnModifyingConfigSubscriber(func(newConfig interface{}) {
		t.Log(newConfig.(*config).Options)
		if newConfig.(*config).Options != "other options" {
			t.Log("No puedo aceptarlo")
			go func() {
				<-time.After(10 * time.Millisecond)
				wg.Done()
			}()
			panic("No puedo aceptarlo")
		} else {
			t.Log("aceptado")
			go func() {
				<-time.After(10 * time.Millisecond)
				wg.Done()
			}()
		}
	})

	configHandler.OnModifiedConfigSubscriber(onModifiedConfigFunc)
	lotroCofnig := &config{
		Options: "no another options",
	}
	lcontent, lerr := json.MarshalIndent(lotroCofnig, "", "\t")
	errorhandler.TryPanic(lerr)
	errorhandler.TryPanic(disk.CreateFile("config.json", lcontent))
	wg.Wait()
	wg.Add(2)

	otherCofnig := &config{
		Options: "other options",
	}
	content, err := json.MarshalIndent(otherCofnig, "", "\t")
	errorhandler.TryPanic(err)
	errorhandler.TryPanic(disk.CreateFile("config.json", content))
	wg.Wait()
}
