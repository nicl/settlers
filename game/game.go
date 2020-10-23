// Package game provides types and helpers for the game. The grid is
// achieved using an axial implementation. See
// https://www.redblobgames.com/grids/hexagons/#coordinates-axial for
// more on this. I've relabelled 'q' and 'r' as 'Column' and 'Row' for
// readability.
package game

// Resource is a type for resources.
type Resource int

// Resources and tile types for the game.
const (
	Brick Resource = iota
	Grain
	Lumber
	Ore
	Wool

	// Pseudo-resources for extra tile types
	Desert
	Sea
)

// Hex contains axial coordinates and other meta for a tile.
type Hex struct {
	Column, Row int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

// IsNeighbour is a predicate that returns true if 'other' is a
// neighbour to the hex. Returns false is a = b.
func (a Hex) IsNeighbour(b Hex) bool {
	neighbours := map[Hex]interface{}{
		{1, 0}:  nil,
		{1, -1}: nil,
		{0, -1}: nil,
		{-1, 0}: nil,
		{-1, 1}: nil,
		{0, 1}:  nil,
	}

	colOffset := a.Column - b.Column
	rowOffset := a.Row - b.Row

	_, ok := neighbours[Hex{Column: colOffset, Row: rowOffset}]

	return ok
}

// FindSharedNeighbour returns the shared neighbour of two neighbouring hexes.
func FindSharedNeighbour(a, b Hex, candidates []Hex) (Hex, bool) {
	for _, candidate := range candidates {
		if a.IsNeighbour(candidate) && b.IsNeighbour(candidate) {
			return candidate, true
		}
	}

	return Hex{}, false
}

// Tile is a hex location as well as resource type and dice number.
type Tile struct {
	Location Hex
	Resource Resource
	Number   int
}

// Road denotes location (lies between two neighbours) and who owns it.
type Road struct {
	A, B   Hex
	Player int
}

// Settlement denotes location (meeting point of three neighbours) and
// owner.
type Settlement struct {
	A, B, C Hex
	Player  int
}

// City denotes location (meeting point of three neighbours) and
// owner.
type City struct {
	A, B, C Hex
	Player  int
}

// DevelopmentCard is a type alias for development cards.
type DevelopmentCard int

// Development card enums.
const (
	Knight DevelopmentCard = iota
	YearOfPlenty
	Monopoly
	RoadBuilding
	VictoryPoint
)

// Player represents a player's hand.
type Player struct {
	Resources struct {
		Brick, Grain, Lumber, Ore, Wool int
	}
	DevelopmentCards struct {
		InHand, Played []DevelopmentCard
	}
}

// Board is the general wrapper type for game state for a board. Scores
// and actions can be derived from this. Actions operate on a board (are
// 'reducers').
type Board struct {
	Robber  Hex
	Players []Player
	Tiles   []Tile
}

// * find (available) adjacent vertices
// * find settlements/cities 'on' a tile

// Implementation
// * use axial and then with neighbours

// * a road connects two neighbours
// * a settlement connects three neighbours (or two if edge of board)
