package generator

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"path"
	"path/filepath"
	"strings"
)
import _ "embed"

//go:embed docs.js
var docScript string

//go:embed normalize.css
var normalize string

//go:embed global.css
var globals string

type FolderNode struct {
	childFolders []*FolderNode
	childFiles   []*FileNode
	DisplayName  string
	Name         string
}

func (f *FolderNode) GetDisplayName() string {
	if len(f.DisplayName) == 0 {
		return f.Name
	}

	return f.DisplayName
}

type FileNode struct {
	Name        string
	DisplayName string
	Data        []byte
}

func (fn *FolderNode) AddFolder(node *FolderNode) {
	fn.childFolders = append(fn.childFolders, node)
}

func (fn *FolderNode) AddFile(file string, displayName string, data []byte) {
	c := FileNode{
		Name:        file,
		Data:        data,
		DisplayName: displayName,
	}
	fn.childFiles = append(fn.childFiles, &c)
}

func NewFolderNode(name string) *FolderNode {
	f := FolderNode{
		childFolders: make([]*FolderNode, 0),
		childFiles:   make([]*FileNode, 0),
		Name:         name,
	}

	return &f
}

type Generator struct {
	Target string
	Tree   *FolderNode
	Nav    *html.Node
}

func NewGenerator(target string) *Generator {
	g := Generator{
		Target: target,
	}
	return &g
}

func (g *Generator) EnterDir(dir string, parent *FolderNode) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, e := range entries {
		f := path.Join(dir, e.Name())

		if e.IsDir() {
			fn := NewFolderNode(e.Name())
			parent.AddFolder(fn)
			g.EnterDir(f, fn)
		} else {

			if e.Name() == "info.md" {
				d, err := os.ReadFile(f)
				if err != nil {
					fmt.Println(err)
					return
				}
				title := extractTitle(d)
				parent.DisplayName = title
				continue
			}

			if path.Ext(e.Name()) == ".md" {
				d, err := os.ReadFile(f)
				if err != nil {
					fmt.Println(err)
					return
				}
				title := extractTitle(d)
				parent.AddFile(e.Name(), title, d)
			}
		}
	}

}

func (g *Generator) Start() error {
	p, err := filepath.Abs(g.Target)
	if err != nil {
		return fmt.Errorf("failed to resolve filepath")
	}

	outPath, err := filepath.Abs("out")
	if err != nil {
		return fmt.Errorf("failed to resolve output filepath")
	}

	os.RemoveAll(outPath)

	fmt.Println(p)

	f := NewFolderNode("out")
	f.DisplayName = "Docs"
	g.Tree = f
	g.EnterDir(p, f)

	return nil
}

func (fn *FolderNode) print(depth int64) {
	for i := int64(0); i < depth; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("%s (%s)\n", fn.DisplayName, fn.Name)
	for _, folder := range fn.childFolders {
		folder.print(depth + 1)
	}

	for _, file := range fn.childFiles {
		for i := int64(0); i < depth+1; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("%s (%s)\n", file.DisplayName, file.Name)

	}
}

func (g *Generator) PrintTree() {
	g.Tree.print(0)
}

func (fn *FolderNode) walk(p string, node *html.Node) {
	cwd := path.Join(p, fn.Name)
	os.Mkdir(cwd, os.ModePerm)
	for _, folder := range fn.childFolders {
		fl := createHTMLElement("ul", nil)
		text := createHTMLElement("span", nil, &html.Node{
			Type: html.TextNode,
			Data: folder.GetDisplayName(),
		})

		f := createHTMLElement("li", map[string]string{"class": "dropdown"}, text, fl)
		node.AppendChild(f)
		folder.walk(cwd, fl)
	}

	for _, file := range fn.childFiles {
		filePath := path.Join(cwd, file.Name)
		atrs := make(map[string]string)
		atrs["data-page"] = strings.ReplaceAll(strings.ReplaceAll(filePath, ".md", ".html"), "out/", "")
		link := createHTMLElement("a", atrs, &html.Node{
			Type: html.TextNode,
			Data: file.DisplayName,
		})
		f := createHTMLElement("li", nil, link)
		node.AppendChild(f)

		convertFile(filePath, file.Data)
	}
}

func (g *Generator) Walk() {

	g.Nav = createHTMLElement("ul", nil)
	g.Tree.walk("", g.Nav)

	head := createHTMLElement("head", nil)
	body := createHTMLElement("body", nil)
	title := createHTMLElement("h1", nil, &html.Node{
		Type: html.TextNode,
		Data: "Docute",
	})

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

	pageAtrs := make(map[string]string)
	pageAtrs["id"] = "navbar"
	navBar := createHTMLElement("div", pageAtrs, title, g.Nav)

	pageAtrs["id"] = "page"
	page := createHTMLElement("div", pageAtrs)

	pageAtrs["id"] = "main"
	container := createHTMLElement("div", pageAtrs, navBar, page)

	body.AppendChild(container)
	body.AppendChild(docsjs)
	body.AppendChild(highlightInit)

	root := createHTMLElement("html", nil, head, body)

	i := path.Join("out", "index.html")
	sr, err := renderHTML(root)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(i, []byte(sr), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}
