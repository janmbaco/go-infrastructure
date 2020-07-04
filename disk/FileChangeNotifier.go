package disk

import (
	"github.com/fsnotify/fsnotify"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/event"
)

const onFileChangedEvent = "onFileChangedEvent"

type FileChangedNotifier struct {
	file           string
	eventPublisher *event.EventPublisher
	isWatchingFile bool
	watcher        *fsnotify.Watcher
	isSubscribing  chan bool
}

func NewFileChangedNotifier(file string) *FileChangedNotifier {
	watcher, err := fsnotify.NewWatcher()
	errorhandler.TryPanic(err)
	errorhandler.TryPanic(watcher.Add(file))
	return &FileChangedNotifier{file: file, watcher: watcher, eventPublisher: event.NewEventPublisher(), isSubscribing: make(chan bool, 1)}
}

func (this *FileChangedNotifier) Subscribe(subscribeFunc func()) {
	this.isSubscribing <- true
	this.eventPublisher.Subscribe(onFileChangedEvent, subscribeFunc)
	if !this.isWatchingFile {
		go this.watchFile()
		this.isWatchingFile = true
	}
	<-this.isSubscribing
}

func (this *FileChangedNotifier) publish() {
	this.eventPublisher.Publish(onFileChangedEvent)
}

func (this *FileChangedNotifier) watchFile() {
	for range this.watcher.Events {
		this.publish()
	}
}
