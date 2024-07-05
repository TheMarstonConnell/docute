package gen

import (
	"embed"
	"fmt"
	"os"
	"path"
	"regexp"
)

//go:embed fonts/*
var fonts embed.FS

func ReplaceMarkdownLinks(text string) string {
	re := regexp.MustCompile(`(\[[^\]]+\]\()([^\)/][^\)]+)(\))`)

	// Function to add leading slash
	return re.ReplaceAllString(text, `$1/$2$3`)
}

func MakeAbsoluteLinks(text string) string {
	re := regexp.MustCompile(`\(([^()]+)\.md\)`)

	return re.ReplaceAllString(text, "($1.html)")
}

func Gen(out string, base string) error {
	_ = os.RemoveAll(out)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return err
	}

	d, err := os.ReadFile("SUMMARY.md")
	if err != nil {
		fmt.Println(err)
		return err
	}
	dd := MakeAbsoluteLinks(string(d))
	// fileData := ReplaceMarkdownLinks(string(d))
	htmlData := mdToHTML([]byte(dd))
	err = Walk(dir, out, htmlData, base)
	if err != nil {
		fmt.Println(err)
		return err
	}
	readme, err := os.ReadFile(path.Join(out, "README.html"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = os.WriteFile(path.Join(out, "index.html"), readme, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fontEntries, err := fonts.ReadDir("fonts")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, entry := range fontEntries {
		fmt.Println(entry.Name())
		fontData, err := fonts.ReadFile(path.Join("fonts", entry.Name()))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		err = os.WriteFile(path.Join(out, entry.Name()), fontData, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	logo, err := os.ReadFile("logo.png")
	if err == nil {
		err = os.WriteFile(path.Join(out, "logo.png"), logo, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func Walk(dir string, out string, summary []byte, base string) error {
	_ = os.MkdirAll(out, os.ModePerm)

	contents, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read %s | %w", dir, err)
	}
	for _, content := range contents {
		p := path.Join(dir, content.Name())
		o := path.Join(out, content.Name())
		fmt.Println(p)
		if content.Name()[0:1] == "." {
			continue
		}

		if content.IsDir() {
			err = os.MkdirAll(o, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err := Walk(p, o, summary, base)
			if err != nil {
				fmt.Println(err)
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
			fmt.Println(err)
			return err
		}

		fileData := ReplaceMarkdownLinks(string(f))
		htmlData := mdToHTML([]byte(fileData))

		newData, err := CreateIndex(summary, htmlData, newPath[5:], base)
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = os.WriteFile(newPath, newData, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
