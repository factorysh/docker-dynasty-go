package dynasty

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Dynasty struct {
	client     *client.Client
	layers     map[string][]byte
	layerNames *Layers
}

func New(cli *client.Client) (*Dynasty, error) {
	d := &Dynasty{
		client:     cli,
		layers:     make(map[string][]byte),
		layerNames: NewLayers(),
	}
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All: false,
	})
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		inspect, _, err := cli.ImageInspectWithRaw(ctx, image.ID)
		if err != nil {
			return nil, err
		}
		d.layers[image.ID] = d.encode_layers(inspect.RootFS.Layers)
		fmt.Println("image", image.ID, image.RepoTags,
			string(d.layers[image.ID]))
	}

	return d, nil
}

func (d *Dynasty) encode_layers(layers []string) []byte {
	e := make([]byte, len(layers)*3)
	for i, layer := range layers {
		code := d.layerNames.layer(layer)
		for j, c := range code {
			e[i*3+j] = c
		}
	}
	return e
}
