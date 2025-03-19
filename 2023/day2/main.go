package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Game struct {
	ID      int
	Rounds  []*Round
	IsValid bool
}

type Round struct {
	Game  *Game
	Red   int
	Green int
	Blue  int
}

var Limits = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
	}
	defer file.Close()

	fmt.Println("Reading game history...")
	games := []*Game{}
	rounds := []Round{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Split(line, ":")        // {"Game ID", rgb_sets}
		temp2 := strings.Fields(temp[0])        // {"Game", "ID"}
		game_id, _ := strconv.Atoi(temp2[1])    // Convert ID to int
		rgb_sets := strings.Split(temp[1], ";") // {rgb, rgb, ...}

		new_game := &Game{ID: game_id, IsValid: true}
		for _, rgb_set := range rgb_sets {
			new_round := Round{Game: new_game, Red: 0, Green: 0, Blue: 0}
			rgb := strings.Split(rgb_set, ",") // {"# red", "# green", "# blue"}

			for _, group := range rgb {
				cube_info := strings.Fields(group) // {"num_cubes", "cube_color"}
				num_cubes, _ := strconv.Atoi(cube_info[0])
				cube_color := cube_info[1]

				if cube_color == "red" {
					new_round.Red = num_cubes
				} else if cube_color == "green" {
					new_round.Green = num_cubes
				} else if cube_color == "blue" {
					new_round.Blue = num_cubes
				}
			}
			rounds = append(rounds, new_round)
			new_game.Rounds = append(new_game.Rounds, &new_round)
		}
		games = append(games, new_game)
	}

	fmt.Println("Validating games...")
	validateGames(&rounds)
	id_sum := 0
	for _, game := range games {
		if game.IsValid {
			id_sum += game.ID
		}
	}
	fmt.Println("Calculating powers...")
	power_sum := 0
	for _, game := range games {
		power_sum += getPower(*game)
	}
	fmt.Printf("[Part 1] The sum of id numbers for valid games is: %d\n", id_sum)
	fmt.Printf("[Part 2] The sum of powers for all the games is: %d\n", power_sum)
}

func validateGames(rounds *[]Round) {
	for _, round := range *rounds {
		if !round.Game.IsValid {
			continue
		}
		if round.Red > Limits["red"] {
			round.Game.IsValid = false
			continue
		}
		if round.Green > Limits["green"] {
			round.Game.IsValid = false
			continue
		}
		if round.Blue > Limits["blue"] {
			round.Game.IsValid = false
			continue
		}
	}
}

func getPower(game Game) int {
	max_red := 0
	max_green := 0
	max_blue := 0

	for _, round := range game.Rounds {
		if round.Red > max_red {
			max_red = round.Red
		}
		if round.Green > max_green {
			max_green = round.Green
		}
		if round.Blue > max_blue {
			max_blue = round.Blue
		}
	}

	power := max_red * max_green * max_blue
	return power
}
