/*
Package generator_2d

2D Tile Renderer

This is a simple library used to generated and manage tile set map and render them in PNG
With this library, you can use X, Y position based on the tile grid
*/
package generator_2d

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
)

/*
Map

Representation of a final map with entities, base tile map, and other functionality.

A map is based on
  - a TileMap (base image for the 2D scene)
  - Lists of Sprite (representation of an image in the 2D scene)
  - Message for in-build notification system
  - Debug is a boolean variable used to display the debugging mod
*/
type Map struct {
	TileMap                   // Representation of an TileMap currently the scene rendered to the user
	Sprites map[string]Sprite // Sprites representations representation of an image in the 2D scene
	Message string            // For in-build image notification system
	Debug   bool
}

const (
	CheckCanBePosedAlreadyHere            = "ALREADY_ENT_HERE_NO_HOVER_BY"
	CheckCanBePosedAlreadyHereBadHovered  = "BAD_HOVERED_ENTITY"
	CheckCanBePosedAlreadyHereAlreadyUsed = "HOVERED_ENTITY_ALREADY_USED"
	CheckCanBePosedNeedHover              = "NO_ENTITY_TO_PUT"
)

func (m *Map) CheckCanBePosed(uuid string, ent Sprite) error {
	if ent.WidthCell > 1 {
		if ent.X+ent.WidthCell >= m.MaxWidth {
			log.Println("error above max witdh", m.MaxWidth, ent.X)
			return errors.New("above maximum width")
		}
	} else {
		if ent.X >= m.MaxWidth {
			log.Println("error above max witdh", m.MaxWidth, ent.X)
			return errors.New("above maximum width")
		}
	}

	if ent.HeightCell > 1 {
		if ent.Y+ent.HeightCell >= m.MaxHeight {
			log.Println("error above max height")
			return errors.New("above maximum height")
		}
	} else {
		if ent.Y >= m.MaxHeight {
			log.Println("error above max height")
			return errors.New("above maximum height")
		}
	}

	// check if an entity is already in that position
	for key, alreadyHereSprite := range m.Sprites {
		// if entity does need to be hover something else
		if len(ent.NeedHoverBy) == 0 {
			if ent.X == alreadyHereSprite.X && ent.Y == alreadyHereSprite.Y && key != uuid {
				return errors.New(CheckCanBePosedAlreadyHere)
			}

			return nil
		} else {
			if ent.X == alreadyHereSprite.X && ent.Y == alreadyHereSprite.Y && key != uuid {
				if !alreadyHereSprite.CanBeHoverBy[ent.UniqueID] {
					return errors.New(CheckCanBePosedAlreadyHereBadHovered)
				}

				if alreadyHereSprite.Hovered {
					return errors.New(CheckCanBePosedAlreadyHereAlreadyUsed)
				}
				alreadyHereSprite.Hovered = true

				return nil
			}
		}
	}

	if len(ent.NeedHoverBy) != 0 {
		return errors.New(CheckCanBePosedNeedHover)
	}

	return nil
}

/*
AddUpdateEnt
Add an entity to a final Map, add it into a Sprite lists
When you are going to use the RenderScene function, all Sprites added wil be generated and added to the scene according to there X,Y position.

Always use AddUpdateEnt to add or update entity, even more if you change the X,Y position.
Every conditional positional check is made in this function and in no other step.

If you don't use it, you will maybe have some clapping entity, if not managed well on your side.
*/
func (m *Map) AddUpdateEnt(uuid string, ent Sprite) error {
	err := m.CheckCanBePosed(uuid, ent)
	if err != nil {
		return err
	}

	m.Sprites[uuid] = ent

	return nil
}

func (m *Map) DeleteEnt(uuid string) error {
	delete(m.Sprites, uuid)
	return nil
}

/*
DrawImage
In build option to draw sprite into final img scene

Take a Sprite in entry and a *gg.Context
*/
func (m *Map) DrawImage(entity Sprite, ctx *gg.Context) {
	entReader := bytes.NewReader(entity.Image)
	entityImage, _ := png.Decode(entReader)

	x := entity.X
	y := entity.Y

	log.Println(x, y)

	xpos, ypos := m.GridToPixel(x, y)

	ctx.Push()
	// Draw the square on the image
	ctx.DrawImage(entityImage, xpos, ypos)
}

/*
GridToPixel
Simple, in-build position, to transform X,Y cells based position, into there pixel equivalent.

  - x : X grid position
  - x : Y grid position
*/
func (m *Map) GridToPixel(x, y int) (int, int) {
	return x * m.GridSize, y * m.GridSize
}

/*
SendNotification
// TODO Simple, in-build image notification system
*/
func (m *Map) SendNotification(message string) {
	m.Message = message

	// TODO what can we do with notification system, does with need to make a more complexe one for now
	// or keep using basic fixed size image ?
}

/*
RenderScene

Main function used to get the final Image png representation of a user Map.

When you are using RenderScene every Sprite's that is present on the Map will be generated.

Linked to their equivalent Sprite.Level if having one.

Or displayed as their base image Sprite.Image if having no level.

You will get a string path, where the generated file belong.

/!\ *The file will NOT BE automatically delete, so you need to handle this part.* /!\
*/
func (m *Map) RenderScene() string {
	reader := bytes.NewReader(m.BaseImage)
	base, _ := png.Decode(reader)
	rgba := image.NewRGBA(base.Bounds())

	m.RGBA = rgba

	dc := gg.NewContextForImage(rgba)
	dc.SetRGB(255, 255, 255)
	dc.Clear()
	dc.DrawImage(base, 0, 0)

	baseX := base.Bounds().Dx()
	baseY := base.Bounds().Dy()

	maxX := baseX / m.GridSize
	maxY := baseY / m.GridSize

	log.Println(maxX, maxY)

	m.renderSprites(dc)

	m.renderMessage(dc)

	m.renderDebug(maxY, dc, maxX)

	largerDC := m.scaleUP(baseX, baseY, dc)

	// Open the PNG image file
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
	}

	_uuid, _ := uuid.NewUUID()
	log.Println(currentDir)
	_filename := fmt.Sprintf(currentDir+"\\tmp_%s.png", _uuid.String())

	err = largerDC.SavePNG(_filename)

	if err != nil {
		return ""
	}

	return _filename
}

/*
scaleUP

Simple DiscordInternal function used to scale UP the final scene by 2.
*/
func (m *Map) scaleUP(baseX int, baseY int, dc *gg.Context) *gg.Context {
	// Create a larger context for scaling up
	largerDC := gg.NewContext(int(float64(baseX)*2), int(float64(baseY)*2))

	// Scale the drawing to the larger context
	largerDC.Scale(2, 2)
	largerDC.DrawImage(dc.Image(), 0, 0)
	return largerDC
}

/*
renderSprites

Simple DiscordInternal function used to render all Sprite of a Map.
*/
func (m *Map) renderSprites(dc *gg.Context) {

	// create a zindex map of sprite
	var ZSpritesMap = make(map[int][]Sprite)

	for _, sprite := range m.Sprites {
		zsprite, finded := ZSpritesMap[sprite.ZIndex]

		if !finded {
			zsprite = []Sprite{}
		}

		zsprite = append(zsprite, sprite)

		ZSpritesMap[sprite.ZIndex] = zsprite
	}

	for i := 0; i < len(ZSpritesMap); i++ {
		zsprites := ZSpritesMap[i+1]
		for _, sprite := range zsprites {
			if sprite.Level > 1 {
				_lvl := sprite.Levels[sprite.Level]
				_lvl.X = sprite.X
				_lvl.Y = sprite.Y
				log.Println(_lvl, "level")
				m.DrawImage(_lvl, dc)
			} else {
				m.DrawImage(sprite, dc)
			}
		}
	}

}

/*
renderMessage

Simple DiscordInternal function to render in-build notification message
*/
func (m *Map) renderMessage(dc *gg.Context) {
	if m.Message != "" {
		base3, _ := LoadImage("\\assets\\gui.png")

		gui := Sprite{
			Image:      base3,
			X:          0,
			Y:          5,
			WidthCell:  0,
			HeightCell: 24,
		}

		m.DrawImage(gui, dc)

		xpos, ypos := m.GridToPixel(2, gui.Y+2)

		// # TODO maybe add a new font for gui purpose
		// ${line}
		// ${fullPath}
		// <!-- epic:"lovely" -->

		//// Open the PNG image file
		//currentDir, err := os.Getwd()
		//if err != nil {
		//	fmt.Println("Failed to get current directory:", err)
		//}
		//if err := dc.LoadFontFace(currentDir+"\\assets\\monogram.ttf", 20); err != nil {
		//	panic(err)
		//}

		dc.DrawString(m.Message, float64(xpos), float64(ypos))

		m.Message = ""
	}
}

/*
renderDebug

Simple DiscordInternal function to the final scene with a debugging grids.
*/
func (m *Map) renderDebug(maxY int, dc *gg.Context, maxX int) {
	if m.Debug {
		for y := 0; y < maxY; y++ {
			fmt.Println(10 + float64(y*16))

			dc.SetRGB(255, 0, 21)
			dc.DrawString(strconv.Itoa(y), 3, 10+float64(y*16))

			for x := 0; x < maxX; x++ {
				if y == 0 {
					dc.SetRGB(255, 0, 21)
					dc.DrawString(strconv.Itoa(x), 3+float64(x*16)-float64(y), 10)
				}
				dc.Push()
				dc.DrawRectangle(float64(x*16), float64(y*16), 15, 15)
				dc.SetLineCap(2)
				dc.SetRGBA(0, 0, 0, 0.5)
				dc.Fill()
			}
		}
	}
}
