package server_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/configuration"
	"github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig"
	infraerrors "github.com/janmbaco/go-infrastructure/v2/errors"
	"github.com/janmbaco/go-infrastructure/v2/eventsmanager"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-infrastructure/v2/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type revisionConfig struct {
	Revision int               `json:"revision"`
	Labels   map[string]string `json:"labels,omitempty"`
}
type capturedWatcher struct{ notify func() }

func (w *capturedWatcher) Subscribe(fn func()) error {
	w.notify = fn
	return nil
}

type fixture struct {
	config  configuration.ConfigHandler
	builder server.ListenerBuilder
	watcher *capturedWatcher
	path    string
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	logger := logs.NewLogger()
	catcher := infraerrors.NewErrorCatcher(nil)
	w := new(capturedWatcher)
	p := filepath.Join(t.TempDir(), "test.json")
	c, err := fileconfig.NewFileConfigHandler(p, &revisionConfig{Revision: 1}, catcher, eventsmanager.NewEventManager(), w, logger)
	require.NoError(t, err)
	return &fixture{config: c, builder: server.NewListenerBuilder(c, logger, catcher), watcher: w, path: p}
}

func testAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := ln.Addr().String()
	require.NoError(t, ln.Close())
	return addr
}

func await[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(8 * time.Second):
		t.Fatal("deadline waiting for listener transition")
		var v T
		return v
	}
}

func checkResponse(t *testing.T, client *http.Client, url, want string) {
	t.Helper()
	require.Eventually(t, func() bool {
		resp, err := client.Get(url)
		if err != nil {
			return false
		}
		b, err := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		return err == nil && closeErr == nil && string(b) == want
	}, 5*time.Second, 10*time.Millisecond)
}

func TestListenerReloadFromItsOwnTLSRequest(t *testing.T) {
	f := newFixture(t)
	addr := testAddr(t)
	cert := httptest.NewTLSServer(nil)
	client := cert.Client()
	client.Timeout = 3 * time.Second
	tlsConfig := &tls.Config{Certificates: cert.TLS.Certificates, MinVersion: tls.VersionTLS12}
	cert.Close()
	returned := make(chan struct{})
	release := make(chan struct{})
	var unblock sync.Once
	defer unblock.Do(func() { close(release) })
	listener, err := f.builder.SetBootstrapper(func(value any, s *server.ServerSetter) error {
		v := value.(*revisionConfig)
		s.Addr = addr
		s.TLSConfig = tlsConfig
		s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/change" {
				if err := f.config.SetConfig(&revisionConfig{Revision: 2}); err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				close(returned)
				<-release
			}
			_, _ = fmt.Fprint(w, v.Revision)
		})
		return nil
	}).GetListener()
	require.NoError(t, err)
	finish := listener.Start()
	defer listener.Stop()
	checkResponse(t, client, "https://"+addr, "1")
	completed := make(chan error, 1)
	go func() {
		res, err := client.Post("https://"+addr+"/change", "text/plain", nil)
		if err == nil {
			body, e := io.ReadAll(res.Body)
			closeErr := res.Body.Close()
			if closeErr != nil {
				err = closeErr
			}
			if e != nil || string(body) != "1" {
				err = fmt.Errorf("change response %q: %v", body, e)
			}
		}
		completed <- err
	}()
	await(t, returned)
	// A watcher and multiple writes while the current HTTP request is draining.
	var writes sync.WaitGroup
	errs := make(chan error, 12)
	for i := 3; i < 15; i++ {
		writes.Add(1)
		go func() {
			defer writes.Done()
			errs <- f.config.SetConfig(&revisionConfig{Revision: i})
			f.watcher.notify()
		}()
	}
	writes.Wait()
	close(errs)
	for writeErr := range errs {
		require.NoError(t, writeErr)
	}
	finalBytes, err := json.Marshal(&revisionConfig{Revision: 99})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(f.path, finalBytes, 0o600))
	f.watcher.notify()
	unblock.Do(func() { close(release) })
	require.NoError(t, await(t, completed))
	checkResponse(t, client, "https://"+addr, "99")
	listener.Stop()
	require.NoError(t, await(t, finish))
}

func TestListenerStopFromHandlerAndFreshRun(t *testing.T) {
	f := newFixture(t)
	addr := testAddr(t)
	var l server.Listener
	var err error
	l, err = f.builder.SetBootstrapper(func(_ any, s *server.ServerSetter) error {
		s.Addr = addr
		s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Stop()
			l.Stop()
			_, _ = fmt.Fprint(w, "stopped")
		})
		return nil
	}).GetListener()
	require.NoError(t, err)
	for range 3 {
		finish := l.Start()
		require.Equal(t, finish, l.Start())
		res, err := http.Get("http://" + addr)
		require.NoError(t, err)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, res.Body.Close())
		require.NoError(t, err)
		require.Equal(t, "stopped", string(body))
		require.NoError(t, await(t, finish))
		l.Stop()
	}
}

func TestListenerStartupFailureReturnsFinishAndReleasesPort(t *testing.T) {
	for _, mode := range []string{"missing-address", "panic", "invalid-tls", "occupied"} {
		t.Run(mode, func(t *testing.T) {
			f := newFixture(t)
			addr := testAddr(t)
			if mode == "occupied" {
				occupied, e := net.Listen("tcp", addr)
				require.NoError(t, e)
				t.Cleanup(func() { require.NoError(t, occupied.Close()) })
			}
			var calls atomic.Int32
			l, err := f.builder.SetBootstrapper(func(_ any, s *server.ServerSetter) error {
				s.Addr = addr
				if calls.Add(1) == 1 {
					return nil
				}
				switch mode {
				case "missing-address":
					s.Addr = ""
				case "panic":
					panic("startup test")
				case "invalid-tls":
					s.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
				}
				return nil
			}).GetListener()
			require.NoError(t, err)
			started := make(chan chan server.ListenerError, 1)
			go func() { started <- l.Start() }()
			terminal := await(t, await(t, started))
			require.Error(t, terminal)
			if mode == "missing-address" {
				require.Equal(t, server.AddressNotConfigured, terminal.GetErrorType())
			}
			l.Stop()
			l.Stop()
			if mode != "occupied" {
				bound, e := net.Listen("tcp", addr)
				require.NoError(t, e)
				require.NoError(t, bound.Close())
			}
		})
	}
}

func TestListenerGRPCGenerationServesAndStops(t *testing.T) {
	f := newFixture(t)
	addr := testAddr(t)
	l, err := f.builder.SetBootstrapper(func(_ any, s *server.ServerSetter) error {
		s.Addr = addr
		s.ServerType = server.GRpcSever
		return nil
	}).SetGrpcDefinitions(func(s *grpc.Server) { healthpb.RegisterHealthServer(s, health.NewServer()) }).GetListener()
	require.NoError(t, err)
	finish := l.Start()
	defer l.Stop()
	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer func() { require.NoError(t, cc.Close()) }()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	resp, err := healthpb.NewHealthClient(cc).Check(ctx, &healthpb.HealthCheckRequest{})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
	l.Stop()
	require.NoError(t, await(t, finish))
}

func TestListenerValidateAndApplyWithoutRebinding(t *testing.T) {
	f := newFixture(t)
	addr := testAddr(t)
	validated := make(chan int, 4)
	applied := make(chan int, 4)
	var value atomic.Int32
	l, err := f.builder.SetBootstrapper(func(v any, s *server.ServerSetter) error {
		s.Addr = addr
		original := v.(*revisionConfig).Revision
		s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = fmt.Fprintf(w, "%d/%d", original, value.Load()) })
		return nil
	}).SetConfigValidatorFunc(func(v any) (bool, error) { r := v.(*revisionConfig).Revision; validated <- r; return r != 2, nil }).SetConfigApplicatorFunc(func(v any, a *server.ConfigApplication) error {
		r := v.(*revisionConfig).Revision
		value.Store(int32(r))
		*a.NeedsRestart = r == 4
		applied <- r
		return nil
	}).GetListener()
	require.NoError(t, err)
	finish := l.Start()
	defer l.Stop()
	require.NoError(t, f.config.SetConfig(&revisionConfig{Revision: 2}))
	require.Equal(t, 2, await(t, validated))
	checkResponse(t, http.DefaultClient, "http://"+addr, "1/0")
	require.NoError(t, f.config.SetConfig(&revisionConfig{Revision: 3}))
	require.Equal(t, 3, await(t, validated))
	require.Equal(t, 3, await(t, applied))
	checkResponse(t, http.DefaultClient, "http://"+addr, "1/3")
	require.NoError(t, f.config.SetConfig(&revisionConfig{Revision: 4}))
	require.Equal(t, 4, await(t, validated))
	require.Equal(t, 4, await(t, applied))
	checkResponse(t, http.DefaultClient, "http://"+addr, "4/4")
	l.Stop()
	require.NoError(t, await(t, finish))
}

func TestListenerRollbackAndNewerConfiguration(t *testing.T) {
	for _, newer := range []bool{false, true} {
		t.Run(fmt.Sprintf("newer=%t", newer), func(t *testing.T) {
			f := newFixture(t)
			addr := testAddr(t)
			failing := make(chan struct{})
			release := make(chan struct{})
			var unblock sync.Once
			defer unblock.Do(func() { close(release) })
			restored := make(chan struct{}, 1)
			onRestored := func() { restored <- struct{}{} }
			f.config.RestoredSubscribe(&onRestored)
			l, err := f.builder.SetBootstrapper(func(v any, s *server.ServerSetter) error {
				r := v.(*revisionConfig).Revision
				s.Addr = addr
				s.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = fmt.Fprint(w, r) })
				if r == 2 {
					close(failing)
					<-release
					return fmt.Errorf("rejected in bootstrap")
				}
				return nil
			}).SetConfigApplicatorFunc(func(v any, a *server.ConfigApplication) error {
				// Revision 2 retires the active server before its bootstrap fails.
				// Later revisions normally hot-apply, but now need a new binding.
				*a.NeedsRestart = v.(*revisionConfig).Revision == 2
				return nil
			}).GetListener()
			require.NoError(t, err)
			finish := l.Start()
			defer l.Stop()
			require.NoError(t, f.config.SetConfig(&revisionConfig{Revision: 2}))
			await(t, failing)
			if newer {
				require.NoError(t, f.config.SetConfig(&revisionConfig{Revision: 3}))
			}
			unblock.Do(func() { close(release) })
			if newer {
				checkResponse(t, http.DefaultClient, "http://"+addr, "3")
				select {
				case <-restored:
					t.Fatal("reverted a newer generation")
				default:
				}
			} else {
				await(t, restored)
				require.Equal(t, 1, f.config.GetConfig().(*revisionConfig).Revision)
				require.NoError(t, f.config.SetConfig(&revisionConfig{Revision: 4}))
				checkResponse(t, http.DefaultClient, "http://"+addr, "4")
			}
			l.Stop()
			require.NoError(t, await(t, finish))
		})
	}
}

func TestConfigurationSnapshotPreservesNestedTypesAndRevision(t *testing.T) {
	f := newFixture(t)
	v := f.config.(interface {
		ConfigSnapshot() (any, uint64, error)
		RestoreIfVersion(uint64) (bool, error)
	})
	require.NoError(t, f.config.SetConfig(&revisionConfig{Revision: 2, Labels: map[string]string{"keep": "yes"}}))
	snapshot, revision, err := v.ConfigSnapshot()
	require.NoError(t, err)
	cfg := snapshot.(*revisionConfig)
	cfg.Labels["keep"] = "mine"
	require.Equal(t, "yes", f.config.GetConfig().(*revisionConfig).Labels["keep"])
	require.NoError(t, f.config.SetConfig(&revisionConfig{Revision: 3, Labels: map[string]string{"extra": "no"}}))
	didRestore, err := v.RestoreIfVersion(revision)
	require.NoError(t, err)
	require.False(t, didRestore)
	_, current, err := v.ConfigSnapshot()
	require.NoError(t, err)
	didRestore, err = v.RestoreIfVersion(current)
	require.NoError(t, err)
	require.True(t, didRestore)
	restored, version, err := v.ConfigSnapshot()
	require.NoError(t, err)
	require.Greater(t, version, current)
	require.Equal(t, &revisionConfig{Revision: 2, Labels: map[string]string{"keep": "yes"}}, restored)
	require.Equal(t, "mine", cfg.Labels["keep"])
}

func TestListenerStopDuringInitialBootstrap(t *testing.T) {
	f := newFixture(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	var calls atomic.Int32
	l, err := f.builder.SetBootstrapper(func(_ any, s *server.ServerSetter) error {
		if calls.Add(1) == 1 {
			return nil
		}
		close(entered)
		<-release
		return fmt.Errorf("bootstrap superseded by Stop")
	}).GetListener()
	require.NoError(t, err)
	started := make(chan chan server.ListenerError, 1)
	go func() { started <- l.Start() }()
	await(t, entered)
	l.Stop()
	l.Stop()
	close(release)
	require.NoError(t, await(t, await(t, started)))
}
