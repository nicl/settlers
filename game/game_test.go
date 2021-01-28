package game

import (
	"testing"
)

func TestIsNeighbour(t *testing.T) {
	hex := Hex{Row: 0, Column: 0}
	candidates := map[Hex]bool{
		{0, 1}:  true,
		{1, 0}:  true,
		{0, -1}: true,
		{1, -1}: true,
		{-1, 1}: true,
		{-1, 0}: true,

		{1, 1}:   false,
		{0, 0}:   false,
		{-1, -1}: false,
		{-2, 2}:  false,
	}

	for h, want := range candidates {
		if got := hex.IsNeighbour(h); got != want {
			t.Errorf("got %t; want %t. For %v", got, want, h)
		}
	}
}

func TestFindSharedNeighbour(t *testing.T) {
	b := Hex{Column: 0, Row: 0}
	a := Hex{Column: -1, Row: 0}
	want := Hex{Column: -1, Row: 1}

	candidates := []Hex{
		{-1, -1},
		{-1, 1},
		{-1, 0},
		{0, 0},
		{-1, -0},
		{-2, 2},
	}

	got, ok := FindSharedNeighbour(a, b, candidates)

	if !ok || got != want {
		t.Errorf("got %v; want %v.", got, want)
	}
}
