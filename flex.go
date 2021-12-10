package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/otiai10/copy"
)

type Flex struct {
	Components map[string]bool
	Source     string
	Target     string
}

func (f *Flex) componentPath(name string) string {
	filename := fmt.Sprintf("%s.html", name)
	return path.Join(f.Source, "components", filename)
}

func (f *Flex) bundle() {
	start := time.Now()
	fmt.Println("[INFO] Bundling...")

	f.ClearTarget()

	entries, err := os.ReadDir(f.Source)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filepath := path.Join(f.Source, entry.Name())
			file := NewFile(f, filepath, "")
			file.Load()
			file.CleanUp()
			file.Write()
		}
	}

	// copy static files
	f.copyStatic()

	duration := time.Since(start).Milliseconds()
	fmt.Printf("[INFO] %d components loaded bundled in %dms.\n", len(f.Components), duration)
}

func (f *Flex) copyStatic() {
	// copy static files to dist dir
	sourcePath := "static"
	targetPath := path.Join("dist", "static")
	if err := copy.Copy(sourcePath, targetPath); err != nil {
		log.Fatal(err)
	}
}

func (f *Flex) loadComponent(name string, props string) *File {
	// add to components set
	f.Components[name] = true

	component := NewFile(f, f.componentPath(name), props)
	component.Load()
	return component
}

func (f *Flex) ClearTarget() {
	// remove target dir
	if err := os.RemoveAll(f.Target); err != nil {
		log.Fatal(err)
	}

	// create target dir
	if err := os.MkdirAll(f.Target, 0755); err != nil {
		log.Fatal(err)
	}
}

func NewFlex(source string, target string) *Flex {
	return &Flex{
		Components: make(map[string]bool),
		Source:     source,
		Target:     target,
	}
}
