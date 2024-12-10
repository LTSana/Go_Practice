package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type task struct {
	id          int
	description string
	completed   bool
	created     int
}

func main() {

	// Check if the user has provided an action
	if len(os.Args) < 2 {
		fmt.Println("Please specify an action")
		return
	}

	// Check the action
	switch os.Args[1] {
	case "add":

		// Check if the user provided a description
		if len(os.Args) < 3 {
			fmt.Println("Please specify a task")
			return
		}

		// Pass the description to the functions parameter
		addTask(os.Args[2])

	case "list":

		// Check if the users action is valid
		if len(os.Args) == 3 {
			switch os.Args[2] {
			case "-a":
				listTasks(true)
			case "--all":
				listTasks(true)
			default:
				listTasks(false)
			}
		} else {
			listTasks(false)
		}

	case "complete":

		// Check if the user provided the task ID
		if len(os.Args) < 3 {
			fmt.Println("Please specify a task ID")
			return
		}

		// Check if the user provided a valid int
		taskId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}

		// Set the task as completed
		completeTask(taskId)

	case "delete":

		// Check if the user provided the task ID
		if len(os.Args) < 3 {
			fmt.Println("Please specify a task ID")
			return
		}

		// Check if the user provided a valid int
		taskId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}

		// Delete the task
		deleteTask(taskId)
	default:
		fmt.Println("Invalid command. Please use add, list or delete")
	}
}

func addTask(s string) error {

	// Open the tasks CSV file to read only
	file2, err := os.OpenFile("tasks.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file2.Close()

	// Read the csv data from the file
	r := csv.NewReader(file2)

	// Read all the records in the CSV file
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Create or open existing CSV file
	file, err := os.Create("tasks.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Prepare to write to the file
	w := csv.NewWriter(file)
	defer w.Flush()
	data := [][]string{}

	// Task ID
	taskId := 0

	// Append all the tasks to the new data
	data = append(data, []string{"id", "Task", "completed", "created"})
	if len(records) > 1 {
		for _, task := range records[1:] {
			data = append(data, []string{task[0], task[1], task[2], task[3]})
		}

		// Get the ID of the last task
		v, err := strconv.Atoi(records[len(records)-1][0])
		if err != nil {
			return err
		}
		taskId = v
	}

	// Append the new task to the list
	data = append(data, []string{strconv.Itoa(taskId + 1), s, "false", time.Now().Format("2006/01/02 15:04:05.00000")})

	// Write the data to the CSV file
	for _, value := range data {
		if err := w.Write(value); err != nil {
			return err
		}
	}

	return nil
}

func listTasks(show bool) error {

	// Open the file if it exists else creates the file
	file, err := os.OpenFile("tasks.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file
	r := csv.NewReader(file)

	// Get all the records from the file
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Create a new tabwriter.Writer instance.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// Write some data to the Writer.
	var wFormatString string
	if show {
		wFormatString = "ID\tTask\tCreated\tDone"
	} else {
		wFormatString = "ID\tTask\tCreated"
	}
	fmt.Fprintln(w, wFormatString)

	// Print out the tasks one by one
	for _, task := range records[1:] {

		// Get the time diff
		t1, err := time.Parse("2006/01/02 15:04:05.00000", time.Now().Format("2006/01/02 15:04:05.00000"))
		t2, err := time.Parse("2006/01/02 15:04:05.00000", task[3])
		if err != nil {
			fmt.Println(err)
			return err
		}
		taskTimeDiff := t1.Sub(t2)

		formattedString := fmt.Sprintf("")
		if taskTimeDiff.Seconds() > 1.0 && taskTimeDiff.Minutes() < 2.0 && taskTimeDiff.Hours() < 1.0 {
			formattedString = fmt.Sprintf("%.0f seconds ago", taskTimeDiff.Seconds())
		} else if taskTimeDiff.Minutes() > 1.0 && taskTimeDiff.Hours() < 2.0 {
			formattedString = fmt.Sprintf("%.0f minutes ago", taskTimeDiff.Minutes())
		} else if taskTimeDiff.Hours() > 1.0 {
			formattedString = fmt.Sprintf("%.0f hours ago", taskTimeDiff.Hours())
		} else {
			formattedString = fmt.Sprintf("Now")
		}

		// Check if the task is completed
		taskCompleted, err := strconv.ParseBool(task[2])
		if err != nil {
			fmt.Println(err)
			return err
		}
		var taskCompletedEmoji string
		if taskCompleted {
			taskCompletedEmoji = "✔️"
		} else {
			taskCompletedEmoji = "❌"
		}

		// Check if the user wants to see all data
		if show {
			fmt.Fprint(w, task[0]+"\t"+task[1]+"\t"+formattedString+"\t"+taskCompletedEmoji+"\n")
		} else {
			fmt.Fprint(w, task[0]+"\t"+task[1]+"\t"+formattedString+"\n")
		}
	}

	// Flush the Writer to ensure all data is written to the output.
	w.Flush()

	return nil
}

func completeTask(id int) error {

	// Open the tasks CSV file to read only
	file2, err := os.OpenFile("tasks.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file2.Close()

	// Read the csv data from the file
	r := csv.NewReader(file2)

	// Read all the records in the CSV file
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Create or open existing CSV file
	file, err := os.Create("tasks.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Prepare to write to the file
	w := csv.NewWriter(file)
	defer w.Flush()
	data := [][]string{}

	// Append all the tasks to the new data
	data = append(data, []string{"id", "Task", "completed", "created"})
	if len(records) > 1 {
		for _, task := range records[1:] {

			// Check if the ID of the task matches the task the user wants to set as complete
			taskID, err := strconv.Atoi(task[0])
			if err != nil {
				return err
			}
			if taskID == id {
				data = append(data, []string{task[0], task[1], "true", task[3]})
			} else {
				data = append(data, []string{task[0], task[1], task[2], task[3]})
			}
		}
	}

	// Write the data to the CSV file
	for _, value := range data {
		if err := w.Write(value); err != nil {
			return err
		}
	}

	return nil
}

func deleteTask(id int) error {

	// Open the tasks CSV file to read only
	file2, err := os.OpenFile("tasks.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file2.Close()

	// Read the csv data from the file
	r := csv.NewReader(file2)

	// Read all the records in the CSV file
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Create or open existing CSV file
	file, err := os.Create("tasks.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Prepare to write to the file
	w := csv.NewWriter(file)
	defer w.Flush()
	data := [][]string{}

	// Append all the tasks to the new data
	data = append(data, []string{"id", "Task", "completed", "created"})
	if len(records) > 1 {
		for _, task := range records[1:] {

			// Check if the ID of the task matches the task the user wants to set as complete
			taskID, err := strconv.Atoi(task[0])
			if err != nil {
				return err
			}
			if taskID != id {
				data = append(data, []string{task[0], task[1], task[2], task[3]})
			}
		}
	}

	// Write the data to the CSV file
	for _, value := range data {
		if err := w.Write(value); err != nil {
			return err
		}
	}

	return nil
}
