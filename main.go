package main

import (
	"flag"
	"fmt"
)

// Is set by the compiler during build
var GitSemverStr string

func main() {

	// Flags
	var flagWatchShort bool
	var flagWatchLong bool
	var flagServeShort bool
	var flagServeLong bool
	var flagVersionShort bool
	var flagVersionLong bool

	flagWatchDesc := "Watch the src directory for changes"
	flagServeDesc := "Serve the dist directory on port 8080"
	flagVersionDesc := "Print version"

	flag.BoolVar(&flagWatchShort, "w", false, flagWatchDesc)
	flag.BoolVar(&flagWatchLong, "watch", false, flagWatchDesc)
	flag.BoolVar(&flagServeShort, "s", false, flagServeDesc)
	flag.BoolVar(&flagServeLong, "serve", false, flagServeDesc)
	flag.BoolVar(&flagVersionShort, "v", false, flagVersionDesc)
	flag.BoolVar(&flagVersionLong, "version", false, flagVersionDesc)

	flag.Parse()

	flagWatch := flagWatchShort || flagWatchLong
	flagServe := flagServeShort || flagServeLong
	flagVersion := flagVersionShort || flagVersionLong

	if flagVersion {
		cmdVersion()
		return
	}

	// Bundle HTML
	b := NewBundler("src", "dist")
	b.bundle()

	// Watch for file changes
	if flagWatch {
		serveInBackground(b.Target)
		fmt.Println("[INFO] Watching file changes...")
		watch(b)
	}

	// Serve dist directory
	if flagServe {
		serve(b.Target)
	}
}

// Print program version
// is passed by the compiler
func cmdVersion() {
	fmt.Println(GitSemverStr)
}
