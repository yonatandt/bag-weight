package main

import (
	"bag-weight/game"
	"bag-weight/player"
	"fmt"
)

const (
	MIN_BAG_WEIGHT     = 40
	MAX_BAG_WEIGHT     = 140
	MAX_GLOBAL_GUESSES = 100
	MAX_GAME_TIME      = 1500
	MIN_NUM_OF_PLAYERS = 2
	MAX_NUM_OF_PLAYERS = 8
)

// getPlayerNumInput gets the number of players in the game, and
// also ensures that the number is in the pre-configured range.
func getPlayerNumInput() int {
	fmt.Println(
		"Enter number of players (",
		MIN_NUM_OF_PLAYERS,
		"-",
		MAX_NUM_OF_PLAYERS,
		"):",
	)
	var numOfPlayers int
	_, err := fmt.Scanln(&numOfPlayers)
	for {
		if err == nil && numOfPlayers >= 2 && numOfPlayers <= 8 {
			break
		}
		fmt.Println(
			"Players number must be a number between ",
			MIN_NUM_OF_PLAYERS,
			"-",
			MAX_NUM_OF_PLAYERS,
		)
		_, err = fmt.Scanln(&numOfPlayers)
	}
	return numOfPlayers
}

// getInputAboutPlayers gets info about the game players from a user's input.
func getInputAboutPlayers(numOfPlayers int) []*player.Player {
	players := []*player.Player{}
	for i := 0; i < numOfPlayers; i++ {
		var playerName string
		var playerTypeInput int
		fmt.Println("Enter player", i+1, "name:")
		fmt.Scanln(&playerName)
		fmt.Println("Enter player", i+1, "type (1-5):")
		fmt.Scanln(&playerTypeInput)
		playerType := player.NumberToTypeOfPlayer[playerTypeInput]
		p := player.NewPlayer(playerName, playerType)
		players = append(players, &p)
	}
	return players

}

// getUserInput gets the required user input to start the game.
func getUserInput() []*player.Player {
	numOfPlayers := getPlayerNumInput()
	return getInputAboutPlayers(numOfPlayers)
}

func main() {
	players := getUserInput()
	g := game.NewGame(players, MIN_BAG_WEIGHT, MAX_BAG_WEIGHT, MAX_GLOBAL_GUESSES, MAX_GAME_TIME)

	g.StartGame()
}
