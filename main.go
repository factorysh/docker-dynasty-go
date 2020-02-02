package main

import (
	"github.com/docker/docker/client"
	"gitlhub.com/factorysh/docker-dynasty-go/dynasty"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = dynasty.New(cli)
	if err != nil {
		panic(err)
	}
}
