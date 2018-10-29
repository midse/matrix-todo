package main

import (
	"regexp"
	"strings"

	ui "github.com/gizak/termui"
)

func setupHandlers(content Content, rows []*ui.Row, blocks []*ui.List) {
	blockEvents := false
	writablePopup := false
	buffer := ""
	var popup *ui.Par

	ui.Handle("t", func(e ui.Event) {
		if !blockEvents {
			blockEvents = true
			writablePopup = true
			popup = ui.NewPar("")
			popup.Height = 3
			popup.Width = 50
			popup.Y = (ui.TermHeight() / 2) - (popup.Height / 2) - 10
			popup.X = (ui.TermWidth() / 2) - (popup.Width / 2)
			popup.BorderLabel = "New task --> " + findCurrentBlock(blocks).BorderLabel

			ui.Render(popup)
		}

	})

	ui.Handle("d", func(e ui.Event) {
		if !blockEvents {
			blockEvents = true
			writablePopup = false

			currentBlock := findCurrentBlock(blocks)
			taskName := findCurrentTask(blocks)

			if taskName == "" {
				blockEvents = false
			} else {
				currentTask := getTask(getCurrentBoard(content), currentBlock.BorderLabel, taskName)
				popup = ui.NewPar("")
				popup.Text += "Created: " + currentTask.CreationDate + "\n"
				popup.Text += "Due to: " + currentTask.DueDate + "\n"
				popup.Text += "Description: " + currentTask.Description + "\n"
				popup.Text += "Status: " + currentTask.Status + "\n"
				popup.Height = 10
				popup.Width = 50
				popup.Y = (ui.TermHeight() / 2) - (popup.Height / 2) - 10
				popup.X = (ui.TermWidth() / 2) - (popup.Width / 2)
				popup.BorderLabel = currentTask.Name

				ui.Render(popup)
			}
		}

	})

	ui.Handle("x", func(e ui.Event) {
		if !blockEvents {
			taskName := findCurrentTask(blocks)
			blockName := findCurrentBlock(blocks).BorderLabel

			if len(taskName) > 0 {

				taskStatus := getTaskStatus(getCurrentBoard(content), blockName, taskName)

				if taskStatus == "TODO" {
					taskStatus = "DONE"
				} else {
					taskStatus = "TODO"
				}
				content.Boards[currentBoard] = setTaskStatus(getCurrentBoard(content), blockName, taskName, taskStatus)

				rows, blocks = rerender(getCurrentBoard(content), rows, blocks, blockName, taskName)
			}
		}
	})

	ui.Handle("r", func(e ui.Event) {
		if !blockEvents {
			taskName := findCurrentTask(blocks)
			blockName := findCurrentBlock(blocks).BorderLabel

			if len(taskName) > 0 {
				content.Boards[currentBoard] = removeTask(getCurrentBoard(content), blockName, taskName)

				rows, blocks = rerender(getCurrentBoard(content), rows, blocks, blockName, taskName)
			}
		}
	})

	ui.Handle("k", func(e ui.Event) {
		if !blockEvents {
			myBoard := getNextBoard(content)
			rows, blocks = rerender(myBoard, rows, blocks, "", "")
		}
	})

	ui.Handle("j", func(e ui.Event) {
		if !blockEvents {
			myBoard := getPreviousBoard(content)
			rows, blocks = rerender(myBoard, rows, blocks, "", "")
		}
	})

	ui.Handle("<Keyboard>", func(e ui.Event) {
		if blockEvents && writablePopup {
			if e.ID == "C-8>" {
				if len(buffer) > 1 {
					buffer = buffer[0 : len(buffer)-1]
				}
			} else if e.ID == "<Enter>" {
				logger.Println("<Enter> received")
				if buffer[1:] != "" {
					taskName := buffer[1:]
					blockName := findCurrentBlock(blocks).BorderLabel
					content.Boards[currentBoard] = addTask(getCurrentBoard(content), blockName, taskName)

					blockEvents = false
					buffer = ""
					popup.Text = ""
					rows, blocks = rerender(getCurrentBoard(content), rows, blocks, blockName, taskName)
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
	})

	ui.Handle("<Up>", func(e ui.Event) {
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
	})

	ui.Handle("<Down>", func(e ui.Event) {
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
	})

	ui.Handle("<Left>", func(e ui.Event) {
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
	})

	ui.Handle("<Right>", func(e ui.Event) {
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
	})

	ui.Handle("<Resize>", func(e ui.Event) {
		payload := e.Payload.(ui.Resize)
		ui.Body.Width = payload.Width
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Handle("<Escape>", func(e ui.Event) {
		if blockEvents {
			ui.Render(ui.Body)
			buffer = ""
			popup.Text = ""
			blockEvents = false
		} else {
			ui.StopLoop()
			saveData(content, &credentials)
		}
	})

	// handle key q pressing
	ui.Handle("q", func(ui.Event) {
		if !blockEvents {
			// press q to quit
			ui.StopLoop()
			saveData(content, &credentials)
		}
	})

}
