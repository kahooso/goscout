// task-01 (A0.1): тип Task с методами и конструктором.
// Подробный разбор: knowledge-base/topics/go/structs.md
package main

import "fmt"

func main() {
	tasks := []*Task{
		NewTask(0, "Купить хлеб"),
		NewTask(1, "Написать код"),
		NewTask(2, "Лечь спать"),
	}
	tasks[0].Complete()
	PrintAll(tasks)
}

type Task struct {
	ID     int
	Title  string
	IsDone bool
}

// NewTask возвращает *Task — методы с pointer receiver требуют адреса.
func NewTask(id int, title string) *Task {
	return &Task{ID: id, Title: title}
}

// Complete: pointer receiver — иначе изменения IsDone теряются на копии.
func (t *Task) Complete() {
	t.IsDone = true
}

// String делает Task реализующим fmt.Stringer — fmt.Println зовёт его автоматически.
func (t *Task) String() string {
	sgn := " "
	if t.IsDone {
		sgn = "x"
	}
	return fmt.Sprintf("[%s] %d. %s", sgn, t.ID, t.Title)
}

func PrintAll(tasks []*Task) {
	for _, t := range tasks {
		fmt.Println(t)
	}
}
