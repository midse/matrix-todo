package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	ui "github.com/gizak/termui"
)

var logger *log.Logger
var parTime *ui.Par
var dataFile = "./data.json"
var focusTaskColor = "fg-white,bg-red"
var focusOnTask = regexp.MustCompile(`\[(.*?)\]\(` + focusTaskColor + `\)`)
var currentBoard int

func findCurrentBlock(blocks []*ui.List) *ui.List {
	for _, block := range blocks {
		if block.BorderFg == ui.ColorYellow {
			return block
		}
	}

	return nil
}

func findCurrentTask(blocks []*ui.List) string {

	block := findCurrentBlock(blocks)
	for _, item := range block.Items {
		if strings.HasSuffix(item, "("+focusTaskColor+")") {
			return focusOnTask.ReplaceAllString(item, "$1")[4:]
		}
	}

	return ""
}

func setFocusOnBlock(blocks []*ui.List, blockName string) {
	for _, block := range blocks {
		if block.BorderFg == ui.ColorYellow {
			block.BorderFg = ui.ColorWhite
		}

		if block.BorderLabel == blockName {
			block.BorderFg = ui.ColorYellow
		}
	}
}

func setFocusOnTask(blocks []*ui.List, blockName string, taskName string) {
	logger.Println("Try to focus on " + taskName)
	for indexBlock, block := range blocks {
		if block.BorderLabel == blockName {
			for indexItem, item := range block.Items {
				if item[4:] == taskName {
					logger.Println("Focus on " + item)
					blocks[indexBlock].Items[indexItem] = "[" + item + "](" + focusTaskColor + ")"
					break
				}
			}

			break
		}
	}
}

func getCurrentBoard(content Content) Board {
	return content.Boards[currentBoard]
}

func getNextBoard(content Content) Board {
	currentBoard++

	if currentBoard >= len(content.Boards) {
		currentBoard = 0
	}
	return content.Boards[currentBoard]
}

func getPreviousBoard(content Content) Board {
	currentBoard--

	if currentBoard < 0 {
		currentBoard = len(content.Boards) - 1
	}

	return content.Boards[currentBoard]
}

func headers(boardName string) *ui.Row {
	parTime = ui.NewPar("Rendering...")
	parTime.Border = false
	par2 := ui.NewPar(boardName)
	par2.Height = 3

	ls1 := ui.NewList()
	ls1.Items = []string{
		"[n] New board          [t] New task",
		"[k] Next board         [r] Remove task",
		"[j] Previous board     [x] Mark task done",
		"[q] Exit               [d] Task details"}
	ls1.Height = 6
	ls1.ItemFgColor = ui.ColorYellow
	ls1.BorderLabel = "Menu"

	return ui.NewRow(ui.NewCol(7, 0, parTime, par2),
		ui.NewCol(5, 0, ls1))
}

func board(board Board) ([]*ui.Row, []*ui.List) {
	blocks := []*ui.List{}
	rows := []*ui.Row{}
	rows = append(rows, headers(board.Name))

	cols := []*ui.Row{}
	for index, block := range board.Blocks {
		ls0 := ui.NewList()

		var tasks []string

		for _, task := range block.Tasks {
			symbol := "[ ]"
			if task.Status == "DONE" {
				symbol = "[X]"
			}
			taskContent := symbol
			taskContent = taskContent + " " + task.Name
			tasks = append(tasks, taskContent)
		}
		ls0.Items = tasks
		ls0.ItemFgColor = ui.ColorYellow

		if index == 0 {
			ls0.BorderFg = ui.ColorYellow
		}
		ls0.BorderLabel = block.Name
		ls0.Height = 20

		cols = append(cols, ui.NewCol(6, 0, ls0))
		blocks = append(blocks, ls0)

		if len(cols) == 2 {
			rows = append(rows, ui.NewRow(cols...))
			cols = []*ui.Row{}
		}

	}

	return rows, blocks
}

func addTask(board Board, blockName string, taskName string) Board {
	logger.Println("New task on " + blockName + " --> " + taskName)
	for index, block := range board.Blocks {
		if block.Name == blockName {
			logger.Println("Add task to board")
			now := time.Now()
			creationDate := now.Format("2006-01-02 15:04:05")
			task := Task{ID: 9999, Name: taskName, Description: "", CreationDate: creationDate, DueDate: "", Status: "TODO"}
			logger.Println(task)
			board.Blocks[index].Tasks = append(board.Blocks[index].Tasks, task)

			break
		}
	}
	return board
}

func removeTask(board Board, blockName string, taskName string) Board {
	logger.Println("Set task status on " + blockName + " --> " + taskName)
	for indexBlock, block := range board.Blocks {
		if block.Name == blockName {
			for indexTask, task := range block.Tasks {
				if task.Name == taskName {
					logger.Println("Remove task " + task.Name)
					board.Blocks[indexBlock].Tasks = append(board.Blocks[indexBlock].Tasks[:indexTask], board.Blocks[indexBlock].Tasks[indexTask+1:]...)
					break
				}
			}

			break
		}
	}

	return board
}

func setTaskStatus(board Board, blockName string, taskName string, status string) Board {
	logger.Println("Set task status on " + blockName + " --> " + taskName)
	for indexBlock, block := range board.Blocks {
		if block.Name == blockName {
			for indexTask, task := range block.Tasks {
				if task.Name == taskName {
					logger.Println("Set status on " + task.Name)
					board.Blocks[indexBlock].Tasks[indexTask].Status = status
					break
				}
			}

			break
		}
	}

	return board
}

func getTaskStatus(board Board, blockName string, taskName string) string {
	logger.Println("Get task status on " + blockName + " --> " + taskName)
	for _, block := range board.Blocks {
		if block.Name == blockName {
			for _, task := range block.Tasks {
				if task.Name == taskName {
					return task.Status
				}
			}

			break
		}
	}

	return ""
}

func getTask(board Board, blockName string, taskName string) *Task {
	for _, block := range board.Blocks {
		if block.Name == blockName {
			for _, task := range block.Tasks {
				if task.Name == taskName {
					return &task
				}
			}

			break
		}
	}

	return nil
}

func saveData(content Content) {
	bytes, err := json.Marshal(content)

	if err != nil {
		logger.Println("Unable to marshal json")
	}

	ioutil.WriteFile(dataFile, bytes, 0644)
}

func loadData() Content {
	var content Content

	jsonData, err := ioutil.ReadFile(dataFile)

	if err != nil {
		panic(err)
	}
	json.Unmarshal(jsonData, &content)

	return content
}

func rerender(boardStruct Board, rows []*ui.Row, blocks []*ui.List, focusedBlock string, focusedTask string) ([]*ui.Row, []*ui.List) {
	ui.Clear()
	ui.Body.Rows = nil
	rows, blocks = board(boardStruct)

	ui.Body.AddRows(rows...)
	ui.Body.Align()

	if focusedBlock != "" {
		setFocusOnBlock(blocks, focusedBlock)
	}

	if focusedTask != "" {
		setFocusOnTask(blocks, focusedBlock, focusedTask)
	}
	ui.Render(ui.Body)

	return rows, blocks
}

func main() {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger = log.New(f, "", log.LstdFlags)

	content := loadData()

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	currentBoard = 0

	rows, blocks := board(getCurrentBoard(content))

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
