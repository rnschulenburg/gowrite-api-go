package ConverterService

import (
	"os"
	"os/exec"
)

type DomToPdf struct{}

func NewDomToPdf() *DomToPdf {
	return &DomToPdf{}
}

func (d *DomToPdf) CreatePdfDocument(
	outputPath string,
	root interface{},
	options map[string]interface{},
) error {

	html := "<html><body>Generated</body></html>"

	tmp := outputPath + ".html"
	err := os.WriteFile(tmp, []byte(html), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("wkhtmltopdf", tmp, outputPath)
	return cmd.Run()
}
