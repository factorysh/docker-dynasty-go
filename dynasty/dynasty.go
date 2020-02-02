package dynasty

import (
	"bytes"
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Dynasty explores all the layers
type Dynasty struct {
	client     *client.Client
	layers     map[string][]byte
	layerNames *Layers
	all        map[string]types.ImageSummary
}

// New Dynasty
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
		d.layers[image.ID] = d.encodeLayers(inspect.RootFS.Layers)
	}

	return d, nil
}

func (d *Dynasty) encodeLayers(layers []string) []byte {
	e := make([]byte, len(layers)*4)
	for i, layer := range layers {
		code := d.layerNames.layer(layer)
		for j, c := range code {
			e[i*4+j] = c
		}
		e[i*4+3] = byte(' ')
	}
	return e
}

// LayerCodeTags Layer, code and tags
type LayerCodeTags struct {
	Layer string
	Code  []byte
	Tags  []string
}

// ByLayerCodeTags sorts LayerCodeTags
type ByLayerCodeTags []LayerCodeTags

func (b ByLayerCodeTags) Len() int           { return len(b) }
func (b ByLayerCodeTags) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByLayerCodeTags) Less(i, j int) bool { return bytes.Compare(b[i].Code, b[j].Code) < 0 }

// Tree return sorted layers
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

// Ancestor returns ancestors of the image
func (d *Dynasty) Ancestor(name string) ([]LayerCodeTags, error) {
	return d.beforeAfter(name, func(a, b []byte) bool {
		return startswith(a, b)
	})
}

// Descendant returns descendant of the image
func (d *Dynasty) Descendant(name string) ([]LayerCodeTags, error) {
	return d.beforeAfter(name, func(a, b []byte) bool {
		return startswith(b, a)
	})
}

func (d *Dynasty) beforeAfter(name string, cmp func(a, b []byte) bool) ([]LayerCodeTags, error) {
	ctx := context.Background()
	inspect, _, err := d.client.ImageInspectWithRaw(ctx, name)
	if err != nil {
		return nil, err
	}
	id := d.layers[inspect.ID]
	r := make(ByLayerCodeTags, 0)
	for k, layer := range d.layers {
		is, ok := d.all[k]
		var tags []string
		if ok {
			tags = is.RepoTags
		} else {
			tags = []string{}
		}
		if cmp(layer, id) {
			r = append(r, LayerCodeTags{
				Code:  layer,
				Layer: k,
				Tags:  tags,
			})
		}
	}
	sort.Sort(r)
	return r, nil
}

// Does haystack starts with needle ?
func startswith(needle, haystack []byte) bool {
	if len(needle) > len(haystack) {
		return false
	}

	for i, b := range needle {
		if haystack[i] != b {
			return false
		}
	}
	return true
}
