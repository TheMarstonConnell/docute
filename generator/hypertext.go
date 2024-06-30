package generator

import "bytes"
import (
	"golang.org/x/net/html"
)

// createHTMLElement creates an HTML element with given tag, attributes, and children
func createHTMLElement(tag string, attrs map[string]string, children ...*html.Node) *html.Node {
	// Create the element node
	elem := &html.Node{
		Type: html.ElementNode,
		Data: tag,
	}

	// Add attributes
	for key, value := range attrs {
		elem.Attr = append(elem.Attr, html.Attribute{Key: key, Val: value})
	}

	// Append children
	for _, child := range children {
		elem.AppendChild(child)
	}

	return elem
}

// renderHTML renders an HTML node to a string
func renderHTML(node *html.Node) (string, error) {
	var buf bytes.Buffer
	err := html.Render(&buf, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
