package main

import (
	"fmt"

	"github.com/rakyll/ticktock"
	"github.com/rakyll/ticktock/t"
)

type PrintJob struct {
	Msg string
}

func (j *PrintJob) Run() error {
	fmt.Println(j.Msg)
	return nil
}
func main() {
	// Prints "Hello world" once in every seconds
	err := ticktock.Schedule(
		"print-hello",
		&PrintJob{Msg: "Hello world"},
		&t.When{Every: t.Every(5).Seconds()})
	if err != nil {
		fmt.Println(err)
	}
	ticktock.Start()
}
