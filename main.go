package main

import (
	"fmt"
	"time"
	"math/rand"

	"github.com/jybateman/goneat"
)

var Frames int

var gRand *rand.Rand

var PopSize int = 200

func FrameCount() {
	for {
		fmt.Println(goneat.GetGeneration(), len(goneat.GetSpecies()))
		Frames = 0
		time.Sleep(time.Second)
	}
}

func main() {
	s := rand.NewSource(time.Now().UnixNano())
	gRand = rand.New(s)
	Genomes = goneat.InitNEAT(PopSize, 3, 1)

	// go FrameCount()

	InitSdl()
	for i := 0; i < PopSize; i++ {
		NewPlayer()
	}
	for {
		StartGame()
		goneat.NextGeneration(Genomes)
		fmt.Println(goneat.GetGeneration(), len(goneat.GetSpecies()))
	}
}

