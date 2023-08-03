package runner

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func Build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory("./coolProject")

	// get `golang` image
	golang := client.Container().From(GOLANG_IMAGE)

	// mount cloned repository into `golang` image
	golang = golang.WithDirectory("/src", src).WithWorkdir("/src")

	// add environment variables
	golang = golang.WithEnvVariable("CGO_ENABLED", "0")
	golang = golang.WithEnvVariable("GOOS", "linux")
	golang = golang.WithEnvVariable("GOARCH", "amd64")

	// define the application build command
	golang = golang.WithExec([]string{"go", "build", "-o", SERVER_OUTPUT, "./cmd/main.go"})

	// get reference to the built binary in the container
	output := golang.File(SERVER_OUTPUT)

	// write the binary from the container to the host
	_, err = output.Export(ctx, SERVER_OUTPUT)
	if err != nil {
		return err
	}
	return nil
}
