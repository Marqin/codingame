/*

   This got me to Gold.

*/

package main

import (
	"fmt"
	"math"
)

type point struct {
	x int
	y int
}

type pointDist struct {
	p point
	d float64
}

const (
	RADIUS      = 600
	FORCE_FIELD = 400
)

func find(points *[]pointDist, f point) int {

	for i, p := range *points {
		if f.x == p.p.x && f.y == p.p.y {
			return i
		}
	}

	return -1
}

func distance(c0, c1 point) float64 {
	x := c0.x - c0.x
	y := c0.y - c1.y
	return math.Sqrt(float64(x*x + y*y))
}

func vector(from, to point) point {
	return point{to.x - from.x, to.y - from.y}
}

func getStrikePoint(checkpoint, enemy point) (strikePoint point) {

	vec := vector(enemy, checkpoint)

	d := distance(checkpoint, vec)

	x := int(float64(FORCE_FIELD*vec.x) / d)
	y := int(float64(FORCE_FIELD*vec.y) / d)

	return addVec(enemy, point{x, y})
}

func addVec(x, y point) point {
	return point{x.x + y.x, x.y + y.y}
}

func main() {

	//boost := 1
	var checkpoints []pointDist
	mapped := false

	var previous, current, starting, best point
	firstLoop := true

	bestDist := -1.0

	boost := 1

	var lastPlace point

	for {
		var nextCheckpointDist, nextCheckpointAngle int
		var curr, next, enemy point

		fmt.Scan(&curr.x, &curr.y, &next.x, &next.y, &nextCheckpointDist, &nextCheckpointAngle)
		fmt.Scan(&enemy.x, &enemy.y)

		if firstLoop {
			lastPlace = curr
			starting = curr
			previous = curr
			current = next

			checkpoints = append(checkpoints, pointDist{starting, 0})

			strike := getStrikePoint(current, enemy)
			fmt.Printf("%d %d BOOST\n", strike.x, strike.y)

			firstLoop = false
			continue
		}

		if current != next {
			previous = current
			current = next
		}

		if mapped {
			if bestDist < 0 {
				for _, p := range checkpoints {
					if p.d > bestDist {
						bestDist = p.d
						best = p.p
					}
				}
			}
		} else {
			d := distance(previous, current)
			if distance(current, starting) < RADIUS {
				i := find(&checkpoints, starting)
				checkpoints[i] = pointDist{current, d}
				mapped = true
			} else {
				checkpoints = append(checkpoints, pointDist{current, d})
			}
		}

		thrust := 100

		if nextCheckpointAngle > 90 || nextCheckpointAngle < -90 {
			thrust = 0
		} else if nextCheckpointAngle < 10 && nextCheckpointAngle > -10 {
			if current == best && boost > 0 {
				thrust = 200
				boost--
			}
		}

		speed := vector(lastPlace, curr)
		lastPlace = curr

		target := point{current.x - speed.x*3, current.y - speed.y*3}

		if thrust <= 100 {
			fmt.Printf("%d %d %d\n", target.x, target.y, thrust)
		} else {
			fmt.Printf("%d %d BOOST\n", current.x, current.y)
		}
	}
}
