package server_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc"
	fileConfigResolver "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc/resolver"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/disk"
	diskIoc "github.com/janmbaco/go-infrastructure/disk/ioc"
	errorsIoc "github.com/janmbaco/go-infrastructure/errors/ioc"
	eventsIoc "github.com/janmbaco/go-infrastructure/eventsmanager/ioc"
	logsIoc "github.com/janmbaco/go-infrastructure/logs/ioc"
	"github.com/janmbaco/go-infrastructure/server"
	serverIoc "github.com/janmbaco/go-infrastructure/server/ioc"
	serverResolver "github.com/janmbaco/go-infrastructure/server/ioc/resolver"
)

// / <summary>
// / Implementa la responsabilidad de probar el comportamiento del Listener con recuperación de pánico.
// / </summary>
type ListenerTests struct {
	configFilePath string
	configHandler  configuration.ConfigHandler
	Resolver       dependencyinjection.Resolver
}

type testConfig struct {
	Address  string `json:"address"`
	Address2 string `json:"address_2"`
}

const (
	testConfigFile = "listener_test_config.json"
	firstAddress   = ":18080"
	secondAddress  = ":18090"
)

func (lt *ListenerTests) setup(t *testing.T) {
	lt.configFilePath = testConfigFile

	// Setup IoC container
	container := dependencyinjection.NewBuilder().
		AddModule(logsIoc.NewLogsModule()).
		AddModule(errorsIoc.NewErrorsModule()).
		AddModule(eventsIoc.NewEventsModule()).
		AddModule(diskIoc.NewDiskModule()).
		AddModule(ioc.NewConfigurationModule()).
		AddModule(serverIoc.NewServerModule()).
		MustBuild()

	resolver := container.Resolver()

	lt.Resolver = resolver
	lt.configHandler = fileConfigResolver.GetFileConfigHandler(resolver,
		lt.configFilePath,
		&testConfig{Address: firstAddress, Address2: secondAddress})
}

func (lt *ListenerTests) teardown() {
	_ = disk.DeleteFile(lt.configFilePath)
}

func (lt *ListenerTests) createListener(name string, addressSelector func(*testConfig) string) (server.Listener, error) {
	builder := serverResolver.GetListenerBuilder(lt.Resolver, lt.configHandler)
	builder.SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) error {
		serverSetter.Name = name
		mux := http.NewServeMux()
		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("started"))
		}))
		serverSetter.Addr = addressSelector(config.(*testConfig))
		serverSetter.Handler = mux
		return nil
	})
	return builder.GetListener()
}

func (lt *ListenerTests) updateConfigToDuplicatePort(t *testing.T) {
	content, err := json.MarshalIndent(&testConfig{
		Address:  firstAddress,
		Address2: firstAddress,
	}, "", "\t")
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}
	_ = disk.CreateFile(lt.configFilePath, content)
}

// / <summary>
// / Listener_WhenPortConflictOccurs_ThenRecoversByRestoringConfig valida que el listener recupere
// / automáticamente de un conflicto de puerto mediante restauración de configuración.
// / </summary>
func TestListener_WhenPortConflictOccurs_ThenRecoversByRestoringConfig(t *testing.T) {
	// Arrange
	lt := &ListenerTests{}
	lt.setup(t)
	defer lt.teardown()

	listener1, err := lt.createListener("Listener1", func(c *testConfig) string { return c.Address })
	if err != nil {
		t.Fatalf("Arrange failed: %v", err)
	}
	listener2, err := lt.createListener("Listener2", func(c *testConfig) string { return c.Address2 })
	if err != nil {
		t.Fatalf("Arrange failed: %v", err)
	}
	finish1 := listener1.Start()
	finish2 := listener2.Start()

	// Act
	lt.updateConfigToDuplicatePort(t)

	go func() {
		time.Sleep(50 * time.Millisecond)
		listener1.Stop()
		listener2.Stop()
	}()

	// Assert
	err1 := <-finish1
	err2 := <-finish2

	if err1 != nil {
		t.Errorf("Assert failed: Listener1 terminó con error: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Assert failed: Listener2 terminó con error: %v", err2)
	}
}
