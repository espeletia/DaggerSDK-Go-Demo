package runner

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

var (
	BUILD           = "build/"
	DOCKER_USERNAME = "DOCKER_USERNAME"
	DOCKER_PASSWORD = "DOCKER_PASSWORD"
	GOLANG_IMAGE    = "golang:latest"
	SERVER_OUTPUT   = "build/start_server"
	SERVER_ENTRY    = "./start_server"
)

func PushRegistry(ctx context.Context) error {

	vars := []string{DOCKER_USERNAME, DOCKER_PASSWORD}
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

	password := client.SetSecret("password", os.Getenv(DOCKER_PASSWORD))
	username := os.Getenv(DOCKER_USERNAME)

	source := client.Host().Directory(".")

	golang := client.Container().From(GOLANG_IMAGE)

	golang = golang.WithDirectory("/src", source).WithWorkdir("/src")

	// add environment variables
	golang = golang.WithEnvVariable("CGO_ENABLED", "0")
	golang = golang.WithEnvVariable("GOOS", "linux")
	golang = golang.WithEnvVariable("GOARCH", "amd64")

	golang = golang.WithExec([]string{"go", "build", "-o", SERVER_OUTPUT, "./cmd/main.go"})
	// get reference to the built binary in the container
	output := golang.File(SERVER_OUTPUT)

	// write the binary from the container to the host
	_, err = output.Export(ctx, SERVER_OUTPUT)
	if err != nil {
		return err
	}
	buildDir := client.Host().Directory(BUILD)
	deploy := client.Container().From(GOLANG_IMAGE).WithDirectory("/app", buildDir).WithWorkdir("/app").WithEntrypoint([]string{SERVER_ENTRY})

	address, err := deploy.WithRegistryAuth("docker.io", username, password).Publish(ctx, fmt.Sprintf("%s/cool-project", username))
	if err != nil {
		return err
	}

	// print image address
	fmt.Println("Image published at:", address)
	return nil
}
