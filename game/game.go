package game

import (
	"bag-weight/player"
	"bag-weight/utils"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Game struct {
	Players           []*player.Player
	MinWeight         int
	MaxWeight         int
	MaxGlobalGuesses  int
	GlobalGuessesNum  int
	MaxGameTime       int
	BagWeight         int
	GuessedNumbers    *utils.SafeHistoryMap
	IsGameOver        bool
	DidSomeoneWin     bool
	ClosestGuess      int
	BestGuesser       *player.Player
	GameStatusChannel chan bool
}

// NewGame creates and returns a new Game
func NewGame(
	players []*player.Player,
	minWeight int,
	maxWeight int,
	maxGlobalGuesses int,
	maxGameTime int,
) Game {
	rand.Seed(time.Now().UnixNano())
	bagWeight := rand.Intn(maxWeight-minWeight) + minWeight
	fmt.Println("The bag weight is ", bagWeight)
	g := Game{
		Players:           players,
		MinWeight:         minWeight,
		MaxWeight:         maxWeight,
		MaxGlobalGuesses:  maxGlobalGuesses,
		GlobalGuessesNum:  0,
		MaxGameTime:       maxGameTime,
		BagWeight:         bagWeight,
		IsGameOver:        false,
		DidSomeoneWin:     false,
		ClosestGuess:      minWeight + maxWeight,
		BestGuesser:       players[0],
		GameStatusChannel: make(chan bool),
	}
	g.InitiateGuessedNumbers()
	for _, p := range players {
		p.InitiatePlayerInGame(minWeight, maxWeight)
	}
	return g
}

// InitiateGuessedNumbers initiate the safe history map of the guessed
// numbers in the game, by initiating all the possible numbers to guess
// to false.
func (g *Game) InitiateGuessedNumbers() {
	guessesMap := make(map[int]bool)
	for i := g.MinWeight; i < g.MaxWeight; i++ {
		guessesMap[i] = false
	}
	g.GuessedNumbers = &utils.SafeHistoryMap{
		HistoryMap: guessesMap,
	}
}

// StartGame starts the guessing game and enforces the roles of the game.
// In case no player guessed the correct weight before the game ended,
// the function prints the player which guessed the closest weight, and
// the weight he had guessed
func (g *Game) StartGame() {
	go g.setDeadline()
	for _, p := range g.Players {
		go g.letPlayerGuess(p)
	}
	didSomeoneWin := <-g.GameStatusChannel
	if !didSomeoneWin {
		fmt.Println(
			"No one guessed the bag's weight",
			g.BagWeight,
			"\nThe closest player was",
			g.BestGuesser.Name,
			"with the guess:",
			g.ClosestGuess,
		)
	}
}

// setDeadline sets up a deadline of time to the game (which was configured in th constructor)
func (g *Game) setDeadline() {
	time.Sleep(time.Duration(g.MaxGameTime) * time.Millisecond)
	g.GameStatusChannel <- false
}

// handlePlayerGuess deals with a player guess. It updates the GuessedNumbers history
// map, and checks wheatear the player was correct, or if not, checks if the player made
// the closest guess in the game, then forces the sleep penalty.
func (g *Game) handlePlayerGuess(p *player.Player, guess int) {
	// fmt.Printf("%v guessed the number %v\n", p.Name, guess)
	g.GlobalGuessesNum += 1
	g.GuessedNumbers.Lock.Lock()
	g.GuessedNumbers.HistoryMap[guess] = true
	g.GuessedNumbers.Lock.Unlock()
	diff := int(math.Abs(float64(guess - g.BagWeight)))

	if diff == 0 {
		g.handlePlayerWin(p)
		return
	}

	if diff < g.ClosestGuess {
		g.ClosestGuess = guess
		g.BestGuesser = p
	}

	if g.GlobalGuessesNum >= g.MaxGlobalGuesses {
		g.IsGameOver = true
		return
	}

	// fmt.Printf("%v is sleeping for %v milliseconds\n", p.Name, diff)
	time.Sleep(time.Duration(diff) * time.Millisecond)
}

// handlePlayerWin finishes the game after a player guessed the correct weight,
// and prints the winner and the total number of guesses.
func (g *Game) handlePlayerWin(p *player.Player) {
	fmt.Println(p.PlayerType, p.Name, "won!\nAfter", g.GlobalGuessesNum, "overall guesses")
	g.IsGameOver = true
	g.ClosestGuess = 0
	g.DidSomeoneWin = true
}

// letPlayerGuess make player guesses in the game while the game is still on.
func (g *Game) letPlayerGuess(p *player.Player) {
	for !g.IsGameOver {
		guess := p.Guess(
			g.MinWeight,
			g.MaxWeight,
			g.GuessedNumbers,
		)
		g.handlePlayerGuess(p, guess)
	}
	g.GameStatusChannel <- g.DidSomeoneWin
}
