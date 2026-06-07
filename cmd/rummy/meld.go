package main

import (
	"cards/common"
	"fmt"
	"sort"
	"strings"
)

/*

  File:    manager.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

type MeldType int

const (
	MeldRun MeldType = iota
	MeldSet
)

type Meld struct {
	Type  MeldType
	Cards []common.Card
}

type MeldResult struct {
	Melds          []Meld
	Deadwood       []common.Card
	DeadwoodPoints int
}

// create string:
// [run1] [run2]... [set1] [set2]... [deadwood] (points)
func (mr MeldResult) toString() string {
	segs := make([]string, 0, len(mr.Melds)+2)

	// runs
	for _, m := range mr.Melds {
		if m.Type == MeldRun {
			segs = append(segs, fmt.Sprintf("[ %s ]", meldToString(m.Cards)))
		}
	}
	// sets
	for _, m := range mr.Melds {
		if m.Type == MeldSet {
			segs = append(segs, fmt.Sprintf("[ %s ]", meldToString(m.Cards)))
		}
	}
	// deadwood
	segs = append(segs, fmt.Sprintf("Deadwood: [ %s ]", meldToString(mr.Deadwood)))
	// points
	segs = append(segs, fmt.Sprintf("Points (%d)", mr.DeadwoodPoints))

	return strings.Join(segs, " ")
}

func meldToString(cards []common.Card) string {
	s := make([]string, 0)
	for _, c := range cards {
		s = append(s, fmt.Sprintf("%3s", c.Icon))
	}
	return strings.Join(s, ", ")
}

// find all the
func computeMelds(hand []common.Card) MeldResult {

	allMelds := append(findSets(hand), findRuns(hand)...)

	best := MeldResult{
		Melds:          nil,
		Deadwood:       copyCards(hand),
		DeadwoodPoints: totalDeadwood(hand),
	}

	n := len(allMelds)

	// depth first search
	var dfs func(idx int, used []common.Card, chosen []Meld)
	dfs = func(idx int, used []common.Card, chosen []Meld) {
		if idx == n {
			// compute deadwood
			dead := remainingDeadwood(hand, used)
			points := totalDeadwood(dead)
			if points < best.DeadwoodPoints {
				best.Melds = copyMelds(chosen)
				best.Deadwood = dead
				best.DeadwoodPoints = points
			}
			return
		}

		// skip this meld
		dfs(idx+1, used, chosen)

		// include this meld if no overlap
		m := allMelds[idx]
		if !overlaps(used, m.Cards) {
			newUsed := append(copyCards(used), m.Cards...)
			newChosen := append(copyMelds(chosen), m)
			dfs(idx+1, newUsed, newChosen)
		}
	}

	dfs(0, nil, nil)
	return best
}

func findSets(hand []common.Card) []Meld {
	byRank := map[common.Rank][]common.Card{}
	for _, card := range hand {
		byRank[card.Rank] = append(byRank[card.Rank], card)
	}
	var melds []Meld
	for _, group := range byRank {
		if len(group) >= 3 {
			melds = append(melds, Meld{
				Type:  MeldSet,
				Cards: copyCards(group),
			})
		}
	}
	return melds
}

func findRuns(hand []common.Card) []Meld {
	bySuit := map[common.Suit][]common.Card{}
	for _, c := range hand {
		bySuit[c.Suit] = append(bySuit[c.Suit], c)
	}

	var melds []Meld
	for _, group := range bySuit {
		if len(group) < 3 {
			continue
		}

		// Sort by Ace-low ordering
		sort.Slice(group, func(i, j int) bool {
			return AceLow(group[i].Rank) < AceLow(group[j].Rank)
		})

		start := 0
		for start < len(group) {
			end := start + 1
			// find longest consecutive segment
			for end < len(group) &&
				AceLow(group[end].Rank) == AceLow(group[end-1].Rank)+1 {
				end++
			}
			segLen := end - start
			if segLen >= 3 {
				// accept ALL contiguous subsequences of length >= 3
				for i := start; i <= end-3; i++ {
					for j := i + 3; j <= end; j++ {
						melds = append(melds,
							Meld{Type: MeldRun, Cards: copyCards(group[i:j])})
					}
				}
			}
			start = end
		}
	}

	return melds
}

func copyMelds(ms []Meld) []Meld {
	out := make([]Meld, len(ms))
	for i, m := range ms {
		out[i].Type = m.Type
		out[i].Cards = copyCards(m.Cards)
	}
	return out
}

func overlaps(used []common.Card, cards []common.Card) bool {
	for _, u := range used {
		for _, card := range cards {
			if sameCard(u, card) {
				return true
			}
		}
	}
	return false
}

func remainingDeadwood(hand, used []common.Card) []common.Card {
	dead := copyCards(hand)
	for _, u := range used {
		dead = removeCard(dead, u)
	}
	return dead
}

func totalDeadwood(cards []common.Card) int {
	sum := 0
	for _, card := range cards {
		sum += RankValue(card)
	}
	return sum
}

func copyCards(cards []common.Card) []common.Card {
	out := make([]common.Card, len(cards))
	copy(out, cards)
	return out
}

func sameCard(a, b common.Card) bool {
	return a.Suit == b.Suit && a.Rank == b.Rank
}

func removeCard(cards []common.Card, card common.Card) []common.Card {
	for i, x := range cards {
		if sameCard(x, card) {
			return append(cards[:i], cards[i+1:]...)
		}
	}
	return cards
}
