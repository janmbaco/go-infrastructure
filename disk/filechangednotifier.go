package disk

import (
	"github.com/fsnotify/fsnotify"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
	"github.com/janmbaco/go-infrastructure/logs"
)

type (
	// FileChangedNotifier defines and object that observe changes of a file
	FileChangedNotifier interface {
		Subscribe(subscribeFunc func()) error
	}

	fileChangedNotifier struct {
		subscriptions  eventsmanager.Subscriptions[FileChangedEvent]
		eventPublisher eventsmanager.Publisher[FileChangedEvent]
		watcher        *fsnotify.Watcher
		file           string
		isWatchingFile bool
	}
)

// NewFileChangedNotifier returns a FileChangedNotifier
func NewFileChangedNotifier(filePath string, eventManager *eventsmanager.EventManager, logger logs.Logger) (FileChangedNotifier, error) {
	subscriptions := eventsmanager.NewSubscriptions[FileChangedEvent]()
	publisher := eventsmanager.NewPublisher(subscriptions, logger)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &fileChangedNotifier{file: filePath, watcher: watcher, subscriptions: subscriptions, eventPublisher: publisher}, nil
}

// Subscribe subscribes a functio to observe changes of a file
func (f *fileChangedNotifier) Subscribe(subscribeFunc func()) error {
	fn := func(FileChangedEvent) { subscribeFunc() }
	if err := f.subscriptions.Add(fn); err != nil {
		return err
	}
	if !f.isWatchingFile {
		if err := f.watcher.Add(f.file); err != nil {
			return err
		}
		go f.watchFile()
		f.isWatchingFile = true
	}
	return nil
}

func (f *fileChangedNotifier) watchFile() {
	for evt := range f.watcher.Events {
		if evt.Op == fsnotify.Write {
			f.eventPublisher.Publish(FileChangedEvent{})
		}
	}
}
