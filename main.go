package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"z41dth3c0d3r/go-todo-list/src"
	"z41dth3c0d3r/go-todo-list/src/utils"
)

func createFileIfNotExists(fileName string) (*os.File, error) {
	file, err := utils.CreateFile(fileName)
	if err != nil {
		return nil, err
	}
	fmt.Println("File created!")
	utils.WaitForEnter()
	return file, nil
}

func getTodoFile() (string, bool) {

	// will store the new file or existing one
	var file string
	var existingFile bool

	// printing the main menu
	src.PrintMainMenu()
	// get the choice
	choice, err := utils.GetOption()
	if err != nil {
		fmt.Println(err)
		utils.WaitForEnter()
		return "", false
	}

	// choice from the main menu
	switch choice {
	case 1:
		location := utils.GetAnswer("Enter the existing todo lists's location: ")
		if !utils.IsFileExists(location) {
			fmt.Println("File not exists")
			utils.WaitForEnter()
		} else {
			file = location
			existingFile = true
		}
	case 2:
		newTodoListFileName := utils.GetAnswer("Enter the todo list name: ")

		if !utils.IsValidPath(newTodoListFileName) {
			fmt.Println("Invalid path")
			utils.WaitForEnter()
		} else {
			if utils.IsFileExists(newTodoListFileName) {

			alreadyExistingFileWipeOrNotQuestionLoop:
				for {
					utils.ClearScreen()
					fmt.Println("Given file already exists! do you want to overwrite?")
					fmt.Println("1. Yes")
					fmt.Println("2. No")
					fmt.Println("3. Go back")

					choice2, err := utils.GetOption()
					if err != nil {
						fmt.Println(err)
						utils.WaitForEnter()
					}

					switch choice2 {
					case 1:
						_file, err := createFileIfNotExists(newTodoListFileName)
						if err != nil {
							fmt.Println(err)
							utils.WaitForEnter()
							return "", false
						}
						defer _file.Close()
						file = newTodoListFileName
						existingFile = false
						break alreadyExistingFileWipeOrNotQuestionLoop
					case 2, 3:
						utils.WaitForEnter()
						break alreadyExistingFileWipeOrNotQuestionLoop
					default:
						fmt.Println("Invalid choice!")
						utils.WaitForEnter()
					}
				}
			} else {
				_file, err := createFileIfNotExists(newTodoListFileName)
				if err != nil {
					fmt.Println(err)
					utils.WaitForEnter()
					return "", false
				}
				defer _file.Close()
				existingFile = false
				file = newTodoListFileName
			}
		}
	case 3:
		fmt.Println("Bye bye!")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice!")
		utils.WaitForEnter()
	}

	return file, existingFile
}

func addTodo(filename string, append bool) {
	var file *os.File
	var err error

	if append {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file for append:", err)
			utils.WaitForEnter()
			return
		}

	} else {
		file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error creating file:", err)
			utils.WaitForEnter()
			return
		}
	}

	defer file.Close()

	name := utils.GetAnswer("What's the name of the todo? : ")

	newTodo := src.Todo{Name: name, IsDone: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	if err := src.SaveTodo(newTodo, file); err != nil {
		fmt.Println(err.Error())
		utils.WaitForEnter()
		return
	}

}

func readTodoFile(filename string) []src.Todo {
	utils.ClearScreen()

	todos, err := src.ReadTodo(filename)
	if err != nil {
		fmt.Println(err.Error())
		utils.WaitForEnter()
		return nil
	}

	for i := 0; i < len(todos); i++ {
		fmt.Printf("ID: %d\n", todos[i].ID)
		fmt.Printf("Name: %s\n", todos[i].Name)
		fmt.Printf("Is done: %t\n", todos[i].IsDone)
		fmt.Printf("Created at: %s\n", todos[i].CreatedAt.Format(src.DatetimeLayout))
		fmt.Printf("Updated at: %s\n\n\n", todos[i].UpdatedAt.Format(src.DatetimeLayout))
	}

	return todos
}

func removeTodo(filename string) {
	todos := readTodoFile(filename)
	id := utils.GetAnswer("Enter the removing todo's id : ")

	convertedId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		utils.WaitForEnter()
		return
	}

	if err := src.RemoveTodo(filename, convertedId, todos); err != nil {
		fmt.Println(err)
		utils.WaitForEnter()
		return
	}
}

func editTodo(filename string) {
	todos := readTodoFile(filename)
	id := utils.GetAnswer("Enter the editing todo's id : ")

	convertedId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		utils.WaitForEnter()
		return
	}

	if err := src.EditTodo(filename, convertedId, todos); err != nil {
		fmt.Println(err)
		utils.WaitForEnter()
		return
	}

}

func operateTodoList(filename string, existingFile bool) bool {
	src.PrintTodoOperationMenu()

	choice, err := utils.GetOption()
	if err != nil {
		fmt.Println(err)
		utils.WaitForEnter()
		return false
	}

	switch choice {
	case 1:
		addTodo(filename, existingFile)
	case 2:
		removeTodo(filename)
	case 3:
		completeTodo(filename)
	case 4:
		editTodo(filename)
	case 5:
		readTodoFile(filename)
		utils.WaitForEnter()
	case 6:
		return false
	default:
		fmt.Println("Invalid choice!")
		utils.WaitForEnter()
	}

	return true
}

func completeTodo(filename string) {
	todos := readTodoFile(filename)
	id := utils.GetAnswer("Enter the completing todo's id : ")

	convertedId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		utils.WaitForEnter()
		return
	}

	if err := src.MarkItAsComplete(filename, convertedId, todos); err != nil {
		fmt.Println(err)
		utils.WaitForEnter()
		return
	}
}

func main() {

	for {
		file, existingFile := getTodoFile()
		if file != "" {
			for {
				response := operateTodoList(file, existingFile)
				if !response {
					break
				}
				existingFile = true
			}
		}
	}
}
