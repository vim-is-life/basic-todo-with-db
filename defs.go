package main

const (
	TODO_FILE_FORMAT = "csv"
	TODO_FILE        = "todos." + TODO_FILE_FORMAT
)

// Todo represents a todo item in the list
type Todo struct {
	todo_id uint
	name    string
	kind    Category
	state   ProgressState
	// desc string
}

// UserChoice represents user's choice for what to do with todo list
type UserChoice int

const (
	ChoiceViewTodos UserChoice = iota + 1
	ChoiceAddTodo
	ChoiceUpdateTodo
	ChoiceDeleteTodo
	ChoiceDeleteAllTodos
)

// Category represents the category of a todo item.
// (ie what type of todo it is)
type Category int

const (
	CategoryProject Category = iota
	CategoryHomework
	CategoryReading
)

// Method to get a string representation of a category
func (c Category) String() string {
	return [...]string{"Project", "Homework", "Reading"}[c]
}

// ProgressState represents the state of a todo item.
// (ie whether it's done, being worked on, or not started)
type ProgressState int

const (
	StateTodo ProgressState = iota
	StateInProgress
	StateDone
)

// Method to get a string representation of a ProgressState
func (p ProgressState) String() string {
	return [...]string{"Todo", "InProgress", "Done"}[p]
}
