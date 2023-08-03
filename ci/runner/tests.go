package runner

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func DoTests(ctx context.Context) error {
	fmt.Println("Testing with Dagger")

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	src := client.Host().Directory(".")

	golang := client.Container().From("golang:latest")

	golang = golang.WithDirectory("/src", src).WithWorkdir("/src")

	golang = golang.WithExec([]string{"go", "test", "./..."})
	_, err = golang.Sync(ctx)
	if err != nil {
		return err
	}
	return nil
}
