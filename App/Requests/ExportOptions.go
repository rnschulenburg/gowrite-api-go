package Requests

type ExportOptions struct {
	FileType string `json:"fileType"`
	H1       bool   `json:"h1"`
	H2       bool   `json:"h2"`
	H3       bool   `json:"h3"`
	H4       bool   `json:"h4"`
	Span     bool   `json:"span"`
}
