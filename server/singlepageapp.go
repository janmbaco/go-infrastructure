package server

import (
	"net/http"
	"os"
	"path/filepath"
)

type singlePageApp struct {
	staticPath string
	indexPath  string
}

// NewSinglePageApp return the handler for a Single Page App
func NewSinglePageApp(staticPath string, indexPath string) *singlePageApp {
	return &singlePageApp{staticPath: staticPath, indexPath: indexPath}
}

func (singlePageApp *singlePageApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(singlePageApp.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(singlePageApp.staticPath, singlePageApp.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(singlePageApp.staticPath)).ServeHTTP(w, r)
}
