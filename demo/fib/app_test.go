package fib

import (
	"context"
	"log"
	"os"
	"os/signal"
	"testing"
)

// doc https://opentelemetry.io/docs/instrumentation/go/getting-started/

func TestFibonacci(t *testing.T) {
	l := log.New(os.Stdout, "", 0)

	sigCh := make(chan os.Signal, 1)
	// Ctrl + C
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := NewApp(os.Stdin, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Printf("\ngoodbye")
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
