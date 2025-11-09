package server_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/configuration/fileconfig"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errors"
	errorsIoc "github.com/janmbaco/go-infrastructure/errors/ioc"
	errorsResolver "github.com/janmbaco/go-infrastructure/errors/ioc/resolver"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
	eventsIoc "github.com/janmbaco/go-infrastructure/eventsmanager/ioc"
	"github.com/janmbaco/go-infrastructure/logs"
	logsIoc "github.com/janmbaco/go-infrastructure/logs/ioc"
	logsResolver "github.com/janmbaco/go-infrastructure/logs/ioc/resolver"
	"github.com/janmbaco/go-infrastructure/server"
)

// / <summary>
// / Implementa la responsabilidad de probar el comportamiento del Listener con recuperación de pánico.
// / </summary>
type ListenerTests struct {
	configFilePath string
	logger         logs.Logger
	errorCatcher   errors.ErrorCatcher
	configHandler  configuration.ConfigHandler
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
		MustBuild()

	resolver := container.Resolver()

	// Resolve services
	lt.logger = logsResolver.GetLogger(resolver)
	lt.errorCatcher = errorsResolver.GetErrorCatcher(resolver)
	eventManager := resolver.Type(new(*eventsmanager.EventManager), nil).(*eventsmanager.EventManager)

	// Create subscriptions and publisher for file change events

	filechangeNotifier, err := disk.NewFileChangedNotifier(lt.configFilePath, eventManager, lt.logger)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	configHandler, err := fileconfig.NewFileConfigHandler(
		lt.configFilePath,
		&testConfig{Address: firstAddress, Address2: secondAddress},
		lt.errorCatcher,
		eventManager,
		filechangeNotifier,
		lt.logger,
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	lt.configHandler = configHandler
}

func (lt *ListenerTests) teardown() {
	_ = disk.DeleteFile(lt.configFilePath)
}

func (lt *ListenerTests) createListener(name string, addressSelector func(*testConfig) string) (server.Listener, error) {
	builder := server.NewListenerBuilder(lt.configHandler, lt.logger, lt.errorCatcher)
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
	time.Sleep(100 * time.Millisecond)

	// Act
	lt.updateConfigToDuplicatePort(t)
	time.Sleep(200 * time.Millisecond)

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
