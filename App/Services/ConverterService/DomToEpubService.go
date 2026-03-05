package ConverterService

import (
	"archive/zip"
	"github.com/rnschulenburg/gowrite-api-go/App/Requests"
	"os"
	"strings"
)

func CreateEpubDocument(
	path string,
	htmlInput string,
	options Requests.ExportOptions,
) error {

	file, err := os.Create(path)
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

	// ---------------- mimetype (MUSS uncompressed sein) ----------------

	mimeHeader := &zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store,
	}

	mimeWriter, err := zipWriter.CreateHeader(mimeHeader)
	if err != nil {
		return err
	}

	_, err2 := mimeWriter.Write([]byte("application/epub+zip"))
	if err2 != nil {
		return err2
	}

	// ---------------- META-INF/container.xml ----------------

	containerXML := `<?xml version="1.0"?>
<container version="1.0"
xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
<rootfiles>
<rootfile
full-path="OEBPS/content.opf"
media-type="application/oebps-package+xml"/>
</rootfiles>
</container>`

	addEpubFile(zipWriter, "META-INF/container.xml", containerXML)

	// ---------------- CSS ----------------

	styleCSS := `
body{
font-family: serif;
line-height:1.5;
margin:0;
padding:0;
}

p{
margin:0 0 1em 0;
}

h1{
font-size:1.6em;
}

h2{
font-size:1.4em;
color:blue;
}

h3{
font-size:1.2em;
font-style:italic;
color:#f8339e;
}

h4{
font-size:1.1em;
font-style:italic;
color:#909096;
}

span{
background:yellow;
color:#ff8000;
font-style:italic;
}
`
	if !options.H1 {
		styleCSS = styleCSS + `
h1{
display:none;
}
`
	}
	if !options.H2 {
		styleCSS = styleCSS + `
h2{
display:none;
}
`
	}
	if !options.H3 {
		styleCSS = styleCSS + `
h3{
display:none;
}
`
	}
	if !options.H4 {
		styleCSS = styleCSS + `
h4{
display:none;
}
`
	}
	if !options.Span {
		styleCSS = styleCSS + `
span{
display:none;
}
`
	}

	addEpubFile(zipWriter, "OEBPS/style.css", styleCSS)

	// ---------------- XHTML Content ----------------

	contentHTML := `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<title>Book</title>
<link rel="stylesheet" type="text/css" href="style.css"/>
</head>
<body>
` + sanitizeHTML(htmlInput) + `
</body>
</html>`

	addEpubFile(zipWriter, "OEBPS/chapter1.xhtml", contentHTML)

	// ---------------- content.opf ----------------

	contentOPF := `<?xml version="1.0" encoding="UTF-8"?>
<package version="3.0"
xmlns="http://www.idpf.org/2007/opf"
unique-identifier="bookid">

<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
<dc:identifier id="bookid">gowrite-book</dc:identifier>
<dc:title>Exported Book</dc:title>
<dc:language>de</dc:language>
</metadata>

<manifest>
<item id="chapter1"
href="chapter1.xhtml"
media-type="application/xhtml+xml"/>

<item id="css"
href="style.css"
media-type="text/css"/>
</manifest>

<spine>
<itemref idref="chapter1"/>
</spine>

</package>`

	addEpubFile(zipWriter, "OEBPS/content.opf", contentOPF)

	return nil
}

func addEpubFile(zipWriter *zip.Writer, name string, content string) {

	f, _ := zipWriter.Create(name)

	_, err := f.Write([]byte(content))
	if err != nil {
		return
	}
}

func sanitizeHTML(input string) string {

	// EPUB braucht sauberes XHTML
	html := strings.ReplaceAll(input, "<br>", "<br/>")

	return html
}
