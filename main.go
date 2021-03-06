package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	docopt "github.com/docopt/docopt-go"
	ui "github.com/gizak/termui"
)

const version = "0.2.0"

var logger *log.Logger
var parTime *ui.Par
var dataFile = "./data.json"
var focusTaskColor = "fg-white,bg-red"
var focusOnTask = regexp.MustCompile(`\[(.*?)\]\(` + focusTaskColor + `\)`)
var currentBoard int
var credentials Credentials

func main() {

	usage := `Matrix Todo - Simplistic todo list app

Usage:
	matrix-todo [ ((-f|--file) <data-file>) ] [ (-e|--encrypt) ]
	matrix-todo help | -h | --help
	matrix-todo version | -v | -V | --version

Options:
	-h --help        Show this screen.
	-v --version     Show version.
	-f --file        Read/write data to this file (default: ./data.json).
	-e --encrypt     Encrypt data file using a password.

Examples:
    # Use a custom location for data and decrypt its content
    $ matrix-todo --file ~/Documents/my_todo_list --decrypt`

	arguments, _ := docopt.ParseDoc(usage)

	displayVersion := checkMultipleBoolArgs(arguments, []string{"version", "-V", "--version"})

	logFile, err := os.OpenFile("matrix-todo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create log file")
		os.Exit(1)
	}
	defer logFile.Close()

	logger = log.New(logFile, "", log.LstdFlags)
	// logger.Println(arguments)

	if displayVersion {
		fmt.Printf("matrix-todo v%s\n", version)
		os.Exit(0)
	}

	if help, _ := arguments.Bool("help"); help {
		fmt.Println(usage)
		os.Exit(0)
	}

	if arguments["<data-file>"] != nil {
		dataFile, _ = arguments.String("<data-file>")
	}

	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		emptyDataContent := "{\"author\": \"Unknown\", \"boards\": [{\"board_name\": \"My First Matrix\",\"board_type\": \"eisenhower_matrix\",\"board_last_update\": \"\", \"board_blocks\": [{\"block_name\": \"Urgent/Important\",\"block_type\": \"list\", \"block_tasks\": [] },{ \"block_name\": \"Not Urgent/Important\",\"block_tasks\": []},{\"block_name\": \"Urgent/Not Important\",\"block_tasks\": []},{\"block_name\": \"Not Urgent/Not Important\",\"block_tasks\": []}]}]}"
		bytes := []byte(emptyDataContent)
		ioutil.WriteFile(dataFile, bytes, 0644)
	}

	if toEncrypt, _ := arguments.Bool("--encrypt"); toEncrypt {
		fmt.Println("*Warning* This action is irreversible! Your data file will be entirely encrypted.")
		fmt.Println("Use ctrl+c to abort")
		password := getPassword()
		derivedKey, salt, err := generateKeyFromPassword(password, nil)

		if err != nil {
			fmt.Println("Unable to generate key!")
			os.Exit(1)
		}

		credentials.derivedKey = derivedKey
		credentials.salt = salt

		content := loadData(nil)
		saveData(content, &credentials)

		fmt.Println("Data successfully encrypted! Launch again matrix-todo to open your data.")
		os.Exit(0)
	}

	if isDataFileEncrypted() {
		password := getPassword()

		salt := readSaltFromDataFile()
		derivedKey, salt, err := generateKeyFromPassword(password, &salt)

		credentials.derivedKey = derivedKey
		credentials.salt = salt

		if err != nil {
			fmt.Println("Unable to generate key!")
			os.Exit(1)
		}
	}

	logger.Printf("matrix-todo v%s\n", version)
	logger.Printf("Loading data from '%s' file\n", dataFile)

	content := loadData(&credentials)

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
