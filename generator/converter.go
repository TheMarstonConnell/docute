package generator

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"os"
	"strings"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func extractText(node ast.Node) string {
	var text string
	ast.WalkFunc(node, func(node ast.Node, entering bool) ast.WalkStatus {
		if textNode, ok := node.(*ast.Text); ok {
			text += string(textNode.Literal)
		}
		return ast.GoToNext
	})
	return text
}

func extractTitle(data []byte) string {
	parser := parser.New()

	// Parse the Markdown content
	root := markdown.Parse(data, parser)
	header := ""
	// Walk through the AST and extract headings
	ast.WalkFunc(root, func(node ast.Node, entering bool) ast.WalkStatus {
		if heading, ok := node.(*ast.Heading); ok && entering {
			if heading.Level > 1 {
				return ast.GoToNext
			}
			headingText := extractText(heading)
			header = headingText
		}
		return ast.GoToNext
	})
	return header
}

func convertFile(path string, data []byte) error {

	ht := mdToHTML(data)

	err := os.WriteFile(fmt.Sprintf("%s.html", strings.ReplaceAll(path, ".md", "")), ht, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write file %s | %w", path, err)
	}
	return nil
}
