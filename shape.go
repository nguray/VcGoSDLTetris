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

	offSet := int(sh.typ * 4)
	for i := 0; i < 4; i++ {
		sh.v[i].x = tetrominos[i+offSet].x
		sh.v[i].y = tetrominos[i+offSet].y
	}

}

func (sh *Shape) Draw(renderer *sdl.Renderer) {

	var (
		l, t int32
		rect sdl.Rect
	)

	renderer.SetDrawColor(sh.color.R, sh.color.G, sh.color.B, sh.color.A)
	a := int32(cellSize - 2)
	for _, v := range sh.v {
		l = v.x + sh.x
		t = v.y + sh.y
		if t >= 0 {
			rect = sdl.Rect{X: int32(l*cellSize + LEFT + 1), Y: int32(t*cellSize + TOP + 1), W: a, H: a}
			renderer.FillRect(&rect)
		}
	}

}

func (sh *Shape) RotateLeft() {
	if sh.typ != 5 {
		var x, y int32
		for i := 0; i < 4; i++ {
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
		for i := 0; i < 4; i++ {
			x = -sh.v[i].y
			y = sh.v[i].x
			sh.v[i].x = x
			sh.v[i].y = y
		}
	}
}

func (sh *Shape) OutBoardLimit() bool {
	//--------------------------------------------------
	for _, v := range sh.v {
		x := v.x + sh.x
		y := v.y + sh.y
		if (x < 0) || x > (NB_COLUMNS-1) || (y > NB_ROWS-1) {
			return true
		}
	}
	return false
}

func (sh *Shape) HitGround(board []int) bool {
	//--------------------------------------------------
	for _, v := range sh.v {
		x := v.x + sh.x
		y := v.y + sh.y
		if (x >= 0) && x < NB_COLUMNS && (y >= 0) && (y < NB_ROWS) {
			v := board[y*NB_COLUMNS+x]
			if v != 0 {
				return true
			}
		}
	}
	return false
}

func (sh *Shape) MinX() int32 {
	var (
		x    int32
		minX int32
	)
	minX = sh.v[0].x + sh.y
	for i := 1; i < 4; i++ {
		x = sh.v[i].x + sh.x
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
	maxX = sh.v[0].x + sh.y
	for i := 1; i < 4; i++ {
		x = sh.v[i].x + sh.x
		if x > maxX {
			maxX = x
		}
	}
	return maxX
}

func (sh *Shape) MaxY() int32 {
	var y int32
	maxY := sh.v[0].y + sh.y
	for i := 1; i < 4; i++ {
		y = sh.v[i].y + sh.y
		if y > maxY {
			maxY = y
		}
	}
	return maxY
}
