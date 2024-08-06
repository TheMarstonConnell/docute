package gen

import (
	"bytes"
	"embed"
	"fmt"
	"path"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)
import _ "embed"

//go:embed script.js
var docScript string

//go:embed devsocket.js
var devSocket string

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
			if attr.Key == "href" && v == href {
				return true
			}
		}
	}
	return false
}

func visit(node *html.Node, marker string, class string) string {
	if node == nil {
		return ""
	}

	if hasMatchingHref(node, marker) {
		addClassToNode(node, class)
		d := node.FirstChild.Data
		return d
	}

	// Recursively visit child nodes
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		s := visit(c, marker, class)
		if len(s) > 1 {
			return s
		}
	}

	return ""
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

func CreateHead(base string, titleText string, color Colors, description string) *html.Node {
	head := createHTMLElement("head", nil)

	meta := createHTMLElement("meta", map[string]string{"name": "viewport", "content": "width=device-width, initial-scale=1, viewport-fit=cover"})
	desc := createHTMLElement("meta", map[string]string{"name": "description", "content": description})

	title := createHTMLElement("title", nil, &html.Node{
		Type: html.TextNode,
		Data: titleText,
	})

	icons := createHTMLElement("link", map[string]string{"rel": "stylesheet", "href": "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css", "defer": ""})
	highlightcss := createHTMLElement("link", map[string]string{"rel": "stylesheet", "href": fmt.Sprintf("%sdefault.min.css", base)})
	highlightjs := createHTMLElement("script", map[string]string{"src": fmt.Sprintf("%shighlight.min.js", base)})
	highlightjsGo := createHTMLElement("script", map[string]string{"src": fmt.Sprintf("%sgo.min.js", base)})

	favicon := createHTMLElement("link", map[string]string{"rel": "icon", "type": "image/png", "href": fmt.Sprintf("%sicon.png", base)})

	b := createHTMLElement("base", map[string]string{"href": base})

	head.AppendChild(meta)
	head.AppendChild(desc)
	head.AppendChild(title)
	head.AppendChild(favicon)
	head.AppendChild(highlightcss)
	head.AppendChild(highlightjs)
	head.AppendChild(highlightjsGo)
	head.AppendChild(icons)

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

	colorData := fmt.Sprintf(`:root {
	--primary: %s;
	--text: %s;
	--background: %s;
	--secondary: %s;
	--title: %s;
}`, color.Primary, color.Text, color.Background, color.Secondary, color.TitleBar)

	head.AppendChild(createHTMLElement("style", nil, &html.Node{
		Type: html.TextNode,
		Data: colorData,
	}))

	for _, entry := range styleEntries {

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

func CreateIndex(summaryData []byte, pageData []byte, bodyText string, marker string, base string, titleText string, color Colors, inDev bool) ([]byte, error) {
	body := createHTMLElement("body", nil)

	n, err := html.Parse(bytes.NewReader(summaryData))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	s := visit(n, marker, "active")

	bodyData, err := html.Parse(bytes.NewReader(pageData))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	nav := CreateNav(n)

	page := createHTMLElement("div", id("page"), bodyData)

	title := createHTMLElement("img", add(add(id("logo"), "src", "logo.png"), "alt", "site logo"))

	buttons := createHTMLElement("div", nil)

	burgerIcon := createHTMLElement("i", add(id("bars"), "class", "fa fa-bars"))

	hamburgerButton := createHTMLElement("button", add(id("menu"), "aria-label", "navigation menu button"), burgerIcon)

	headerBar := createHTMLElement("div", id("header"), title, buttons, hamburgerButton)

	snackbar := createHTMLElement("div", id("snackbar"))

	main := createHTMLElement("div", id("main"), headerBar, nav, page, snackbar)

	body.AppendChild(main)

	docsJs := createHTMLElement("script", nil, &html.Node{
		Type: html.TextNode,
		Data: docScript,
	})
	body.AppendChild(docsJs)

	if inDev {
		webSocket := createHTMLElement("script", nil, &html.Node{
			Type: html.TextNode,
			Data: devSocket,
		})
		body.AppendChild(webSocket)
	}

	root := createHTMLElement("html", add(id("docs-root"), "lang", "en"), CreateHead(base, fmt.Sprintf("%s - %s", s, titleText), color, getDesc(bodyText)), body)

	r, err := renderHTML(root)
	if err != nil {
		return nil, err
	}

	return []byte(replaceHintTags(string(r))), nil
}

func getDesc(page string) string {
	re := regexp.MustCompile(`(?m)^#{1,6} .*\n`)

	// Replace all matches with an empty string
	cleanedText := re.ReplaceAllString(page, "")

	re = regexp.MustCompile(`\[(.*?)\]\(.*?\)|\[(.*?)\]\[.*?\]|\[.*?\]: .*?\n`)

	// Replace all matches with the captured group
	cleanedText = re.ReplaceAllStringFunc(cleanedText, func(match string) string {
		// Find matches for inline links and reference links
		inlineMatch := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
		referenceMatch := regexp.MustCompile(`\[(.*?)\]\[.*?\]`)

		// Check which type of match it is and return the text inside the square brackets
		if inlineMatch.MatchString(match) {
			return inlineMatch.ReplaceAllString(match, "$1")
		} else if referenceMatch.MatchString(match) {
			return referenceMatch.ReplaceAllString(match, "$1")
		}
		// Return an empty string for reference definitions
		return ""
	})

	sentences := strings.Split(cleanedText, ".")

	return sentences[0] + "."
}

func replaceHintTags(content string) string {
	// Define the regex pattern to match the hint tags
	re := regexp.MustCompile(`(?s){% hint style=[“"](\w+)[”"] %}(.*?){% endhint %}`)

	// Replace the hint tags with the appropriate div
	return re.ReplaceAllStringFunc(content, func(s string) string {
		matches := re.FindStringSubmatch(s)
		if len(matches) != 3 {
			return s
		}
		style, hintContent := matches[1], matches[2]
		return fmt.Sprintf(`<div class="hint %s">%s</div>`, style, hintContent)
	})
}
