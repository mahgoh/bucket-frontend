package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

type Pure struct {
	Config     *Config
	Components map[string]bool
}

func (p *Pure) entryPath() string {
	return path.Join(p.Config.Source, "index.html")
}

func (p *Pure) outputPath() string {
	return path.Join(p.Config.Output, "index.html")
}

func (p *Pure) componentPath(name string) string {
	filename := fmt.Sprintf("%s.html", name)
	return path.Join(p.Config.Source, "components", filename)
}

func (p *Pure) bundle() {
	start := time.Now()
	fmt.Println("[INFO] Bundling...")

	entry := NewFile(p, p.entryPath())
	entry.Load()
	entry.cleanUp()

	os.WriteFile(p.outputPath(), []byte(entry.Out), 0666)

	duration := time.Since(start).Milliseconds()
	fmt.Printf("[INFO] %d components loaded bundled in %dms.\n", len(p.Components), duration)
}

func (p *Pure) loadComponent(name string) *File {
	// add to components set
	p.Components[name] = true

	component := NewFile(p, p.componentPath(name))
	component.Load()
	return component
}

func NewPure(config *Config) *Pure {
	return &Pure{
		Config:     config,
		Components: make(map[string]bool),
	}
}
