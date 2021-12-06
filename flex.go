package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

type Flex struct {
	Components map[string]bool
	Source     string
	Target     string
}

func (f *Flex) entryPath() string {
	return path.Join(f.Source, "index.html")
}

func (f *Flex) outputPath() string {
	return path.Join(f.Target, "index.html")
}

func (f *Flex) componentPath(name string) string {
	filename := fmt.Sprintf("%s.html", name)
	return path.Join(f.Source, "components", filename)
}

func (f *Flex) bundle() {
	start := time.Now()
	fmt.Println("[INFO] Bundling...")

	entry := NewFile(f, f.entryPath())
	entry.Load()
	entry.cleanUp()

	os.WriteFile(f.outputPath(), []byte(entry.Out), 0666)

	duration := time.Since(start).Milliseconds()
	fmt.Printf("[INFO] %d components loaded bundled in %dms.\n", len(f.Components), duration)
}

func (f *Flex) loadComponent(name string) *File {
	// add to components set
	f.Components[name] = true

	component := NewFile(f, f.componentPath(name))
	component.Load()
	return component
}

func NewFlex(source string, target string) *Flex {
	return &Flex{
		Components: make(map[string]bool),
		Source:     source,
		Target:     target,
	}
}
