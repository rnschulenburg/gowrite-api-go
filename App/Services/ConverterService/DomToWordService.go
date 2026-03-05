package ConverterService

import (
	"archive/zip"
	"github.com/rnschulenburg/gowrite-api-go/App/Requests"
	"golang.org/x/net/html"
	"os"
	"strconv"
	"strings"
)

type DomToWord struct{}

func CreateWordDocument(
	filename string,
	htmlInput string,
	options Requests.ExportOptions,
) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	zipWriter := zip.NewWriter(file)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {

		}
	}(zipWriter)

	root, err := parseHTML(htmlInput)
	if err != nil {
		return err
	}

	body := replaceHtmlToWord(root, options)

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

	stylesXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">

<w:docDefaults>

<w:rPrDefault>
<w:rPr>
<w:rFonts w:ascii="Courier New" w:hAnsi="Courier New"/>
<w:sz w:val="24"/>
<w:spacing w:val="3"/>
</w:rPr>
</w:rPrDefault>

<w:pPrDefault>
<w:pPr>
<w:spacing w:line="350" w:lineRule="auto"/>
</w:pPr>
</w:pPrDefault>

</w:docDefaults>

<w:style w:type="paragraph" w:default="1" w:styleId="Normal">
<w:name w:val="Normal"/>
<w:qFormat/>
</w:style>

<w:style w:type="paragraph" w:styleId="Heading1">
<w:name w:val="heading 1"/>
<w:basedOn w:val="Normal"/>

<w:pPr>
<w:spacing w:line="300" w:lineRule="auto"/>
</w:pPr>

<w:rPr>
<w:rFonts w:ascii="Segoe UI" w:hAnsi="Segoe UI"/>
<w:sz w:val="28"/>
</w:rPr>

</w:style>

<w:style w:type="paragraph" w:styleId="Heading2">
<w:name w:val="heading 2"/>
<w:basedOn w:val="Normal"/>

<w:pPr>
<w:spacing w:line="300" w:lineRule="auto"/>
</w:pPr>

<w:rPr>
<w:rFonts w:ascii="Segoe UI" w:hAnsi="Segoe UI"/>
<w:sz w:val="26"/>
<w:color w:val="0000FF"/>
</w:rPr>

</w:style>

<w:style w:type="paragraph" w:styleId="Heading3">
<w:name w:val="heading 3"/>
<w:basedOn w:val="Normal"/>

<w:pPr>
<w:spacing w:line="300" w:lineRule="auto"/>
</w:pPr>

<w:rPr>
<w:rFonts w:ascii="Segoe UI" w:hAnsi="Segoe UI"/>
<w:sz w:val="24"/>
<w:i/>
<w:color w:val="F8339E"/>
</w:rPr>

</w:style>

<w:style w:type="paragraph" w:styleId="Heading4">
<w:name w:val="heading 4"/>
<w:basedOn w:val="Normal"/>

<w:pPr>
<w:spacing w:line="300" w:lineRule="auto"/>
</w:pPr>

<w:rPr>
<w:rFonts w:ascii="Segoe UI" w:hAnsi="Segoe UI"/>
<w:sz w:val="22"/>
<w:i/>
<w:color w:val="909096"/>
</w:rPr>

</w:style>

</w:styles>`

	addFile(zipWriter, "word/styles.xml", stylesXML)

	documentRels := `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1"
    Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"
    Target="styles.xml"/>
</Relationships>`

	addFile(zipWriter, "word/_rels/document.xml.rels", documentRels)

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
	_, err := f.Write([]byte(content))
	if err != nil {
		return
	}
}

func parseHTML(htmlStr string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func replaceHtmlToWord(n *html.Node, options Requests.ExportOptions) string {

	var sb strings.Builder

	if n.Data == "h1" && !options.H1 ||
		n.Data == "h2" && !options.H2 ||
		n.Data == "h3" && !options.H3 ||
		n.Data == "h4" && !options.H4 ||
		n.Data == "span" && !options.Span {
		return ""
	}

	switch n.Type {

	case html.ElementNode:

		switch n.Data {

		case "h1", "h2", "h3", "h4":
			//if n.Data == "h1" && op

			level := int(n.Data[1] - '0')

			sb.WriteString("<w:p>")
			sb.WriteString("<w:pPr>")
			sb.WriteString(`<w:pStyle w:val="Heading` + strconv.Itoa(level) + `"/>`)
			sb.WriteString("</w:pPr>")

			sb.WriteString(buildRuns(n, options))

			sb.WriteString("</w:p>")

			return sb.String()

		case "p":

			sb.WriteString("<w:p>")
			sb.WriteString("<w:pPr>")
			sb.WriteString(`<w:spacing w:after="160"/>`)
			sb.WriteString("</w:pPr>")

			sb.WriteString(buildRuns(n, options))

			sb.WriteString("</w:p>")

			return sb.String()

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

	if n.Type == html.TextNode {

		text := strings.TrimSpace(n.Data)

		if text != "" {

			sb.WriteString("<w:r>")
			sb.WriteString("<w:t>" + html.EscapeString(text) + "</w:t>")
			sb.WriteString("</w:r>")
		}

		return sb.String()
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(replaceHtmlToWord(c, options))
	}

	return sb.String()
}

func buildRuns(n *html.Node, options Requests.ExportOptions) string {

	var sb strings.Builder

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		if c.Type == html.ElementNode && c.Data == "span" {
			sb.WriteString(replaceHtmlToWord(c, options))
			continue
		}

		sb.WriteString(replaceHtmlToWord(c, options))
	}

	return sb.String()
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
