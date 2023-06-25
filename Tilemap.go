package generator_2d

import "image"

/*
TileMap

A TileMap is the entry point of that 2D library.

A TileMap is build with an []byte array of an BaseImage.
A GridSize like 16 (=16x per grid cell)

Every other field are DiscordInternal one, and will be automatically populated.
*/
type TileMap struct {
	BaseImage []byte `json:"baseImage,omitempty"`
	GridSize  int

	RGBA      *image.RGBA // read-only
	MaxHeight int         // read-only
	MaxWidth  int         // read-only
	Height    int         // read-only
	Width     int         // read-only
}
