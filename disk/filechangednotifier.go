package disk

import (
	"github.com/fsnotify/fsnotify"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

type (
	// FileChangedNotifier defines and object that observe changes of a file
	FileChangedNotifier interface {
		Subscribe(subscribeFunc func())
	}

	fileChangedNotifier struct {
		file           string
		subscriptions  eventsmanager.Subscriptions
		eventPublisher eventsmanager.Publisher
		isWatchingFile bool
		watcher        *fsnotify.Watcher
	}
)

// NewFileChangedNotifier returns a FileChangedNotifier
func NewFileChangedNotifier(filePath string, subscriptions eventsmanager.Subscriptions, publisher eventsmanager.Publisher) FileChangedNotifier {
	errorschecker.CheckNilParameter(map[string]interface{}{ "subscriptions": subscriptions, "publisher": publisher})
	watcher, err := fsnotify.NewWatcher()
	errorschecker.TryPanic(err)
	return &fileChangedNotifier{file: filePath, watcher: watcher, subscriptions: subscriptions, eventPublisher: publisher}
}

// Subscribe subscribes a functio to observe changes of a file
func (f *fileChangedNotifier) Subscribe(subscribeFunc func()) {
	errorschecker.CheckNilParameter(map[string]interface{}{"subscribeFunc": subscribeFunc})
	f.subscriptions.Add(&fileChangedEvent{}, subscribeFunc)
	if !f.isWatchingFile {
		errorschecker.TryPanic(f.watcher.Add(f.file))
		go f.watchFile()
		f.isWatchingFile = true
	}
}

func (f *fileChangedNotifier) watchFile() {
	for evt := range f.watcher.Events {
		if evt.Op == fsnotify.Write {
			f.eventPublisher.Publish(&fileChangedEvent{})
		}
	}
}
