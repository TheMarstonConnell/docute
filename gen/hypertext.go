package gen

import (
	"bytes"
	"regexp"
	"strings"
)

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

func processHintTags(content string) string {
	re := regexp.MustCompile(`{% hint style="(\w+)" %}(.*?){% endhint %}`)
	return re.ReplaceAllStringFunc(content, func(s string) string {
		matches := re.FindStringSubmatch(s)
		if len(matches) != 3 {
			return s
		}
		style, hintContent := matches[1], matches[2]
		hintContent = strings.TrimSpace(hintContent)
		return `<div class="hint ` + style + `">` + hintContent + `</div>`
	})
}

// renderHTML renders an HTML node to a string
func renderHTML(node *html.Node) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("<!DOCTYPE html>\n") // Prepend the DOCTYPE declaration

	err := html.Render(&buf, node)
	if err != nil {
		return nil, err
	}
	return []byte(processHintTags(buf.String())), nil
}
