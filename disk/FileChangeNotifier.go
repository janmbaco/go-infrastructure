package disk

import (
	"github.com/fsnotify/fsnotify"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
)

const onFileChangedEvent = "onFileChangedEvent"

// FileChangedNotifier defines and object that observe changes of a file
type FileChangedNotifier interface {
	Subscribe(subscribeFunc *func())
}

type fileChangedNotifier struct {
	file           string
	eventPublisher events.Publisher
	isWatchingFile bool
	watcher        *fsnotify.Watcher
	isBusy         chan bool
	isPublishing   bool
}

// NewFileChangedNotifier returns a FileChangedNotifier
func NewFileChangedNotifier(file string) FileChangedNotifier {
	watcher, err := fsnotify.NewWatcher()
	errorhandler.TryPanic(err)
	return &fileChangedNotifier{file: file, watcher: watcher, eventPublisher: events.NewPublisher(), isPublishing: false, isBusy: make(chan bool, 1)}
}

// Subscribe subscribes a functio to observe changes of a file
func (this *fileChangedNotifier) Subscribe(subscribeFunc *func()) {
	this.eventPublisher.Subscribe(onFileChangedEvent, subscribeFunc)
	if !this.isWatchingFile {
		errorhandler.TryPanic(this.watcher.Add(this.file))
		go this.watchFile()
		this.isWatchingFile = true
	}
}

func (this *fileChangedNotifier) publish(isPublising bool) {
	this.isBusy <- true
	if !isPublising {
		this.isPublishing = true
		this.eventPublisher.Publish(onFileChangedEvent)
		this.isPublishing = false
	}
	<-this.isBusy
}

func (this *fileChangedNotifier) watchFile() {
	for evt := range this.watcher.Events {
		if evt.Op == fsnotify.Write {
			go this.publish(this.isPublishing)
		}
	}
}
