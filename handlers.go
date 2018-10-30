package main

import (
	"regexp"
	"strings"

	ui "github.com/gizak/termui"
)

var handlerMapping = map[string]func(ui.Event, *Content, []*ui.Row, []*ui.List){
	"t":          newTaskHandler,
	"n":          newBoardHandler,
	"d":          taskDetailHandler,
	"r":          removeTaskHandler,
	"x":          markTaskHandler,
	"k":          nextBoardHandler,
	"j":          previousBoardHandler,
	"<Keyboard>": genericKeyboardHandler,
	"<Up>":       previousTaskHandler,
	"<Down>":     nextTaskHandler,
	"<Right>":    nextBlockHandler,
	"<Left>":     previousBlockHandler,
	"<Resize>":   resizeHandler,
	"<Escape>":   escapeHandler,
	"q":          quitHandler,
}

var blockEvents = false
var writablePopup = false
var popupType = ""
var buffer = ""
var popup *ui.Par

func newTaskHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		blockEvents = true
		writablePopup = true
		popupType = "NEWTASK"
		popup = createSimplePopup("New task --> " + findCurrentBlock(blocks).BorderLabel)

		ui.Render(popup)
	}
}

func newBoardHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		blockEvents = true
		writablePopup = true
		popupType = "NEWBOARD"
		popup = createSimplePopup("New Board")

		ui.Render(popup)
	}
}

func taskDetailHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		blockEvents = true
		writablePopup = false

		currentBlock := findCurrentBlock(blocks)
		taskName := findCurrentTask(blocks)

		if taskName == "" {
			blockEvents = false
		} else {
			currentTask := getTask(getCurrentBoard(*content), currentBlock.BorderLabel, taskName)
			popup = createSimplePopup(currentTask.Name)
			popup.Text += "Created: " + currentTask.CreationDate + "\n"
			popup.Text += "Due to: " + currentTask.DueDate + "\n"
			popup.Text += "Description: " + currentTask.Description + "\n"
			popup.Text += "Status: " + currentTask.Status + "\n"
			popup.Height = 10

			ui.Render(popup)
		}
	}
}

func markTaskHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		taskName := findCurrentTask(blocks)
		blockName := findCurrentBlock(blocks).BorderLabel

		if len(taskName) > 0 {

			taskStatus := getTaskStatus(getCurrentBoard(*content), blockName, taskName)

			if taskStatus == "TODO" {
				taskStatus = "DONE"
			} else {
				taskStatus = "TODO"
			}
			content.Boards[currentBoard] = setTaskStatus(getCurrentBoard(*content), blockName, taskName, taskStatus)

			rows, blocks = rerender(getCurrentBoard(*content), rows, blocks, blockName, taskName)
		}
	}
}

func removeTaskHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		taskName := findCurrentTask(blocks)
		blockName := findCurrentBlock(blocks).BorderLabel

		if len(taskName) > 0 {
			content.Boards[currentBoard] = removeTask(getCurrentBoard(*content), blockName, taskName)

			rows, blocks = rerender(getCurrentBoard(*content), rows, blocks, blockName, taskName)
		}
	}
}

func previousBoardHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		myBoard := getPreviousBoard(*content)
		rows, blocks = rerender(myBoard, rows, blocks, "", "")
	}
}

func nextBoardHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		myBoard := getNextBoard(*content)
		rows, blocks = rerender(myBoard, rows, blocks, "", "")
	}
}

func genericKeyboardHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if blockEvents && writablePopup {
		if e.ID == "C-8>" {
			if len(buffer) > 1 {
				buffer = buffer[0 : len(buffer)-1]
			}
		} else if e.ID == "<Enter>" {
			logger.Println("<Enter> received")
			if buffer[1:] != "" {
				taskName := ""
				blockName := ""

				if popupType == "NEWTASK" {
					taskName = buffer[1:]
					blockName = findCurrentBlock(blocks).BorderLabel
					content.Boards[currentBoard] = addTask(getCurrentBoard(*content), blockName, taskName)
				}

				if popupType == "NEWBOARD" {
					boardName := buffer[1:]
					newBoard := createNewBoard(boardName)
					content.Boards = append(content.Boards, newBoard)

					// Go to newly created board
					currentBoard = len(content.Boards) - 1
				}

				blockEvents = false
				buffer = ""
				popup.Text = ""
				popupType = ""
				rows, blocks = rerender(getCurrentBoard(*content), rows, blocks, blockName, taskName)
			}
		} else {
			input := e.ID
			if input != "<Space>" {
				r := regexp.MustCompile(`<[a-zA-Z0-9-]+>`)
				input = r.ReplaceAllString(input, "")
			}
			buffer += strings.Replace(input, "<Space>", " ", 1)
		}

		if len(buffer) > 0 {
			popup.Text = buffer[1:]
			ui.Render(popup)
		}
	}
}

func previousTaskHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		for _, block := range blocks {
			if block.BorderFg == ui.ColorYellow {
				found := false
				for index, item := range block.Items {
					if strings.HasSuffix(item, "("+focusTaskColor+")") {
						previousIndex := (index - 1)

						if previousIndex < 0 {
							previousIndex = len(block.Items) - 1
						}

						block.Items[index] = focusOnTask.ReplaceAllString(item, "$1")
						block.Items[previousIndex] = "[" + block.Items[previousIndex] + "](" + focusTaskColor + ")"
						found = true
						break
					}
				}

				if !found && len(block.Items) > 0 {
					block.Items[len(block.Items)-1] = "[" + block.Items[len(block.Items)-1] + "](" + focusTaskColor + ")"
				}
			}
		}

		ui.Render(ui.Body)
	}
}

func nextTaskHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		for _, block := range blocks {
			if block.BorderFg == ui.ColorYellow {
				found := false
				for index, item := range block.Items {
					if strings.HasSuffix(item, "("+focusTaskColor+")") {
						nextIndex := (index + 1) % len(block.Items)

						block.Items[index] = focusOnTask.ReplaceAllString(item, "$1")
						block.Items[nextIndex] = "[" + block.Items[nextIndex] + "](" + focusTaskColor + ")"
						found = true
						break
					}
				}

				if !found && len(block.Items) > 0 {
					block.Items[0] = "[" + block.Items[0] + "](" + focusTaskColor + ")"
				}
			}
		}

		ui.Render(ui.Body)
	}
}

func previousBlockHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		for index, block := range blocks {
			if block.BorderFg == ui.ColorYellow {
				block.BorderFg = ui.ColorWhite

				for indexBlock, item := range block.Items {
					if strings.HasSuffix(item, "("+focusTaskColor+")") {
						block.Items[indexBlock] = focusOnTask.ReplaceAllString(item, "$1")
						break
					}
				}

				index--
				if index < 0 {
					index = len(blocks) - 1
				}
				blocks[index].BorderFg = ui.ColorYellow
				break
			}
		}

		ui.Render(ui.Body)
	}
}

func nextBlockHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		for index, block := range blocks {
			if block.BorderFg == ui.ColorYellow {
				block.BorderFg = ui.ColorWhite

				for indexBlock, item := range block.Items {
					if strings.HasSuffix(item, "("+focusTaskColor+")") {
						block.Items[indexBlock] = focusOnTask.ReplaceAllString(item, "$1")
						break
					}
				}

				blocks[(index+1)%len(blocks)].BorderFg = ui.ColorYellow
				break
			}
		}

		ui.Render(ui.Body)
	}
}

func resizeHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	payload := e.Payload.(ui.Resize)
	ui.Body.Width = payload.Width
	ui.Body.Align()
	ui.Clear()
	ui.Render(ui.Body)
}

func escapeHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if blockEvents {
		ui.Render(ui.Body)
		buffer = ""
		popup.Text = ""
		blockEvents = false
	} else {
		ui.StopLoop()
		saveData(*content, &credentials)
	}
}

func quitHandler(e ui.Event, content *Content, rows []*ui.Row, blocks []*ui.List) {
	if !blockEvents {
		// press q to quit
		ui.StopLoop()
		saveData(*content, &credentials)
	}
}

func setupHandlers(content Content, rows []*ui.Row, blocks []*ui.List) {

	// for item := range handlerMapping {
	// 	logger.Println("Mapping " + item)

	// 	myFunc := func(e ui.Event) {
	// 		logger.Println(item)
	// 		handlerMapping[item](e, &content, rows, blocks)
	// 	}
	// 	ui.Handle(item, myFunc)
	// }

	ui.Handle("t", func(e ui.Event) { newTaskHandler(e, &content, rows, blocks) })
	ui.Handle("n", func(e ui.Event) { newBoardHandler(e, &content, rows, blocks) })
	ui.Handle("d", func(e ui.Event) { taskDetailHandler(e, &content, rows, blocks) })
	ui.Handle("r", func(e ui.Event) { removeTaskHandler(e, &content, rows, blocks) })
	ui.Handle("x", func(e ui.Event) { markTaskHandler(e, &content, rows, blocks) })
	ui.Handle("k", func(e ui.Event) { nextBoardHandler(e, &content, rows, blocks) })
	ui.Handle("j", func(e ui.Event) { previousBoardHandler(e, &content, rows, blocks) })
	ui.Handle("<Keyboard>", func(e ui.Event) { genericKeyboardHandler(e, &content, rows, blocks) })
	ui.Handle("<Up>", func(e ui.Event) { previousTaskHandler(e, &content, rows, blocks) })
	ui.Handle("<Down>", func(e ui.Event) { nextTaskHandler(e, &content, rows, blocks) })
	ui.Handle("<Right>", func(e ui.Event) { nextBlockHandler(e, &content, rows, blocks) })
	ui.Handle("<Left>", func(e ui.Event) { previousBlockHandler(e, &content, rows, blocks) })
	ui.Handle("<Resize>", func(e ui.Event) { resizeHandler(e, &content, rows, blocks) })
	ui.Handle("<Escape>", func(e ui.Event) { escapeHandler(e, &content, rows, blocks) })
	ui.Handle("q", func(e ui.Event) { quitHandler(e, &content, rows, blocks) })
}
