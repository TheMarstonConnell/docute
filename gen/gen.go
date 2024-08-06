package gen

import (
	"embed"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed fonts/*
var fonts embed.FS

//go:embed thirdparty/default.min.css
var defaultcss []byte

//go:embed thirdparty/font-awesome.min.css
var fontawesome []byte

//go:embed thirdparty/go.min.js
var gomin []byte

//go:embed thirdparty/highlight.min.js
var highlight []byte

//go:embed empty.png
var empty []byte

func ReplaceMarkdownLinks(text string) string {
	re := regexp.MustCompile(`(\[[^\]]+\]\()([^\)/][^\)]+)(\))`)

	// Function to add leading slash
	t := re.ReplaceAllString(text, `$1/$2$3`)
	return strings.ReplaceAll(t, "/http", "http")
}

func MakeAbsoluteLinks(text string) string {
	re := regexp.MustCompile(`\(([^()]+)\.md\)`)

	return re.ReplaceAllString(text, "($1.html)")
}

func Gen(out string, base string, titleText string, inDev bool) error {
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

	color := DefaultColors()
	colorFile, err := os.ReadFile("colors.yaml")
	if err == nil {
		err = yaml.Unmarshal(colorFile, &color)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		fmt.Println(err)
	}

	// fileData := ReplaceMarkdownLinks(string(d))
	htmlData := mdToHTML([]byte(dd))
	err = Walk(dir, out, htmlData, base, titleText, color, inDev)
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
	if err != nil {
		logo = empty
	}

	err = os.WriteFile(path.Join(out, "logo.png"), logo, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	favi, err := os.ReadFile("icon.png")
	if err == nil {
		err = os.WriteFile(path.Join(out, "icon.png"), favi, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	err = os.WriteFile(path.Join(out, "default.min.css"), defaultcss, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.WriteFile(path.Join(out, "font-awesome.min.css"), fontawesome, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.WriteFile(path.Join(out, "go.min.js"), gomin, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.WriteFile(path.Join(out, "highlight.min.js"), highlight, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Walk(dir string, out string, summary []byte, base string, titleText string, color Colors, inDev bool) error {
	_ = os.MkdirAll(out, os.ModePerm)

	contents, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read %s | %w", dir, err)
	}
	for _, content := range contents {
		p := path.Join(dir, content.Name())
		o := path.Join(out, content.Name())
		if content.Name()[0:1] == "." {
			continue
		}

		if content.IsDir() {
			err = os.MkdirAll(o, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err := Walk(p, o, summary, base, titleText, color, inDev)
			if err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}

		ext := path.Ext(p)
		if ext != ".md" {
			f, err := os.ReadFile(p)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err = os.WriteFile(o, f, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}

		newPath := o[:len(o)-2] + "html"
		f, err := os.ReadFile(p)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// fileData := ReplaceMarkdownLinks(string(f))
		fileData := MakeAbsoluteLinks(string(f))
		htmlData := mdToHTML([]byte(fileData))

		newData, err := CreateIndex(summary, htmlData, fileData, newPath[5:], base, titleText, color, inDev)
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
