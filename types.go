package main

// Task type
type Task struct {
	ID           int    `json:"task_id"`
	Name         string `json:"task_name"`
	Description  string `json:"task_description"`
	CreationDate string `json:"task_creation_date"`
	DueDate      string `json:"task_due_date"`
	Status       string `json:"task_status"`
}

// Block type
type Block struct {
	Name  string `json:"block_name"`
	Type  string `json:"block_type"`
	Tasks []Task `json:"block_tasks"`
}

// Board type
type Board struct {
	Name       string  `json:"board_name"`
	LastUpdate string  `json:"board_last_update"`
	Type       string  `json:"board_type"`
	Blocks     []Block `json:"board_blocks"`
}

// Content type
type Content struct {
	Boards []Board `json:"boards"`
	Author string  `json:"author"`
	Mail   string  `json:"mail"`
}
