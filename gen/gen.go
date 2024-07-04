package gen

import (
	"fmt"
	"os"
	"path"
	"regexp"
)

func ReplaceMarkdownLinks(text string) string {
	re := regexp.MustCompile(`\(([^()]+)\.md\)`)

	return re.ReplaceAllString(text, "($1.html)")
}

func ListFiles(dir string, out string) error {
	p := path.Join(dir, "SUMMARY.md")
	data, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("cannot read %s | %w", p, err)
	}

	fmt.Println(ReplaceMarkdownLinks(string(data)))

	return nil
}

func Gen(dir string, out string) error {

	_ = os.RemoveAll(out)

	err := Walk(dir, out)
	if err != nil {
		return err
	}

	d, err := os.ReadFile(path.Join(out, "SUMMARY.html"))
	if err != nil {
		return err
	}

	index, err := CreateIndex(d)
	if err != nil {
		return err
	}

	_ = os.WriteFile(path.Join(out, "index.html"), index, os.ModePerm)
	return nil
}

func Walk(dir string, out string) error {
	contents, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read %s | %w", dir, err)
	}

	for _, content := range contents {
		p := path.Join(dir, content.Name())
		o := path.Join(out, content.Name())

		if content.IsDir() {
			_ = os.MkdirAll(o, os.ModePerm)
			err := Walk(p, o)
			if err != nil {
				return err
			}

			continue
		}

		ext := path.Ext(p)
		if ext != ".md" {
			continue
		}

		newPath := o[:len(o)-2] + "html"
		f, err := os.ReadFile(p)
		if err != nil {
			return err
		}

		fileData := ReplaceMarkdownLinks(string(f))
		htmlData := mdToHTML([]byte(fileData))

		_ = os.WriteFile(newPath, htmlData, os.ModePerm)
	}

	return nil
}
