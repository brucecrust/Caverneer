package main

import "fmt"

type MapLength struct {
	xLength, yLength int
}

type MapPoint struct {
	x, y int
}

type Map struct {
	playerPosition MapPoint
	worldMap       [][]int8
}

type Entity struct {
	health, damage int
}

type CombatEntity interface {
	attack(CombatEntity)
	takeDamage(int)
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

func editPlayerPosition(worldMap *Map, newPlayerPosition MapPoint) {
	worldMap.playerPosition = newPlayerPosition
	worldMap.worldMap[worldMap.playerPosition.x][worldMap.playerPosition.y] = 1
}

func printWorldMap(worldMap [][]int8) {
	for i := range worldMap {
		fmt.Println(worldMap[i])
	}
	fmt.Printf("\n")
}

func main() {
	mapLength := MapLength{
		xLength: 10,
		yLength: 10,
	}

	worldMap := &Map{
		playerPosition: MapPoint{
			0, 0,
		},
		worldMap: createWorldMap(mapLength.xLength, mapLength.yLength),
	}
	printWorldMap(worldMap.worldMap)

	player := &Entity{
		health: 10,
		damage: 5,
	}

	playerPosition := MapPoint{
		x: 0,
		y: 0,
	}

	enemy := &Entity{
		health: 10,
		damage: 5,
	}

	editPlayerPosition(worldMap, playerPosition)
	fmt.Printf("Adding player value to world map as `1`, at pos (0, 0)\n")
	printWorldMap(worldMap.worldMap)

	fmt.Printf("Starting combat...")
	fmt.Println("Attacking...\n", player.health, enemy.health)
	player.attack(enemy)
	fmt.Println("Damage taken...\n", player.health, enemy.health)
}
