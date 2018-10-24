# matrix-todo - Simplistic todo list app (CLI only)

[![Go Report Card](https://goreportcard.com/badge/github.com/midse/matrix-todo)](https://goreportcard.com/report/github.com/midse/matrix-todo)
[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/midse)

Built just for fun using Golang and [https://github.com/gizak/termui](termui), almost too simplistic to be useful (for now). I was tired of sketching [Eisenhower Matrices](https://en.wikipedia.org/wiki/Time_management#The_Eisenhower_Method) on my notebook.

## To do

+ Task edition : description, due date
+ Allow scrolling within a block (infinite number of tasks should not be allowed)
+ Allow scrolling when terminal's height is too small
+ Board creation
+ Handle CLI args
+ Introduce new board types (for now only eisenhower_matrix is supported)
+ Reference task within a commit (custom task id per board?)
+ Monitor n git repositories per board
+ Encrypt local storage?

Feel free to contribute. :)