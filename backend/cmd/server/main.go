package main

import (
	"fmt"

	"emitrr_assignment/backend/internal/engine"
	"emitrr_assignment/backend/internal/game"
)

func main() {
	player := &game.Player{ID: "1", Username: "Human", Symbol: 'X'}
	botPlayer := &game.Player{ID: "2", Username: "Bot", Symbol: 'O'}

	g := game.NewGame("test", player)
	g.AddSecondPlayer(botPlayer)

	for {
		var col int
		fmt.Print("Enter column (1-7): ")
        fmt.Scan(&col)

        col = col - 1 // convert to 0-based for user friendly input


		result := engine.PlayTurn(g, col)

		printBoard(g)

		if result != "continue" {
			fmt.Println("Game result:", result)
			break
		}
	}
}

func printBoard(g *game.Game) {
	fmt.Println()
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%c ", cell)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
