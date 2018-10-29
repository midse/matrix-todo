# matrix-todo - Simplistic todo list app (CLI only)

[![Build Status](https://travis-ci.org/midse/matrix-todo.svg?branch=master)](https://travis-ci.org/midse/matrix-todo)
[![Go Report Card](https://goreportcard.com/badge/github.com/midse/matrix-todo)](https://goreportcard.com/report/github.com/midse/matrix-todo)
[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/midse)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/midse/matrix-todo/blob/master/LICENSE)

Built just for fun using Golang and [termui](https://github.com/gizak/termui), almost too simplistic to be useful (for now). I was tired of sketching [Eisenhower Matrices](https://en.wikipedia.org/wiki/Time_management#The_Eisenhower_Method) on my notebook.

If you want to test this app, you can download the binaries provided in the [Releases](https://github.com/midse/matrix-todo/releases) section.

Encrypt your data file using AES GCM (key is derived from password using pbkdf2)

```bash
$ ./matrix-todo --encrypt
Data successfully encrypted! Launch again matrix-todo to open your data.
```

## To do

+ Task edition : description, due date
+ Allow scrolling within a block (infinite number of tasks should not be allowed)
+ Allow scrolling when terminal's height is too small
+ Board creation
+ Introduce new board types (for now only eisenhower_matrix is supported)
+ Reference task within a commit (custom task id per board?)
+ Monitor n git repositories per board

## Changelog

### v0.1.0 (2018-10-29)

+ Basic CLI args support
+ Data file encryption support (prototype)

### v0.1.1 (2018-10-29)

+ Small display improvement when --encrypt is used

### v0.1.2 (2018-10-29)

+ Base64 is not used anymore after encryption

### v0.1.3 (2018-10-29)

+ Provide binaries for Linux, Windows and Osx (feedback needed)

### v0.2.0 (2018-10-29)

+ Remove --decrypt option
+ Detect automatically if data file is encrypted
+ Create an empty data file with a single board when data file doesn't exist

## Contributing

Feel free to contribute. :)
Feedback is always welcome!