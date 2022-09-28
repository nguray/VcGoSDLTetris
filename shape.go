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

	y := sh.y + TOP + 1

	for _, v := range sh.v {
		l = v.x + sh.x
		if t >= 0 {
			rect = sdl.Rect{X: int32(l*cellSize + LEFT + 1), Y: int32(v.y*cellSize + y), W: a, H: a}
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

func (sh *Shape) OutBoardLimit1() bool {
	//--------------------------------------------------

	//-- Offset to have the bottom
	b := sh.y + cellSize - 1

	iy := int32(b / cellSize)
	for _, v := range sh.v {
		x := v.x + sh.x
		y := v.y + iy
		if (x < 0) || x > (NB_COLUMNS-1) || (y > NB_ROWS-1) {
			return true
		}
	}
	return false
}

func (sh *Shape) HitGround1(renderer *sdl.Renderer, board []int) int32 {
	var (
		iy int32
		//rect sdl.Rect
	)
	//--------------------------------------------------

	t := sh.y

	//renderer.SetDrawColor(255, 0, 0, 255)
	//-- Top
	iy = int32(t / cellSize)
	for _, v := range sh.v {
		x := v.x + sh.x
		y := v.y + iy

		//rect = sdl.Rect{X: int32(x*cellSize + LEFT + 1), Y: int32(y*cellSize + TOP + 1), W: 5, H: 5}
		//renderer.FillRect(&rect)

		//renderer.DrawLine(x*cellSize-5, y*cellSize, x*cellSize+5, y*cellSize)
		//renderer.DrawLine(x*cellSize, y*cellSize-5, x*cellSize, y*cellSize+5)
		if (x >= 0) && x < NB_COLUMNS && (y >= 0) && (y < NB_ROWS) {
			iHit := y*NB_COLUMNS + x
			v := board[iHit]
			if v != 0 {
				return iHit
			}
		}
	}

	//renderer.SetDrawColor(0, 0, 255, 255)
	//-- Bottom
	t += cellSize - 1
	iy = int32(t / cellSize)
	for _, v := range sh.v {
		x := v.x + sh.x
		y := v.y + iy

		//rect = sdl.Rect{X: int32(x*cellSize + LEFT + 1), Y: int32(y*cellSize + TOP + 1), W: 5, H: 5}
		//renderer.FillRect(&rect)
		//renderer.DrawLine(x*cellSize-5, y*cellSize, x*cellSize+5, y*cellSize)
		//renderer.DrawLine(x*cellSize, y*cellSize-5, x*cellSize, y*cellSize+5)
		if (x >= 0) && x < NB_COLUMNS && (y >= 0) && (y < NB_ROWS) {
			iHit := y*NB_COLUMNS + x
			v := board[iHit]
			if v != 0 {
				return iHit
			}
		}
	}

	return -1
}

func (sh *Shape) MinX() int32 {
	var (
		x    int32
		minX int32
	)
	minX = sh.v[0].x + sh.x
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
	maxX = sh.v[0].x + sh.x
	for i := 1; i < 4; i++ {
		x = sh.v[i].x + sh.x
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
	iy := int32(sh.y / cellSize)
	maxY := sh.v[0].y + iy
	for i := 1; i < 4; i++ {
		y = sh.v[i].y + iy
		if y > maxY {
			maxY = y
		}
	}
	return maxY
}
