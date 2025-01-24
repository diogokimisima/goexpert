package main

import (
	"fmt"

	"github.com/diogokimisima/fcutils-secret/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	fmt.Println(ed)
}
