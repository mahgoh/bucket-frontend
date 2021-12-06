package main

type Config struct {
	Source string
	Output string
}

func NewConfig(source string, output string) *Config {
	return &Config{
		Source: source,
		Output: output,
	}
}
