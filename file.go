package main

import (
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type File struct {
	Bundler *Bundler
	Path    string
	Props   string
	In      string
	Out     string
}

func NewFile(bundler *Bundler, path string, props string) *File {
	return &File{
		Bundler: bundler,
		Path:    path,
		Props:   props,
	}
}

func (f *File) Load() {
	f.Read()
	f.LoadProps()
	f.LoadDependencies()
}

// reads file contents
func (f *File) Read() {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		log.Fatal(err)
	}
	f.In = string(data)
	f.Out = string(data)
}

// writes output contents to file in target directory
func (f *File) Write() {
	sourcePath := strings.Split(f.Path, "/")
	filename := sourcePath[len(sourcePath)-1]
	filepath := path.Join(f.Bundler.Target, filename)
	if err := os.WriteFile(filepath, []byte(f.Out), 0644); err != nil {
		log.Fatal(err)
	}
}

// replace prop placeholders with provided props
func (f *File) LoadProps() {
	pattern := regexp.MustCompile(`{{[\s]?\$([0-9]+)[\s]?}}`)
	matches := pattern.FindAllStringSubmatch(f.Out, -1)

	if matches == nil {
		return
	}

	for _, match := range matches {
		props := strings.Split(f.Props, ";")
		index, err := strconv.ParseInt(match[1], 10, 0)
		if err != nil {
			log.Fatal(err)
		}

		if len(props) > int(index) {
			f.Out = strings.ReplaceAll(f.Out, match[0], props[index])
		}
	}
}

// load additional components
func (f *File) LoadDependencies() {
	pattern := regexp.MustCompile(`<component[\s\n\r]+name="([a-z-]+)"[\s\n\r]+(?:props="([a-zA-Z0-9;.\s!-]+)"[\s\n\r]+)?/>`)
	matches := pattern.FindAllStringSubmatch(f.Out, -1)

	if matches == nil {
		return
	}

	for _, match := range matches {
		component := f.Bundler.loadComponent(match[1], match[2])
		f.Out = strings.ReplaceAll(f.Out, match[0], component.Out)
	}
}

// remove unnecessaty whitespace
func (f *File) CleanUp() {
	// remove whitespace between tags
	f.Out = regexp.MustCompile(`>\s+<`).ReplaceAllString(f.Out, "><")

	// remove whitespace before tags
	f.Out = regexp.MustCompile(`\s+<`).ReplaceAllString(f.Out, "<")

	// remove whitespace after tags
	f.Out = regexp.MustCompile(`>\s+`).ReplaceAllString(f.Out, ">")

	// trim the rest of the whitespace
	f.Out = regexp.MustCompile(`\s+`).ReplaceAllString(f.Out, " ")
}
