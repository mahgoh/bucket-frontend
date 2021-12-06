package main

import (
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

type File struct {
	Flex *Flex
	Path string
	In   string
	Out  string
}

func NewFile(flex *Flex, path string) *File {
	return &File{
		Flex: flex,
		Path: path,
	}
}

func (f *File) Load() {
	f.Read()
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

func (f *File) LoadDependencies() {
	pattern := regexp.MustCompile("<!--c:([a-z]+)-->")
	matches := pattern.FindAllStringSubmatch(f.In, -1)

	if matches == nil {
		return
	}

	for _, match := range matches {
		component := f.Flex.loadComponent(match[1])
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
