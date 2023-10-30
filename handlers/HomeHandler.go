package handlers

import (
	"net/http"
	"path/filepath"
)



func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    // Define the path to your index.html file
    path := filepath.Join("views", "index.html")
    http.ServeFile(w, r, path)
}