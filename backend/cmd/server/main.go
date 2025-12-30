package main

import (
	"fmt"

	"emitrr_assignment/backend/internal/bot"
	"emitrr_assignment/backend/internal/game"
)

func main() {
	player := &game.Player{ID: "1", Username: "Human", Symbol: 'X'}
	botPlayer := &game.Player{ID: "2", Username: "Bot", Symbol: 'O'}

	g := game.NewGame("test", player)
	g.AddSecondPlayer(botPlayer)

	col := bot.DecideMove(g, 'O', 'X')
	fmt.Println("Bot decided to play column:", col)
}
