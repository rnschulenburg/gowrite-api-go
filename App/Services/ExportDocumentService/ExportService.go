package ExportDocumentService

import (
	"errors"
	"github.com/rnschulenburg/gowrite-api-go/App/Services/ConverterService"

	"log"
	"path/filepath"
	"strings"
)

//type ExportDocument struct {
//	DomToWord  *ConverterService.DomToWord
//	DomToPdf   *ConverterService.DomToPdf
//	WordToPdf  *ConverterService.WordToPdf
//	RuntimeDir string
//}
//
//func NewExportDocument(
//	runtimeDir string,
//	d2w *ConverterService.DomToWord,
//	d2p *ConverterService.DomToPdf,
//	w2p *ConverterService.WordToPdf,
//) *ExportDocument {
//	return &ExportDocument{
//		RuntimeDir: runtimeDir,
//		DomToWord:  d2w,
//		DomToPdf:   d2p,
//		WordToPdf:  w2p,
//	}
//}

func Export(
	runtimeDir string,
	projectName string,
	htmlContent string,
	options map[string]interface{},
) (string, error) {

	log.Println("Toll mal scheuen")
	projectName = strings.TrimSpace(projectName)
	//root, err := parseHTML(htmlContent)
	//if err != nil {
	//	return "", err
	//}
	//log.Println("root", root)
	wordPath := filepath.Join(runtimeDir, projectName+".docx")
	err := ConverterService.CreateWordDocument(wordPath, htmlContent, options)
	if err != nil {
		return "", err
	}

	docPdfPath := filepath.Join(runtimeDir, projectName+".doc.pdf")
	//err = e.WordToPdf.CreatePdfDocument(wordPath, docPdfPath)
	//if err != nil {
	//	return "", err
	//}
	//
	domPdfPath := filepath.Join(runtimeDir, projectName+".dom.pdf")
	//err = e.DomToPdf.CreatePdfDocument(domPdfPath, root, options)
	//if err != nil {
	//	return "", err
	//}

	switch options["fileType"] {
	case "dom-pdf":
		return domPdfPath, nil
	case "doc-pdf":
		return docPdfPath, nil
	case "word":
		return wordPath, nil
	}

	return "", errors.New("invalid fileType")
}
