package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/otiai10/copy"
)

type Bundler struct {
	Components map[string]bool
	Source     string
	Target     string
}

// utility to resolve component path by name
func (f *Bundler) componentPath(name string) string {
	filename := fmt.Sprintf("%s.html", name)
	return path.Join(f.Source, "components", filename)
}

// for each file in the source directory
// 		read the file
//		load linked components, replace props, cleanup file
// 		write file to target directory
// copy static files to target directory
func (f *Bundler) bundle() {
	start := time.Now()
	fmt.Println("[INFO] Loading components...")

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
	fmt.Printf("[INFO] %d components loaded in %dms.\n", len(f.Components), duration)
}

// copy static files to target dir
func (f *Bundler) copyStatic() {
	sourcePath := "static"
	targetPath := path.Join(f.Target, "static")

	if err := copy.Copy(sourcePath, targetPath); err != nil {
		log.Fatal(err)
	}
}

// load component by (file)name and return it
func (f *Bundler) loadComponent(name string, props string) *File {
	// add to components set
	f.Components[name] = true

	component := NewFile(f, f.componentPath(name), props)
	component.Load()
	return component
}

// remove and recreate target directory
func (f *Bundler) ClearTarget() {
	// remove target dir
	if err := os.RemoveAll(f.Target); err != nil {
		log.Fatal(err)
	}

	// create target dir
	if err := os.MkdirAll(f.Target, 0755); err != nil {
		log.Fatal(err)
	}
}

func NewBundler(source string, target string) *Bundler {
	return &Bundler{
		Components: make(map[string]bool),
		Source:     source,
		Target:     target,
	}
}
