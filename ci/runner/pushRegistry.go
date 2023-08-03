package runner

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

var (
	BUILD           = "coolProject/build/"
	DOCKER_USERNAME = "DOCKER_USERNAME"
	DOCKER_PASSWORD = "DOCKER_PASSWORD"
	GOLANG_IMAGE    = "golang:latest"
	SERVER_OUTPUT   = "coolProject/build/start_server"
	SERVER_ENTRY    = "./start_server"
)

func PushRegistry(ctx context.Context) error {

	vars := []string{DOCKER_USERNAME, DOCKER_PASSWORD}
	for _, v := range vars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set", v)
		}
	}

	// Build the go binary
	Build(ctx)

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// define ENV variables
	password := client.SetSecret("password", os.Getenv(DOCKER_PASSWORD))
	username := os.Getenv(DOCKER_USERNAME)

	// get reference to the local project
	buildDir := client.Host().Directory(BUILD)
	deploy := client.Container().From(GOLANG_IMAGE).WithDirectory("/app", buildDir).WithWorkdir("/app").WithEntrypoint([]string{SERVER_ENTRY})

	// push to registry
	address, err := deploy.WithRegistryAuth("docker.io", username, password).Publish(ctx, fmt.Sprintf("%s/cool-project", username))
	if err != nil {
		return err
	}

	// print image address
	fmt.Println("Image published at:", address)
	return nil
}
