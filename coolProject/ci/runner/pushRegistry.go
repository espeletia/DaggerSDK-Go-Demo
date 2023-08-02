package runner

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

func PushRegistry(ctx context.Context) error {

	// check for Docker Hub registry credentials in host environment
	vars := []string{"DOCKER_USERNAME", "DOCKER_PASSWORD"}
	for _, v := range vars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set", v)
		}
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// set registry password as secret for Dagger pipeline
	password := client.SetSecret("password", os.Getenv("DOCKERHUB_PASSWORD"))
	username := os.Getenv("DOCKERHUB_USERNAME")

	// get reference to source code directory
	source := client.Host().Directory(".")

	golang := client.Container().From("golang:latest")

	// mount cloned repository into `golang` image
	golang = golang.WithDirectory("/src", source).WithWorkdir("/src")

	// add environment variables
	golang = golang.WithEnvVariable("CGO_ENABLED", "0")
	golang = golang.WithEnvVariable("GOOS", "linux")
	golang = golang.WithEnvVariable("GOARCH", "amd64")

	// define the application build command
	serverOutputPath := "build/start_server"
	golang = golang.WithExec([]string{"go", "build", "-o", serverOutputPath, "./cmd/main.go"})
	// get reference to the built binary in the container
	output := golang.File(serverOutputPath)

	// write the binary from the container to the host
	_, err = output.Export(ctx, serverOutputPath)
	if err != nil {
		return err
	}

	// publish image to registry
	address, err := golang.WithRegistryAuth("docker.io", username, password).
		Publish(ctx, fmt.Sprintf("%s/cool-project", username))
	if err != nil {
		return err
	}

	// print image address
	fmt.Println("Image published at:", address)
	return nil
}
