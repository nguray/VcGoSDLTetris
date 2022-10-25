package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Vector2i struct {
	x int32
	y int32
}

var (
	tetrominos []Vector2i
	colors     []sdl.Color
)

type Shape struct {
	typ   int32
	x     int32
	y     int32
	v     [4]Vector2i
	color sdl.Color
}

func InitTetrominos() {

	tetrominos = []Vector2i{
		{0, 0}, {0, 0}, {0, 0}, {0, 0},
		{0, -1}, {0, 0}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 0}, {1, 0}, {1, 1},
		{0, -1}, {0, 0}, {0, 1}, {0, 2},
		{-1, 0}, {0, 0}, {1, 0}, {0, 1},
		{0, 0}, {1, 0}, {0, 1}, {1, 1},
		{-1, -1}, {0, -1}, {0, 0}, {0, 1},
		{1, -1}, {0, -1}, {0, 0}, {0, 1}}

	colors = []sdl.Color{
		{R: 0, G: 0, B: 0, A: 0xFF},
		{R: 0xFF, G: 0x60, B: 0x60, A: 0xFF},
		{R: 0x60, G: 0xFF, B: 0x60, A: 0xFF},
		{R: 0x60, G: 0x60, B: 0xFF, A: 0xFF},
		{R: 0xCC, G: 0xCC, B: 0x60, A: 0xFF},
		{R: 0xCC, G: 0x60, B: 0xCC, A: 0xFF},
		{R: 0x60, G: 0xCC, B: 0xCC, A: 0xFF},
		{R: 0xDA, G: 0xAA, B: 0x00, A: 0xFF}}

}

func ShapeNew(typ, x, y int32) *Shape {

	shape := &Shape{typ, x, y, [4]Vector2i{}, sdl.Color{R: 0xFF, G: 0, B: 0, A: 0xFF}}
	shape.InitGfx()
	shape.color = colors[shape.typ]
	//--
	return shape
}

func (sh *Shape) InitGfx() {

	offSet := int(sh.typ) * len(sh.v)
	for i := 0; i < len(sh.v); i++ {
		sh.v[i].x = tetrominos[i+offSet].x
		sh.v[i].y = tetrominos[i+offSet].y
	}

}

func (sh *Shape) Draw(renderer *sdl.Renderer) {

	var (
		x, y int32
		rect sdl.Rect
	)

	renderer.SetDrawColor(sh.color.R, sh.color.G, sh.color.B, sh.color.A)
	a := int32(cellSize - 2)

	for _, v := range sh.v {
		x = v.x*cellSize + sh.x + LEFT + 1
		y = v.y*cellSize + sh.y + TOP + 1
		if y >= TOP {
			rect = sdl.Rect{X: int32(x), Y: int32(y), W: a, H: a}
			renderer.FillRect(&rect)
		}
	}

}

func (sh *Shape) RotateLeft() {
	if sh.typ != 5 {
		var x, y int32
		for i := 0; i < len(sh.v); i++ {
			x = sh.v[i].y
			y = -sh.v[i].x
			sh.v[i].x = x
			sh.v[i].y = y
		}
	}
}

func (sh *Shape) RotateRight() {
	if sh.typ != 5 {
		var x, y int32
		for i := 0; i < len(sh.v); i++ {
			x = -sh.v[i].y
			y = sh.v[i].x
			sh.v[i].x = x
			sh.v[i].y = y
		}
	}
}

func (sh *Shape) MinX() int32 {
	var (
		x    int32
		minX int32
	)
	minX = sh.v[0].x
	for i := 1; i < len(sh.v); i++ {
		x = sh.v[i].x
		if x < minX {
			minX = x
		}
	}
	return minX
}

func (sh *Shape) MaxX() int32 {
	var (
		x    int32
		maxX int32
	)
	maxX = sh.v[0].x
	for i := 1; i < len(sh.v); i++ {
		x = sh.v[i].x
		if x > maxX {
			maxX = x
		}
	}
	return maxX
}

func (sh *Shape) MaxY() int32 {
	var (
		y int32
	)
	maxY := sh.v[0].y
	for i := 1; i < len(sh.v); i++ {
		y = sh.v[i].y
		if y > maxY {
			maxY = y
		}
	}
	return maxY
}

func (sh *Shape) Column() int32 {
	return int32(sh.x / cellSize)
}

func (sh *Shape) IsOutLeftBoardLimit() bool {
	l := sh.MinX()*cellSize + sh.x
	return l < 0
}

func (sh *Shape) IsOutRightBoardLimit() bool {
	r := sh.MaxX()*cellSize + cellSize + sh.x
	return r > NB_COLUMNS*cellSize
}

func (sh *Shape) IsAlwaysOutBoardLimit() bool {
	return true
}

func (sh *Shape) IsOutBottomLimit() bool {
	//--------------------------------------------------
	b := sh.MaxY()*cellSize + cellSize + sh.y
	return b > NB_ROWS*cellSize
}

func (sh *Shape) HitGround(board []int) bool {

	//--------------------------------------------------

	Hit := func(x int32, y int32) bool {
		ix := int32(x / cellSize)
		iy := int32(y / cellSize)
		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			v := board[iy*NB_COLUMNS+ix]
			if v != 0 {
				return true
			}
		}
		return false
	}

	for _, v := range sh.v {

		x := v.x*cellSize + sh.x + 1
		y := v.y*cellSize + sh.y + 1
		if Hit(x, y) {
			return true
		}

		x = v.x*cellSize + cellSize - 1 + sh.x
		y = v.y*cellSize + sh.y + 1
		if Hit(x, y) {
			return true
		}

		x = v.x*cellSize + cellSize - 1 + sh.x
		y = v.y*cellSize + cellSize - 1 + sh.y
		if Hit(x, y) {
			return true
		}

		x = v.x*cellSize + sh.x + 1
		y = v.y*cellSize + cellSize - 1 + sh.y
		if Hit(x, y) {
			return true
		}

	}

	return false
}
