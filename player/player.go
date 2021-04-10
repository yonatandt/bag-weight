package player

import (
	"bag-weight/utils"
	"math/rand"
	"time"
)

type TypeOfPlayer string

const (
	Random           TypeOfPlayer = "Random"
	Remembering      TypeOfPlayer = "Remembering"
	Thorough         TypeOfPlayer = "Through"
	Cheating         TypeOfPlayer = "Cheating"
	ThoroughCheating TypeOfPlayer = "Thorugh Cheating"
)

var NumberToTypeOfPlayer = map[int]TypeOfPlayer{
	1: Random,
	2: Remembering,
	3: Thorough,
	4: Cheating,
	5: ThoroughCheating,
}

type Player struct {
	Name       string
	PlayerType TypeOfPlayer
	Memory     []int
	LastGuess  int
}

// NewPlayer creates and returns a new player
func NewPlayer(name string, playerType TypeOfPlayer) Player {
	p := Player{
		Name:       name,
		PlayerType: playerType,
		Memory:     []int{},
	}
	return p
}

// InitiatePlayerInGame initate information relevant to the player by its type
// Currently relevant just for the remembering type which hold the history of his guesses.
func (p *Player) InitiatePlayerInGame(min int, max int) {
	switch p.PlayerType {
	case Remembering:
		slice := []int{}
		for i := min; i < max; i++ {
			slice = append(slice, i)
		}
		p.Memory = slice
	}
}

// Guess returns a player's guess depending on his type
func (p *Player) Guess(min int, max int, globalMemory *utils.SafeHistoryMap) int {
	rand.Seed(time.Now().UnixNano())
	switch p.PlayerType {
	case Random:
		return rand.Intn(max-min) + min
	case Remembering:
		index := rand.Intn(len(p.Memory))
		value := p.Memory[index]
		p.Memory = append(p.Memory[:index], p.Memory[index+1:]...)
		return value
	case Thorough:
		if p.LastGuess < min {
			p.LastGuess = min
			return min
		}
		p.LastGuess += 1
		return p.LastGuess
	case Cheating:
		notGuessedYet := []int{}
		globalMemory.Lock.RLock()
		for number, isGuessed := range globalMemory.HistoryMap {
			if !isGuessed {
				notGuessedYet = append(notGuessedYet, number)
			}
		}
		globalMemory.Lock.RUnlock()
		return notGuessedYet[rand.Intn(len(notGuessedYet))]
	case ThoroughCheating:
		globalMemory.Lock.RLock()
		var guess int
		for number, isGuessed := range globalMemory.HistoryMap {
			if !isGuessed {
				guess = number
			}
		}
		globalMemory.Lock.RUnlock()
		return guess
	}
	// default value in case of a wrong type of player.
	return min
}
