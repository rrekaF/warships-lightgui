/*
Warships-LightGUI provides an easy to use graphical user interface
for the `Warships Online` game.
*/
package board

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type pos int

const (
	Left  pos = iota // indicates player's board (on the left side)
	Right            // indicates enemy's board (on the right side)
)

type state int

const (
	empty state = iota
	Hit
	Miss
	Ship
	border
	ruler
)

const (
	delimiter   = 5
	boardWidth  = 10
	boardHeight = 10
	maxX        = boardWidth + delimiter + boardWidth + 1
	maxY        = boardHeight
)

// Board represents a game board, including both player's (Left)
// and enemy's (Right) sides.
type Board struct {
	b [maxX][maxY]state
	c *Config
}

func (b *Board) printChar(s state) string {
	switch s {
	case empty:
		return string(b.c.EmptyChar)
	case Ship:
		return string(b.c.ShipChar)
	case Hit:
		return string(b.c.HitChar)
	case Miss:
		return string(b.c.MissChar)
	case border:
		return string(b.c.BorderChar)
	case ruler:
		return ""
	default:
		return string('.')
	}
}

/*
Import imports player's ships from a slice of coordinates
(as returned by the game server) and places them on the
Left board.

If any of the coordinates is invalid, the function returnes ErrInvalidCoord.
*/
func (b *Board) Import(coords []string) error {
	for _, c := range coords {
		err := b.Set(Left, c, Ship)
		if err != nil {
			return err
		}
	}
	return nil
}

// Export exports ships from either Left (player's) or Right (enemy's)
// board. The return value is a slice of ship coordinates (using format
// expected by the game server).
func (b *Board) Export(p pos) []string {
	var result []string

	if p == Left {
		for y := maxY - 1; y >= 0; y-- {
			for x := 0; x < boardHeight; x++ {
				if b.b[x][y] == Ship {
					coords := fmt.Sprintf("%c%d", 'A'+byte(x-1), y+1)
					result = append(result, coords)
				}
			}
		}
	}

	if p == Right {
		for y := maxY - 1; y >= 0; y-- {
			for x := boardWidth + delimiter; x < maxX; x++ {
				if b.b[x][y] == Ship {
					coords := fmt.Sprintf("%c%d", 'A'+byte(x-boardWidth-delimiter-1), y+1)
					result = append(result, coords)
				}
			}
		}
	}

	return result
}

/*
HitOrMiss updates and returns the state of a coordinate on the board,
depending on the previous state:
  - Empty -> Miss
  - Ship, Hit  -> Hit

If the coordinate is invalid, the function returnes ErrInvalidCoord.

Parameters:
  - p (pos): Left or Right board
  - coord (string): a string representing the coordinate (e.g. "A1", "B2")

Returns:
  - s (state): updated state value (Empty, Miss, or Hit)
*/
func (b *Board) HitOrMiss(p pos, coord string) (state, error) {
	var s state

	x, y, err := b.stringCoordToInt(coord)
	if err != nil {
		return empty, err
	}

	if p == Left {
		s = b.b[x][y]
	} else {
		s = b.b[x+boardWidth+delimiter][y]
	}

	switch s {
	case Ship:
		b.Set(p, coord, Hit)
		return Hit, nil
	case Hit:
		return Hit, nil
	default:
		b.Set(p, coord, Miss)
		return Miss, nil
	}
}

var ErrInvalidCoord = errors.New("invalid coordinate")

func (b *Board) stringCoordToInt(coord string) (int, int, error) {
	if len(coord) != 2 && len(coord) != 3 {
		return 0, 0, ErrInvalidCoord
	}

	x := strings.ToUpper(coord)[0] - 'A'

	if x > 10 {
		return 0, 0, ErrInvalidCoord
	}

	y, err := strconv.Atoi(coord[1:])
	if err != nil {
		return 0, 0, ErrInvalidCoord
	}

	y--

	if y < 0 || y > maxY-1 {
		return 0, 0, ErrInvalidCoord
	}

	x++

	return int(x), y, nil
}

/*
Set updates the state of a coordinate on the board.

For the Left board, the function validates the state of the coordinate based on the following logic:
  - If the state is Miss and the previous state is not Empty, it does not update the state.
  - If the state is Hit and the previous state is not Ship, it does not update the state.
  - If the state is Ship and the previous state is not Empty, it does not update the state.

For the Right board, the function does not update the state if the previous state is not Empty.

If the coordinate is invalid, the function returnes ErrInvalidCoord.

Parameters:
  - p (pos): Left or Right board
  - coord (string): a string representing the coordinate (e.g. "A1", "B2")
  - s (state): the state to update the coordinate to (Empty, Miss, or Ship)
*/
func (b *Board) Set(p pos, coord string, s state) error {
	x, y, err := b.stringCoordToInt(coord)
	if err != nil {
		return err
	}

	if p == Left {
		switch s {
		case Miss, Ship:
			if b.b[x][y] != empty {
				return nil
			}
		case Hit:
			if b.b[x][y] != Ship {
				return nil
			}
		}
	}

	if p == Right {
		x = x + boardWidth + delimiter
		// if b.b[x][y] != empty {
		// 	return nil
		// }
	}

	b.b[x][y] = s
	return nil
}

// New returns a new Board.
func New(c *Config) *Board {
	b := &Board{
		c: c,
	}

	for y := maxY - 1; y >= 0; y-- {
		b.b[0][y] = ruler
		b.b[boardWidth+delimiter][y] = ruler
	}

	return b
}
