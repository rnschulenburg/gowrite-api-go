package ProjectController

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Save(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	projectName := strings.TrimSpace(vars["projectName"])

	basePath := os.Getenv("GeneratedHtmlPath")
	file := filepath.Join(basePath, "books", projectName+".html")

	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = os.WriteFile(file, content, 0644)
	if err != nil {
		http.Error(w, "cannot write file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}
