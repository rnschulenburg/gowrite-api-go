package ConverterService

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"strings"
)

type DocxToDom struct{}

type Document struct {
	Body Body `xml:"body"`
}

type Body struct {
	Paragraphs []Paragraph `xml:"p"`
}

type Paragraph struct {
	Style PStyle `xml:"pPr>pStyle"`
	Runs  []Run  `xml:"r"`
}

type PStyle struct {
	Val string `xml:"val,attr"`
}

type Run struct {
	Properties *RunProperties `xml:"rPr"`
	Text       string         `xml:"t"`
}

type RunProperties struct {
	Highlight *Highlight `xml:"highlight"`
}

type Highlight struct {
	Val string `xml:"val,attr"`
}

func ConvertDocxBytesToHTML(data []byte) (string, error) {

	reader := bytes.NewReader(data)

	zipReader, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", err
	}

	var documentXML []byte

	for _, f := range zipReader.File {

		if f.Name == "word/document.xml" {

			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			documentXML, err = io.ReadAll(rc)
			if err != nil {
				return "", err
			}
		}
	}

	// hier dein bestehender Parser
	doc, err := parseDocument(documentXML)
	if err != nil {
		return "", err
	}

	return renderHTML(doc), nil
}

func readDocumentXML(path string) ([]byte, error) {

	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {

		println("ZIP FILE:", f.Name)

		if f.Name == "word/document.xml" {

			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			println("document.xml size:", len(data))

			return data, nil
		}
	}

	return nil, errors.New("document.xml not found")
}

func parseDocument(data []byte) (*Document, error) {

	var doc Document

	decoder := xml.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func renderHTML(doc *Document) string {

	var sb strings.Builder

	for _, p := range doc.Body.Paragraphs {

		tag := mapStyleToTag(p.Style.Val)

		var runBuffer strings.Builder
		renderRuns(&runBuffer, p.Runs)

		content := strings.TrimSpace(runBuffer.String())

		if content == "" {
			continue
		}

		sb.WriteString("<" + tag + ">")
		sb.WriteString(content)
		sb.WriteString("</" + tag + ">")
	}

	return sb.String()
}

func mapStyleToTag(style string) string {

	s := strings.ToLower(style)

	switch {

	case strings.Contains(s, "heading1"),
		strings.Contains(s, "heading 1"),
		s == "title":
		return "h1"

	case strings.Contains(s, "heading2"),
		strings.Contains(s, "heading 2"),
		s == "subtitle":
		return "h2"

	case strings.Contains(s, "heading3"),
		strings.Contains(s, "heading 3"):
		return "h3"

	case strings.Contains(s, "heading4"),
		strings.Contains(s, "heading 4"):
		return "h4"
	}

	return "p"
}

func renderRuns(sb *strings.Builder, runs []Run) {

	for _, r := range runs {

		text := strings.TrimSpace(r.Text)

		if text == "" {
			continue
		}

		escaped := html.EscapeString(text)

		if isHighlighted(r) {

			sb.WriteString("<span>")
			sb.WriteString(escaped)
			sb.WriteString("</span>")

		} else {

			sb.WriteString(escaped)
		}
	}
}

func isHighlighted(r Run) bool {

	if r.Properties == nil {
		return false
	}

	if r.Properties.Highlight == nil {
		return false
	}

	return true
}
