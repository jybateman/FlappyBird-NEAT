package main

import (
	"log"
	"time"

	"github.com/jybateman/goneat"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
)

const BaseHeight int32 = 112
const BaseWidth int32 = 336
const ScreenHeight int32 = 512
const ScreenWidth int32 = 288
const PipeWidth int32 = 52
const PipeHeight int32 = 320
var BaseXPos int32 = 0

var FPS int = 30
var FrameRate int64 = int64(1.0/float64(FPS)*1000000000)
var FramesPerNano int64 = 33333333

const PipeVelX int32 = -4

type pipePos struct {
	x int32
	y int32
}
var uPipes []pipePos
var bPipes []pipePos

var gRenderer *sdl.Renderer
var gWindow *sdl.Window

var BirdTextCount int = 0
var BirdNumberAni int = 4
var BirdTexts [3]*sdl.Texture
var PipeText *sdl.Texture
var BGText *sdl.Texture
var BaseText *sdl.Texture

var AllDead int

func AddPipe() {
	// uPipes = append(uPipes, pipePos{ScreenWidth+10, 100-PipeHeight})
	uPipes = append(uPipes, pipePos{ScreenWidth+10, int32(gRand.Intn(142)+80)-PipeHeight})
	bPipes = append(bPipes, pipePos{ScreenWidth+10, uPipes[len(uPipes)-1].y+PipeHeight+100})
}

func IsDead(pIdx int) bool {
	if (Players[pIdx].x < uPipes[0].x + PipeWidth &&
		Players[pIdx].x + PlayerWidth > uPipes[0].x &&
		Players[pIdx].y < uPipes[0].y + PipeHeight &&
		PlayerHeight + Players[pIdx].y > uPipes[0].y) {
		return true
	}
	if (Players[pIdx].x < bPipes[0].x + PipeWidth &&
		Players[pIdx].x + PlayerWidth > bPipes[0].x &&
		Players[pIdx].y < bPipes[0].y + PipeHeight &&
		PlayerHeight + Players[pIdx].y > bPipes[0].y) {
		return true
	}

	if (Players[pIdx].x < uPipes[1].x + PipeWidth &&
		Players[pIdx].x + PlayerWidth > uPipes[1].x &&
		Players[pIdx].y < uPipes[1].y + PipeHeight &&
		PlayerHeight + Players[pIdx].y > uPipes[1].y) {
		return true
	}
	if (Players[pIdx].x < bPipes[1].x + PipeWidth &&
		Players[pIdx].x + PlayerWidth > bPipes[1].x &&
		Players[pIdx].y < bPipes[1].y + PipeHeight &&
		PlayerHeight + Players[pIdx].y > bPipes[1].y) {
		return true
	}

	if Players[pIdx].y+Players[pIdx].vel+PlayerHeight > ScreenHeight-BaseHeight {
		return true
	}

	return false
}

func Update(elapse int64) {
	for i := 0; i < len(uPipes); i++ {
		if uPipes[i].x < -PipeWidth {
			uPipes = append(uPipes[:i], uPipes[i+1:]...)
			bPipes = append(bPipes[:i], bPipes[i+1:]...)
			i--
		} else {
			uPipes[i].x += PipeVelX
			bPipes[i].x += PipeVelX
		}
	}
	BaseXPos += PipeVelX
	if BaseXPos <= -48 {
		BaseXPos = 0
	}
	if uPipes[0].x > 0 && bPipes[0].x < 5 {
		AddPipe()
	}
	for i := 0; i < len(Players); i++ {
		if !Players[i].dead {
			if Players[i].flap {
				// if Players[i].vel+PlayerFlapAcc > PlayerMaxAsc {
				// 	Players[i].vel += PlayerFlapAcc
				// } else {
				// 	Players[i].vel = PlayerMaxAsc
				// }
				if (Players[i].score >= 10) {
					Players[i].score -= 10
				}
				Players[i].rot = PlayerRotMax
				Players[i].vel = PlayerMaxAsc
				Players[i].flap = false
			} else {
				if Players[i].vel+PlayerDescAcc < PlayerMaxDesc {
					Players[i].vel += PlayerDescAcc
				} else {
					Players[i].vel = PlayerMaxDesc
				}
			}
			if Players[i].y+Players[i].vel+PlayerHeight < ScreenHeight-BaseHeight && Players[i].y+Players[i].vel > 0 {
				Players[i].y += Players[i].vel
			}
			if Players[i].rot+PlayerRotVel< PlayerRotMin {
				Players[i].rot += PlayerRotVel
			} else {
				Players[i].rot = PlayerRotMin
			}
			if IsDead(i) {
				Players[i].dead = true
				Genomes[i].SetFitness(int(Players[i].score))
				AllDead++
				// log.Println("DEAD!!")
			}
			if (bPipes[0].x <= 50 && Players[i].y <= bPipes[1].y-20 && Players[i].y >= uPipes[1].y+PipeHeight+20) {
				Players[i].score++
			} else if (Players[i].y <= bPipes[0].y-20 && Players[i].y >= uPipes[0].y+PipeHeight+20) {
				Players[i].score++
			}
		}
	}
}

func Draw() {
	//gRenderer.SetDrawColor(255, 255, 255, 255)
	gRenderer.Copy(BGText, nil, nil)
	//gRenderer.Clear()
	//gRenderer.SetDrawColor(0, 0, 0, 255)
	for i := range uPipes {
		gRenderer.CopyEx(PipeText, nil, &sdl.Rect{uPipes[i].x, uPipes[i].y, PipeWidth, PipeHeight}, 0, nil, sdl.FLIP_VERTICAL)
		gRenderer.Copy(PipeText, nil, &sdl.Rect{bPipes[i].x, bPipes[i].y, PipeWidth, PipeHeight})
	}
	for i := 0; i < len(Players); i++ {
		if !Players[i].dead {
			gRenderer.CopyEx(BirdTexts[BirdTextCount/BirdNumberAni], nil, &sdl.Rect{Players[i].x, Players[i].y, PlayerWidth, PlayerHeight}, Players[i].rot, nil, sdl.FLIP_NONE)
		}
	}
	BirdTextCount++
	if BirdTextCount == BirdNumberAni*3 {
		BirdTextCount = 0
	}

	if (bPipes[0].x <= 50) {
		gRenderer.DrawLine(0, bPipes[1].y-20, ScreenWidth, bPipes[1].y-20)
		gRenderer.DrawLine(0, uPipes[1].y+PipeHeight+20, ScreenWidth, uPipes[1].y+PipeHeight+20)
	} else {
		gRenderer.DrawLine(0, bPipes[0].y-20, ScreenWidth, bPipes[0].y-20)
		gRenderer.DrawLine(0, uPipes[0].y+PipeHeight+20, ScreenWidth, uPipes[0].y+PipeHeight+20)
	}
	
	gRenderer.Copy(BaseText, nil, &sdl.Rect{BaseXPos, ScreenHeight-BaseHeight, BaseWidth, BaseHeight})
	gRenderer.Present()
}

func ProcessInput() {
	id, err := gWindow.GetID()
	if err != nil {
		log.Fatalln(err)
	}
	pEvent := &sdl.UserEvent{sdl.USEREVENT, sdl.GetTicks(), id, 1331, nil, nil}

	_, err = sdl.PushEvent(pEvent) // Here's where the event is actually pushed
	if err != nil {
		log.Fatalln("Failed to push event:", err)
	}


	sdl.PollEvent()
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			if t.Type == sdl.KEYDOWN && t.Repeat == 0 && t.Keysym.Sym == sdl.K_SPACE {
				Players[0].flap = true
			}
		}
	}
}

func ProcessAIInput() {
	for i := 0; i < len(Players); i++ {
		// Genomes[i].SetInputs([]float64{float64(Players[i].x), float64(uPipes[0].x), float64(uPipes[0].y), float64(bPipes[0].x), float64(bPipes[0].y), float64(uPipes[1].x), float64(uPipes[1].y), float64(bPipes[1].x), float64(bPipes[1].y)})
		if (bPipes[0].x <= 50) {
			Genomes[i].SetInputs([]float64{float64(Players[i].y), float64(uPipes[1].y+PipeHeight), float64(bPipes[1].y)})
		} else {
			Genomes[i].SetInputs([]float64{float64(Players[i].y), float64(uPipes[0].y+PipeHeight), float64(bPipes[0].y)})
		}
		out := goneat.GetOutput(&Genomes[i])
		Players[i].flap = out[0] > 0.0
	}
}

func StartGame() {
	var elapse int64
	var lag int64 = 0

	uPipes = nil
	bPipes = nil
	AllDead = 0
	for i := 0; i < len(Players); i++ {
		Players[i] = Player{x:56, y:200, vel:0, rot:0, flap:false, score:0, dead:false}
	}
	AddPipe()
	AddPipe()
	uPipes[0].x = ScreenWidth + 200
	bPipes[0].x = ScreenWidth + 200
	uPipes[1].x = ScreenWidth + 200 + ScreenWidth / 2
	bPipes[1].x = ScreenWidth + 200 + ScreenWidth / 2
	previous := time.Now()
	for AllDead < len(Players) {
		current := time.Now()
		elapse = current.UnixNano() - previous.UnixNano()
		previous = current
		lag += elapse

		// ProcessInput()
		ProcessAIInput()

		Draw()
		Frames++
		// for lag >= FramesPerNano {
		Update(elapse)
		// 	lag -= FramesPerNano
		// }
		if elapse > 0 {
			elapse = 0
		}
		time.Sleep(time.Duration(FrameRate+elapse)*time.Nanosecond)
	}
}

func InitSdl() {
	var err error

	gWindow, err = sdl.CreateWindow("FlappyBird", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, ScreenWidth, ScreenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalln(err)
	}
	gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		log.Fatalln(err)
	}
	BirdTexts[0], err = img.LoadTexture(gRenderer, "./assets/sprites/bluebird-downflap.png")
	if err != nil {
		log.Fatalln(err)
	}
	BirdTexts[1], err = img.LoadTexture(gRenderer, "./assets/sprites/bluebird-midflap.png")
	if err != nil {
		log.Fatalln(err)
	}
	BirdTexts[2], err = img.LoadTexture(gRenderer, "./assets/sprites/bluebird-upflap.png")
	if err != nil {
		log.Fatalln(err)
	}
	PipeText, err = img.LoadTexture(gRenderer, "./assets/sprites/pipe-green.png")
	if err != nil {
		log.Fatalln(err)
	}
	BGText, err = img.LoadTexture(gRenderer, "./assets/sprites/background-day.png")
	if err != nil {
		log.Fatalln(err)
	}
	BaseText, err = img.LoadTexture(gRenderer, "./assets/sprites/base.png")
	if err != nil {
		log.Fatalln(err)
	}
}
