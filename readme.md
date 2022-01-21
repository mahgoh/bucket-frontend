<p align="center">
  <img src="./bucket_logo.svg" width="600" />
</p>
<p align="center">  
  Manage your bucket list in a modern web experience.<br>
  <a href="https://github.com/972C8/bucket-webapp">Main Repository</a> | <a href="https://bucket-webapp.herokuapp.com">Demo</a>
</p>

# Bucket Frontend

This repository serves as the development environment for the frontend of the [Bucket](https://github.com/972C8/bucket-webapp) webapp. The frontend is developed using a simple custom HTML bundler written in [Go](https://go.dev) that produces static HTML files which are then integrated in the Bucket backend.

Read the [Introduction](#introdution) to learn more about how the bundler works and how to create components. For instructions on how to use bundler CLI, consult [Usage](#usage). Interactivity on the frontend is achieved with the use of custom JavaScript utility modules. For detailed information visit [JavaScript Modules Reference](#javascript-modules-reference). If you want to build the bundler yourself (e.g. for a different platform), see [Build](#build).

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

## JavaScript Modules Reference

Interactivity of the frontend is achieved with the use of JavaScript. Utility functions are grouped into modules that provide a unified way to tackle various tasks. Find below the detailed reference.

### `Request`

The request module groups functions to make HTTP requests using the promise-based [Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API). Required headers such as `X-XSRF-TOKEN` or `Content-Type` are set and the functions return a Promise. `endpoint` means the path of an URL (according to the [URL spec](https://url.spec.whatwg.org/#url-representation)) including the root `/` (e.g., `/api/buckets`)

- `Request.HEAD(endpoint)`  
  Makes an HTTP request with the method `HEAD` to the specified `endpoint`.
- `Request.GET(endpoint)`  
  Makes an HTTP request with the method `GET` to the specified `endpoint`.
- `Request.POST(endpoint, payload, json?, jsonFeedback?)`  
  Makes an HTTP request with the method `POST` to the specified `endpoint`. The `payload` is stringified to JSON by default - this can be disabled by setting the third argument `json` to `false` (is required if using FormData). Additionally, the response body is expected to be JSON and is parsed before returned by default - this can be disabled by setting the fourth argument `jsonFeedback` to false.
- `Request.PUT(endpoint, payload, jsonFeedback?)`  
  Makes an HTTP request with the method `PUT` to the specified `endpoint`. The `payload` is stringified to JSON by default. Additionally, the response body is expected to be JSON and is parsed before returned by default - this can be disabled by setting the fourth argument `jsonFeedback` to false.
- `Request.DELETE(endpoint)`  
  Makes an HTTP request with the method `DELETE` to the specified `endpoint`.

### `API`

The API module groups functions to make specific HTTP requests using the promise-based `Request` functions.

- `API.register(payload)`  
  Makes a request to create a new avatar with the provided `payload`.
- `API.login(payload)`  
  Makes a request to login the avatar with the provided `payload`.
- `API.getProfile()`  
  Makes a request to get the profile of the logged in avatar.
- `API.putProfile(profile)`  
  Makes a request to update the profile of the logged in avatar with the provided `profile`.
- `API.postAvatarPicture(data)`  
  Makes a request to update the profile picture of the logged in avatar with the provided image `data`.
- `API.postBucketItem(bucketItem)`  
  Makes a request to create the provided `bucketItem`.
- `API.getBucketItem(bucketItemID)`  
  Makes a request to get the bucket item with the provided `bucketItemID`.
- `API.getBucketItems(params)`  
  Makes a request to get all bucket items of the logged in avatar. Additionally, params can be provided that are sent with the request.
- `API.putBucketItem(bucketItemID, bucketItem)`  
  Makes a request to update the bucket item with the provided `bucketItemID` and the new `bucketItem`.
- `API.deleteBucketItem(bucketItemID)`  
  Makes a request to delete the bucket item with the provided `bucketItemID`.
- `API.postBucket(bucket)`  
  Makes a request to create the provided `bucket`.
- `API.getBucket(bucketID)`  
  Makes a request to get the bucket with the provided `bucketID`.
- `API.getBuckets()`  
  Makes a request to get all buckets of the logged in avatar.
- `API.putBucket(bucketID, bucket)`  
  Makes a request to update the bucket with the provided `bucketID` and the new `bucket`.
- `API.deleteBucket(bucketID)`  
  Makes a request to delete the bucket with the provided `bucketID`.
- `API.postLabel(label)`  
  Makes a request to create the provided `label`.
- `API.getLabel(labelID)`  
  Makes a request to get the label with the provided `labelID`.
- `API.getLabels()`  
  Makes a request to get all labels of the logged in avatar.
- `API.putLabel(labelID, label)`  
  Makes a request to update the label with the provided `labelID` and the new `label`.
- `API.deleteLabel(labelID)`  
  Makes a request to delete the label with the provided `labelID`.
- `API.postBucketImage(data)`  
  Makes a request to create a bucket image with the provided image `data`.
- `API.postLocation(location)`  
  Makes a request to create the provided `location`.
- `API.getLocation(locationID)`  
  Makes a request to get the location with the provided `locationID`.
- `API.getLocations()`  
  Makes a request to get all locations of the logged in avatar.

### Shortcuts

The shortcuts module provides functions for shorter syntax. As a result, these functions are not grouped into a unified namespace.

- `param(name)`  
  Shortcut for getting a query param by `name`.
- `query(identifier)`  
  Shortcut for `document.querySelector(identifier)`.
- `queryAll(identifier)`  
  Shortcut for `document.querySelectorAll(identifier)`.
- `ready(callback)`  
  Shortcut for when document is ready.
- `redirect(path)`  
  Shortcut for redirecting to `path`.
- `val(identifier, value)`  
  Shortcut for getting the value of an element or setting the value if `value` is provided.

### `Animation` (static)

A static module that provides simple animation functions.

- `Animation.fadeIn(identifier)`  
  Fades in the element with provided `identifier`.
- `Animation.fadeOut(identifier)`
  Fades out the element with provided `identifier`.
- `Animation.slideIn(identifier)`
  Slides in the element with provided `identifier`.
- `Animation.slideOut(identifier)`
  Slides out the element with provided `identifier`.

### `Filter`

Processes filtering and sorting query params and generates new params based off of these.

Initialize a filter with `const filter = new Filter(options)` where `options` define the allowed params with default values.

#### Options

An object with key-value pairs where the key is the query param key to bind with and value a `FilterParam`.

##### FilterParam

Defines what query params to look for and how to process them. There are three different types: boolean, number, and sort.

- `boolean`  
  True if the query param key exists neglecting a potential value.
- `number`  
  Parses the query param value as an integer or null if inavlid.
- `sort`  
  Checks if the value with an optional preceding `-` is included in the allowed list. If it is not included, the default value is set. If the preceding `-` exists, ordering is set to descending, otherwise ascending.

##### Example

Complete options object with all `FilterParam` types included.

```js
{
  completed: {
    type: 'boolean',
    default: null,
  },
  bucketId: {
    type: 'number',
    default: null,
  },
  sort: {
    type: 'sort',
    default: { key: 'dateToAccomplish', desc: false },
    allowed: [
      'title',
      'bucket',
      'dateToAccomplish',
      'dateAccomplishedOn',
    ],
  },
}
```

- `filter.query(key?, value?)`
  Returns a query param string of the present filter settings and setting the provided `key` with the provided `value`. This can be used to set the `href` a button. If not key and value is provided, it simply returns the current filter settings.
- `filter.sortQuery(key)`  
  Returns a query param string of the present filter settings and set the sorting to the provided `key`. Automatically reverses the ordering if the provided `key` is already the active sorting.
- `filter.toggleQuery(key)`  
  Returns a query param string of the present filter settings and toggles the value of the provided `key` between true, false and null.
- `filter.get(key)`  
  Returns the value for the provided `key`.

### `Form` (static)

A static module that provides functions to validate form elements. An array of `watchers` and the identifier of the submit button can be passed to `Form.watch(watchers, '#submit')` to setup all required validation listeners.

#### Watcher

- `identifier`  
  A string with a DOM selector.
- `type`  
  Type of the error styling: text or submit.
- `required`  
  If the element is required.
- `message`  
  Error message to display if invalid.
- `event?`  
  The event that fires the validation. Can be a string or an array of strings. By default `keyup`.
- `colorize?`  
  If the text color of the element should adapt to the value. By default `false`.
- `validator(value)`  
  A function that validates the provided `value` which is called when the event is fired. Should return true or false.

#### Functions

- `Form.watch(watchers, submit)`  
  Sets up all event listeners for the provided `watchers`. Automatically, validates the form elements upon input and enables or disables the `submit` button based on the validity of the form inputs.
- `Form.valid()`  
  Returns true or false based on the validity of the form inputs.

### `Format` (static)

A static module that provides functions to format values.

- `Format.camelToKebab(camelCase)`  
  Converts a provided `camelCase` string into a kebab-case string.
- `Format.capitalize(value)`  
  Capitalizes the first letter of the `value`.
- `Format.decodeURIComponent(value)`  
  Decodes a URIComponent string and replaces `+` with spaces, as decodeURIComponent expects spaces to be represented as `%20`.
- `Format.parseDate(dateString, short?)`  
  Parses given ISO date string and returns a formatted date as string in the format DD Month YYYY or if short is set to true in the format DD Mon YYYY.
- `Format.sort(items, key, desc?)`  
  Sorts an array of bucket items by the provided `key` in ascending or descending order (if `desc` is set).
- `Format.dateToInt(dateString, low)`  
  Converts a provided `dateString` into an integer. If the dateString is `null` it is set to 0 if `low` or otherwise 99999999. This ensures that if elements with a date set are always displayed first.
- `Format.groupByYear(array, key)`  
  Splits an array into multiple arrays, each for a separate year.

### `Notification` (static)

A static module that provides functions to create notifications.

- `Notification.notify(title, message)`  
  Creates a notification with the provided `title` and `message` that is displayed for 10s or until the user dismisses it.

### `Modal` (static)

A static module that provides functions to create modals.

- `Modal.confirm(title, message, callback)`  
  Creates a confirm modal with the provided `title` and `message` that has to be confirmed. Once confirmed, the callback is executed.
- `Modal.prompt(title, message, callback)`  
  Creates a prompt modal with the providec `title` and `message` that requries a date input. Once submitted, the callback is executed.

### `Random` (static)

A static module that provides functions to generate random hex strings.

- `Random.hex(length)`  
  Generates a random hexadecimal string of the provided `length`. By default the length is 8.

### `Regex` (static)

A static module that provides regular expression presets.

- `Regex.preset(preset, min?, max?, flags?)`  
  Returns a RegExp of the given preset or the default preset if `preset` does not exist. In addition, it compiles the preset with a provided `min` and `max` if applicable and provided `flags`.
- `Regex.validate(regex, value)`
  Validates if the provided `value` matches the provided `regex`.

### `Store` (static)

A static module that provides a simple state management and provides resources.

- `Store.get(key)`  
  Get the value for a specified `key`.
- `Store.set(key, value)`  
  Set the `value` for a specified `key`.

### `Template` (static)

A static module that provides functions to render templates.

- `Template.render(identifier, rules, parent)`  
  Renders the template with the provided `identifier` and the provided `rules` and append it to the `parent`.

#### Rules

Rules allow to describe how to manipulate elements inside the template. The corresponding elements must have a data attribute with a specified key e.g., `data-title`. This key is required in any rule. A rule can than have multiple properties that define how to process the element.

- `key`  
  A string that matches `data-{key}`.
- `textContent`  
  Sets the textContent of the element to the specified value.
- `html`  
  Sets the innerHTML of the element to the specified value.
- `classes`  
  An array of strings. Adds each string to the classList of the element.
- `src`
  Sets the src attribute of the element to the specified value.
- `href`  
  Sets the href attribute of the element to the specified value.
- `attr`  
  An object of key-value pairs. For each key, set the attribute with the name key to the specified value.
- `listener`  
  Sets the onClick EventListener to the specified function.
- `remove`  
  Removes the element if its textContent is empty.

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

## Dependencies

To watch for file changes, the third-party development dependency `fsnotify` is used. However, this is only required for the watch mode.

- [fsnotify](github.com/fsnotify/fsnotify) ([BSD-3-Clause License](https://github.com/fsnotify/fsnotify/blob/master/LICENSE))
