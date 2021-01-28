package game

import (
	"log"
	"math/rand"
)

// Effects

// AddResources represents a player gaining resources
type AddResources struct {
	PlayerID  PlayerID
	Resources map[Resource]int
}

// RemoveResources represents loss of resources
type RemoveResources struct {
	PlayerID  PlayerID
	Resources map[Resource]int
}

// AddDevCard represents a player gaining a dev card.
type AddDevCard struct {
	PlayerID PlayerID
	Card     DevelopmentCard
}

// AddRoad represents a player gaining a road.
type AddRoad struct {
	PlayerID PlayerID
	Road     Road
}

// AddSettlement represents a player gaining a settlement.
type AddSettlement struct {
	PlayerID   PlayerID
	Settlement Settlement
}

// AddCity represents a player gaining a city.
type AddCity struct {
	PlayerID PlayerID
	City     City
}

// Effect is the interface for all state changes that can then be
// reduced onto a board. They also can be stored and undone/replayed.
type Effect interface{ isEffect() }

func (e AddResources) isEffect()    {}
func (e RemoveResources) isEffect() {}
func (e AddDevCard) isEffect()      {}
func (e AddRoad) isEffect()         {}
func (e AddSettlement) isEffect()   {}
func (e AddCity) isEffect()         {}

// Actions

// Roll is a dice roll.
type Roll struct {
	A, B int
}

// Rob represents a crime.
type Rob struct {
	Robber, Victim PlayerID
}

// Offer is one-side of a trade.
type Offer struct {
	PlayerID  PlayerID
	Resources map[Resource]int
}

// Trade is an exchange of resources between two players.
type Trade struct {
	Party, Counterparty Offer
}

// BuyDevCard is a dev card purchase. (Card selection handled separately.)
type BuyDevCard struct {
	PlayerID PlayerID
	Card     DevelopmentCard
}

// BuildRoad represents a new road action.
type BuildRoad struct {
	PlayerID PlayerID
	Road     Road
}

// BuildSettlement represents a new settlement action.
type BuildSettlement struct {
	PlayerID   PlayerID
	Settlement Settlement
}

// BuildCity represents a new city action.
type BuildCity struct {
	PlayerID PlayerID
	City     City
}

// Action is something that a player does that will affect board state.
type Action interface{ isAction() }

func (a Roll) isAction()            {}
func (a Rob) isAction()             {}
func (a Trade) isAction()           {}
func (a BuyDevCard) isAction()      {}
func (a BuildRoad) isAction()       {}
func (a BuildSettlement) isAction() {}
func (a BuildCity) isAction()       {}

// Functions

func pay(playerID PlayerID, r Resource, n int) Effect {
	if n < 1 {
		log.Fatal("Attempted to pay less than 1!")
	}

	return AddResources{
		PlayerID:  playerID,
		Resources: map[Resource]int{r: n},
	}
}

func resourcesAsSlice(p Player) []Resource {
	var rs []Resource

	add := func(n int, r Resource) {
		for i := 0; i < n; i++ {
			rs = append(rs, r)
		}
	}

	add(p.Resources.Brick, Brick)
	add(p.Resources.Grain, Grain)
	add(p.Resources.Lumber, Lumber)
	add(p.Resources.Ore, Ore)
	add(p.Resources.Wool, Wool)

	return rs
}

// DoRoll returns effects for a roll action.
func DoRoll(b Board, r Roll) []Effect {
	var effects []Effect

	sum := r.A + r.B

	for _, t := range b.Tiles {
		if !(t.Number == sum) {
			break
		}

		for _, p := range b.Players {
			cities := FindAdjacentCities(t.Location, p.Cities)
			if len(cities) > 0 {
				effects = append(effects, pay(p.ID, t.Resource, len(cities)*2))
			}

			settlements := FindAdjacentSettlements(t.Location, p.Settlements)
			if len(settlements) > 0 {
				effects = append(effects, pay(p.ID, t.Resource, len(settlements)))
			}
		}
	}

	return effects
}

// DoRob returns effects for a rob action.
func DoRob(b Board, r Rob) []Effect {
	var effects []Effect

	for _, p := range b.Players {
		if p.ID == r.Victim {
			resources := resourcesAsSlice(p)
			choice := rand.Intn(len(resources) - 1)
			robbedResource := resources[choice]

			effects = append(effects,
				RemoveResources{PlayerID: r.Victim, Resources: map[Resource]int{robbedResource: 1}},
				AddResources{PlayerID: r.Robber, Resources: map[Resource]int{robbedResource: 1}},
			)

			break
		}
	}

	return effects
}

// DoTrade returns effects for a trade action.
func DoTrade(b Board, t Trade) []Effect {
	// TODO validate each side has the resources

	return []Effect{
		AddResources{PlayerID: t.Party.PlayerID, Resources: t.Counterparty.Resources},
		RemoveResources{PlayerID: t.Party.PlayerID, Resources: t.Party.Resources},

		AddResources{PlayerID: t.Counterparty.PlayerID, Resources: t.Party.Resources},
		RemoveResources{PlayerID: t.Counterparty.PlayerID, Resources: t.Counterparty.Resources},
	}
}

// DoBuyDevCard returns effects for a dev card purchase.
func DoBuyDevCard(b Board, buy BuyDevCard) []Effect {
	player, ok := b.FindPlayer(buy.PlayerID)
	if !ok {
		log.Fatalf("Invalid Player ID: %d", buy.PlayerID)
	}

	if player.Resources.Grain < 1 || player.Resources.Ore < 1 || player.Resources.Wool < 1 {
		log.Fatalf("Player has insufficient resources to buy dev card.")
	}

	return []Effect{
		AddDevCard{PlayerID: buy.PlayerID, Card: buy.Card},
		RemoveResources{PlayerID: buy.PlayerID, Resources: map[Resource]int{Grain: 1, Wool: 1, Ore: 1}},
	}
}

// DoBuildRoad returns effects for road purchase.
func DoBuildRoad(b Board, road BuildRoad) []Effect {
	player, ok := b.FindPlayer(road.PlayerID)
	if !ok {
		log.Fatalf("Invalid Player ID: %d", road.PlayerID)
	}

	if player.Resources.Lumber < 1 || player.Resources.Brick < 1 {
		log.Fatalf("Player has insufficient resources to buy road.")
	}

	return []Effect{
		AddRoad{PlayerID: road.PlayerID, Road: road.Road},
		RemoveResources{PlayerID: road.PlayerID, Resources: map[Resource]int{Lumber: 1, Brick: 1}},
	}
}

// DoBuildSettlement returns effects for settlement purchase.
func DoBuildSettlement(b Board, s BuildSettlement) []Effect {
	// TODO validate settlement hexes are adjacent

	player, ok := b.FindPlayer(s.PlayerID)
	if !ok {
		log.Fatalf("Invalid Player ID: %d", s.PlayerID)
	}

	if player.Resources.Lumber < 1 || player.Resources.Brick < 1 || player.Resources.Grain < 1 || player.Resources.Wool < 1 {
		log.Fatalf("Player has insufficient resources to buy settlement.")
	}

	return []Effect{
		AddSettlement{PlayerID: s.PlayerID, Settlement: s.Settlement},
		RemoveResources{PlayerID: s.PlayerID, Resources: map[Resource]int{Lumber: 1, Brick: 1, Grain: 1, Wool: 1}},
	}
}

// DoBuildCity returns effects for city purchase.
func DoBuildCity(b Board, c BuildCity) []Effect {
	player, ok := b.FindPlayer(c.PlayerID)
	if !ok {
		log.Fatalf("Invalid Player ID: %d", c.PlayerID)
	}

	if player.Resources.Ore < 3 || player.Resources.Grain < 2 {
		log.Fatalf("Player has insufficient resources to buy settlement.")
	}

	return []Effect{
		AddCity{PlayerID: c.PlayerID, City: c.City},
		RemoveResources{PlayerID: c.PlayerID, Resources: map[Resource]int{Ore: 3, Grain: 2}},
	}
}
