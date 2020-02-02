package dynasty

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Dynasty struct {
	client *client.Client
	layers map[string]interface{}
}

func New(cli *client.Client) (*Dynasty, error) {
	d := &Dynasty{
		client: cli,
		layers: make(map[string]interface{}),
	}
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		d.layers[image.ID] = nil
	}

	return d, nil
}
