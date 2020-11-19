package main

import (
	"github.com/jybateman/goneat"
)

type Player struct {
	x int32
	y int32
	vel int32
	rot float64
	flap bool
	score int64
	dead bool
}
var Players []Player

var Genomes []goneat.Genome

const PlayerWidth int32 = 34
const PlayerHeight int32 = 24
const PlayerMaxDesc int32 = 10
const PlayerMaxAsc int32 = -8
const PlayerFlapAcc int32 = -9
const PlayerDescAcc int32 = 1
const PlayerRotVel float64 = 3
const PlayerRotMax float64 = -30
const PlayerRotMin float64 = 90

func NewPlayer() {
	// x: ScreenWidth*0.2
	Players = append(Players, Player{x:56, y:200, vel:0, rot:0, flap:false, score:0, dead:false})
}
