package server_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/server"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"

	fileConfigResolver "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc/resolver"
	errorsResolver "github.com/janmbaco/go-infrastructure/errors/ioc/resolver"
	serverResolver "github.com/janmbaco/go-infrastructure/server/ioc/resolver"
)

type configFile struct {
	Address  string `json:"address"`
	Address2 string `json:"address_2"`
}

var filePath = "config.json"
var address = ":8080"
var address2 = ":8090"

func TestNewListener(t *testing.T) {
	
	errorsResolver.GetErrorCatcher().TryFinally(func() {

		builder := serverResolver.GetListenerBuilder(
			fileConfigResolver.GetFileConfigHandler(
				filePath,  
				&configFile{Address: address, Address2: address2}, 
			),
		) 

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
		PutTheSamePortInConfig()
		go func() {
			<-time.After(50 * time.Millisecond)
			listener.Stop()
			listener2.Stop()
		}()
		err1 := <-finishListener
		err2 := <-finishListener2
		errorschecker.TryPanic(err1)
		errorschecker.TryPanic(err2)

	}, func() {
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
