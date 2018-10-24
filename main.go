package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	ui "github.com/gizak/termui"
)

var logger *log.Logger
var parTime *ui.Par
var dataFile = "./data.json"
var focusTaskColor = "fg-white,bg-red"
var focusOnTask = regexp.MustCompile(`\[(.*?)\]\(` + focusTaskColor + `\)`)
var currentBoard int

func main() {
	logFile, err := os.OpenFile("matrix-todo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create log file")
		os.Exit(1)
	}
	defer logFile.Close()

	logger = log.New(logFile, "", log.LstdFlags)

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
