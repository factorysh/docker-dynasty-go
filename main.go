package main

import (
	"fmt"

	"github.com/docker/docker/client"
	"gitlhub.com/factorysh/docker-dynasty-go/dynasty"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	d, err := dynasty.New(cli)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
}
