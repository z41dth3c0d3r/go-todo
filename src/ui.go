package src

import (
	"fmt"
	"z41dth3c0d3r/go-todo-list/src/utils"
)

func PrintMainMenu() {
	utils.ClearScreen()
	fmt.Println("1. Open existing todo list.")
	fmt.Println("2. Create todo list")
	fmt.Println("3. Exit")
}

func PrintTodoOperationMenu() {
	utils.ClearScreen()
	fmt.Println("1. Add todo")
	fmt.Println("2. Remove todo")
	fmt.Println("3. Complete todo")
	fmt.Println("4. Edit todo")
	fmt.Println("5. Read todo list")
	fmt.Println("6. Go back")
}
