package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	docopt "github.com/docopt/docopt-go"
	ui "github.com/gizak/termui"
)

var logger *log.Logger
var parTime *ui.Par
var dataFile = "./data.json"
var focusTaskColor = "fg-white,bg-red"
var focusOnTask = regexp.MustCompile(`\[(.*?)\]\(` + focusTaskColor + `\)`)
var currentBoard int

func main() {

	usage := `Matrix Todo - Simplistic todo list app

Usage:
	matrix-todo [ (-f|--file) <data-file> ] [ (-e|--encrypt) | (-d|--decrypt)]
	matrix-todo [ help | -h | --help ]
	matrix-todo [ version | -v | -V | --version ]

Options:
	-h --help        Show this screen.
	-v --version     Show version.
	-f --file        Read/write data to this file (default: ./data.json).
	-e --encrypt     Encrypt data file using a password.
	-d --decrypt     Decrypt data file using a password.

Examples:
    # Use a custom location for data and decrypt its content
    $ matrix-todo --file ~/Documents/my_todo_list --decrypt`

	arguments, _ := docopt.ParseDoc(usage)

	logFile, err := os.OpenFile("matrix-todo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create log file")
		os.Exit(1)
	}
	defer logFile.Close()

	logger = log.New(logFile, "", log.LstdFlags)

	logger.Println(arguments)

	if arguments["<data-file>"] != nil {
		dataFile, _ = arguments.String("<data-file>")
	}

	logger.Printf("Loading data from '%s' file", dataFile)

	content := loadData()

	err = ui.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialize UI")
		os.Exit(1)
	}
	defer ui.Close()

	currentBoard = 0

	rows, blocks := renderBoard(getCurrentBoard(content))

	ui.Body.AddRows(rows...)
	ui.Body.Align()

	ui.Render(ui.Body)

	drawTicker := time.NewTicker(time.Second)
	go func() {
		for {
			now := time.Now()
			parTime.Text = now.Format("2006-01-02 15:04:05")
			ui.Render(parTime)

			<-drawTicker.C
		}
	}()

	setupHandlers(content, rows, blocks)

	ui.Loop()
}
