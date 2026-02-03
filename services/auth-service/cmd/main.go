package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrxacker/inventory-management-system/services/auth-service/internal/app"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := app.StartApp(ctx)
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}
