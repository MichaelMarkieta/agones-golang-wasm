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
	"strconv"
	"strings"
	"time"
)

const (
	screenWidth   = 1024
	screenHeight  = 1024
	spriteScaling = 0.25
	spriteSpeed   = 5
	spriteWidth   = 284
	spriteHeight  = 285
	spriteCountY  = 6
	bottomBound = screenHeight-(spriteHeight*spriteScaling)
	rightBound = screenWidth-(spriteWidth*spriteScaling)
)

var (
	spritesheet *ebiten.Image
	spritesheetMap = map[int]*image.Rectangle{}
	localPlayer string
)

type Game struct {
	spritesheet *ebiten.Image
	players     map[string]*Player
	wsconn      *websocket.Conn
	ctx         context.Context
}

type Player struct {
	x float64
	y float64
	sprite int
}

func (g *Game) Update(screen *ebiten.Image) error {
	keypressed := false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			keypressed = true
			// Moving Down
			if ebiten.Key.String(k) == "Down" && g.players[localPlayer].y <= bottomBound {
				if g.players[localPlayer].y + 1 * spriteSpeed > bottomBound {
					g.players[localPlayer].y = bottomBound
				} else {
					g.players[localPlayer].y += 1 * spriteSpeed
				}
			}
			// Moving Up
			if ebiten.Key.String(k) == "Up" && g.players[localPlayer].y >= 0 {
				if g.players[localPlayer].y - 1 * spriteSpeed < 0 {
					g.players[localPlayer].y = 0
				} else {
					g.players[localPlayer].y -= 1 * spriteSpeed
				}
			}
			// Moving Left
			if ebiten.Key.String(k) == "Left" && g.players[localPlayer].x > 0 {
				if g.players[localPlayer].x - 1 * spriteSpeed < 0 {
					g.players[localPlayer].x = 0
				} else {
					g.players[localPlayer].x -= 1 * spriteSpeed
				}
			}
			// Moving Right
			if ebiten.Key.String(k) == "Right" && g.players[localPlayer].x <= rightBound {
				if g.players[localPlayer].x + 1 * spriteSpeed > rightBound {
					g.players[localPlayer].x = rightBound
				} else {
					g.players[localPlayer].x += 1 * spriteSpeed
				}
			}
		}
	}

	if keypressed {
		err := wsjson.Write(g.ctx, g.wsconn, fmt.Sprintf("POSITION %s %f %f %d", localPlayer, g.players[localPlayer].x, g.players[localPlayer].y, g.players[localPlayer].sprite))
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	err1 := screen.Fill(color.NRGBA{0xf0, 0xf0, 0xf0, 0xff}); if err1 != nil {
		log.Fatalf("Cannot fill screen: %s", err1)
	}

	for _, p := range g.players {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(spriteScaling, spriteScaling)
		op.GeoM.Translate(p.x, p.y)
		err2 := screen.DrawImage(spritesheet.SubImage(*spritesheetMap[p.sprite]).(*ebiten.Image), op); if err2 != nil {
			log.Fatalf("Cannot draw sprite: %s", err2)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newPlayer() *Player {
	p := &Player{}
	p.x = float64(rand.Intn(screenWidth- spriteWidth))
	p.y = float64(rand.Intn(screenHeight- spriteHeight))
	p.sprite = rand.Intn(spriteCountY)+1
	return p
}

func isNewPlayer(g *Game, id string) bool {
	log.Printf("Checking if %s is new", id)
	exists := false
	_, exists = g.players[id]
	return !exists
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// SPRITESHEET SETUP
	f, err := ebitenutil.OpenFile("round_nodetails_outline.png")
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	spritesheet, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	spritesheetMap[1] = &image.Rectangle{image.Point{0,((1*spriteHeight)+1)-spriteHeight},image.Point{spriteWidth,(1*spriteHeight)+1}}
	spritesheetMap[2] = &image.Rectangle{image.Point{0,((2*spriteHeight)+1)-spriteHeight},image.Point{spriteWidth,(2*spriteHeight)+1}}
	spritesheetMap[3] = &image.Rectangle{image.Point{0,((3*spriteHeight)+1)-spriteHeight},image.Point{spriteWidth,(3*spriteHeight)+1}}
	spritesheetMap[4] = &image.Rectangle{image.Point{0,((4*spriteHeight)+1)-spriteHeight},image.Point{spriteWidth,(4*spriteHeight)+1}}
	spritesheetMap[5] = &image.Rectangle{image.Point{0,((5*spriteHeight)+1)-spriteHeight},image.Point{spriteWidth,(5*spriteHeight)+1}}
	spritesheetMap[6] = &image.Rectangle{image.Point{0,((6*spriteHeight)+1)-spriteHeight},image.Point{spriteWidth,(6*spriteHeight)+1}}
	spritesheetMap[7] = &image.Rectangle{image.Point{0,((7*spriteHeight)+1)-spriteHeight},image.Point{spriteWidth,(7*spriteHeight)+1}}

	// GAME SETUP
	g := &Game{}
	localPlayer = uuid.New().String()
	p := newPlayer()
	g.players = make(map[string]*Player)
	g.players[localPlayer] = p

	// WEBSOCKET SETUP
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wsconn, _, err := websocket.Dial(ctx, "ws://gameserver.michaelmarkieta.com:7777", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer wsconn.Close(websocket.StatusInternalError, "")

	go func() {
		for {
			_, message, err := wsconn.Read(ctx);
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Message receieved: %s", message)
			parts := strings.Split(strings.Replace(strings.Trim(string(message), "\""), "\"\n\"", " ", -1), " ")
			log.Print(parts)
			switch parts[0] {
			case "POSITION":
				x, err := strconv.ParseFloat(parts[2], 64); if err != nil {log.Fatal(err)}
				y, err := strconv.ParseFloat(parts[3], 64); if err != nil {log.Fatal(err)}
				sprite, err := strconv.Atoi(parts[4]); if err != nil {log.Fatal(err)}
				if isNewPlayer(g, parts[1]) {
					log.Printf("Add remote player")
					g.players[parts[1]] = &Player{x: x, y: y, sprite: sprite}
				}
				if parts[1] != localPlayer {
					log.Printf("Move remote player")
					g.players[parts[1]].x = x
					g.players[parts[1]].y = y
					g.players[parts[1]].sprite = sprite
				}
			case "PLAYER_LEAVE":
				log.Printf("Remove player [PLAYER_LEAVE UUID]")
			}
		}
	}()

	// RUN GAME
	g.wsconn = wsconn
	g.ctx = ctx
	err4 := wsjson.Write(g.ctx, g.wsconn, fmt.Sprintf("POSITION %s %f %f %d", localPlayer, g.players[localPlayer].x, g.players[localPlayer].y, g.players[localPlayer].sprite))
	if err4 != nil {
		log.Fatal(err4)
	}
	ebiten.SetRunnableOnUnfocused(true)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}