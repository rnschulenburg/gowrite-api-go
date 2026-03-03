package ExportDocumentController

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rnschulenburg/gowrite-api-go/App/Services/ExportDocumentService"
	"io"
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	projectName := mux.Vars(r)["projectName"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	optionsHeader := r.Header.Get("x-options")
	var options map[string]interface{}
	err = json.Unmarshal([]byte(optionsHeader), &options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filePath, err := ExportDocumentService.Export(
		"/tmp",
		projectName,
		string(body),
		options,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileType := options["fileType"].(string)

	w.Header().Set("Content-Type", getMimeByFileType(fileType))
	w.Header().Set("Content-Disposition",
		"attachment; filename="+projectName+"."+strings.ReplaceAll(fileType, "-", "."))

	http.ServeFile(w, r, filePath)
}

func getMimeByFileType(fileType string) string {
	switch fileType {
	case "dom-pdf", "doc-pdf":
		return "application/pdf"
	case "word":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}
