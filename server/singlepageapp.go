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
func NewSinglePageApp(staticPath string, indexPath string) http.Handler {
	return &singlePageApp{staticPath: staticPath, indexPath: indexPath}
}

func (sap *singlePageApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(filepath.Join(sap.staticPath, r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(sap.staticPath, sap.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(sap.staticPath)).ServeHTTP(w, r)
}
