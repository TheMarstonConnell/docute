package gen

import (
	"bytes"
	"golang.org/x/net/html"
)
import _ "embed"

//go:embed docs.js
var docScript string

//go:embed normalize.css
var normalize string

//go:embed global.css
var globals string

func id(i string) map[string]string {
	return map[string]string{"id": i}
}

func CreateIndex(summaryData []byte) ([]byte, error) {
	head := createHTMLElement("head", nil)
	body := createHTMLElement("body", nil)

	highlightcss := createHTMLElement("link", map[string]string{"rel": "stylesheet", "href": "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/default.min.css"})

	highlightjs := createHTMLElement("script", map[string]string{"src": "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"})
	highlightjsGo := createHTMLElement("script", map[string]string{"src": "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/go.min.js"})

	highlightInit := createHTMLElement("script", nil, &html.Node{
		Type: html.TextNode,
		Data: "hljs.highlightAll();",
	})

	head.AppendChild(highlightcss)
	head.AppendChild(highlightjs)
	head.AppendChild(highlightjsGo)

	docsjs := createHTMLElement("script", nil, &html.Node{
		Type: html.TextNode,
		Data: docScript,
	})

	normalizecss := createHTMLElement("style", nil, &html.Node{
		Type: html.TextNode,
		Data: normalize,
	})
	head.AppendChild(normalizecss)

	globalscss := createHTMLElement("style", nil, &html.Node{
		Type: html.TextNode,
		Data: globals,
	})
	head.AppendChild(globalscss)

	n, err := html.Parse(bytes.NewReader(summaryData))
	if err != nil {
		return nil, err
	}

	nav := createHTMLElement("div", id("navbar"), n)

	page := createHTMLElement("div", id("page"))

	main := createHTMLElement("div", id("main"), nav, page)

	body.AppendChild(main)
	body.AppendChild(docsjs)
	body.AppendChild(highlightInit)

	root := createHTMLElement("html", nil, head, body)
	return renderHTML(root)

}
