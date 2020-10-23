package main

// file:///home/nicl/Downloads/catan_base_rules_2020_200707.pdf

// 19 terrain hexes (one desert)
// 6 sea frame pieces
// 9 harbour pieces
// 18 number tokens
// 95 resource cards
// - brick, grain, lumber, ore, wool (19 each)
// 25 development cards
// - 14 knights, 5 victory points, 2 monopoly, 2 year of plenty, 2 road building
// longest road, largest army
// 4 cities/player
// 5 settlements/player
// 15 roads/player
// 2 dice
// 1 robber
//

type DevelopmentCards struct {
	Knight, YearOfPlenty, Monopoly, RoadBuilding int
}

type Player struct {
	Resources struct {
		Brick, Grain, Lumber, Ore, Wool int
	}
	DevelopmentCards struct {
		InHand, Played DevelopmentCards
	}
}

// Actions are types that operate over a game board

// https://www.redblobgames.com/grids/hexagons/

// Turn is a pure function of state
// model grid
func main() {
	println("foo")
}
