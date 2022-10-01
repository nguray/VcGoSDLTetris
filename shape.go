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

func (sh *Shape) CheckLeftBoardLimit(renderer *sdl.Renderer) bool {
	l := sh.MinX1()*cellSize + sh.x
	return l < 0
}

func (sh *Shape) CheckRightBoardLimit(renderer *sdl.Renderer) bool {
	r := sh.MaxX1()*cellSize + cellSize + sh.x
	return r > NB_COLUMNS*cellSize
}

func (sh *Shape) CheckBottomLimit(renderer *sdl.Renderer) bool {
	//--------------------------------------------------
	b := sh.MaxY1()*cellSize + cellSize + sh.y
	return b > NB_ROWS*cellSize
}

func (sh *Shape) HitGround1(renderer *sdl.Renderer, board []int) int32 {
	var (
		ix int32
		iy int32
	)
	//--------------------------------------------------

	renderer.SetDrawColor(255, 0, 0, 255)
	for _, v := range sh.v {
		x := v.x*cellSize + sh.x + 1
		y := v.y*cellSize + sh.y + 1
		ix = int32(x / cellSize)
		iy = int32(y / cellSize)

		//rect = sdl.Rect{X: int32(x*cellSize + LEFT + 1), Y: int32(y*cellSize + TOP + 1), W: 5, H: 5}
		//renderer.FillRect(&rect)

		// renderer.DrawLine(LEFT+ix*cellSize, 0, LEFT+ix*cellSize, NB_ROWS*cellSize)
		// renderer.DrawLine(0, TOP+iy*cellSize, NB_COLUMNS*cellSize, TOP+iy*cellSize)
		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			iHit := iy*NB_COLUMNS + ix
			v := board[iHit]
			if v != 0 {
				return iHit
			}
		}

		x = v.x*cellSize + cellSize - 1 + sh.x
		y = v.y*cellSize + sh.y + 1
		ix = int32(x / cellSize)
		iy = int32(y / cellSize)
		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			iHit := iy*NB_COLUMNS + ix
			v := board[iHit]
			if v != 0 {
				return iHit
			}
		}

		x = v.x*cellSize + cellSize - 1 + sh.x
		y = v.y*cellSize + cellSize - 1 + sh.y
		ix = int32(x / cellSize)
		iy = int32(y / cellSize)
		// renderer.DrawLine(LEFT+ix*cellSize, 0, LEFT+ix*cellSize, NB_ROWS*cellSize)
		// renderer.DrawLine(0, TOP+iy*cellSize, NB_COLUMNS*cellSize, TOP+iy*cellSize)

		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			iHit := iy*NB_COLUMNS + ix
			v := board[iHit]
			if v != 0 {
				return iHit
			}
		}

		x = v.x*cellSize + sh.x + 1
		y = v.y*cellSize + cellSize - 1 + sh.y
		ix = int32(x / cellSize)
		iy = int32(y / cellSize)

		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			iHit := iy*NB_COLUMNS + ix
			v := board[iHit]
			if v != 0 {
				return iHit
			}
		}

	}

	return -1
}

// func (sh *Shape) MinX() int32 {
// 	var (
// 		x    int32
// 		minX int32
// 	)
// 	ix := int32(sh.x / cellSize)
// 	minX = sh.v[0].x + ix
// 	for i := 1; i < 4; i++ {
// 		x = sh.v[i].x + ix
// 		if x < minX {
// 			minX = x
// 		}
// 	}
// 	return minX
// }

func (sh *Shape) MinX1() int32 {
	var (
		x    int32
		minX int32
	)
	minX = sh.v[0].x
	for i := 1; i < 4; i++ {
		x = sh.v[i].x
		if x < minX {
			minX = x
		}
	}
	return minX
}

// func (sh *Shape) MaxX() int32 {
// 	var (
// 		x    int32
// 		maxX int32
// 	)
// 	ix := int32(sh.x / cellSize)
// 	maxX = sh.v[0].x + ix
// 	for i := 1; i < 4; i++ {
// 		x = sh.v[i].x + ix
// 		if x > maxX {
// 			maxX = x
// 		}
// 	}
// 	return maxX
// }

func (sh *Shape) MaxX1() int32 {
	var (
		x    int32
		maxX int32
	)
	maxX = sh.v[0].x
	for i := 1; i < 4; i++ {
		x = sh.v[i].x
		if x > maxX {
			maxX = x
		}
	}
	return maxX
}

// func (sh *Shape) MaxY() int32 {
// 	var (
// 		y int32
// 	)
// 	iy := int32(sh.y / cellSize)
// 	maxY := sh.v[0].y + iy
// 	for i := 1; i < 4; i++ {
// 		y = sh.v[i].y + iy
// 		if y > maxY {
// 			maxY = y
// 		}
// 	}
// 	return maxY
// }

func (sh *Shape) MaxY1() int32 {
	var (
		y int32
	)
	maxY := sh.v[0].y
	for i := 1; i < 4; i++ {
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
