package main

import (
	"fmt"
	"io"
)

type Event struct {
	Id   string
	Data string
	Type string
}

func (e Event) String() string {
	str := ""
	if e.Type != "" {
		str += fmt.Sprintf("event: %s\n", e.Type)
	}
	str += fmt.Sprintf("data: %s\n", e.Data)
	if e.Id != "" {
		str += fmt.Sprintf("id: %s\n", e.Id)
	}
	return str + "\n"
}

func (e Event) Publish(writer io.Writer) {
	writer.Write([]byte(e.String()))
}
