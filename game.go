package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type KeyChar struct {
	keycode sdl.Keycode
	c       string
}

type Game struct {
	velX          int32
	fDrop         bool
	fFastDown     bool
	curMode       GameMode
	curScore      int
	curTetromino  *Shape
	nextTetromino *Shape
	board         []int
	highScores    []HightScore
	idHighScore   int
	userName      string
	tblKeyChars   []KeyChar
	fQuitGame     bool
}

func GameNew() *Game { //int32(myRand.Intn(7)+1)
	game := &Game{0, false, false, STANDBY, 0, nil, ShapeNew(int32(myRand.Intn(7)+1),
		NB_COLUMNS+3, 10*cellSize), make([]int, NB_ROWS*NB_COLUMNS), make([]HightScore, 10), -1, "", make([]KeyChar, 1), false}
	for i := 0; i < len(game.highScores); i++ {
		game.highScores[i].name = "--------"
		game.highScores[i].score = 0
	}

	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_a, c: "A"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_b, c: "B"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_c, c: "C"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_d, c: "D"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_e, c: "E"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_f, c: "F"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_g, c: "G"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_h, c: "H"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_i, c: "I"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_j, c: "J"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_k, c: "K"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_l, c: "L"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_m, c: "M"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_n, c: "N"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_o, c: "O"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_p, c: "P"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_q, c: "Q"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_r, c: "R"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_s, c: "S"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_t, c: "T"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_u, c: "U"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_v, c: "V"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_w, c: "W"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_x, c: "X"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_y, c: "Y"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_z, c: "Z"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_0, c: "0"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_1, c: "1"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_2, c: "2"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_3, c: "3"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_4, c: "4"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_5, c: "5"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_6, c: "6"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_7, c: "7"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_8, c: "8"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_9, c: "9"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_0, c: "0"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_1, c: "1"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_2, c: "2"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_3, c: "3"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_4, c: "4"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_5, c: "5"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_6, c: "6"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_7, c: "7"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_8, c: "8"})
	game.tblKeyChars = append(game.tblKeyChars, KeyChar{keycode: sdl.K_KP_9, c: "9"})

	return game
}

func (ga *Game) getChar(keycode sdl.Keycode) string {
	for _, kc := range ga.tblKeyChars {
		if keycode == kc.keycode {
			return kc.c
		}
	}
	return ""
}

func (ga *Game) DrawBoard(renderer *sdl.Renderer) {
	//----------------------------------------------------------------
	var (
		rect sdl.Rect
		x    int32
		y    int32
		l, c int32
	)
	a := cellSize - 2
	for l = 0; l < NB_ROWS; l++ {
		for c = 0; c < NB_COLUMNS; c++ {
			v := ga.board[l*NB_COLUMNS+c]
			if v != 0 {
				x = int32(c*cellSize + LEFT + 1)
				y = int32(l*cellSize + TOP + 1)
				rect = sdl.Rect{X: x, Y: y, W: int32(a), H: int32(a)}
				col := colors[v]
				renderer.SetDrawColor(col.R, col.G, col.B, col.A)
				renderer.FillRect(&rect)

			}
		}
	}

}

func (ga *Game) DrawScore(renderer *sdl.Renderer) {
	x := LEFT
	y := (NB_ROWS + 1) * cellSize
	textScore := fmt.Sprintf("Score : %06d", ga.curScore)
	surfScore, err := tt_font.RenderUTF8Blended(textScore, sdl.Color{R: 255, G: 255, B: 0, A: 255})
	if err == nil {
		textureScore, err := renderer.CreateTextureFromSurface(surfScore)
		if err == nil {
			_, _, width, height, _ := textureScore.Query()
			renderer.Copy(textureScore, nil, &sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height})
			textureScore.Destroy()
		}
		surfScore.Free()
	}
}

func (ga *Game) DrawHightScores(renderer *sdl.Renderer) {
	var (
		x      int32
		y      int32
		width  int32
		height int32
	)
	y = TOP + cellSize
	strTitle := fmt.Sprintf("HIGH SCORES")
	surfTitle, err := tt_font.RenderUTF8Blended(strTitle, sdl.Color{R: 255, G: 255, B: 0, A: 255})
	if err == nil {
		textureTitle, err := renderer.CreateTextureFromSurface(surfTitle)
		if err == nil {
			_, _, width, height, _ = textureTitle.Query()
			x = LEFT + (NB_COLUMNS/2)*cellSize - int32(width/2)
			renderer.Copy(textureTitle, nil, &sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height})
			textureTitle.Destroy()
			y += int32(3 * height)
		}
		surfTitle.Free()
	}

	xCol0 := LEFT + cellSize
	xCol1 := LEFT + (NB_COLUMNS/2+2)*cellSize
	for _, h := range ga.highScores {

		surfName, err := tt_font.RenderUTF8Blended(h.name, sdl.Color{R: 255, G: 255, B: 0, A: 255})
		if err == nil {
			textureName, err := renderer.CreateTextureFromSurface(surfName)
			if err == nil {
				_, _, width, height, _ = textureName.Query()
				renderer.Copy(textureName, nil, &sdl.Rect{X: int32(xCol0), Y: int32(y), W: width, H: height})
				textureName.Destroy()
			}
			surfName.Free()
		}

		strScore := fmt.Sprintf("%06d", h.score)
		surfScore, err := tt_font.RenderUTF8Blended(strScore, sdl.Color{R: 255, G: 255, B: 0, A: 255})
		if err == nil {
			textureScore, err := renderer.CreateTextureFromSurface(surfScore)
			if err == nil {
				_, _, width, height, _ = textureScore.Query()
				renderer.Copy(textureScore, nil, &sdl.Rect{X: int32(xCol1), Y: int32(y), W: width, H: height})
				textureScore.Destroy()
			}
			surfScore.Free()
		}

		//--
		y += int32(height + 8)

	}
}

func (ga *Game) DrawStandBy(renderer *sdl.Renderer) {
	var (
		x int32
		y int32
	)
	y = TOP + (NB_ROWS/4)*cellSize
	strTitle := fmt.Sprintf("GoLang Tetris in SDL2")
	surfTitle, err := tt_font.RenderUTF8Blended(strTitle, sdl.Color{R: 255, G: 255, B: 0, A: 255})
	if err == nil {
		textureTitle, err := renderer.CreateTextureFromSurface(surfTitle)
		if err == nil {
			_, _, width, height, _ := textureTitle.Query()
			x = LEFT + (NB_COLUMNS/2)*cellSize - int32(width/2)
			renderer.Copy(textureTitle, nil, &sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height})
			textureTitle.Destroy()
			y += int32(2*height + 4)
		}
		surfTitle.Free()
	}

	strMsg := fmt.Sprintf("Press SPACE to Start")
	surfMsg, err := tt_font.RenderUTF8Blended(strMsg, sdl.Color{R: 255, G: 255, B: 0, A: 255})
	if err == nil {
		textureMsg, err := renderer.CreateTextureFromSurface(surfMsg)
		if err == nil {
			_, _, width, height, _ := textureMsg.Query()
			x = LEFT + (NB_COLUMNS/2)*cellSize - int32(width/2)
			renderer.Copy(textureMsg, nil, &sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height})
			textureMsg.Destroy()
			//y += int(height + 4)
		}
		surfMsg.Free()
	}

}

func (ga *Game) NewTetromino() {
	//--------------------------------------------------
	ga.curTetromino = ga.nextTetromino
	ga.curTetromino.x = 6
	ga.curTetromino.y = 0
	ga.curTetromino.y = -ga.curTetromino.MaxY()
	ga.nextTetromino = ShapeNew(TetrisRandomizer(), NB_COLUMNS+3, 10*cellSize)

}

func (ga *Game) InitGame() {
	//--------------------------------------------------
	ga.curScore = 0
	for i := 0; i < NB_ROWS*NB_COLUMNS; i++ {
		ga.board[i] = 0
	}
	ga.curTetromino = nil
	ga.nextTetromino = ShapeNew(int32(myRand.Intn(7)+1), NB_COLUMNS+3, 10*cellSize)

}

func (ga *Game) IsGameOver() bool {
	//------------------------------------------------------
	for i := 0; i < NB_COLUMNS; i++ {
		if ga.board[i] != 0 {
			return true
		}
	}
	return false
}

func (ga *Game) FreezeCurTetramino() {
	//--------------------------------------------------
	if ga.curTetromino != nil {
		for _, v := range ga.curTetromino.v {
			x := v.x + ga.curTetromino.x
			y := v.y + ga.curTetromino.y
			if x >= 0 && x < NB_COLUMNS && y >= 0 && y < NB_ROWS {
				ga.board[y*NB_COLUMNS+x] = int(ga.curTetromino.typ)
			}
		}
		//--
		nbLines := ga.EraseCompletedLines()
		if nbLines > 0 {
			ga.curScore += ComputeScore(nbLines)
			succes_sound.Play(-1, 0)
		}

	}
}

func (ga *Game) FreezeCurTetramino1() {
	//--------------------------------------------------
	if ga.curTetromino != nil {
		iy := int32((ga.curTetromino.y + 1) / cellSize)
		for _, v := range ga.curTetromino.v {
			x := v.x + ga.curTetromino.x
			y := v.y + iy
			if x >= 0 && x < NB_COLUMNS && y >= 0 && y < NB_ROWS {
				ga.board[y*NB_COLUMNS+x] = int(ga.curTetromino.typ)
			}
		}
		//--
		nbLines := ga.EraseCompletedLines()
		if nbLines > 0 {
			ga.curScore += ComputeScore(nbLines)
			succes_sound.Play(-1, 0)
		}

	}
}

func (ga *Game) EraseCompletedLines() int {
	//--------------------------------------------------
	nbLines := 0
	fCompleted := false
	for r := 0; r < NB_ROWS; r++ {
		fCompleted = true
		for c := 0; c < NB_COLUMNS; c++ {
			if ga.board[r*NB_COLUMNS+c] == 0 {
				fCompleted = false
				break
			}
		}
		if fCompleted {
			nbLines++
			//-- Décaler d'une ligne le plateau
			for r1 := r; r1 > 0; r1-- {
				for c1 := 0; c1 < NB_COLUMNS; c1++ {
					ga.board[r1*NB_COLUMNS+c1] = ga.board[(r1-1)*NB_COLUMNS+c1]
				}
			}
		}
	}
	//fmt.Println("Nbre Erased Lines ", nbLines)
	return nbLines
}

func (ga *Game) ProcessEventsStandBy(renderer *sdl.Renderer) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			ga.fQuitGame = true
			println("Quit")
			return false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym
			if t.State == sdl.PRESSED && t.Repeat == 0 {
				switch keyCode {
				case sdl.K_SPACE:
					ga.curMode = PLAY
					processEvents = ga.ProcessEventsPlay
					ga.NewTetromino()
				case sdl.K_ESCAPE:
					ga.fQuitGame = true
					return false
				}
			}

		}
	}
	return true
}

func (ga *Game) ProcessEventsPlay(renderer *sdl.Renderer) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			ga.fQuitGame = true
			println("Quit")
			return false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym

			//keys := ""
			if t.State == sdl.PRESSED && t.Repeat == 0 {
				switch keyCode {
				case sdl.K_LEFT:
					ga.velX = -1
				case sdl.K_RIGHT:
					ga.velX = 1
				case sdl.K_UP:
					if ga.curTetromino != nil {
						ga.curTetromino.RotateLeft()

						idHit := ga.curTetromino.HitGround1(renderer, ga.board)

						if idHit >= 0 {
							//-- Undo Rotate
							ga.curTetromino.RotateRight()

						} else if ga.curTetromino.OutBoardLimit1() {
							max_pos_x := ga.curTetromino.MaxX()
							if max_pos_x >= NB_COLUMNS {
								dx := max_pos_x - (NB_COLUMNS - 1)
								ga.curTetromino.x -= dx
								idHit := ga.curTetromino.HitGround1(renderer, ga.board)
								if idHit >= 0 {
									ga.curTetromino.x += dx
									//-- Undo Rotate
									ga.curTetromino.RotateRight()
								}
							} else {
								min_pos_x := ga.curTetromino.MinX()
								if min_pos_x < 0 {
									dx := min_pos_x
									ga.curTetromino.x -= dx
									idHit := ga.curTetromino.HitGround1(renderer, ga.board)
									if idHit >= 0 {
										ga.curTetromino.x += dx
										//-- Undo Rotate
										ga.curTetromino.RotateRight()
									}
								}
							}

						}

					}
				case sdl.K_DOWN:
					ga.fFastDown = true
				case sdl.K_SPACE:
					if ga.curTetromino != nil {
						//-- Drop current Tetromino
						ga.fDrop = true
					}
				case sdl.K_ESCAPE:
					return false
				}
			} else if t.State == sdl.RELEASED {
				switch keyCode {
				case sdl.K_LEFT:
					ga.velX = 0
				case sdl.K_RIGHT:
					ga.velX = 0
				case sdl.K_DOWN:
					ga.fFastDown = false
				}

			}

		}
	}
	return true
}

func (ga *Game) IsHightScore(newscore int) int {
	//--------------------------------------------------
	for i, v := range ga.highScores {
		if newscore > v.score {
			return i
		}
	}
	return -1
}

func (ga *Game) InsertHightScore(id int, name string, score int) {
	//--------------------------------------------------
	ga.highScores = append(ga.highScores[:id+1], ga.highScores[id:]...)
	ga.highScores[id] = HightScore{name: name, score: score}
	ga.idHighScore = id
	ga.userName = name

}

func (ga *Game) ProcessEventsHightScores(renderer *sdl.Renderer) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			ga.fQuitGame = true
			println("Quit")
			return false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym
			if t.State == sdl.PRESSED && t.Repeat == 0 {
				switch keyCode {
				case sdl.K_BACKSPACE:
					sz := len(ga.userName)
					if sz > 0 {
						ga.userName = ga.userName[:sz-1]
						ga.highScores[ga.idHighScore].name = ga.userName
					}
				case sdl.K_ESCAPE:
					ga.SaveHighScores("HighScores.txt")
					ga.curMode = STANDBY
					processEvents = ga.ProcessEventsStandBy
				case sdl.K_RETURN:
					ga.SaveHighScores("HighScores.txt")
					ga.curMode = STANDBY
					processEvents = ga.ProcessEventsStandBy
				default:
					c := ga.getChar(keyCode)
					if c != "" && ga.idHighScore >= 0 {
						if len(ga.userName) < 10 {
							ga.userName += c
							ga.highScores[ga.idHighScore].name = ga.userName
						}
					}
				}
			}
		}
	}
	return true
}

func (ga *Game) SaveHighScores(fileName string) {
	//------------------------------------------------------

	var (
		str1 string
	)

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < 10; i++ {
		h := ga.highScores[i]
		if h.name == "" {
			h.name = "XXXX"
		}
		str1 = fmt.Sprintf("%s %d\n", h.name, h.score)
		_, _ = f.WriteString(str1)

	}

}

func (ga *Game) LoadHighScores(fileName string) {
	//------------------------------------------------------

	f, err := os.Open(fileName)

	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for nbL := 0; nbL < 10 && scanner.Scan(); nbL++ {

		//--
		strLineVal := scanner.Text()

		wordBreakDown := strings.Fields(strLineVal)

		ga.highScores[nbL].name = wordBreakDown[0]
		val, _ := strconv.ParseInt(wordBreakDown[1], 10, 32)
		ga.highScores[nbL].score = int(val)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
