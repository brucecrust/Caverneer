package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Map struct {
	areaOfMap []int
	worldMap  [][]int8
}

type Entity struct {
	name           string
	position       []int
	graphicChar    int8
	health, damage int
}

type CombatEntity interface {
	attack(CombatEntity)
	takeDamage(int)
	editPosition(*Map, []int, []*Entity)
	entityOnBlankTile(*Map, []int) bool
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

func (e *Entity) combat(worldMap *Map, enemy *Entity, enemiesOnMap []*Entity) {
	playerTurn := true
	for e.health > 0 && enemy.health > 0 {
		if playerTurn {
			e.attack(enemy)
			enemy.takeDamage(e.damage)
			playerTurn = !playerTurn
		} else {
			enemy.attack(e)
			e.takeDamage(enemy.damage)
			playerTurn = !playerTurn
		}
	}
	if e.health <= 0 {
		fmt.Printf("You have died...\n")
	}

	if enemy.health <= 0 {
		fmt.Printf("Enemy has died...\n")
		for x := range enemiesOnMap {
			if enemiesOnMap[x].name == enemy.name {
				worldMap.worldMap[enemy.position[0]][enemy.position[1]] = 0
				enemiesOnMap = removeEnemyIndex(enemiesOnMap, x)
				fmt.Println("Combat", len(enemiesOnMap))
			}
		}
	}
}

func removeEnemyIndex(a []*Entity, i int) []*Entity {
	a[i] = a[len(a)-1]
	return a[:len(a)-1]
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

func (e *Entity) entityOnBlankTile(worldMap *Map, movementArray []int) bool {
	if worldMap.worldMap[e.position[0]+movementArray[0]][e.position[1]+movementArray[1]] != 0 {
		return false
	}

	return true
}

func (e *Entity) mapCollisionDetection(worldMap *Map, movementArray []int) bool {
	if e.position[0]+movementArray[0] >= worldMap.areaOfMap[0] || e.position[0]+movementArray[0] < 0 {
		return false
	}

	if e.position[1]+movementArray[1] >= worldMap.areaOfMap[1] || e.position[1]+movementArray[1] < 0 {
		return false
	}

	return true
}

func (e *Entity) editPosition(worldMap *Map, movementArray []int, enemiesOnMap []*Entity) {
	canMove := e.mapCollisionDetection(worldMap, movementArray)
	for !canMove {
		if e.name == "player" {
			fmt.Printf("Cannot move this way...\n")
		}
		return
	}

	if canMove == true {
		worldMap.worldMap[e.position[0]][e.position[1]] = 0
		e.position[0] += movementArray[0]
		e.position[1] += movementArray[1]
		worldMap.worldMap[e.position[0]][e.position[1]] = e.graphicChar
		if e.name == "player" {
			for x := range enemiesOnMap {
				enemy := enemiesOnMap[x]
				if enemy.position[0] == e.position[0] && enemy.position[1] == e.position[1] {
					e.combat(worldMap, enemiesOnMap[x], enemiesOnMap)
				}
			}
		}
	}
}

func (e *Entity) randomizeMovement(worldMap *Map) []int {
	acceptableMovement := []int{-1, 0, 1}
	x := acceptableMovement[rand.Intn(len(acceptableMovement))]
	y := acceptableMovement[rand.Intn(len(acceptableMovement))]
	movementArray := []int{x, y}
	canMove := e.mapCollisionDetection(worldMap, movementArray)
	for !canMove {
		acceptableMovement := []int{-1, 0, 1}
		x := acceptableMovement[rand.Intn(len(acceptableMovement))]
		y := acceptableMovement[rand.Intn(len(acceptableMovement))]
		movementArray = []int{x, y}
		canMove = e.mapCollisionDetection(worldMap, movementArray)
	}
	entityOnBlankTile := e.entityOnBlankTile(worldMap, movementArray)
	if !entityOnBlankTile {
		e.randomizeMovement(worldMap)
	}

	return []int{x, y}
}

func createEnemies(worldMap *Map, iterationCount int) *Entity {
	enemy := &Entity{
		name:        "enemy " + strconv.Itoa(iterationCount),
		position:    []int{rand.Intn(worldMap.areaOfMap[0]), rand.Intn(worldMap.areaOfMap[1])},
		health:      10,
		damage:      5,
		graphicChar: 2,
	}

	return enemy
}

func main() {
	enemiesOnMap := []*Entity{}

	xLength := 5
	yLength := 5
	worldMap := &Map{
		areaOfMap: []int{xLength, yLength},
		worldMap:  createWorldMap(xLength, yLength),
	}

	player := &Entity{
		name:        "player",
		position:    []int{0, 0},
		health:      10,
		damage:      6,
		graphicChar: 1,
	}

	fmt.Printf("Adding enemy values to world map as `2`\n")
	for x := 0; x < 4; x++ {
		enemy := createEnemies(worldMap, x)
		enemiesOnMap = append(enemiesOnMap, enemy)
	}

	for x := 0; x < len(enemiesOnMap); x++ {
		fmt.Println("ENEMY NAME: ", enemiesOnMap[x].name)
		enemyPos := enemiesOnMap[x].randomizeMovement(worldMap)
		enemiesOnMap[x].editPosition(worldMap, enemyPos, enemiesOnMap)
	}

	fmt.Printf("Adding player value to world map as `1`, at pos (0, 0)\n")
	player.editPosition(worldMap, []int{0, 0}, enemiesOnMap)

	fmt.Printf("Starting loop...\n")
	fmt.Printf("\n")
	printWorldMap(worldMap)
	for {
		input := userInput([]string{"north", "south", "east", "west"})
		if input == "north" {
			player.editPosition(worldMap, []int{-1, 0}, enemiesOnMap)
		} else if input == "south" {
			player.editPosition(worldMap, []int{1, 0}, enemiesOnMap)
		} else if input == "east" {
			player.editPosition(worldMap, []int{0, 1}, enemiesOnMap)
		} else if input == "west" {
			player.editPosition(worldMap, []int{0, -1}, enemiesOnMap)
		}
		for x := 0; x < len(enemiesOnMap); x++ {
			enemyPos := enemiesOnMap[x].randomizeMovement(worldMap)
			enemiesOnMap[x].editPosition(worldMap, enemyPos, enemiesOnMap)
		}
		printWorldMap(worldMap)
	}
}
