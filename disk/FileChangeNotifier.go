package disk

import (
	"github.com/fsnotify/fsnotify"

	"github.com/janmbaco/go-infrastructure/errorhandler"
)

type FileChangedNotifier struct {
	file                     string
	fileChangedSubscriptions []func()
	isWatchingFile           bool
	watcher                  *fsnotify.Watcher
}

func NewFileChangedNotifier(file string) *FileChangedNotifier {
	watcher, err := fsnotify.NewWatcher()
	errorhandler.TryPanic(err)
	errorhandler.TryPanic(watcher.Add(file))
	return &FileChangedNotifier{file: file, watcher: watcher}
}

func (this *FileChangedNotifier) Subscribe(subscribeFunc func()) {
	onFileChangedFunc := func() {
		subscribeFunc()
	}
	this.fileChangedSubscriptions = append(this.fileChangedSubscriptions, onFileChangedFunc)
	if !this.isWatchingFile {
		go this.watchFile()
	}
}

func (this *FileChangedNotifier) publish() {
	for _, f := range this.fileChangedSubscriptions {
		errorhandler.OnErrorContinue(f)
	}
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
