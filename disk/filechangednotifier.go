package disk

import (
	"github.com/fsnotify/fsnotify"
	"github.com/janmbaco/go-infrastructure/errors"
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
func NewFileChangedNotifier(file string, errorCatcher errors.ErrorCatcher, errorThrower errors.ErrorThrower) FileChangedNotifier {
	errors.CheckNilParameter(map[string]interface{}{"errorThrower": errorThrower})
	watcher, err := fsnotify.NewWatcher()
	errors.TryPanic(err)
	subscriptions := eventsmanager.NewSubscriptions(errorThrower)
	return &fileChangedNotifier{file: file, watcher: watcher, subscriptions: subscriptions, eventPublisher: eventsmanager.NewPublisher(subscriptions, errorCatcher)}
}

// Subscribe subscribes a functio to observe changes of a file
func (f *fileChangedNotifier) Subscribe(subscribeFunc func()) {
	errors.CheckNilParameter(map[string]interface{}{"subscribeFunc": subscribeFunc})
	f.subscriptions.Add(&fileChangedEvent{}, subscribeFunc)
	if !f.isWatchingFile {
		errors.TryPanic(f.watcher.Add(f.file))
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
