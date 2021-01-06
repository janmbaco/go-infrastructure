package Config_test

import (
	"encoding/json"
	config2 "github.com/janmbaco/go-infrastructure/config"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"sync"
	"testing"
)

func TestNewFileConfigHandler(t *testing.T) {
	_ = disk.DeleteFile("config.json")
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
		if myConfig.Options != "other options" {
			t.Error("Options does not changed")
		}
		wg.Done()
	}
	configHandler.OnModifiedConfigSubscriber(&onModifiedConfigFunc)
	otherCofnig := &config{
		Options: "other options",
	}
	content, err := json.Marshal(otherCofnig)
	errorhandler.TryPanic(err)
	errorhandler.TryPanic(disk.CreateFile("config.json", content))
	wg.Wait()
}
