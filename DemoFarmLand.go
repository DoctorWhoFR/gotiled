package generator_2d

import (
	"fmt"
	"image"
	"os"
	"strconv"
)

/*
CreateDemoFarmLand
  - Base : base.png
  - GridSize : 16px
  - Theme : basic

Create a basic FarmLand and returning a working Map.

TODO: Do we need to put this on bot part ? or keep it on 2D lib ?
*/
func CreateDemoFarmLand() (Map, error) {
	base, imagec, err := LoadImage("\\assets\\base.png")

	if err != nil {
		fmt.Println(err)
		return Map{}, err
	}

	baseX := imagec.Bounds().Dx()
	baseY := imagec.Bounds().Dy()

	maxX := baseX / 16
	maxY := baseY / 16

	baseLand := TileMap{
		BaseImage: base,
		GridSize:  16,
		MaxWidth:  maxX,
		MaxHeight: maxY,
		Height:    imagec.Bounds().Dy(),
		Width:     imagec.Bounds().Dx(),
		RGBA:      image.NewRGBA(imagec.Bounds()),
	}

	isDebug, _ := strconv.ParseBool(os.Getenv("LAND_DEBUGGING"))

	return Map{
		TileMap: baseLand,
		Sprites: map[string]Sprite{},
		Debug:   isDebug,
	}, nil
}
