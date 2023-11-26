// TODO refactor to put writing and saving in own file? maybe myio.go or sum
// TODO add input validation where appropriate

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"time"
	// "unicode"
)

func main() {
	// todoList := []Todo{}
	// readTodos(&todoList)
	db := initDb()
	defer db.Close()

	var choice UserChoice
	for {
		printMenu()
		fmt.Scan(&choice)

		// if !unicode.IsDigit(rune(choice)) {
		// 	os.Exit(0)
		// }

		switch choice {
		case ChoiceViewTodos:
			viewTodos(db)

		case ChoiceAddTodo:
			addTodo(db)
			viewTodos(db)

		case ChoiceUpdateTodo:
			updateTodo(db, getDesiredTodoIdFromUser())
			viewTodos(db)

		case ChoiceDeleteTodo:
			deleteOneTodo(db, getDesiredTodoIdFromUser())
			viewTodos(db)

		case ChoiceDeleteAllTodos:
			deleteAllTodos(db)
			viewTodos(db)

		default:
			// fmt.Println("Invalid option, please try again.")
			// continue
			os.Exit(0)
		}

		// switch choice {
		// case ChoiceAddTodo, ChoiceUpdateTodo, ChoiceDeleteTodo, ChoiceDeleteAllTodos:
		// 	saveTodos(todoList)
		// }
	}
}

// checkErr logs the error and quits if it is not nil.
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func initDb() *sql.DB {
	migrationStr := "CREATE TABLE IF NOT EXISTS `TodoList` (\n" +
		"`todo_id`     INTEGER  NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,\n" +
		"`name`        TEXT     NOT NULL UNIQUE,\n" +
		"`kind`        INTEGER  NOT NULL DEFAULT 0,\n" +
		"`state`       INTEGER  NOT NULL DEFAULT 0\n" +
		");"
	db, err := sql.Open("sqlite3", "./todolist.sqlite")
	checkErr(err)
	db.Exec(migrationStr)
	// defer db.Close()
	// stmt, err := db.Prepare(migrationStr)
	// stmt.Exec()
	return db
}

// printMenu outputs the menu text to the user.
func printMenu() {
	menuText := `What would you like to do?
	1. View list of todos
	2. Add to the list
	3. Update an item in the list
	4. Delete an item from the list
	5. Delete all items from the list

Any other input quits the program. Your choice?: `
	fmt.Print(menuText)
}

// getDesiredTodoIdFromUser will prompt for the ID of todo user wants to change
// TODO validate input to make sure id user gives is within range
// returns the id that the user wanted
func getDesiredTodoIdFromUser() uint {
	fmt.Print("Id to modify? ")
	var userGivenId uint
	fmt.Scan(&userGivenId)
	return userGivenId
}

// viewTodos will display all todos in todoList.
func viewTodos(db *sql.DB) {
	fmt.Println()
	fmt.Printf("%5s%40s\t%10s\t%10s\n", "id", "name", "kind", "state")
	fmt.Printf("%5s%40s\t%10s\t%10s\n", "==", "====", "====", "=====")

	rows, err := db.Query("SELECT * FROM TodoList")
	checkErr(err)
	for rows.Next() {
		// load current todo
		var currentTodo Todo
		err = rows.Scan(
			&currentTodo.todo_id,
			&currentTodo.name,
			&currentTodo.kind,
			&currentTodo.state)
		checkErr(err)

		// print current todo
		fmt.Printf("%4d)%40s\t%10s\t%10s\n",
			currentTodo.todo_id,
			currentTodo.name,
			currentTodo.kind,
			currentTodo.state)
	}

	fmt.Println()
}

// addTodo will add one todo item to todoList
func addTodo(db *sql.DB) {
	var newTodoItem Todo
	var sc = bufio.NewScanner(os.Stdin)
	// we actually don't need to do anything with an id because db will give us
	// that automagically, so set it to dummy value
	newTodoItem.todo_id = 0

	// get the name of the todo
	fmt.Print("Todo name?\n> ")
	sc.Scan()
	newTodoItem.name = sc.Text()

	// get the type of todo
	fmt.Printf(
		"Todo type? (%d for project, %d for homework, %d for reading)\n> ",
		CategoryProject, CategoryHomework, CategoryReading)
	//TODO input validation
	fmt.Scan(&newTodoItem.kind)

	// set the state of the todo
	newTodoItem.state = StateTodo

	// put stuff in db
	stmt, err := db.Prepare(`INSERT INTO TodoList (name, kind, state)` +
		`VALUES (?, ?, ?);`)
	checkErr(err)
	stmt.Exec(newTodoItem.name, newTodoItem.kind, newTodoItem.state)
}

// updateTodo will modify the entry in todoList with identifier id.
func updateTodo(db *sql.DB, todo_id uint) {
	sqlQueryTempl := `UPDATE TodoList SET %s = ? where todo_id = ?;`
	var finalQuery string
	var newValToSet int // trick so

	// switch on choice forever until valid input
	for {
		var choice int
		fmt.Print("Update how? (1 for type, 2 for state) ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			finalQuery = fmt.Sprintf(sqlQueryTempl, "kind")
			// fmt.Sscanf(&finalQuery, sqlQueryTempl, "kind")
			fmt.Printf(
				"New type? (%d for project, %d for homework, %d for reading) ",
				CategoryProject, CategoryHomework, CategoryReading)
		case 2:
			finalQuery = fmt.Sprintf(sqlQueryTempl, "state")
			// fmt.Sscanf(&finalQuery, sqlQueryTempl, "state")
			fmt.Printf(
				"New state? (%d for todo, %d for in progress, %d for done) ",
				StateTodo, StateInProgress, StateDone)
		}

		// save here out of the switch because we have same range of values from
		// both enums right now
		fmt.Scan(&newValToSet)

		if !(0 <= newValToSet && newValToSet <= 2) {
			fmt.Println("Please enter a value between 0 and 2")
		} else {
			break
		}
	}

	fmt.Println(finalQuery)
	stmt, err := db.Prepare(finalQuery)
	checkErr(err)
	stmt.Exec(newValToSet, todo_id)
}

// updateTodo will delete the entry in todoList with identifier id.
func deleteOneTodo(db *sql.DB, todo_id uint) {
	sqlQueryTempl := `DELETE FROM TodoList WHERE todo_id = ?`
	stmt, err := db.Prepare(sqlQueryTempl)
	checkErr(err)
	stmt.Exec(todo_id)
}

// deleteAllTodos removes all entries in todoList (in-place).
// The function will, however, ask for user confirmation and wait 2 seconds
// before allowing a response, so as to prevent accidental deletion.
func deleteAllTodos(db *sql.DB) {
	fmt.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)

	var userResp string
	fmt.Print("Are you sure? [y/N] ")
	fmt.Scan(&userResp)

	if strings.ToUpper(userResp) == "Y" {
		stmt, err := db.Prepare(`DELETE FROM TodoList`)
		checkErr(err)
		stmt.Exec()
	} else {
		fmt.Println("Aborting...")
	}
}
