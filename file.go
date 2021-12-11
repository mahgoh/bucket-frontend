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
	Flex  *Flex
	Path  string
	Props string
	In    string
	Out   string
}

func NewFile(flex *Flex, path string, props string) *File {
	return &File{
		Flex:  flex,
		Path:  path,
		Props: props,
	}
}

func (f *File) Load() {
	f.Read()
	f.LoadProps()
	f.LoadDependencies()
}

func (f *File) Read() {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		log.Fatal(err)
	}
	f.In = string(data)
	f.Out = string(data)
}

func (f *File) Write() {
	sourcePath := strings.Split(f.Path, "/")
	filename := sourcePath[len(sourcePath)-1]
	filepath := path.Join(f.Flex.Target, filename)
	if err := os.WriteFile(filepath, []byte(f.Out), 0644); err != nil {
		log.Fatal(err)
	}
}

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

func (f *File) LoadDependencies() {
	pattern := regexp.MustCompile(`<component name="([a-z-]+)" (?:props="([a-zA-Z0-9;.\s!-]+)" )?/>`)
	matches := pattern.FindAllStringSubmatch(f.Out, -1)

	if matches == nil {
		return
	}

	for _, match := range matches {
		component := f.Flex.loadComponent(match[1], match[2])
		f.Out = strings.ReplaceAll(f.Out, match[0], component.Out)
	}
}

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
