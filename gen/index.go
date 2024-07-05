package gen

import (
	"bytes"
	"embed"
	"fmt"
	"path"

	"golang.org/x/net/html"
)
import _ "embed"

////go:embed docs.js
//var docScript string

//go:embed normalize.css
var normalize string

//go:embed styles/*
var styles embed.FS

func id(i string) map[string]string {
	return map[string]string{"id": i}
}

func add(atrs map[string]string, field, val string) map[string]string {
	atrs[field] = val
	return atrs
}

// Function to check if a node has a matching href attribute
func hasMatchingHref(node *html.Node, href string) bool {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			v := attr.Val
			fmt.Printf("%s = %s ? \n", href, v)
			if attr.Key == "href" && v == href {
				return true
			}
		}
	}
	return false
}

func visit(node *html.Node, marker string, class string) {
	if node == nil {
		return
	}

	if hasMatchingHref(node, marker) {
		addClassToNode(node, class)
	}

	// Recursively visit child nodes
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		visit(c, marker, class)
	}
}

// Function to add a class to an HTML node
func addClassToNode(node *html.Node, classToAdd string) {
	for i, attr := range node.Attr {
		if attr.Key == "class" {
			node.Attr[i].Val = attr.Val + " " + classToAdd
			return
		}
	}
	node.Attr = append(node.Attr, html.Attribute{Key: "class", Val: classToAdd})
}

func CreateHead(base string) *html.Node {
	head := createHTMLElement("head", nil)

	highlightcss := createHTMLElement("link", map[string]string{"rel": "stylesheet", "href": "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/default.min.css"})
	highlightjs := createHTMLElement("script", map[string]string{"src": "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"})
	highlightjsGo := createHTMLElement("script", map[string]string{"src": "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/go.min.js"})

	b := createHTMLElement("base", map[string]string{"href": base})

	head.AppendChild(highlightcss)
	head.AppendChild(highlightjs)
	head.AppendChild(highlightjsGo)

	head.AppendChild(b)

	normalizecss := createHTMLElement("style", nil, &html.Node{
		Type: html.TextNode,
		Data: normalize,
	})
	head.AppendChild(normalizecss)

	styleEntries, err := styles.ReadDir("styles")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, entry := range styleEntries {
		fmt.Println(entry.Name())

		cssData, err := styles.ReadFile(path.Join("styles", entry.Name()))
		if err != nil {
			fmt.Println(err)
			return nil
		}

		globalscss := createHTMLElement("style", nil, &html.Node{
			Type: html.TextNode,
			Data: string(cssData),
		})
		head.AppendChild(globalscss)
	}

	return head
}

func CreateNav(summary *html.Node) *html.Node {
	nav := createHTMLElement("div", id("navbar"), summary) // creating Nav bar

	return nav
}

func CreateIndex(summaryData []byte, pageData []byte, marker string, base string) ([]byte, error) {
	body := createHTMLElement("body", nil)

	n, err := html.Parse(bytes.NewReader(summaryData))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	visit(n, marker, "active")

	bodyData, err := html.Parse(bytes.NewReader(pageData))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	nav := CreateNav(n)

	page := createHTMLElement("div", id("page"), bodyData)

	title := createHTMLElement("img", add(id("logo"), "src", "logo.png"))

	buttons := createHTMLElement("div", nil)
	headerBar := createHTMLElement("div", id("header"), title, buttons)

	main := createHTMLElement("div", id("main"), headerBar, nav, page)

	body.AppendChild(main)
	highlightInit := createHTMLElement("script", nil, &html.Node{
		Type: html.TextNode,
		Data: "hljs.highlightAll();",
	})
	body.AppendChild(highlightInit)

	root := createHTMLElement("html", nil, CreateHead(base), body)
	return renderHTML(root)
}
