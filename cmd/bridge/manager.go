package main

import (
	"games/common"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
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

// HandManager is the shared controller for the 13 tricks in a Deal.
type HandManager struct {
	cards          []common.Card   // logical cards (rank/suit)
	stack          *fyne.Container // 13 *CardImage widgets
	trickSlot      *canvas.Image   // where the played card is shown
	commitButton   *widget.Button  // enabled after a selection is made
	trickSelection int             // -1 means none selected this trick
	active         int             // the active position
	// "virtual functions"
	deal   func()                         // function to set up the position cards
	renig  func(common.Suit)              // function to eliminate invalid cards
	play   func(common.Suit, common.Suit) // function to play a card for this trick
	reveal func()                         // function to show all Cards
}

func (h *HandManager) Deal() {

	// per - position setup
	h.deal()
	for i := 0; i < 13; i++ {
		h.cards[i].Index = i
	}
	h.trickSelection = -1
}

func (h *HandManager) Play(position, leadPlayer Player, contractSuit, leadSuit common.Suit) {

	//fmt.Printf("Play - player is %s, lead is %s. Suits - contract is %s, lead suit is %s \n",
	//	playerToPosition(position), playerToPosition(leadPlayer),
	//	contractSuit, common.SuitCode(leadSuit))

	if position == leadPlayer {
		// Leader: all unplayed cards are legal
		for i := 0; i < 13; i++ {
			ci := h.stack.Objects[i].(*common.CardImage)
			if h.cards[i].Selected {
				ci.SetDimmed(true)
			} else {
				ci.SetDimmed(false)
			}
		}
	} else {
		// eliminate invalid cards to play on this trick
		// and dim out
		h.renig(leadSuit)
	}
	// initial states
	h.play(contractSuit, leadSuit)

}

func (h *HandManager) RevealAll() {
	// per - position reveal the cards
	h.reveal()
}

func (h *HandManager) setCardFace(idx int, card common.Card) *common.CardImage {
	cardImage := h.stack.Objects[idx].(*common.CardImage)
	cardImage.Img.Resource = card.Resource
	fyne.Do(func() {
		cardImage.Img.Refresh()
	})
	return cardImage
}

func (h *HandManager) CommitSelection() {

	//fmt.Printf("COMMIT: trickSelection=%d (%s)\n",
	//	h.trickSelection,
	//	common.Icon(h.cards[h.trickSelection].Rank,
	//		h.cards[h.trickSelection].Suit))

	if h.trickSelection == -1 {
		return // no card chosen
	}
	h.cards[h.trickSelection].Selected = true
	h.active = -1
	renigCommon(h, common.NoSuit)
}

//
// ------------- Position Managers ------------------------ //
//

func NewSouthHandManager(cards []common.Card) *HandManager {
	h := &HandManager{
		cards:        cards,
		stack:        southStack,
		trickSlot:    southTrick,
		commitButton: southCommit,
	}

	h.deal = func() {
		for i := 0; i < 13; i++ {
			ci := h.setCardFace(i, cards[i])
			ci.Enabled = false
			ci.SetDimmed(false)
			ci.OnTap = func(tapped *common.CardImage) {
				if h.active == 1 {

					//fmt.Printf("TAP SOUTH: trickSelection=%d → %d (%s)\n",
					//	h.trickSelection, tapped.ID,
					//	common.Icon(h.cards[tapped.ID].Rank,
					//		h.cards[tapped.ID].Suit))

					if h.trickSelection != -1 {
						previous := h.stack.Objects[h.trickSelection].(*common.CardImage)
						previous.SetDimmed(false)
					}
					h.trickSelection = tapped.ID
					// Dim the new card
					tapped.SetDimmed(true)
					// Show in trick slot
					h.trickSlot.Resource = tapped.Img.Resource
					h.trickSlot.Refresh()
					h.commitButton.Enable()
				}
			}

		}
	}

	h.renig = func(leadSuit common.Suit) {
		renigCommon(h, leadSuit)
	}

	h.play = func(_, _ common.Suit) { h.active = 1 }

	h.reveal = func() {
		for i := 0; i < 13; i++ {
			h.setCardFace(i, cards[i])
		}
	}

	return h
}

func NewNorthHandManager(cards []common.Card) *HandManager {
	h := &HandManager{
		cards:        cards,
		stack:        northStack,
		trickSlot:    northTrick,
		commitButton: northCommit,
	}

	h.deal = func() {
		for i := 0; i < 13; i++ {
			ci := h.setCardFace(i, cards[i])
			ci.Enabled = false
			ci.SetDimmed(false)
			ci.OnTap = func(tapped *common.CardImage) {
				if h.active == 2 {

					//fmt.Printf("TAP NORTH: trickSelection=%d → %d (%s)\n",
					//	h.trickSelection, tapped.ID,
					//	common.Icon(h.cards[tapped.ID].Rank,
					//		h.cards[tapped.ID].Suit))

					if h.trickSelection != -1 {
						previous := h.stack.Objects[h.trickSelection].(*common.CardImage)
						previous.SetDimmed(false)
					}
					h.trickSelection = tapped.ID
					// Dim the new card
					tapped.SetDimmed(true)
					// Show in trick slot
					h.trickSlot.Resource = tapped.Img.Resource
					h.trickSlot.Refresh()
					h.commitButton.Enable()
				}
			}

		}
	}

	h.renig = func(leadSuit common.Suit) {
		renigCommon(h, leadSuit)
	}

	h.play = func(_, _ common.Suit) { h.active = 2 }

	h.reveal = func() {
		for i := 0; i < 13; i++ {
			h.setCardFace(i, cards[i])
		}
	}

	return h
}

/*

	renig - Allow Tappable on available cards
        and Disallow & dim the unavailable.

// disable any previously used card
// on the opening lead, any card is OK
// not opening lead and card is of that suit, it's OK
// has lead suit but this card is other -- it's NOT OK
// if NO lead suit, any card is available

*/

func renigCommon(h *HandManager, leadSuit common.Suit) {

	hasLeadSuit := false
	for i := 0; i < 13; i++ {
		if !h.cards[i].Selected && h.cards[i].Suit == leadSuit {
			hasLeadSuit = true
			break
		}
	}

	for i := 0; i < 13; i++ {
		ci := h.stack.Objects[i].(*common.CardImage)
		switch {

		case h.cards[i].Selected:
			// disable any previously used card
			ci.SetDimmed(true)
			ci.Enabled = false
			ci.Replace(cardMap["reverse"].Resource)

		case leadSuit == common.NoSuit:
			// on opening lead, any card is OK
			ci.SetDimmed(false)
			ci.Enabled = true

		case hasLeadSuit && h.cards[i].Suit == leadSuit:
			// follow suit
			ci.SetDimmed(false)
			ci.Enabled = true

		case hasLeadSuit:
			// has lead suit but this card is other -- it's NOT OK
			ci.SetDimmed(true)
			ci.Enabled = false

		default:
			// if NO lead suit, any card is available
			ci.SetDimmed(false)
			ci.Enabled = true
		}
	}
}

func NewEastHandManager(cards []common.Card) *HandManager {
	h := &HandManager{
		cards:        cards,
		stack:        eastStack,
		trickSlot:    eastTrick,
		commitButton: eastCommit,
	}

	h.deal = func() {
		for i := 0; i < 13; i++ {
			ci := h.stack.Objects[i].(*common.CardImage)
			if common.BoolGet(showEast) {
				ci.SetDimmed(false)
				h.setCardFace(i, cards[i])
			} else {
				h.setCardFace(i, cardMap["reverse"])
			}
		}
	}

	// East doesn't allow "Enable" for TAP
	h.renig = func(leadSuit common.Suit) {
		// but, only if east cards are visible
		if common.BoolGet(showEast) {
			renigCommon(h, leadSuit)
		}
	}

	h.play = func(contractSuit, leadSuit common.Suit) {
		//fmt.Printf("play East - lead is %s \n", common.SuitCode(leadSuit))
		if leadSuit == common.NoSuit {
			// opening lead logic
			h.autoLead(contractSuit)

		} else {
			h.autoPlay(contractSuit, leadSuit)
		}

		// Dim the played card
		ci := h.stack.Objects[h.trickSelection].(*common.CardImage)
		if common.BoolGet(showEast) {
			ci.SetDimmed(true)
		}
		// Show in trick slot
		h.trickSlot.Resource = h.cards[h.trickSelection].Resource
		fyne.Do(func() {
			h.trickSlot.Refresh()
		})
		// and "tap" it
		h.commitButton.OnTapped()

	}

	h.reveal = func() {
		if !common.BoolGet(showEast) {
			for i := 0; i < 13; i++ {
				h.setCardFace(i, cards[i])
			}
		}
	}

	return h
}

func NewWestHandManager(cards []common.Card) *HandManager {
	h := &HandManager{
		cards:        cards,
		stack:        westStack,
		trickSlot:    westTrick,
		commitButton: westCommit,
	}

	h.deal = func() {
		for i := 0; i < 13; i++ {
			ci := h.stack.Objects[i].(*common.CardImage)
			if common.BoolGet(showWest) {
				ci.SetDimmed(false)
				h.setCardFace(i, cards[i])
			} else {
				h.setCardFace(i, cardMap["reverse"])
			}
		}
	}

	// West doesn't allow "Enable" for TAP
	h.renig = func(leadSuit common.Suit) {
		// but, only is west cards are visible
		if common.BoolGet(showWest) {
			renigCommon(h, leadSuit)
		}
	}

	h.play = func(contractSuit, leadSuit common.Suit) {
		//fmt.Printf("play West - lead is %s \n", common.SuitCode(leadSuit))
		if leadSuit == common.NoSuit {
			// opening lead logic
			h.autoLead(contractSuit)

		} else {
			h.autoPlay(contractSuit, leadSuit)
		}

		// Dim the played card
		ci := h.stack.Objects[h.trickSelection].(*common.CardImage)
		if common.BoolGet(showEast) {
			ci.SetDimmed(true)
		}
		// Show in trick slot
		h.trickSlot.Resource = h.cards[h.trickSelection].Resource
		fyne.Do(func() {
			h.trickSlot.Refresh()
		})
		// and "tap" it
		h.commitButton.OnTapped()

	}

	h.reveal = func() {
		if !common.BoolGet(showWest) {
			for i := 0; i < 13; i++ {
				h.setCardFace(i, cards[i])
			}
		}
	}

	// at end of 13 tricks, show the unused cards
	h.reveal = func() {
		for i := 0; i < 13; i++ {
			h.setCardFace(i, cards[i])
		}
	}
	return h
}

/*

	EAST/WEST auto play --  FOLLOWING  ---

	East/West has suit and a card that can win,
		play the smallest card that still wins.

	East/West has suit but cannot win,
        play the smallest card in that suit.

	East/West has no led suit cards,
		play the smallest trump, if the contract has a trump suit.

	East/West has no lead or trump
		play the smallest card in the hand.


*/

// contractSuit = contract suit (or NoTrump)
// leadSuit = suit led to the trick
func (h *HandManager) autoPlay(contractSuit, leadSuit common.Suit) {

	// find highest lead-suit rank already played
	var highestLeadRank common.Rank = -1

	// by South
	if southManager.trickSelection != -1 {
		c := southManager.cards[southManager.trickSelection]
		if c.Suit == leadSuit && c.Rank > highestLeadRank {
			highestLeadRank = c.Rank
		}
	}

	// and / or North
	if northManager.trickSelection != -1 {
		c := northManager.cards[northManager.trickSelection]
		if c.Suit == leadSuit && c.Rank > highestLeadRank {
			highestLeadRank = c.Rank
		}
	}

	// follow the lead suit
	lowestIdx := -1
	coverIdx := -1
	var coverRank common.Rank = 99

	for idx := 0; idx < len(h.cards); idx++ {
		c := h.cards[idx]
		if c.Selected || c.Suit != leadSuit {
			continue
		}

		// lowest card of the suit
		if lowestIdx == -1 || c.Rank < h.cards[lowestIdx].Rank {
			lowestIdx = idx
		}

		// honor-cover logic
		if isHonor(highestLeadRank) &&
			isHonor(c.Rank) &&
			c.Rank > highestLeadRank &&
			c.Rank < coverRank {
			coverRank = c.Rank
			coverIdx = idx
		}
	}

	if coverIdx != -1 {
		h.trickSelection = coverIdx
		return
	}

	// third hand high - use the lowest winning card
	winningIdx := -1
	var winningRank common.Rank = 99

	for idx := 0; idx < len(h.cards); idx++ {
		c := h.cards[idx]
		if c.Selected || c.Suit != leadSuit {
			continue
		}

		if c.Rank > highestLeadRank && c.Rank < winningRank {
			winningRank = c.Rank
			winningIdx = idx
		}
	}

	if winningIdx != -1 {
		h.trickSelection = winningIdx
		return
	}

	// lowest losing card
	if lowestIdx != -1 {
		h.trickSelection = lowestIdx
		return
	}

	// void in the lead duit, trump ?
	if contractSuit != common.NoTrump {
		lowestTrump := -1
		for idx := 0; idx < len(h.cards); idx++ {
			c := h.cards[idx]
			if !c.Selected && c.Suit == contractSuit {
				if lowestTrump == -1 || c.Rank < h.cards[lowestTrump].Rank {
					lowestTrump = idx
				}
			}
		}
		if lowestTrump != -1 {
			h.trickSelection = lowestTrump
			return
		}
	}

	// pitch the lowest card
	lowestAny := -1
	for idx := 0; idx < len(h.cards); idx++ {
		c := h.cards[idx]
		if !c.Selected {
			if lowestAny == -1 || c.Rank < h.cards[lowestAny].Rank {
				lowestAny = idx
			}
		}
	}

	h.trickSelection = lowestAny
}

func isHonor(rank common.Rank) bool {
	return rank >= common.Jack
}

/*

	Auto Play --  LEADING  --

	No Trump
     Lead the longest suit, 4th‑best, unless a strong sequence
       — then lead the top of the sequence.

    Suit contract
     Lead top of sequence, singleton, doubleton,
      passive low card from a safe suit
       - avoid leading away from an Ace.

*/

// Opening lead:  defined by leadSuit == NoSuit
func (h *HandManager) autoLead(contractSuit common.Suit) {

	// Group cards by suit
	suits := map[common.Suit][]*common.Card{}
	for i := 0; i < 13; i++ {
		if !h.cards[i].Selected {
			suits[h.cards[i].Suit] = append(suits[h.cards[i].Suit], &h.cards[i])
		}
	}

	// find sequences (KQJ, QJ10, J109, KQ, QJ, J10)
	hasSeq := func(cards []*common.Card) *common.Card {
		if len(cards) < 2 {
			return nil
		}
		// Sort high→low
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Rank >
				cards[j].Rank
		})
		for i := 0; i < len(cards)-1; i++ {
			r1 := cards[i].Rank
			r2 := cards[i+1].Rank
			if r1 == r2+1 { // touching honors
				return cards[i] // top of sequence
			}
		}
		return nil
	}

	// lowest card in a suit
	lowest := func(cards []*common.Card) *common.Card {
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Rank <
				cards[j].Rank
		})
		return cards[0]
	}

	// 4th best (or lowest if <4)
	fourthBest := func(cards []*common.Card) *common.Card {
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Rank <
				cards[j].Rank
		})
		if len(cards) >= 4 {
			return cards[len(cards)-4]
		}
		return cards[0]
	}

	//
	// No Trump
	//
	if contractSuit == common.NoTrump {
		// Find longest suit
		var bestSuit common.Suit
		bestLen := -1
		for s, cs := range suits {
			if len(cs) > bestLen {
				bestLen = len(cs)
				bestSuit = s
			}
		}

		cards := suits[bestSuit]

		// Sequence?
		if top := hasSeq(cards); top != nil {
			h.trickSelection = top.Index
			return
		}

		// 4th best
		chosen := fourthBest(cards)
		h.trickSelection = chosen.Index
		return
	}

	//
	// Suit
	//
	trump := contractSuit

	// top of sequence
	for s, cs := range suits {
		if s == trump {
			continue
		}
		if top := hasSeq(cs); top != nil {
			h.trickSelection = top.Index
			return
		}
	}

	// singleton (non‑trump)
	for s, cs := range suits {
		if s == trump {
			continue
		}
		if len(cs) == 1 {
			h.trickSelection = cs[0].Index
			return
		}
	}

	// doubleton (non‑trump)
	for s, cs := range suits {
		if s == trump {
			continue
		}
		if len(cs) == 2 {
			// Lead top of doubleton
			sort.Slice(cs, func(i, j int) bool {
				return cs[i].Rank > cs[j].Rank
			})
			h.trickSelection = cs[0].Index
			return
		}
	}

	// lowest card from longest non‑trump suit
	var bestSuit common.Suit
	bestLen := -1
	for s, cs := range suits {
		if s == trump {
			continue
		}
		if len(cs) > bestLen {
			bestLen = len(cs)
			bestSuit = s
		}
	}

	cards := suits[bestSuit]
	chosen := lowest(cards)
	h.trickSelection = chosen.Index
	return

}
