package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

type Stringer interface {
	String() string
}

func (l *List) String() string {
	formatted := ""
	for i, t := range *l {
		prefix := "  "
		if t.Done {
			prefix = "X "
		}
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, i+1, t.Task)
	}
	return formatted
}

func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, t)
}

func (l *List) Complete(index int) error {
	if 0 >= index || index > len(*l) {
		return fmt.Errorf("item %d doesn't exist", index)
	} else {
		(*l)[index-1].Done = true
		(*l)[index-1].CompletedAt = time.Now()
	}
	return nil
}

func (l *List) Delete(index int) error {
	if index <= 0 || index > len(*l) {
		return fmt.Errorf("item %d doesn't exist", index)
	}
	*l = append((*l)[:index-1], (*l)[index:]...)
	return nil
}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, l)
}
