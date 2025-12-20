package assets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var cache map[string]*ebiten.Image = make(map[string]*ebiten.Image)

func Get(path string) (*ebiten.Image, error) {
	if img, ok := cache[path]; ok {
		return img, nil
	}

	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}

	cache[path] = img
	return img, nil
}

func MustGet(path string) *ebiten.Image {
	img, err := Get(path)
	if err != nil {
		panic(err)
	}
	return img
}

func ClearCache() {
	for path := range cache {
		delete(cache, path)
	}
}

func LoadAssets(paths []string) error {
	for _, path := range paths {
		if _, err := Get(path); err != nil {
			return err
		}
	}
	return nil
}
