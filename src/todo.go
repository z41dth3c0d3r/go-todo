package src

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"z41dth3c0d3r/go-todo-list/src/utils"
)

var DatetimeLayout = "2006-01-02 15:04:05"

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func EditTodo(filename string, editingTodoId int, todos []Todo) error {
	for i := 0; i < len(todos); i++ {
		if todos[i].ID == editingTodoId {
			todos[i].Name = utils.GetAnswer("Enter the new name: ")
			todos[i].UpdatedAt = time.Now()
			break
		}
	}

	if err := WriteTodosToJson(filename, todos); err != nil {
		return err
	}

	return nil
}

func MarkItAsComplete(filename string, markingTodoId int, todos []Todo) error {
	for i := 0; i < len(todos); i++ {
		if todos[i].ID == markingTodoId {
			todos[i].IsDone = true
			todos[i].UpdatedAt = time.Now()
			break
		}
	}

	if err := WriteTodosToJson(filename, todos); err != nil {
		return err
	}

	return nil
}

func WriteTodosToJson(filename string, todos any) error {
	data, err := json.MarshalIndent(todos, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling todos: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func RemoveTodo(filename string, removingTodoId int, todos []Todo) error {
	newTodo := make([]Todo, 0, len(todos))

	for i := 0; i < len(todos); i++ {

		if todos[i].ID != removingTodoId {
			newTodo = append(newTodo, todos[i])
		}
	}

	if err := WriteTodosToJson(filename, newTodo); err != nil {
		return fmt.Errorf("error reading existing todos: %v", err)
	}

	return nil
}

func SaveTodo(todo Todo, file *os.File) error {
	// get exsiting todos
	todos, err := ReadTodo(file.Name())
	if err != nil {
		return fmt.Errorf("error reading existing todos: %v", err)
	}

	// get the last ID
	todo.ID = 1
	if len(todos) > 0 {
		maxID := todos[0].ID
		for _, t := range todos {
			if t.ID > maxID {
				maxID = t.ID
			}
		}
		todo.ID = maxID + 1
	}

	todos = append(todos, todo)

	if err := WriteTodosToJson(file.Name(), todos); err != nil {
		return fmt.Errorf("error reading existing todos: %v", err)
	}

	return nil
}

func ReadTodo(filename string) ([]Todo, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Todo{}, nil
		}
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if len(data) == 0 {
		return []Todo{}, nil
	}

	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, fmt.Errorf("error unmarshaling todos: %v", err)
	}

	return todos, nil
}
