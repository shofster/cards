package main

import "cards/common"

/*

  File:    move.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

type Move struct {
	From      Location
	FromIndex int
	To        Location
	ToIndex   int
	Card      common.Card
	SeqLen    int // 1 for single card, len(seq) for sequences
}

func GenerateLegalMoves(gs *GameState) []Move {
	var moves []Move
	for col := 0; col < 8; col++ {
		stack := gs.CardColumns[col]
		if len(stack) == 0 {
			continue
		}

		// CardColumn → Foundation / FreeCell / CardColumn / EmptyColumn
		top := stack[len(stack)-1]
		// CardColumn → Foundation
		if canMoveToFoundation(gs, top) {
			moves = append(moves, Move{
				From:      CardColumn,
				FromIndex: col,
				To:        Foundation,
				ToIndex:   int(top.Suit),
				Card:      top,
				SeqLen:    1,
			})
		}
		// CardColumn → FreeCell
		for i := 0; i < 4; i++ {
			if gs.FreeCells[i] == nil {
				moves = append(moves, Move{
					From:      CardColumn,
					FromIndex: col,
					To:        FreeCell,
					ToIndex:   i,
					Card:      top,
					SeqLen:    1,
				})
				break
			}
		}
		// CardColumn → CardColumn (single card)
		for dest := 0; dest < 8; dest++ {
			if dest == col {
				continue
			}
			if canPlaceOnCardColumn(gs, top, dest) {
				moves = append(moves, Move{
					From:      CardColumn,
					FromIndex: col,
					To:        CardColumn,
					ToIndex:   dest,
					Card:      top,
					SeqLen:    1,
				})
			}
		}
		// CardColumn → Empty Column (single card)
		if top.Rank == common.King {
			for dest := 0; dest < 8; dest++ {
				if len(gs.CardColumns[dest]) == 0 && dest != col {
					moves = append(moves, Move{
						From:      CardColumn,
						FromIndex: col,
						To:        EmptyColumn,
						ToIndex:   dest,
						Card:      top,
						SeqLen:    1,
					})
				}
			}
		}
		// CardColumn → CardColumn (multi-card sequence)
		moves = append(moves, generateSequenceMoves(gs, col)...)
	}

	// FreeCell → Foundation / CardColumn / EmptyColumn
	for i := 0; i < 4; i++ {
		c := gs.FreeCells[i]
		if c == nil {
			continue
		}
		//FreeCell → Foundation
		if canMoveToFoundation(gs, *c) {
			moves = append(moves, Move{
				From:      FreeCell,
				FromIndex: i,
				To:        Foundation,
				ToIndex:   int(c.Suit),
				Card:      *c,
				SeqLen:    1,
			})
		}
		// FreeCell → CardColumn
		for dest := 0; dest < 8; dest++ {
			if canPlaceOnCardColumn(gs, *c, dest) {
				moves = append(moves, Move{
					From:      FreeCell,
					FromIndex: i,
					To:        CardColumn,
					ToIndex:   dest,
					Card:      *c,
					SeqLen:    1,
				})
			}
		}
		// FreeCell → Empty Column
		if c.Rank == common.King {
			for dest := 0; dest < 8; dest++ {
				if len(gs.CardColumns[dest]) == 0 {
					moves = append(moves, Move{
						From:      FreeCell,
						FromIndex: i,
						To:        EmptyColumn,
						ToIndex:   dest,
						Card:      *c,
						SeqLen:    1,
					})
				}
			}
		}
	}

	return moves
}

func canMoveToFoundation(gs *GameState, c common.Card) bool {
	needed := gs.Foundation[c.Suit] + 1
	return freeValue(c.Rank) == needed
}

func canPlaceOnCardColumn(gs *GameState, c common.Card, dest int) bool {
	col := gs.CardColumns[dest]
	if len(col) == 0 {
		return false
	}

	top := col[len(col)-1]
	if sameColor(top.Suit, c.Suit) {
		return false
	}
	return c.Rank+1 == top.Rank
}

func generateSequenceMoves(gs *GameState, from int) []Move {
	var moves []Move
	col := gs.CardColumns[from]
	if len(col) < 2 {
		return moves
	}

	// find longest descending alternating sequence
	start := len(col) - 1
	for start > 0 {
		a := col[start]
		b := col[start-1]
		if sameColor(a.Suit, b.Suit) || b.Rank+1 != a.Rank {
			break
		}
		start--
	}
	seq := col[start:]
	maxMoveable := maxMovable(gs)
	if len(seq) > maxMoveable {
		seq = seq[len(seq)-maxMoveable:]
		start = len(col) - len(seq)
	}

	// try placing sequence
	for dest := 0; dest < 8; dest++ {
		if dest == from {
			continue
		}
		if len(gs.CardColumns[dest]) == 0 {
			if seq[0].Rank == common.King {
				moves = append(moves, Move{
					From:      CardColumn,
					FromIndex: from,
					To:        EmptyColumn,
					ToIndex:   dest,
					Card:      seq[0],
					SeqLen:    len(seq),
				})
			}
			continue
		}
		top := gs.CardColumns[dest][len(gs.CardColumns[dest])-1]
		if !sameColor(top.Suit, seq[0].Suit) && seq[0].Rank+1 == top.Rank {
			moves = append(moves, Move{
				From:      CardColumn,
				FromIndex: from,
				To:        CardColumn,
				ToIndex:   dest,
				Card:      seq[0],
				SeqLen:    len(seq),
			})
		}
	}

	return moves
}

func maxMovable(gs *GameState) int {
	free := gs.FreeCellCount()
	empty := 0
	for _, col := range gs.CardColumns {
		if len(col) == 0 {
			empty++
		}
	}
	// FreeCell formula: (freeCells + 1) * 2^emptyColumns
	n := free + 1
	for i := 0; i < empty; i++ {
		n *= 2
	}
	return n
}

/* --------------- Effect Moves -------------- */

func ApplyMove(gs *GameState, mv Move) {
	var cards common.Cards
	// remove from source
	switch mv.From {
	case CardColumn:
		col := gs.CardColumns[mv.FromIndex]
		cards = append(cards, col[len(col)-mv.SeqLen:]...)
		gs.CardColumns[mv.FromIndex] = col[:len(col)-mv.SeqLen]
	case FreeCell:
		cards = append(cards, *gs.FreeCells[mv.FromIndex])
		gs.FreeCells[mv.FromIndex] = nil
	default:
	}
	// place at destination
	switch mv.To {
	case CardColumn, EmptyColumn:
		gs.CardColumns[mv.ToIndex] = append(gs.CardColumns[mv.ToIndex], cards...)
	case FreeCell:
		gs.FreeCells[mv.ToIndex] = &cards[0]
	case Foundation:
		gs.Foundation[mv.ToIndex] = freeValue(cards[0].Rank)
	}
}

/*// revert the last move
func UndoMove(gs *GameState, ue UndoEntry) {
	mv := ue.Move
	seq := ue.Seq
	// remove from destination
	switch mv.To {
	case CardColumn, EmptyColumn:
		col := gs.CardColumns[mv.ToIndex]
		gs.CardColumns[mv.ToIndex] = col[:len(col)-mv.SeqLen]
	case FreeCell:
		gs.FreeCells[mv.ToIndex] = nil
	case Foundation:
		gs.Foundation[mv.ToIndex] = ue.PrevFound
	}
	// place at source
	switch mv.From {
	case CardColumn:
		gs.CardColumns[mv.FromIndex] = append(gs.CardColumns[mv.FromIndex], seq...)
	case FreeCell:
		gs.FreeCells[mv.FromIndex] = &seq[0]
	default:
	}
}
*/
