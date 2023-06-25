package generator_2d

/*
Sprite

Representation of an image in a TileMap 2D Scene.

A sprite, is not linked to a specific TileMap scene.

A sprite is based on grid X and Y position.

So you don't need to do any pixels conversions at this point.

A sprite can have multiples images representations, every image after the base one is linked to a Levels.

A basic leveled sprite look like :

	basePath := "\\assets\\farms\\pumpinks\\"

	base, _ := generator_2d.LoadImage(basePath + "\\pum_2.png")

	_sprite := generator_2d.Sprite{
		Image:      base,
		Levels:     make(map[int]generator_2d.Sprite, 0),
		X:          x,
		Y:          y,
		WidthCell:  1,
		HeightCell: 1,
		Level:      level,
	}

	for i := 1; i < maxLevel+1; i++ {
		base, _ := generator_2d.LoadImage(basePath + fmt.Sprintf("\\pum_%d.png", i))

		_sprite.Levels[i] = generator_2d.Sprite{
			Image:      base,
			Levels:     nil,
			X:          x,
			Y:          y,
			WidthCell:  1,
			HeightCell: 1,
			Level:      0,
		}
	}
*/
type Sprite struct {
	Image        []byte
	Levels       map[int]Sprite
	CanBeHoverBy map[string]bool // things that can be on top of this entity
	NeedHoverBy  map[string]bool // things where that entity can be on top off
	Hovered      bool
	X            int
	Y            int
	WidthCell    int
	HeightCell   int
	Level        int
	Type         int
	UniqueID     string
	ZIndex       int
}
