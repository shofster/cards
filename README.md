## Game Summary

This project includes six single‑player card games:

- **Jacks or Better** – Classic video‑poker style game with hold/draw play. Two outcomes only: player win or dealer win.
- **Stud Duel** – 5‑card stud with one‑card reveal streets and optional betting.
- **Texas Hold’em** – Simplified Hold’em duel with configurable street betting.
- **Jokers Wild** – Points‑based duel with special scoring rules shown in Settings.
- **Bridge** – South is the bidder and declarer; North is partner; East/West are AI defenders.
- **Rummy** – Single‑player Gin/Rummy with an AI opponent.
- **Cribbage** – Single‑player Cribbage with an AI opponent.
- **Freecell** - Freecell Solitaire.


All functions are covered by a MIT license and free for use.

---

## GAMES

---

### ♠ ♣ ♥ ♦  —  ♠ ♣ ♥ ♦

## 1. (Jacks) Jacks or Better – "Video Poker"

A video‑poker style game.  
The player receives 5 cards, chooses which to discard, and draws replacements.   
The final hand is evaluated using standard Jacks or Better rankings.<br>
Changing payout table ("8/5" or "9/6") and bet amount are available in the "setting" tab.

### Gameplay

- Deal: Player antes.
- Player receives 5 cards.
- Player may discard one to five cards and draw new ones.
- The hand is evaluated using standard Jacks or Better rules.

### Payouts

**For all hands:**
Jacks or Better uses a fixed paytable with a 4000 payout for a
5‑bet Royal Flush. Other hands scale (mostly) linearly with the bet amount.

The PayoutTable may be either "8/5" or "9/6" (Full House/Flush)


### ♠ ♣ ♥ ♦  —  ♠ ♣ ♥ ♦

## 2. (Stud5) 5‑Card Stud

A simple 5‑card stud duel with optional betting on each reveal street.  
Ante, bet amounts, and wallet adjustments are available in **Settings**.<br>
Dealer qualification (Ace‑King high or better) is manually configurable.

### Gameplay

- Deal: Player antes.
- Player and dealer each receive 5 cards face‑down. 
- Cards are revealed one at a time (Street 1 through Street 5).
- Optional betting may be enabled or disabled via manual configuration.
- After all 5 cards are revealed, the better 5‑card hand wins.


### Payouts

Payouts depend on qualification:

**If dealer does NOT qualify:**
payout = ante * 2 ; Bets push

**If dealer qualifies and player wins:**
payout = (ante + Bets) * 2

**If dealer qualifies and player loses:**
payout = 0

**Tie:**
Push. All wagers (ante + bets) are returned to the wallet.



### ♠ ♣ ♥ ♦  —  ♠ ♣ ♥ ♦

## 3. (Holdem) “Texas” Hold’em

A simplified Texas Hold’em variant with a fixed ante and optional fixed bet on each street.  
Betting on Streets 1–3 is configurable.  
Ante, bet amounts, and wallet adjustments are available in **Settings**.

### Gameplay

- Deal: Player antes.
- Player receives 2 cards.
- Dealer receives 2 cards (face‑down).
- Board is revealed in three streets:
    - Street 1: Flop (3 cards)
    - Street 2: Turn (1 card)
    - Street 3: River (1 card)
- Optional betting may be enabled or disabled via manual configuration.
- Showdown: final bet is enabled.
- Hands are evaluated using standard Texas Hold’em rules.

### Payouts

**Dealer wins**  
All wagers (ante + bets) have already been deducted from the wallet. No additional action is taken.

**Player wins**  
The player receives:  
2 × (ante + bets)

**Tie**  
Push. All wagers (ante + bets) are returned to the wallet.


### ♠ ♣ ♥ ♦  —  ♠ ♣ ♥ ♦


## 4. (Joker) Jokers Wild

A 5‑card points duel with special scoring rules.  
Ante, bet amounts, and wallet adjustments are available in **Settings**.  
All scoring rules (including special card effects such as the Seven and Joker) are displayed in the **Settings** panel.
**Settings** displays a cumulative total of "Dealer" points and "Player" points.

### Gameplay

- Deal: Player antes.
- Player and dealer each receive 5 cards face‑down.
- Cards are revealed one at a time (Street 1 through Street 5).
- Optional betting may be enabled or disabled via configuration.
- Each revealed card updates the running point total according to the scoring rules shown in **Settings**.
- After all 5 cards are revealed, the higher total wins.

### Payouts

**Dealer wins**  
All wagers (ante + bets) have already been deducted from the wallet. No additional action is taken.

**Player wins**  
The player receives:  
2 × (ante + bets)

**Tie**  
Push. All wagers (ante + bets) are returned to the wallet.

### ♠ ♣ ♥ ♦  —  ♠ ♣ ♥ ♦

## 5. (Bridge) Single Player vs 2 AI Opponents Bridge

South is always the declarer.  
North is partner (dummy).  
East and West are AI defenders.  
The contract is set by South’s bid; play proceeds clockwise through 13 tricks.

### Gameplay

- South selects the contract.
- South always makes the opening lead (standard Bridge uses West).
- North’s hand becomes dummy after the opening lead.
- South plays both South and North hands.
- East/West defend and lead automatically when they win a trick.
- Play follows standard Bridge rules.
- Deal may be selected to start a new hand.

### Scoring

**For all contracts:**  
Score is based on contract level, suit (or notrump), and number of tricks taken.  
Overtricks and undertricks follow standard duplicate‑style scoring.
Player: below‑the‑line + above‑the‑line.  
Opponents: above‑the‑line only (no contract points).
Vulnerability is not used.

### ♠ ♣ ♥ ♦  —  ♠ ♣ ♥ ♦

## 6. (Rummy) Single Player vs AI Gin/Rummy

A single‑player Gin/Rummy variant.  
The player discards and draws forming melds; the AI defender plays a simplified opponent strategy.  
Deadwood count determines knock and gin outcomes.

### Gameplay

- Deal: Each side receives 10 cards.
- Player selects a card for discard.
- Player draws from either the stock or discard pile.
- Melds may be runs or sets; deadwood is the sum of the unmatched.
- Player may knock at any time - gin at 0.
- Lay off (playing on opponent's melds) is automatic for knocks - not gin.
- Dealer knocks when deadwood <= 10 or gin.

### Scoring

**For all hands:**
- **Knock:** The knocker scores the difference in deadwood (knockee − knocker).
  If the knockee has lower deadwood, this becomes an undercut.
- **Gin:** The ginner scores opponent deadwood + 25 bonus points.
- **Undercut:** If the opponent has lower deadwood, the opponent scores the difference in deadwood + 25 bonus points.


### ♠ ♣ ♥ ♦  —  ♠ ♣ ♥ ♦

## 7. (Cribbage) Single Player vs AI Cribbage

A single‑player Cribbage variant.
The PLAYER is the logged-in user.
AI is Sir John Suckling.
A 2 track spiral peg board (121 holes, each) displays scores.

### Gameplay

- The first Dealer is the lowest card in a "deck cut".
- Dealer alternates on each hand.
- Deal: Each side receives 6 cards.
- Each sides selects 2 cards for discard to the (hidden) Crib.
- A single card (the starter) is chosen (and shown).
- The non-dealer plays a card, followed by the dealer.
- Cards are alternately player and scored (Pegging)
- Cards are played and scored until the count is >= 31.
- Scored and Total points are visually displayed for PLAYER and AI.
- The hand is complete (and scored) when neither can play.
- The hand(s) are scored including the starter card,
with the dealer also getting the Crib's hand score.
- A new hand is dealt until the first to 121 pegs.



### Scoring

**Pegging Scoring:**
- 2 Points for count of 15.
- 2 Points for count of 31.
- 2 or more for pairs.
- 3 or more for runs.
- 1 Point for card played after Opponent is unable.

**Hand Scoring:**
Score hands in strict order:
1. NON-DEALER hand
2. DEALER hand
3. CRIB (belongs to current dealer)

For each hand:
- Count 15s
- Count pairs
- Count runs
- Flush (4 or 5)
- His Nobs (a Jack matching starter suit)



### ♠ ♣ ♥ ♦  —  ♠ ♣ ♥ ♦

## 8. (Freecell) Patience / Solitaire

A single‑player Freecell game.
Objective is to gather consecutive cards from Ace to King of the same suit.
The board consists of three areas - Card Columns, Free, and Foundation.
- Card Columns (8 card columns) are all of the cards not in the Free or Foundation areas.
- The Free area (4 cells) is a temporary "parking" spot to free up buried
cards in the Card Columns.
- The Foundation piles (4 suits - ♣,♦,♥,♠, each in ascending rank).
The game is "won" when all 4 suits are complete.

A "sequence" is 1 or more cards in descending rank and alternating colors -
clubs and spades are black, diamonds and hearts are red.


### Gameplay

- Clicking the "DEAL" button presents a new board.
- All 52 cards are dealt face up. 7 in the first 4 columns and 6 in the last 4 columns.
- Clicking the "UNDO" button reverses the last valid move.
- Setting the "Hints" to "on" provides a useful hint for the next move.
- Cards are moved via drag-and-drop:

Move the mouse to the card to select a card or sequence.
Pressing the Primary mouse button selects the cell.
Selecting the first of a sequence selects the entire sequence.

Moving the mouse while holding the button down drags the card's icon (i.e. 10♠).

While moving, any valid destination cell is highlighted and any invalid destination
 is dimmed. Locations outside cards are not shown.

Releasing the mouse button ends the move. If it is in a valid location, the card(s)
 are removed from the original location and placed on the destination.

**Specific Rules:**
- Card Columns must be built with descending rank and alternating colors.
- Foundations are built ascending by suit.
- Any Free card or top card of any card column may be moved to
  build on another column, moved to an empty Free card, or its Foundation.
- Any sequence from a card column may be moved to another card column whose top has a rank
one greater and opposite color as the rank/color of that sequence.
- Only a single card may be moved from the top of a Card Column to an
empty Free or Foundation cell.
- Only a King (or King sequence) may be placed on an empty Card Column.


--------------------------------------------------------------

## Tech Stuff

Minumum GO version is 1.22 -- requires "math/rand/v2"

Builds scripts are provided for Windows, MacOS, and LINUX OS'es
 for both AMD and ARM architectures.

Each app is a standalone Go executable:
No external dependencies beyond Fyne.IO.
Shared components:
- Card (rank, suit, resource)
- CardMap (canonical lookup)
- Deck (52‑card slice with random shuffle)
- Deal (n-cards for m-hands, remaining deck)
- assets (PNG card faces)
- other helper functions.

Everything else — rules, payouts, UI, and state machines —
are isolated per application.

Drag-and-Drop is implemented using a full container overlay for tracking
 and animation.
Card image high-lighting is achieved via gray/yellow rectangle overlays.

Each app has the same structure:
- XXXX.go - the main program.
- prefs.go - optional and manually configurable parameters.
- ui.go - the user interface (fyne.io).
- game.go - the games' physical engine.
- evaluate.go - compute the resolutions of a situation.
- splash.png and splash.go -- image resource files.

Application configuration files ("preferences.json") are located in the
standard fyne storage directory with the directory name "com.scsi.XXXX",
where XXXX is "bridge", "rummy" "jacks", "stud5", "joker", "cribbage"
or "holdem".<br>
Windows: ~\AppData\Roaming\fyne\APP<br>
MacOS:    ~/Library/Preferences/fyne/APP<br>
Linux:   ~/.config/fyne/APP<br>
Where:
- ~ is your login (home) directory (folder) and APP is either:
- com.scsi.jacks
- com.scsi.stud5
- com.scsi.holdem
- com.scsi.joker
- com.scsi.bridge
- com.scsi.rummy
- com.scsi.cribbage
- com.scsi.freecell


---

## Shared Engine

All four apps use a unique card engine:

### Card

Defines card RankCode, SuitCode, and Icon

### cardMap

Canonical lookup from string Icons ("A♠", "10♦", "7♣", etc.) to Card.

### deck

Slice of 52 card codes with "math/rand/v2" random shuffle.
Each deal is shuffled using seeded random "rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))"<br>
"joker" and "reverse" card images are included.

### assets

PNG card faces embedded as Fyne resources.
Images from Byron Knoll.
https://www.byronknoll.com/projects.html



