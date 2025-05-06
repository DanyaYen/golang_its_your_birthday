package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var tasks []Task
var nextID int

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

const dataFile = "tasks.json"

func main() {
	loadedTasks, loadedNextID, err := loadTasks(dataFile)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error while loading tasks: %v\n", err)
		os.Exit(1)
	}
	tasks = loadedTasks
	nextID = loadedNextID

	if len(os.Args) < 2 {
		fmt.Println("Usage: todo <command> [arguments]")
		fmt.Println("Available commands: add, list, done, remove")
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		name := addCmd.String("name", "", "Task name (optional)")
		description := addCmd.String("desc", "", "Task description (required)")
		addCmd.Parse(args)

		if *description == "" {
			fmt.Println("Error: Task description is required.")
			addCmd.Usage()
			return
		}

		addTask(*name, *description)
		fmt.Printf("Task with ID %d added.\n", nextID-1)
		if err := saveTasks(dataFile, tasks); err != nil {
			fmt.Fprintf(os.Stderr, "Error while saving tasks: %v\n", err)
			os.Exit(1)
		}

	case "remove":
		removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
		taskID := removeCmd.Int("id", 0, "ID of task to remove")
		removeCmd.Parse(args)

		if *taskID == 0 {
			fmt.Println("Error: Task ID is required for remove command.")
			removeCmd.Usage()
			return
		}

		if removed := removeTask(*taskID); removed {
			if err := saveTasks(dataFile, tasks); err != nil {
				fmt.Fprintf(os.Stderr, "Error while saving tasks: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error: Task with ID %d not found.\n", *taskID)
		}

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		listCmd.Parse(args)
		listTasks(tasks)

	case "done":
		doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
		taskID := doneCmd.Int("id", 0, "ID of task to mark as done")
		doneCmd.Parse(args)

		if *taskID == 0 {
			fmt.Println("Error: Task ID is required for done command.")
			doneCmd.Usage()
			return
		}

		if marked := markTaskDone(*taskID); marked {
			if err := saveTasks(dataFile, tasks); err != nil {
				fmt.Fprintf(os.Stderr, "Error while saving tasks: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error: Task with ID %d not found.\n", *taskID)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: add, list, done, remove")
		return
	}
}

func addTask(name, description string) {
	newTask := Task{ID: nextID, Name: name, Description: description, Done: false}
	nextID++
	tasks = append(tasks, newTask)
}

func listTasks(taskList []Task) {
	if len(taskList) == 0 {
		fmt.Println("Task list is empty.")
		return
	}
	fmt.Println("Task List:")
	for _, task := range taskList {
		status := "[ ]"
		if task.Done {
			status = "[x]"
		}
		fmt.Printf("%d: %s %s", task.ID, status, task.Description)
		if task.Name != "" {
			fmt.Printf(" (%s)", task.Name)
		}
		fmt.Println()
	}
}

func markTaskDone(id int) bool {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			fmt.Printf("Task with ID %d marked as done.\n", id)
			return true
		}
	}
	return false
}

func removeTask(id int) bool {
	for i := range tasks {
		if tasks[i].ID == id {
			removedTaskID := tasks[i].ID
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Printf("Task with ID %d was removed.\n", removedTaskID)
			return true
		}
	}
	return false
}

func saveTasks(filename string, taskList []Task) error {
	jsonData, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		return fmt.Errorf("json encoding error: %w", err)
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("writing file %s error: %w", filename, err)
	}

	return nil
}

func loadTasks(filename string) ([]Task, int, error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, 1, nil
		}
		return nil, 0, fmt.Errorf("reading file %s error: %w", filename, err)
	}

	var loadedTasks []Task
	err = json.Unmarshal(jsonData, &loadedTasks)
	if err != nil {
		return nil, 0, fmt.Errorf("json decoding error from %s: %w", filename, err)
	}

	maxID := 0
	for _, task := range loadedTasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	calculatedNextID := maxID + 1

	return loadedTasks, calculatedNextID, nil
}