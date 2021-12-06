package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func main() {
	c := NewConfig("src", "dist")
	p := NewPure(c)
	p.bundle()
}

type Pure struct {
	Config *PureConfig
}

type PureConfig struct {
	Source string
	Output string
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

func NewPure(config *PureConfig) *Pure {
	return &Pure{
		Config: config,
	}
}

func NewConfig(source string, output string) *PureConfig {
	return &PureConfig{
		Source: source,
		Output: output,
	}
}

func (p *Pure) bundle() {
	start := time.Now()
	fmt.Println("[INFO] Bundling...")

	index, err := os.Open(p.entryPath())
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	scanner := bufio.NewScanner(index)

	out, err := os.Create(p.outputPath())
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	for scanner.Scan() {
		line := p.readLine(scanner.Text())
		if _, err := out.WriteString(line); err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	duration := time.Since(start).Milliseconds()
	fmt.Printf("[INFO] Bundle completed in %dms.\n", duration)
}

func (p *Pure) readLine(line string) string {
	match := p.searchLine(line)
	if match == "" {
		return line
	}

	return p.loadComponent(match)
}

func (p *Pure) searchLine(line string) string {
	pattern := regexp.MustCompile("<!--c:([a-z]+)-->")

	if res := pattern.FindStringSubmatch(line); res != nil {
		return res[1]
	}

	return ""
}

func (p *Pure) loadComponent(name string) string {
	file, err := os.Open(p.componentPath(name))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return strings.Join(lines, "")
}
