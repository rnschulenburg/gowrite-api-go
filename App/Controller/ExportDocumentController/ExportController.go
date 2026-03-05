package ExportDocumentController

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rnschulenburg/gowrite-api-go/App/Requests"
	"github.com/rnschulenburg/gowrite-api-go/App/Services/ExportDocumentService"
	"golang.org/x/net/html/charset"
	"io"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	projectName := mux.Vars(r)["projectName"]

	reader, err := charset.NewReader(r.Body, r.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	options, err := getOptionsHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	w.Header().Set("Content-Type", getMimeByFileType(options.FileType))
	w.Header().Set("Content-Disposition",
		"attachment; filename="+projectName+"."+options.FileType)

	log.Println(filePath)
	log.Println(getMimeByFileType(options.FileType))

	http.ServeFile(w, r, filePath)
}

func getMimeByFileType(fileType string) string {
	switch fileType {
	case "pdf":
		return "application/pdf"
	case "word":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "epub":
		return "application/epub+zip"
	default:
		return "application/octet-stream"
	}
}

func getOptionsHeader(r *http.Request) (Requests.ExportOptions, error) {
	optionsHeader := r.Header.Get("x-options")
	var options Requests.ExportOptions

	if optionsHeader == "" {
		return options, errors.New("missing x-options header")
	}

	err := json.Unmarshal([]byte(optionsHeader), &options)
	if err != nil {
		return options, err
	}

	return options, nil
}
