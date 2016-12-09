/*

   Very ugly code, written fast, for contest.

*/

package main

import (
	"fmt"
	"math"
)

const (
	WIZARD_RADIUS = 400.0
	MAX_ID        = 13
)

var POINT_ZERO = point{0, 0, 0, 0}

type point struct {
	x  int
	y  int
	vx int
	vy int
}

type wizard struct {
	holding   bool
	pos       point
	player    int
	id        int
	target    point
	snaffleID int
}

type snaffle struct {
	taken bool
	pos   point
	id    int
}

type bludger struct {
	pos point
	id  int
}

type strct struct {
	turn    int
	snaffle int
}

func distance(c0, c1 point) float64 {
	x := c0.x - c1.x
	y := c0.y - c1.y

	return math.Sqrt(float64(x*x + y*y))
}

func filterGrabbed(snaffles []snaffle, wizards []wizard) []snaffle {

	grabbed := make([]bool, len(snaffles))

	for j, w := range wizards {
		if w.holding {
			for i, s := range snaffles {
				if distance(w.pos, s.pos) <= WIZARD_RADIUS-1 {
					//snaffles[i].grabbed = true
					wizards[j].snaffleID = s.id
					grabbed[i] = true
				}
			}
		}
	}

	var f []snaffle

	for i, s := range snaffles {
		if !grabbed[i] {
			f = append(f, s)
		}
	}

	return f
}

func findTarget(w wizard, snaffles []snaffle) point {
	var target = 0
	closest := -1.0

	for i, s := range snaffles {
		dist := distance(w.pos, s.pos)

		if !s.taken && (dist < closest || closest < 0) {
			closest = dist
			target = i
		}
	}

	if closest == -1.0 {
		return POINT_ZERO
	}

	snaffles[target].taken = true

	return snaffles[target].pos
}

func targetize(wizards []wizard, snaffles []snaffle, myID int) {

	closestWizards := make([]int, MAX_ID)
	for i := 0; i < MAX_ID; i++ {
		closestWizards[i] = -1
	}

	for i, s := range snaffles {
		closest := -1.0
		closestID := -1
		for j, w := range wizards {
			if w.player == myID {
				d := distance(s.pos, w.pos)
				if d < closest || closest < 0 {
					closest = d
					closestID = j
				}
			}
		}

		closestWizards[i] = closestID
	}

	for i, w := range wizards {
		if w.player == myID {
			closest := -1.0
			closestID := -1
			for sID, wID := range closestWizards {
				if wID == i {
					d := distance(snaffles[sID].pos, w.pos)
					if d < closest || closest < 0 {
						closest = d
						closestID = sID
					}
				}
			}

			if closestID == -1 {
				wizards[i].target = findTarget(w, snaffles)
			} else {
				wizards[i].target = snaffles[closestID].pos
			}

		}
	}
}

/**
 * Grab Snaffles and try to throw them through the opponent's goal!
 * Move towards a Snaffle and use your team id to determine where you need to throw it.
 **/

func main() {
	magicPower := 0
	turn := 1

	// myTeamId: if 0 you need to score on the right of the map, if 1 you need to score on the left
	var myTeamId int
	enemyTeamId := 0
	fmt.Scan(&myTeamId)

	goal := point{0, 3750, 0, 0}

	if myTeamId == 0 {
		enemyTeamId = 1
		goal = point{16000, 3750, 0, 0}
	}

	justThrown := make([]strct, MAX_ID)
	for i := 0; i < MAX_ID; i++ {
		justThrown[i] = strct{-10, -1}
	}

	for {
		// entities: number of entities still in game
		var entities int
		fmt.Scan(&entities)

		var wizards []wizard
		var snaffles, tmp_snaffles []snaffle

		var bludgers []bludger

		for i := 0; i < entities; i++ {
			// entityId: entity identifier
			// entityType: "WIZARD", "OPPONENT_WIZARD" or "SNAFFLE" (or "BLUDGER" after first league)
			// x: position
			// y: position
			// vx: velocity
			// vy: velocity
			// state: 1 if the wizard is holding a Snaffle, 0 otherwise
			var entityId int
			var entityType string
			var x, y, vx, vy, state int
			fmt.Scan(&entityId, &entityType, &x, &y, &vx, &vy, &state)

			switch entityType {
			case "WIZARD":
				holding := false
				if state == 1 {
					holding = true
				}
				wizards = append(wizards, wizard{holding, point{x, y, vx, vy}, myTeamId, entityId, POINT_ZERO, -1})
			case "OPPONENT_WIZARD":
				holding := false
				if state == 1 {
					holding = true
				}
				wizards = append(wizards, wizard{holding, point{x, y, vx, vy}, enemyTeamId, entityId, POINT_ZERO, -1})
			case "SNAFFLE":
				tmp_snaffles = append(tmp_snaffles, snaffle{false, point{x, y, vx, vy}, entityId})
			case "BLUDGER":
				bludgers = append(bludgers, bludger{point{x, y, vx, vy}, entityId})
			}

		}

		snaffles = filterGrabbed(tmp_snaffles, wizards)

		targetize(wizards, snaffles, myTeamId)

		for _, w := range wizards {
			if w.player == myTeamId {

				if w.holding {

					fmt.Printf("THROW %v %v 500\n", goal.x-w.pos.vx, goal.y-w.pos.vy)
					justThrown[w.id] = strct{turn, w.snaffleID}
				} else {

					if magicPower >= 20 && justThrown[w.id].snaffle >= 0 {

						fmt.Printf("FLIPENDO %v AVADA KEDAVRA!\n", justThrown[w.id].snaffle)
						magicPower -= 20

					} else {
						fmt.Printf("MOVE %v %v 150\n", w.target.x, w.target.y)
					}

				}

			}
		}

		for i := 0; i < MAX_ID; i++ {
			if justThrown[i].turn < turn {
				justThrown[i] = strct{-10, -1}
			}
		}

		magicPower++
		turn++
	}
}
