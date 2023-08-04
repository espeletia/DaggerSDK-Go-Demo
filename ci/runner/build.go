package runner

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func Build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// :3
	UwU, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer UwU.Close()

	// :3
	OnO := UwU.Host().Directory("./coolProject")

	// :#
	OwO := UwU.Container().From(GOLANG_IMAGE)

	// :3
	OwO = OwO.WithDirectory("/src", OnO).WithWorkdir("/src")

	// :3
	OwO = OwO.WithEnvVariable("CGO_ENABLED", "0")
	OwO = OwO.WithEnvVariable("GOOS", "linux")
	OwO = OwO.WithEnvVariable("GOARCH", "amd64")

	// :3
	OwO = OwO.WithExec([]string{"go", "build", "-o", SERVER_OUTPUT, "./cmd/main.go"})

	// :3
	Nya := OwO.File(SERVER_OUTPUT)

	// :3
	_, err = Nya.Export(ctx, SERVER_OUTPUT)
	if err != nil {
		return err
	}
	return nil
}
