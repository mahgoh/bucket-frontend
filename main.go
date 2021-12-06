package main

func main() {
	c := NewConfig("src", "dist")
	p := NewPure(c)
	p.bundle()
}
