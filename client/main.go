package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

const (
	screenWidth   = 1024
	screenHeight  = 1024
	frameOX       = 0
	frameOY       = 6
	spriteScaling = 0.25
	spriteSpeed   = 5
	spriteWidth   = 284
	spriteHeight  = 285
	bottomBound = screenHeight-(spriteHeight*spriteScaling)
	rightBound = screenWidth-(spriteWidth*spriteScaling)
)

var (
	spritesheet *ebiten.Image
	playerYellow = image.Rectangle{image.Point{0,0},image.Point{284,285}}
)

type Game struct {
	spritesheet *ebiten.Image
	pressed []ebiten.Key
	players []*Player
	c *websocket.Conn
	ctx context.Context
}

type Player struct {
	id string
	x float64
	y float64
	sprite image.Rectangle
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.pressed = nil
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			g.pressed = append(g.pressed, k)

			// Moving Down
			if ebiten.Key.String(k) == "Down" && g.players[0].y <= bottomBound {
				if g.players[0].y + 1 * spriteSpeed > bottomBound {
					g.players[0].y = bottomBound
				} else {
					g.players[0].y += 1 * spriteSpeed
				}
			}

			// Moving Up
			if ebiten.Key.String(k) == "Up" && g.players[0].y >= 0 {
				if g.players[0].y - 1 * spriteSpeed < 0 {
					g.players[0].y = 0
				} else {
					g.players[0].y -= 1 * spriteSpeed
				}
			}

			// Moving Left
			if ebiten.Key.String(k) == "Left" && g.players[0].x > 0 {
				if g.players[0].x - 1 * spriteSpeed < 0 {
					g.players[0].x = 0
				} else {
					g.players[0].x -= 1 * spriteSpeed
				}
			}

			// Moving Right
			if ebiten.Key.String(k) == "Right" && g.players[0].x <= rightBound {
				if g.players[0].x + 1 * spriteSpeed > rightBound {
					g.players[0].x = rightBound
				} else {
					g.players[0].x += 1 * spriteSpeed
				}
			}
		}
	}

	err := wsjson.Write(g.ctx, g.c, fmt.Sprintf("%s %f %f", g.players[0].id, g.players[0].x, g.players[0].y))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	err1 := screen.Fill(color.NRGBA{0xff, 0xff, 0xff, 0xff}); if err1 != nil {
		log.Fatalf("Cannot fill screen: %s", err1)
	}

	for _, p := range g.players {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(spriteScaling, spriteScaling)
		op.GeoM.Translate(p.x, p.y)
		err2 := screen.DrawImage(spritesheet.SubImage(p.sprite).(*ebiten.Image), op); if err2 != nil {
			log.Fatalf("Cannot draw sprite: %s", err2)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newPlayer() *Player {
	rand.Seed(time.Now().UnixNano())
	p := &Player{}
	p.id = uuid.New().String()
	p.x = float64(rand.Intn(screenWidth- spriteWidth))
	p.y = float64(rand.Intn(screenHeight- spriteHeight))
	p.sprite = playerYellow
	return p
}

func main() {
	// WEBSOCKET SETUP
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://34.95.7.42:7777", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close(websocket.StatusInternalError, "")

	// SPRITESHEET SETUP
	f, err := ebitenutil.OpenFile("round_nodetails_outline.png")
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	spritesheet, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)

	// GLOBAL VARIABLES SETUP
	g := &Game{}
	p := newPlayer()
	g.players = append(g.players, p)
	g.c = c
	g.ctx = ctx

	err3 := wsjson.Write(g.ctx, g.c, fmt.Sprintf("%s %f %f", g.players[0].id, g.players[0].x, g.players[0].y))
	if err3 != nil {
		log.Fatal(err3)
	}

	// RUN GAME
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}