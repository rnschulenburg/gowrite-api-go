package ProjectController

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Get
// return empty
func Get(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	projectName := strings.TrimSpace(vars["projectName"])

	basePath := os.Getenv("GeneratedHtmlPath")
	file := filepath.Join(basePath, "books", projectName+".html")

	content := "<h1>Neuer Roman</h1>"

	if data, err := os.ReadFile(file); err == nil {
		content = string(data)
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(content))
}
