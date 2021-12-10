package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("F L E X")
	fmt.Println("- - - -")
	watchFlag := flag.Bool("watch", false, "Watch the src directory for changes and serves dist on 8080.")

	flag.Parse()

	f := NewFlex("src", "dist")
	f.bundle()

	if *watchFlag {
		serve()
		fmt.Println("[INFO] Watching file changes...")
		watch(f)
	}
}
