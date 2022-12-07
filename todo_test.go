package todo

import (
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	l := List{}
	taskName := "New Task"
	l.Add(taskName)
	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}
}

func TestComplete(t *testing.T) {
	l := List{}
	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}
	if l[0].Done {
		t.Errorf("New task isn't done yet")
	}
	l.Complete(1)
	if !l[0].Done {
		t.Errorf("New task should've been completed by now")
	}
}

func TestDelete(t *testing.T) {
	l := List{}
	tasks := []string{"Task 1", "Task 2", "Task 3"}
	for _, task := range tasks {
		l.Add(task)
	}
	l.Delete(2)
	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d instead", 2, len(l))
	}
	if l[1].Task != tasks[2] {
		t.Errorf("Expected %q, got %q instead", tasks[2], l[1].Task)
	}
}

func TestSaveGet(t *testing.T) {
	l1 := List{}
	l2 := List{}
	taskName := "New Task"
	l1.Add(taskName)

	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating a temp file: %s", err)
	}

	defer os.Remove(tf.Name())

	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}
	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Couldn't read file to list: %s", err)
	}
	if len(l1) != len(l2) {
		t.Errorf("L1 length must equal L2 length")
	}
	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q should match task %q", l1[0].Task, l2[0].Task)
	}
}
