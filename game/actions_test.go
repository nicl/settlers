package game

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestDoRoll(t *testing.T) {
	roll := Roll{1, 3}
	board := Board{
		Robber: Hex{0, -2},
		Players: []Player{
			{
				ID: 1,
				Cities: []City{
					{A: Hex{0, 0}, B: Hex{0, 1}, C: Hex{-1, 1}},
				},
			},
			{
				ID: 2,
				Settlements: []Settlement{
					{A: Hex{0, 0}, B: Hex{0, -1}, C: Hex{1, -1}},
					{A: Hex{0, 2}, B: Hex{0, 3}, C: Hex{1, 3}},
				},
			},
		},
		Tiles: []Tile{{Location: Hex{0, 0}, Resource: Brick, Number: 4}},
	}

	got := DoRoll(board, roll)
	want := []Effect{
		AddResources{PlayerID: 1, Resources: map[Resource]int{Brick: 2}},
		AddResources{PlayerID: 2, Resources: map[Resource]int{Brick: 1}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestDoRob(t *testing.T) {
	robberID := PlayerID(1)
	victimID := PlayerID(2)

	board := Board{
		Robber: Hex{0, -2},
		Players: []Player{
			{
				ID:        robberID,
				Resources: Resources{Ore: 3},
			},
			{
				ID:        victimID,
				Resources: Resources{Brick: 2, Grain: 4},
			},
		},
		Tiles: []Tile{{Location: Hex{0, 0}, Resource: Brick, Number: 4}},
	}

	rob := Rob{Robber: robberID, Victim: victimID}
	rand.Seed(1) // make test predictable

	got := DoRob(board, rob)
	want := []Effect{
		RemoveResources{PlayerID: victimID, Resources: map[Resource]int{Brick: 1}},
		AddResources{PlayerID: robberID, Resources: map[Resource]int{Brick: 1}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}

//

func TestDoTrade(t *testing.T) {
	partyID := PlayerID(1)
	counterpartyID := PlayerID(2)

	board := Board{
		Players: []Player{
			{
				ID:        partyID,
				Resources: Resources{Ore: 3},
			},
			{
				ID:        counterpartyID,
				Resources: Resources{Brick: 2, Grain: 4},
			},
		},
	}

	trade := Trade{
		Party:        Offer{PlayerID: partyID, Resources: map[Resource]int{Brick: 2, Grain: 3}},
		Counterparty: Offer{PlayerID: counterpartyID, Resources: map[Resource]int{Ore: 3}},
	}

	got := DoTrade(board, trade)

	want := []Effect{
		AddResources{PlayerID: partyID, Resources: map[Resource]int{Ore: 3}},
		RemoveResources{PlayerID: partyID, Resources: map[Resource]int{Brick: 2, Grain: 3}},

		AddResources{PlayerID: counterpartyID, Resources: map[Resource]int{Brick: 2, Grain: 3}},
		RemoveResources{PlayerID: counterpartyID, Resources: map[Resource]int{Ore: 3}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestDoBuyDevCard(t *testing.T) {
	buyer := PlayerID(1)
	card := VictoryPoint

	board := Board{
		Players: []Player{
			{
				ID:        buyer,
				Resources: Resources{Grain: 1, Wool: 1, Ore: 3},
			},
		},
	}

	buy := BuyDevCard{
		PlayerID: buyer,
		Card:     card,
	}

	got := DoBuyDevCard(board, buy)

	want := []Effect{
		AddDevCard{PlayerID: buyer, Card: card},
		RemoveResources{PlayerID: buy.PlayerID, Resources: map[Resource]int{Grain: 1, Wool: 1, Ore: 1}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestDoBuildRoad(t *testing.T) {
	buyer := PlayerID(1)
	road := Road{A: Hex{0, 0}, B: Hex{0, 1}}

	board := Board{
		Players: []Player{
			{ID: buyer, Resources: Resources{Brick: 1, Lumber: 1}},
		},
	}

	got := DoBuildRoad(board, BuildRoad{PlayerID: buyer, Road: road})

	want := []Effect{
		AddRoad{PlayerID: buyer, Road: road},
		RemoveResources{PlayerID: buyer, Resources: map[Resource]int{Brick: 1, Lumber: 1}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestDoBuildSettlement(t *testing.T) {
	buyer := PlayerID(1)
	settlement := Settlement{A: Hex{0, 0}, B: Hex{0, 1}, C: Hex{1, 1}}

	board := Board{
		Players: []Player{
			{ID: buyer, Resources: Resources{Brick: 1, Lumber: 1, Grain: 1, Wool: 2}},
		},
	}

	got := DoBuildSettlement(board, BuildSettlement{PlayerID: buyer, Settlement: settlement})

	want := []Effect{
		AddSettlement{PlayerID: buyer, Settlement: settlement},
		RemoveResources{PlayerID: buyer, Resources: map[Resource]int{Brick: 1, Lumber: 1, Grain: 1, Wool: 1}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}
