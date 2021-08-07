package server_test

import (
	"encoding/json"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	dependencyinjection_test "github.com/janmbaco/go-infrastructure/dependencyinjection.test"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/server"
	"net/http"
	"testing"
	"time"
)

type configFile struct {
	Address  string `json:"address"`
	Address2 string `json:"address_2"`
}

var filePath = "config.json"
var address = ":8080"
var address2 = ":8090"

func TestNewListener(t *testing.T) {
	container := dependencyinjection.NewContainer()
	dependencyinjection_test.Registerfacade(container.Register())

	errorCatcher := container.Resolver().Type(new(errors.ErrorCatcher), nil).(errors.ErrorCatcher)
	errorCatcher.TryFinally(func() {
		builder := container.Resolver().Type(
			new(server.ListenerBuilder),
			map[string]interface{}{
				"filePath": filePath,
				"defaults": &configFile{Address: address, Address2: address2},
			},
		).(server.ListenerBuilder)

		builder.SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) {
			serverSetter.Name = "Primero"
			mux := http.NewServeMux()
			mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("started..."))
			}))
			serverSetter.Addr = config.(*configFile).Address
			serverSetter.Handler = mux
		})

		listener := builder.GetListener()
		finishListener := listener.Start()

		builder.SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) {
			serverSetter.Name = "Segundo"
			mux := http.NewServeMux()
			mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("started..."))
			}))
			serverSetter.Addr = config.(*configFile).Address2
			serverSetter.Handler = mux
		})
		listener2 := builder.GetListener()
		finishListener2 := listener2.Start()
		go func() {
			PutTheSamePortInConfig()
			<-time.After(800 * time.Millisecond)
			listener.Stop()
			listener2.Stop()
		}()
		<-finishListener
		<-finishListener2
	}, func() {
		<-time.After(5 * time.Millisecond)
		disk.DeleteFile(filePath)
	})
}

func PutTheSamePortInConfig() {
	lcontent, lerr := json.MarshalIndent(&configFile{
		Address:  address,
		Address2: address,
	}, "", "\t")
	errorschecker.TryPanic(lerr)
	disk.CreateFile(filePath, lcontent)
}
