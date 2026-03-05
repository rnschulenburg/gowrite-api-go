package ExportDocumentService

import (
	"errors"
	"github.com/rnschulenburg/gowrite-api-go/App/Requests"
	"github.com/rnschulenburg/gowrite-api-go/App/Services/ConverterService"
	"path/filepath"
	"strings"
)

func Export(
	runtimeDir string,
	projectName string,
	htmlContent string,
	options Requests.ExportOptions,
) (string, error) {

	projectName = strings.TrimSpace(projectName)

	wordPath := filepath.Join(runtimeDir, projectName+".docx")
	if options.FileType == "word" {
		err := ConverterService.CreateWordDocument(wordPath, htmlContent, options)
		if err != nil {
			return "", err
		}
	}

	docPdfPath := filepath.Join(runtimeDir, projectName+".pdf")
	if options.FileType == "pdf" {
		err := ConverterService.CreatePdfDocument(docPdfPath, htmlContent, options)
		if err != nil {
			return "", err
		}
	}

	epubPath := filepath.Join(runtimeDir, projectName+".epub")
	if options.FileType == "epub" {
		err := ConverterService.CreateEpubDocument(epubPath, htmlContent, options)
		if err != nil {
			return "", err
		}
	}

	switch options.FileType {
	case "pdf":
		return docPdfPath, nil
	case "word":
		return wordPath, nil
	case "epub":
		return epubPath, nil
	}

	return "", errors.New("invalid fileType")
}
