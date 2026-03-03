package ConverterService

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type WordToPdf struct{}

func NewWordToPdf() *WordToPdf {
	return &WordToPdf{}
}

func (w *WordToPdf) CreatePdfDocument(docxPath, pdfPath string) error {

	outDir := filepath.Dir(pdfPath)

	cmd := exec.Command(
		"libreoffice",
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outDir,
		docxPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("libreoffice error: %s", string(output))
	}

	generated := filepath.Join(
		outDir,
		filepath.Base(docxPath[:len(docxPath)-5])+".pdf",
	)

	return exec.Command("mv", generated, pdfPath).Run()
}
