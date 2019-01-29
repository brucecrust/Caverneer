package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Map struct {
	areaOfMap []int
	worldMap  [][]int8
}

type Entity struct {
	position       []int
	graphicChar    int8
	health, damage int
}

type CombatEntity interface {
	attack(CombatEntity)
	takeDamage(int)
	editPosition(*Map, []int) bool
}

type Enemy interface {
	randomizeMovement(*Map, *Entity)
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

func (e *Entity) editPosition(worldMap *Map, xyMove []int) bool {
	entityCanMove := true

	if e.position[1]+xyMove[1] > worldMap.areaOfMap[1]-1 || e.position[1]+xyMove[1] < 0 {
		entityCanMove = false
		fmt.Printf("Cannot move this way...\n")
	}

	if e.position[0]+xyMove[0] > worldMap.areaOfMap[0]-1 || e.position[0]+xyMove[0] < 0 {
		entityCanMove = false
		fmt.Printf("Cannot move this way...\n")
	}

	if entityCanMove == true {
		worldMap.worldMap[e.position[0]][e.position[1]] = 0
		e.position[0] += xyMove[0]
		e.position[1] += xyMove[1]
		worldMap.worldMap[e.position[0]][e.position[1]] = e.graphicChar
	}
	return entityCanMove
}

func printWorldMap(worldMap *Map) {
	for i := range worldMap.worldMap {
		fmt.Println(worldMap.worldMap[i])
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

func (e *Entity) randomizeMovement(worldMap *Map, player *Entity) []int {
	movementArray := []int{-1, 0, 1}
	x := movementArray[rand.Intn(len(movementArray))]
	y := movementArray[rand.Intn(len(movementArray))]

	for {
		if e.position[0]+x < worldMap.areaOfMap[0]-1 || e.position[0]+x > 0 {
			if e.position[1]+y < worldMap.areaOfMap[1]-1 || e.position[0]+y > 0 {
				break
			} else {
				y = movementArray[rand.Intn(len(movementArray))]
				continue
			}
		} else {
			x = movementArray[rand.Intn(len(movementArray))]
			continue
		}
	}

	fmt.Println(e.position[0]+x, e.position[1]+y)
	return []int{x, y}
}

func main() {
	xLength := 5
	yLength := 5
	worldMap := &Map{
		areaOfMap: []int{xLength, yLength},
		worldMap:  createWorldMap(xLength, yLength),
	}

	player := &Entity{
		position:    []int{0, 0},
		health:      10,
		damage:      5,
		graphicChar: 1,
	}

	enemy := &Entity{
		position:    []int{1, 1},
		health:      10,
		damage:      5,
		graphicChar: 2,
	}

	fmt.Printf("Adding player value to world map as `1`, at pos (0, 0)\n")
	player.editPosition(worldMap, []int{0, 0})

	fmt.Printf("Adding enemy value to world map as `2`\n")
	enemyPos := enemy.randomizeMovement(worldMap, player)
	enemy.editPosition(worldMap, enemyPos)

	fmt.Printf("Starting combat...")
	fmt.Println("Attacking...\n", player.health, enemy.health)
	player.attack(enemy)
	fmt.Println("Damage taken...\n", player.health, enemy.health)

	fmt.Printf("Starting loop...\n")
	fmt.Printf("\n")
	printWorldMap(worldMap)
	for {
		enemyPos = enemy.randomizeMovement(worldMap, player)
		enemy.editPosition(worldMap, enemyPos)

		input := userInput([]string{"north", "south", "east", "west"})
		if input == "north" {
			player.editPosition(worldMap, []int{-1, 0})
		} else if input == "south" {
			player.editPosition(worldMap, []int{1, 0})
		} else if input == "east" {
			player.editPosition(worldMap, []int{0, 1})
		} else if input == "west" {
			player.editPosition(worldMap, []int{0, -1})
		}
		printWorldMap(worldMap)
	}
}
