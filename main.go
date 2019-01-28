package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MapLength struct {
	xLength, yLength int
}
type Map struct {
	worldMap [][]int8
}

type Entity struct {
	position       []int
	graphicalChar  int8
	health, damage int
}

type CombatEntity interface {
	attack(CombatEntity)
	takeDamage(int)
	editPosition(*Map, int8, []int)
}

func (e *Entity) attack(target CombatEntity) {
	target.takeDamage(e.damage)
}

func (e *Entity) takeDamage(damageTaken int) {
	e.health -= damageTaken
}

func createWorldMap(xLength int, yLength int) [][]int8 {
	worldMap := make([][]int8, 5)
	for i := range worldMap {
		worldMap[i] = make([]int8, 5)
	}
	return worldMap
}

func (e *Entity) editPosition(worldMap *Map, graphicChar int8, xyMove []int) {
	worldMap.worldMap[e.position[0]][e.position[1]] = 0
	e.position[0] += xyMove[0]
	e.position[1] += xyMove[1]
	worldMap.worldMap[e.position[0]][e.position[1]] = graphicChar
}

func printWorldMap(worldMap [][]int8) {
	for i := range worldMap {
		fmt.Println(worldMap[i])
	}
	fmt.Printf("\n")
}

func userInput(acceptableInput []string) string {
	for {
		var inputList []string
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		for x := range acceptableInput {
			if strings.HasPrefix(acceptableInput[x], strings.ToLower(input)) {
				inputList = append(inputList, acceptableInput[x])
			}
		}

		// Enter just enough input to distinguish between items in the acceptableInput array.
		if len(inputList) == 0 {
			fmt.Printf("Please input another command, I do not understand.\n")
			continue

		} else if len(inputList) == 1 {
			answer := inputList[0]
			stringAnswer := string(answer)
			fmt.Printf("Correct! Your answer was: %v\n", stringAnswer)
			return stringAnswer

		} else {
			fmt.Printf("Which of the following commands did you mean?: %v \n", acceptableInput)
			continue
		}
	}
}

func main() {
	mapLength := MapLength{
		xLength: 10,
		yLength: 10,
	}

	worldMap := &Map{
		worldMap: createWorldMap(mapLength.xLength, mapLength.yLength),
	}
	printWorldMap(worldMap.worldMap)

	player := &Entity{
		position:      []int{0, 0},
		health:        10,
		damage:        5,
		graphicalChar: 1,
	}

	enemy := &Entity{
		position:      []int{0, 0},
		health:        10,
		damage:        5,
		graphicalChar: 2,
	}

	player.editPosition(worldMap, player.graphicalChar, []int{0, 1})
	fmt.Printf("Adding player value to world map as `1`, at pos (0, 0)\n")
	printWorldMap(worldMap.worldMap)

	fmt.Printf("Starting combat...")
	fmt.Println("Attacking...\n", player.health, enemy.health)
	player.attack(enemy)
	fmt.Println("Damage taken...\n", player.health, enemy.health)

	fmt.Printf("Starting loop...\n")

	for {
		input := userInput([]string{"north", "south", "east", "west"})
		if input == "north" {
			player.editPosition(worldMap, player.graphicalChar, []int{-1, 0})
		} else if input == "south" {
			player.editPosition(worldMap, player.graphicalChar, []int{1, 0})
		} else if input == "east" {
			player.editPosition(worldMap, player.graphicalChar, []int{0, 1})
		} else if input == "west" {
			player.editPosition(worldMap, player.graphicalChar, []int{0, -1})
		}
		printWorldMap(worldMap.worldMap)
	}
}
