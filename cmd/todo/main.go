package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"todo"
)

var todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "ToDo tool. Developed for the pragmatic bookshelf.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		flag.PrintDefaults()
	}

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	//* flags:
	add := flag.Bool("add", false, "Add a task to the list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be marked as completed")
	delete := flag.Int("delete", 0, "Delete an item off the list")
	verbose := flag.Bool("verbose", false, "Toggle in verbosity to display extra details")
	all := flag.Bool("all", false, "List all tasks, including the completed ones")
	flag.Parse()

	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch {
	case *list && *all && *verbose:
		for index, item := range *l {
			if item.Done {
				msg := fmt.Sprintf("X %d: %s [Created at: %v | Completed at: %v]", index+1, item.Task, item.CreatedAt, item.CompletedAt)
				fmt.Println(msg)
			} else {
				msg := fmt.Sprintf("  %d: %s [Created at: %v]", index+1, item.Task, item.CreatedAt)
				fmt.Println(msg)
			}
		}

	case *list && *verbose:
		for index, item := range *l {
			if !item.Done {
				msg := fmt.Sprintf("  %d: %s [Created at: %v]", index+1, item.Task, item.CreatedAt)
				fmt.Println(msg)
			}
		}
	case *list && *all:
		// fmt.Print(l)
		for index, item := range *l {
			if item.Done {
				msg := fmt.Sprintf("X %d: %s", index+1, item.Task)
				fmt.Println(msg)
			} else {
				msg := fmt.Sprintf("  %d: %s", index+1, item.Task)
				fmt.Println(msg)
			}
		}

	case *list:
		for index, item := range *l {
			if !item.Done {
				msg := fmt.Sprintf("  %d: %s", index+1, item.Task)
				fmt.Println(msg)
			}
		}

	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	// case *task != "":
	// 	l.Add(*task)
	// 	if err := l.Save(todoFileName); err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 		os.Exit(1)
	// 	}

	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		// fmt.Fprintln(os.Stderr, "Invalid option")
		flag.Usage()
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(scanner.Text()) == 0 {
		return "", fmt.Errorf("a task can't be blank")
	}
	return scanner.Text(), nil
}
