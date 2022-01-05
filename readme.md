# bucket-frontend

This repository serves as the development environment for the frontend of the [Bucket](https://github.com/972C8/bucket-webapp) webapp. The frontend is developed using a simple custom HTML bundler written in [Go](https://go.dev) that produces static HTML files which are then integrated in the Bucket backend.

Read the [Introduction](#introdution) to learn more about how the bundler works and how to create components. For instructions on how to use bundler CLI, consult [Usage](#usage). If you want to build the bundler yourself (e.g. for a different platform), see [Build](#build).

Both the introduction and usage are prerequisites in order to work with the bundler.

There is also a short [section](#bundle-tailwindcss) on how to bundle tailwindcss for production. This step only needs to be done once at the end of development.

## Introduction

This section explains how the bundler works and how to write components and reference them.

Each file in the `src` directory serves as an entry point. These entry points are simple HTML files that can contain additional custom components. These components are located in the subdirectory `components` and are also just HTML files.

The bundler reads the entry files and loads referenced components recursively to generate a single HTML output file. It simply injects the components contents into the entry point file. However, components can reference other components endlessly.

The generated output per entry point file is written to an output file in the `dist` directory.

In addition, values can be passed to a component by the parent who referenced the component. More on this [later](#props).

### Using components

Let's see how such components are written and how they are referenced.
We assume the following source structure, where we have 2 entry files and 4 components.

```bash
# Source Structure
src
  components
    comp-a.html
    comp-b.html
    comp-c.html
    comp-d.html
  entry-1.html
  entry-2.html
```

```html
<!-- entry-1.html -->
<html>
  <component name="comp-a" />
</html>

<!-- entry-2.html -->
<html>
  <component name="comp-b" />
</html>

<!-- comp-a.html -->
<p>Component A simply contains text</p>

<!-- comp-b.html -->
<p>Component B references to component C</p>
<component name="comp-c" />

<!-- comp-c.html -->
<section>Component C contains a section</section>
```

For `entry-1`, `comp-a` is read and injected into the entry file. It simply replaces the component element with the contents of the component. In `entry-2`, `comp-b` itself references another component `comp-c`. `comp-c` is injected into `comp-b`, which is injected into `entry-2`. `comp-d` is never referenced and, therefore, ignored by the bundler. This generates the following output files:

```bash
# Dist Structure
dist
  entry-1.html
  entry-2.html
```

```html
<!-- entry-1.html -->
<html>
  <p>Component A simply contains text</p>
</html>

<!-- entry-2.html -->
<html>
  <p>Component B references to component C</p>
  <section>Component C contains a section</section>
</html>
```

### Props

Additionally, values can be passed to a component by the parent that references the component. These values are called props and essentially are an array of strings. A second attribute `props` on the component element contains this array - the values are separated by a `;`.

The values can be accessed by the receiving child with the syntax `{{$i}}` where `i` is the zero-based index in the array that is passed. The order of the props matters, obviously.

In the following example, the entry file passes two values to the child component: `Home` and `Ben`. These values are then accessed in the child component.

```html
<!-- src/entry.html -->
<html>
  <component name="pagetitle" props="Home;Ben" />
</html>

<!-- src/components/pagetitle.html -->
<h1>{{$0}}</h1>
<p>Hello {{$1}}!</p>
```

The bundler injects the `pagetitle` component into the entry file and replaces the placeholders with the `props` values. This generates the following output:

```html
<!-- dist/entry.html -->
<html>
  <h1>Home</h1>
  <p>Hello Ben!</p>
</html>
```

## Usage

The bundler is a single binary that is executed. It looks for entry files in the `src` directory relative to the current working directory. The directory `src/components` contains all components.

Simply execute the bundler binary to generate the outputs to the `dist` directory - relative to the current working directory.

Pre-built binaries for Windows X64, Intel-based macOS and ARM-based macOS are located in the `bin` directory or are attached to the associated GitHub release.

```bash
# Run on Windows (64 Bit)
$ .\bin\windows_x64.exe

# Run on macOS (Intel-based)
$ ./bin/darwin_x64

# Run on macOS (ARM-based)
$ ./bin/darwin_arm
```

Alternatively, the repository can be cloned and executed by using the `go run` command (requires [Go](https.//go.dev) to be installed).

```bash
# Run with go run
$ go run .
```

### Flags

Furthermore, the program provides the following flags. During development, it is recommended to use `-w` and `-s`.

```bash
$ ./bin/darwin_arm -w -s
```

#### -s -serve

The flag `-s` or `-serve` serves the `dist` directory on port 8080.

#### -w -watch

The flag `-w` or `-watch` enables watch mode. In watch mode, the source directory is watched for file changes. Upon change, all entry files are re-bundled automatically.

#### -h -help

The flag `-h` or `-help` prints information about the usage of the CLI program. Ignores all other flags.

#### -v -version

The flag `-v` or `-version` prints the version of the CLI program.

## Bundle Tailwindcss

In development, the play CDN of [tailwindcss](https://tailwindcss.com) is used. While this is a just-in-time CDN and is remarkably fast, it should not be used in production. Therefore, we generate a bundle containing only the used tailwindcss classes and link it instead of the CDN. This step only needs to be done once at the end of development.

```html
<!-- Development -->
<script src="https://cdn.tailwindcss.com?plugins=forms"></script>

<!-- Production -->
<link rel="stylesheet" href="/bundle.css" />
```

### Prerequisites

The bundle is generated by the tailwindcss CLI using [Node.js](https://nodejs.org/) and [npm](https://www.npmjs.com/) - so these are both required.

### Generate

The npm start script is setup to run the tailwindcss CLI using the provided `tailwind.config.js` config and the `-m` minify flag. It produces `./dist/static/bundle.css` containing all tailwindcss classes that are used in the files in the `dist` directory.

```bash
# Install dependencies (only once)
$ npm install

# Run generation
$ npm run start
```

## Build

### Prerequisites

In order to build the bundler, [Go](https.//go.dev) is required. To use the build script (Makefile) make is required as well.

### Makefile

The `Makefile` contains instructions to build for Windows X64, Intel-based macOS and ARM-based macOS. Additionally, it injects the current semver string into the program.

```bash
$ make
```

> **Note:** Tested on macOS 12.1 (ARM) with GNU Make 3.81 and Go 1.17.2

### Custom

Go includes tooling to build for a variety of operating systems and architectures. Run `go tool dist list` to see a list of all available platforms in the format `os/arch`.

The Go compiler builds for the current platform or, if specified, the platform defined in `GOOS` and `GOARCH`. For example, to build for Linux X64:

```bash
# set environment variables
$ export GOOS=linux GOARCH=amd64

# build for linux x64
$ go build .
```

To inject the current version, we can pass the sematic version string to the compiler with the flag `-ldflags`. First, set the environment variable `GIT_SEMVER_STR` to the result of `git describe --tags --dirty`. The run the build command and pass the environment variable to the compiler.

```bash
# set semver string as env var
$ export GIT_SEMVER_STR=$(git describe --tags --dirty)

# build with -ldflags and pass semver string to compiler
$ go build -ldflags "-X main.GitSemverStr=$GIT_SEMVER_STR" .
```

> **Note:** Tested on macOS 12.1 (ARM) with Go 1.17.2
