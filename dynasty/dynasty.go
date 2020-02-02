package dynasty

import (
	"bytes"
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Dynasty struct {
	client     *client.Client
	layers     map[string][]byte
	layerNames *Layers
	all        map[string]types.ImageSummary
}

func New(cli *client.Client) (*Dynasty, error) {
	d := &Dynasty{
		client:     cli,
		layers:     make(map[string][]byte),
		layerNames: NewLayers(),
		all:        make(map[string]types.ImageSummary),
	}
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All: false,
	})
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		ctx = context.Background()
		inspect, _, err := cli.ImageInspectWithRaw(ctx, image.ID)
		if err != nil {
			return nil, err
		}
		d.all[image.ID] = image
		d.layers[image.ID] = d.encode_layers(inspect.RootFS.Layers)
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

type LayerCodeTags struct {
	Layer string
	Code  []byte
	Tags  []string
}

type ByLayerCodeTags []LayerCodeTags

func (b ByLayerCodeTags) Len() int           { return len(b) }
func (b ByLayerCodeTags) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByLayerCodeTags) Less(i, j int) bool { return bytes.Compare(b[i].Code, b[j].Code) < 0 }

func (d *Dynasty) Tree() []LayerCodeTags {
	size := len(d.layers)
	r := make(ByLayerCodeTags, size)
	i := 0
	for id, code := range d.layers {
		r[i] = LayerCodeTags{
			Code:  code,
			Layer: id,
			Tags:  d.all[id].RepoTags,
		}
		i++
	}
	sort.Sort(r)
	return r
}
