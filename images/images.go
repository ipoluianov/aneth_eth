package images

import (
	"errors"
	"sync"
)

type Img struct {
	Data    []byte
	AltText string
}

type Images struct {
	mtx   sync.Mutex
	items map[string]*Img
}

var Instance *Images

func init() {
	Instance = NewImages()
}

func NewImages() *Images {
	var c Images
	c.items = make(map[string]*Img)
	return &c
}

func (c *Images) Get(code string) (*Img, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if item, ok := c.items[code]; ok {
		return item, nil
	}
	return nil, errors.New("image not found")
}

func (c *Images) Set(code string, data []byte, altText string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var item Img
	item.Data = data
	item.AltText = altText
	c.items[code] = &item
}
