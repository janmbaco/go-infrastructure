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
}

func NewFileChangedNotifier(file string) *FileChangedNotifier {
	watcher, err := fsnotify.NewWatcher()
	errorhandler.TryPanic(err)
	errorhandler.TryPanic(watcher.Add(file))
	return &FileChangedNotifier{file: file, watcher: watcher, eventPublisher: event.NewEventPublisher()}
}

func (this *FileChangedNotifier) Subscribe(subscribeFunc func()) {
	this.eventPublisher.Subscribe(onFileChangedEvent, subscribeFunc)
	if !this.isWatchingFile {
		go this.watchFile()
	}
}

func (this *FileChangedNotifier) publish() {
	this.eventPublisher.Publish(onFileChangedEvent)
}

func (this *FileChangedNotifier) watchFile() {
	for {
		select {
		case event, ok := <-this.watcher.Events:
			if ok {
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Chmod == fsnotify.Chmod {
					this.publish()
				}
			}
		}
	}
}
