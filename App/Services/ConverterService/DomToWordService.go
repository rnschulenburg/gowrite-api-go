package ConverterService

import (
	"archive/zip"
	"golang.org/x/net/html"
	"os"
	"strconv"
	"strings"
)

type DomToWord struct{}

func NewDomToWord() *DomToWord {
	return &DomToWord{}
}

func CreateWordDocument(
	path string,
	root string,
	options map[string]interface{},
) error {
	return createDocx(path, root)
}

func createDocx(filename string, htmlInput string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// HTML → DOM
	root, err := parseHTML(htmlInput)
	if err != nil {
		return err
	}

	body := replaceHtmlToWord(root)

	// =========================
	// 1️⃣ document.xml
	// =========================
	documentXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    ` + body + `
    <w:sectPr>
      <w:pgSz w:w="11906" w:h="16838"/>
      <w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440"/>
    </w:sectPr>
  </w:body>
</w:document>`

	addFile(zipWriter, "word/document.xml", documentXML)

	// =========================
	// 2️⃣ styles.xml  (WICHTIG!)
	// =========================
	stylesXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">

  <!-- ================= GLOBAL DEFAULTS ================= -->
  <w:docDefaults>

    <!-- Default Run Properties (Schrift global) -->
    <w:rPrDefault>
      <w:rPr>
        <w:rFonts w:ascii="Courier New" w:hAnsi="Courier New"/>
        <w:sz w:val="24"/>              <!-- 12pt (24 half-points) -->
        <w:spacing w:val="3"/>          <!-- letter-spacing -->
      </w:rPr>
    </w:rPrDefault>

    <!-- Default Paragraph Properties -->
    <w:pPrDefault>
      <w:pPr>
        <w:spacing w:line="360" w:lineRule="auto"/> <!-- 1.6 line-height -->
      </w:pPr>
    </w:pPrDefault>

  </w:docDefaults>

  <!-- ================= NORMAL ================= -->
  <w:style w:type="paragraph" w:default="1" w:styleId="Normal">
    <w:name w:val="Normal"/>
    <w:qFormat/>
  </w:style>

  <!-- ================= H1 ================= -->
  <w:style w:type="paragraph" w:styleId="Heading1">
    <w:name w:val="heading 1"/>
    <w:basedOn w:val="Normal"/>
    <w:rPr>
      <w:rFonts w:ascii="Segoe UI" w:hAnsi="Segoe UI"/>
      <w:sz w:val="28"/>  <!-- 14pt -->
    </w:rPr>
  </w:style>

  <!-- ================= H2 ================= -->
  <w:style w:type="paragraph" w:styleId="Heading2">
    <w:name w:val="heading 2"/>
    <w:basedOn w:val="Normal"/>
    <w:rPr>
      <w:rFonts w:ascii="Segoe UI" w:hAnsi="Segoe UI"/>
      <w:sz w:val="26"/>
      <w:color w:val="0000FF"/>
    </w:rPr>
  </w:style>

  <!-- ================= H3 ================= -->
  <w:style w:type="paragraph" w:styleId="Heading3">
    <w:name w:val="heading 3"/>
    <w:basedOn w:val="Normal"/>
    <w:rPr>
      <w:rFonts w:ascii="Segoe UI" w:hAnsi="Segoe UI"/>
      <w:sz w:val="24"/>
      <w:i/>
      <w:color w:val="F8339E"/>
    </w:rPr>
  </w:style>

  <!-- ================= H4 ================= -->
  <w:style w:type="paragraph" w:styleId="Heading4">
    <w:name w:val="heading 4"/>
    <w:basedOn w:val="Normal"/>
    <w:rPr>
      <w:rFonts w:ascii="Segoe UI" w:hAnsi="Segoe UI"/>
      <w:sz w:val="22"/>
      <w:i/>
      <w:color w:val="909096"/>
    </w:rPr>
  </w:style>

</w:styles>`

	addFile(zipWriter, "word/styles.xml", stylesXML)

	// =========================
	// 3️⃣ document relationships
	// =========================
	documentRels := `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1"
    Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"
    Target="styles.xml"/>
</Relationships>`

	addFile(zipWriter, "word/_rels/document.xml.rels", documentRels)

	// =========================
	// 4️⃣ Content Types
	// =========================
	contentTypes := `<?xml version="1.0" encoding="UTF-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
</Types>`

	addFile(zipWriter, "[Content_Types].xml", contentTypes)

	// =========================
	// 5️⃣ Root rels
	// =========================
	rootRels := `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1"
    Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
    Target="word/document.xml"/>
</Relationships>`

	addFile(zipWriter, "_rels/.rels", rootRels)

	return nil
}

func addFile(zipWriter *zip.Writer, name, content string) {
	f, _ := zipWriter.Create(name)
	f.Write([]byte(content))
}

func heading(n *html.Node, level int) string {

	text := extractText(n)

	return `
<w:p>
  <w:pPr>
    <w:pStyle w:val="Heading` + strconv.Itoa(level) + `"/>
  </w:pPr>
  <w:r>
    <w:t>` + escapeXML(text) + `</w:t>
  </w:r>
</w:p>`
}
func extractText(n *html.Node) string {

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

func escapeXML(s string) string {
	return html.EscapeString(s)
}

func parseHTML(htmlStr string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func buildRuns(n *html.Node) string {

	var sb strings.Builder

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		// span separat behandeln
		if c.Type == html.ElementNode && c.Data == "span" {
			sb.WriteString(replaceHtmlToWord(c))
			continue
		}

		// normale Verarbeitung
		sb.WriteString(replaceHtmlToWord(c))
	}

	return sb.String()
}

func replaceHtmlToWord(n *html.Node) string {

	var sb strings.Builder

	switch n.Type {

	case html.ElementNode:

		switch n.Data {

		// -------- HEADINGS --------
		case "h1", "h2", "h3", "h4":

			level := int(n.Data[1] - '0') // h1 → 1

			sb.WriteString("<w:p>")
			sb.WriteString("<w:pPr>")
			sb.WriteString(`<w:pStyle w:val="Heading` + strconv.Itoa(level) + `"/>`)
			sb.WriteString("</w:pPr>")

			sb.WriteString(buildRuns(n))

			sb.WriteString("</w:p>")
			return sb.String()

		// -------- PARAGRAPH --------
		case "p":

			sb.WriteString("<w:p>")
			sb.WriteString(buildRuns(n))
			sb.WriteString("</w:p>")
			return sb.String()

		// -------- SPAN (Highlight) --------
		case "span":

			text := extractText(n)

			sb.WriteString("<w:r>")
			sb.WriteString("<w:rPr>")
			sb.WriteString(`<w:highlight w:val="yellow"/>`)
			sb.WriteString("</w:rPr>")
			sb.WriteString("<w:t>" + html.EscapeString(text) + "</w:t>")
			sb.WriteString("</w:r>")

			return sb.String()
		}
	}

	// -------- TEXT NODE --------
	if n.Type == html.TextNode {

		text := strings.TrimSpace(n.Data)
		if text != "" {
			sb.WriteString("<w:r>")
			sb.WriteString("<w:t>" + html.EscapeString(text) + "</w:t>")
			sb.WriteString("</w:r>")
		}
		return sb.String()
	}

	// -------- RECURSION --------
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(replaceHtmlToWord(c))
	}

	return sb.String()
}
