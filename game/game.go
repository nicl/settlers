// Package game provides types and helpers for the game. The grid is
// achieved using an axial implementation. See
// https://www.redblobgames.com/grids/hexagons/#coordinates-axial for
// more on this. I've relabelled 'q' and 'r' as 'Column' and 'Row' for
// readability. It ranges from +-2 in each axis, but note the sum never
// exceeds 2.
package game

import "log"

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

// String returns a sensible string representation.
func (r Resource) String() string {
	switch r {
	case Brick:
		return "Brick"
	case Grain:
		return "Grain"
	case Lumber:
		return "Lumber"
	case Ore:
		return "Ore"
	case Wool:
		return "Wool"
	default:
		log.Fatalf("Unknown resource: %d", r)
		return ""
	}
}

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
func (h Hex) IsNeighbour(b Hex) bool {
	neighbours := map[Hex]interface{}{
		{1, 0}:  nil,
		{1, -1}: nil,
		{0, -1}: nil,
		{-1, 0}: nil,
		{-1, 1}: nil,
		{0, 1}:  nil,
	}

	colOffset := h.Column - b.Column
	rowOffset := h.Row - b.Row

	_, ok := neighbours[Hex{Column: colOffset, Row: rowOffset}]

	return ok
}

// ContainsNeighbour returns true if any candidate is neighbour to hex
func (h Hex) ContainsNeighbour(candidates ...Hex) bool {
	for _, c := range candidates {
		if h.IsNeighbour(c) {
			return true
		}
	}

	return false
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

func FindAdjacentCities(location Hex, cities []City) []City {
	var adjacent []City

	for _, c := range cities {
		if location.ContainsNeighbour(c.A, c.B, c.C) {
			adjacent = append(adjacent, c)
		}
	}

	return adjacent
}

func FindAdjacentSettlements(location Hex, settlements []Settlement) []Settlement {
	var adjacent []Settlement

	for _, s := range settlements {
		if location.ContainsNeighbour(s.A, s.B, s.C) {
			adjacent = append(adjacent, s)
		}
	}

	return adjacent
}

// Tile is a hex location as well as resource type and dice number.
type Tile struct {
	Location Hex
	Resource Resource
	Number   int
}

// Road denotes location (lies between two neighbours) and who owns it.
type Road struct {
	A, B Hex
}

// Settlement denotes location (meeting point of three neighbours) and
// owner.
type Settlement struct {
	A, B, C Hex
}

// City denotes location (meeting point of three neighbours) and
// owner.
type City struct {
	A, B, C Hex
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

// PlayerID uniquely identifies a player.
type PlayerID int

// Resources type to make instantiating Player easier
type Resources struct {
	Brick, Grain, Lumber, Ore, Wool int
}

// Player represents a player's hand.
type Player struct {
	ID               PlayerID
	Resources        Resources
	DevelopmentCards struct {
		InHand, Played []DevelopmentCard
	}
	Settlements []Settlement
	Cities      []City
}

// Board is the general wrapper type for game state for a board. Scores
// and actions can be derived from this. Actions operate on a board (are
// 'reducers').
type Board struct {
	Robber           Hex
	Players          []Player
	Tiles            []Tile
	DevelopmentCards []DevelopmentCard
}

// FindPlayer returns player with given ID.
func (b Board) FindPlayer(playerID PlayerID) (Player, bool) {
	for _, p := range b.Players {
		if p.ID == playerID {
			return p, true
		}
	}

	return Player{}, false
}
