package ConverterService

import (
	"errors"
	"github.com/go-pdf/fpdf"
	"github.com/rnschulenburg/gowrite-api-go/App/Requests"
	"golang.org/x/net/html"
	"log"
	"strings"
)

func CreatePdfDocument(
	path string,
	htmlInput string,
	options Requests.ExportOptions,
) error {

	root, err := parseHTMLPdf(htmlInput)
	if err != nil {
		return err
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetFontLocation("Resources/fonts")

	// Word uses 2.54cm margins
	pdf.SetMargins(25.4, 25.4, 25.4)

	pdf.AddPage()

	// ---------------- Fonts ----------------

	pdf.AddUTF8Font("Body", "", "LiberationMono-Regular.ttf")
	pdf.AddUTF8Font("Body", "B", "LiberationMono-Bold.ttf")
	pdf.AddUTF8Font("Body", "I", "LiberationMono-Italic.ttf")
	pdf.AddUTF8Font("Body", "BI", "LiberationMono-BoldItalic.ttf")

	pdf.AddUTF8Font("Heading", "", "Inter_18pt-Regular.ttf")
	pdf.AddUTF8Font("Heading", "B", "Inter_18pt-Bold.ttf")
	pdf.AddUTF8Font("Heading", "I", "Inter_18pt-Italic.ttf")
	pdf.AddUTF8Font("Heading", "BI", "Inter_18pt-BoldItalic.ttf")

	if pdf.Err() {
		log.Println("PDF ERROR:", pdf.Err())
		return errors.New("PDF ERROR")
	}

	pdf.SetFont("Body", "", 12)

	pdf.SetAutoPageBreak(true, 25.4)

	walkPDF(pdf, root, options)

	return pdf.OutputFileAndClose(path)
}

func walkPDF(pdf *fpdf.Fpdf, n *html.Node, options Requests.ExportOptions) {

	if n.Type == html.ElementNode {

		switch n.Data {

		// ---------------- H1 ----------------

		case "h1":
			if !options.H1 {
				return
			}

			pdf.Ln(5)

			pdf.SetFont("Heading", "B", 16)

			pdf.MultiCell(
				0,
				5.7,
				extractTextPdf(n, pdf),
				"",
				"L",
				false,
			)

			pdf.SetFont("Body", "", 12)

			pdf.Ln(2)

			return

		// ---------------- H2 ----------------

		case "h2":
			if !options.H2 {
				return
			}

			pdf.Ln(4)

			pdf.SetFont("Heading", "B", 14)

			pdf.SetTextColor(0, 0, 255)

			pdf.MultiCell(
				0,
				5.0,
				extractTextPdf(n, pdf),
				"",
				"L",
				false,
			)

			pdf.SetTextColor(0, 0, 0)

			pdf.SetFont("Body", "", 12)

			pdf.Ln(2)

			return

		// ---------------- H3 ----------------

		case "h3":
			if !options.H3 {
				return
			}

			pdf.Ln(3)

			pdf.SetFont("Heading", "I", 12)

			pdf.SetTextColor(248, 51, 158)

			pdf.MultiCell(
				0,
				4.6,
				extractTextPdf(n, pdf),
				"",
				"L",
				false,
			)

			pdf.SetTextColor(0, 0, 0)

			pdf.SetFont("Body", "", 12)

			pdf.Ln(1.5)

			return

		// ---------------- H4 ----------------

		case "h4":
			if !options.H4 {
				return
			}

			pdf.Ln(2)

			pdf.SetFont("Heading", "I", 11)

			pdf.SetTextColor(144, 144, 150)

			pdf.MultiCell(
				0,
				4.0,
				extractTextPdf(n, pdf),
				"",
				"L",
				false,
			)

			pdf.SetTextColor(0, 0, 0)

			pdf.SetFont("Body", "", 12)

			pdf.Ln(1)

			return

		// ---------------- Paragraph ----------------

		case "p":

			renderParagraph(pdf, n, options)

			// Word paragraph spacing ≈ 2.8mm
			pdf.Ln(2.9)

			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkPDF(pdf, c, options)
	}
}

func renderParagraph(pdf *fpdf.Fpdf, n *html.Node, options Requests.ExportOptions) {

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		// ---------- TEXT ----------

		if c.Type == html.TextNode {

			text := extractTextPdf(c, pdf)

			if text != "" {

				// 1.5 line height for 12pt text
				pdf.MultiCell(
					0,
					6.75,
					text,
					"",
					"L",
					false,
				)
			}
		}

		// ---------- SPAN (Highlight) ----------

		if c.Type == html.ElementNode && c.Data == "span" {
			if !options.Span {
				return
			}

			pdf.SetTextColor(255, 128, 0)

			pdf.SetFont("Body", "I", 12)

			pdf.MultiCell(
				0,
				6.75,
				extractTextPdf(c, pdf),
				"",
				"L",
				false,
			)

			pdf.SetTextColor(0, 0, 0)

			pdf.SetFont("Body", "", 12)
		}
	}
}

func extractTextPdf(n *html.Node, pdf *fpdf.Fpdf) string {

	var sb strings.Builder

	var walker func(*html.Node)

	walker = func(node *html.Node) {

		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}

	walker(n)

	return strings.TrimSpace(sb.String())
}

func parseHTMLPdf(input string) (*html.Node, error) {
	return html.Parse(strings.NewReader(input))
}
