package main

import (
	"context"
	"log"

	"gitlab.com/g6834/team31/tasks/internal/application"
)

//TODO
// добавить индексы в коллекцию tasks
func main() {
	ctx := context.Background()
	log.Println("if local -> export CFG_PATH=config/config_debug.yaml")
	application.Start(ctx)
}
