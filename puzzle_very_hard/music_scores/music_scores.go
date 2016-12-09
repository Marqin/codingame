/*

  This is not perfect, but gets 100%.

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type image struct {
	w            int
	h            int
	data         []bool
	firstBlack   pixel
	staveHeight  int
	staveDistace int
}

type pixel struct {
	x int
	y int
}

func (i *image) setStaveInfo() {
	x := i.firstBlack.x
	y := i.firstBlack.y
	for y < i.h {
		if !i.get(pixel{x, y}) {
			i.staveHeight = y - i.firstBlack.y
			break
		}
		y++
	}

	heightStart := y

	for y < i.h {
		if i.get(pixel{x, y}) {
			i.staveDistace = y - heightStart
			break
		}
		y++
	}
}

func (i *image) get(p pixel) bool {
	return i.data[p.x+p.y*i.w]
}

func (i *image) set(p pixel, value bool) {
	i.data[p.x+p.y*i.w] = value
}

func (i *image) nextPixel(p pixel) pixel {
	if p.x+1 < i.w {
		return pixel{p.x + 1, p.y}
	} else {
		return pixel{0, p.y + 1}
	}
}

func (i *image) onLine(p pixel) bool {

	lineStart := i.firstBlack.y
	lineEnd := i.firstBlack.y + i.staveHeight - 1
	for x := 0; x < 6; x++ {
		if p.y >= lineStart && p.y <= lineEnd {
			return true
		}

		lineStart += i.staveDistace + i.staveHeight
		lineEnd += i.staveDistace + i.staveHeight
	}

	return false
}

func (i *image) onStripe(p pixel) bool {

	sum := 0

	nxt := i.firstBlack.y + (i.nextLine(p))*(i.staveDistace+i.staveHeight)

	for y := p.y; y < nxt; y++ {
		if i.get(pixel{p.x, y}) {
			break
		} else {
			sum++
		}
	}
	for y := p.y - 1; y >= i.firstBlack.y-i.staveDistace; y-- {
		if i.get(pixel{p.x, y}) {
			break
		} else {
			sum++
		}
	}

	return (sum >= i.staveDistace)
}

func (i *image) closeToLine(p pixel) bool {

	var SEARCH_DISTANCE = i.staveHeight

	for y := p.y - SEARCH_DISTANCE; y <= p.y+SEARCH_DISTANCE && y < i.h; y++ {
		if i.onLine(pixel{p.x, y}) {
			return true
		}
	}
	return false
}

func (i *image) nextLine(p pixel) int {

	lineStart := i.firstBlack.y
	lineEnd := i.firstBlack.y + i.staveHeight - 1
	for x := 0; x < 6; x++ {

		if p.y <= lineStart {
			return x
		}

		lineStart += i.staveDistace + i.staveHeight
		lineEnd += i.staveDistace + i.staveHeight
	}

	return 6
}

func (i *image) getNote(p pixel) string {

	note := ""

	close := i.closeToLine(p)
	nearestLine := i.nextLine(p)

	if close {
		switch nearestLine {
		case 0:
			note = "F"
		case 1:
			note = "D"
		case 2:
			note = "B"
		case 3:
			note = "G"
		case 4:
			note = "E"
		case 5:
			note = "C"
		}
	} else {
		switch nearestLine {
		case 0:
			note = "G"
		case 1:
			note = "E"
		case 2:
			note = "C"
		case 3:
			note = "A"
		case 4:
			note = "F"
		case 5:
			note = "D"
		}
	}

	if i.get(p) {
		note += "Q"
	} else {
		note += "H"
	}

	return note
}

func (i *image) setFirstBlack() {
	for x := 0; x < i.w; x++ {
		for y := 0; y < i.h; y++ {
			p := pixel{x, y}
			if i.get(p) {
				i.firstBlack = p
				return
			}
		}
	}
}

func readImage(data string, w, h int) image {
	img := image{w, h, make([]bool, w*h), pixel{-1, -1}, 0, 0}

	var firstFree pixel

	letter := ""
	for _, s := range strings.Split(data, " ") {
		if letter == "" {
			letter = s
		} else {
			number, _ := strconv.Atoi(s)
			for x := 0; x < number; x++ {
				if letter == "W" {
					img.set(firstFree, false)
				} else {
					img.set(firstFree, true)
				}
				firstFree = img.nextPixel(firstFree)
			}
			letter = ""
		}
	}

	img.setFirstBlack()
	img.setStaveInfo()

	return img
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var W, H int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &W, &H)

	scanner.Scan()
	img := readImage(scanner.Text(), W, H)

	noteStart := 0

	var notes []string

	for x := 0; x < img.w; x++ {
		for y := 0; y < img.h; y++ {
			firstPixel := pixel{x, y}
			if img.get(firstPixel) && !img.onLine(firstPixel) {
				noteStart = x
				noteEnd := x
				for z := x; z < img.w; z++ {
					lastPixel := pixel{z, y}
					if !img.get(lastPixel) && img.onStripe(lastPixel) {
						noteEnd = z
						break
					}
				}

				mid := noteStart + (noteEnd-noteStart)/2
				notes = append(notes, img.getNote(pixel{mid, y}))

				x = noteEnd
				break
			}
		}
	}

	fmt.Println(strings.Join(notes, " "))
}
