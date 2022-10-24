/*--------------------------------------------*\
		Simple Tetris using sdl2
                 2022
			Raymond NGUYEN THANH
\*--------------------------------------------*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type GameMode int

const (
	STANDBY GameMode = iota
	PLAY
	GAMEPAUSE
	GAMEOVER
	HIGHSCORES
)

const (
	LEFT       = 10
	TOP        = 10
	NB_ROWS    = 20
	NB_COLUMNS = 12
	WIN_WIDTH  = 480
	WIN_HEIGHT = 560
	TITLE      = "Go SDL2 Tetris"
)

type HightScore struct {
	name  string
	score int
}

func HightScoreNew(userName string, scoreVal int) *HightScore {

	score := &HightScore{name: userName, score: scoreVal}
	//--
	return score
}

type ProcessEvents func(renderer *sdl.Renderer) bool

var (
	cellSize        int32
	myRand          *rand.Rand
	processEvents   ProcessEvents
	tt_font         *ttf.Font
	succes_sound    *mix.Chunk
	idtetrominosBag int
	tetrominosBag   []int32
	game            *Game
	curTetromino    *Shape
	nextTetromino   *Shape
)

func ComputeScore(nbLines int) int {
	var score int
	switch nbLines {
	case 0:
		score = 0
	case 1:
		score = 40
	case 2:
		score = 100
	case 3:
		score = 300
	case 4:
		score = 1200
	default:
		score = 2000
	}
	return score
}

func TetrisRandomizer() int32 {

	var (
		iSrc int32
		ityp int32
	)

	if idtetrominosBag < 14 {
		ityp = tetrominosBag[idtetrominosBag]
		idtetrominosBag += 1
	} else {
		//-- Shuttle bag
		for i := 0; i < 14; i++ {
			iSrc = int32(myRand.Intn(14))
			ityp = tetrominosBag[iSrc]
			tetrominosBag[iSrc] = tetrominosBag[0]
			tetrominosBag[0] = ityp
		}
		ityp = tetrominosBag[0]
		idtetrominosBag = 1
	}

	return ityp
}

func IsOutLeftBoardLimit(tetro *Shape) bool {
	l := tetro.MinX1()*cellSize + tetro.x
	return l < 0
}

func IsOutRightBoardLimit(tetro *Shape) bool {
	r := tetro.MaxX1()*cellSize + cellSize + tetro.x
	return r > NB_COLUMNS*cellSize
}

func IsOutBottomLimit(tetro *Shape) bool {
	//--------------------------------------------------
	b := tetro.MaxY1()*cellSize + cellSize + tetro.y
	return b > NB_ROWS*cellSize
}

func HitGround(tetro *Shape, board []int) bool {
	var (
		ix int32
		iy int32
	)
	//--------------------------------------------------

	for _, v := range tetro.v {

		x := v.x*cellSize + tetro.x + 1
		y := v.y*cellSize + tetro.y + 1
		ix = int32(x / cellSize)
		iy = int32(y / cellSize)
		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			iHit := iy*NB_COLUMNS + ix
			v := board[iHit]
			if v != 0 {
				return true
			}
		}

		x = v.x*cellSize + cellSize - 1 + tetro.x
		y = v.y*cellSize + tetro.y + 1
		ix = int32(x / cellSize)
		iy = int32(y / cellSize)
		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			iHit := iy*NB_COLUMNS + ix
			v := board[iHit]
			if v != 0 {
				return true
			}
		}

		x = v.x*cellSize + cellSize - 1 + tetro.x
		y = v.y*cellSize + cellSize - 1 + tetro.y
		ix = int32(x / cellSize)
		iy = int32(y / cellSize)
		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			iHit := iy*NB_COLUMNS + ix
			v := board[iHit]
			if v != 0 {
				return true
			}
		}

		x = v.x*cellSize + tetro.x + 1
		y = v.y*cellSize + cellSize - 1 + tetro.y
		ix = int32(x / cellSize)
		iy = int32(y / cellSize)

		if (ix >= 0) && ix < NB_COLUMNS && (iy >= 0) && (iy < NB_ROWS) {
			iHit := iy*NB_COLUMNS + ix
			v := board[iHit]
			if v != 0 {
				return true
			}
		}

	}

	return false
}

func NewTetromino() {
	//--------------------------------------------------
	curTetromino = nextTetromino
	curTetromino.x = 6 * cellSize
	curTetromino.y = 0
	curTetromino.y = -curTetromino.MaxY1() * cellSize
	nextTetromino = ShapeNew(TetrisRandomizer(), (NB_COLUMNS+3)*cellSize, 10*cellSize)

}

func ProcessEventsPlay(renderer *sdl.Renderer) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			game.fQuitGame = true
			return false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym

			//keys := ""
			if t.State == sdl.PRESSED && t.Repeat == 0 {
				switch keyCode {
				case sdl.K_p:
					game.fPause = !game.fPause
				case sdl.K_LEFT:
					game.velX = -1
				case sdl.K_RIGHT:
					game.velX = 1
				case sdl.K_UP:
					if curTetromino != nil {
						curTetromino.RotateLeft()

						if HitGround(curTetromino, game.board) {
							//-- Undo Rotate
							curTetromino.RotateRight()

						} else if IsOutRightBoardLimit(curTetromino) {
							backupX := curTetromino.x
							//-- Move tetromino inside board
							for IsOutRightBoardLimit(curTetromino) {
								curTetromino.x--
							}
							if HitGround(curTetromino, game.board) {
								curTetromino.x = backupX
								//-- Undo Rotate
								curTetromino.RotateRight()
							}

						} else if IsOutLeftBoardLimit(curTetromino) {

							backupX := curTetromino.x
							//-- Move tetromino inside board
							for IsOutLeftBoardLimit(curTetromino) {
								curTetromino.x++
							}
							if HitGround(curTetromino, game.board) {
								curTetromino.x = backupX
								//-- Undo Rotate
								curTetromino.RotateRight()
							}

						}

					}
				case sdl.K_DOWN:
					game.fFastDown = true
				case sdl.K_SPACE:
					if curTetromino != nil {
						//-- Drop current Tetromino
						game.fDrop = true
					}
				case sdl.K_ESCAPE:
					return false
				}
			} else if t.State == sdl.RELEASED {
				switch keyCode {
				case sdl.K_LEFT:
					game.velX = 0
				case sdl.K_RIGHT:
					game.velX = 0
				case sdl.K_DOWN:
					game.fFastDown = false
				}

			}

		}
	}
	return true
}

func ProcessEventsStandBy(renderer *sdl.Renderer) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			game.fQuitGame = true
			println("Quit")
			return false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym
			if t.State == sdl.PRESSED && t.Repeat == 0 {
				switch keyCode {
				case sdl.K_SPACE:
					game.curMode = PLAY
					processEvents = ProcessEventsPlay
					NewTetromino()
				case sdl.K_ESCAPE:
					game.fQuitGame = true
					return false
				}
			}

		}
	}
	return true
}

func ProcessEventsGameOver(renderer *sdl.Renderer) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			game.fQuitGame = true
			println("Quit")
			return false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym
			if t.State == sdl.PRESSED && t.Repeat == 0 {
				switch keyCode {
				case sdl.K_SPACE:
					game.curMode = STANDBY
					processEvents = ProcessEventsStandBy
				case sdl.K_ESCAPE:
					game.fQuitGame = true
					return false
				}
			}

		}
	}
	return true
}

func ProcessEventsHightScores(renderer *sdl.Renderer) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			game.fQuitGame = true
			println("Quit")
			return false
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym
			if t.State == sdl.PRESSED && t.Repeat == 0 {
				switch keyCode {
				case sdl.K_BACKSPACE:
					sz := len(game.userName)
					if sz > 0 {
						game.userName = game.userName[:sz-1]
						game.highScores[game.idHighScore].name = game.userName
					}
				case sdl.K_ESCAPE:
					game.SaveHighScores("HighScores.txt")
					game.curMode = STANDBY
					processEvents = ProcessEventsStandBy
				case sdl.K_RETURN:
					game.SaveHighScores("HighScores.txt")
					game.curMode = STANDBY
					processEvents = ProcessEventsStandBy
				default:
					c := game.getChar(keyCode)
					if c != "" && game.idHighScore >= 0 {
						if len(game.userName) < 10 {
							game.userName += c
							game.highScores[game.idHighScore].name = game.userName
						}
					}
				}
			}
		}
	}
	return true
}

func FreezeCurTetramino1() {
	//--------------------------------------------------
	if curTetromino != nil {
		ix := int32((curTetromino.x + 1) / cellSize)
		iy := int32((curTetromino.y + 1) / cellSize)
		for _, v := range curTetromino.v {
			x := v.x + ix
			y := v.y + iy
			if x >= 0 && x < NB_COLUMNS && y >= 0 && y < NB_ROWS {
				game.board[y*NB_COLUMNS+x] = int(curTetromino.typ)
			}
		}
		//--
		game.nbCompledLines = game.ComputeCompletedLines()
		if game.nbCompledLines > 0 {
			game.curScore += ComputeScore(game.nbCompledLines)
		}

	}
}

func main() {

	var renderer *sdl.Renderer

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIN_WIDTH, WIN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	ttf.Init()
	defer ttf.Quit()

	curDir, _ := os.Getwd()
	fullPathName := filepath.Join(curDir, "resources", "sansation.ttf")
	tt_font, err = ttf.OpenFont(fullPathName, 18)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load Font: %s\n", err)
		panic(err)
	}
	defer tt_font.Close()
	tt_font.SetStyle(ttf.STYLE_ITALIC | ttf.STYLE_BOLD)

	fullPathName = filepath.Join(curDir, "resources", "Tetris.wav")
	mix.OpenAudio(44100, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, 1024)
	tetris_music, err := mix.LoadMUS(fullPathName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load Music : %s\n", err)
		panic(err)
	}
	defer tetris_music.Free()
	tetris_music.Play(-1)
	mix.VolumeMusic(20)

	fullPathName = filepath.Join(curDir, "resources", "109662__grunz__success.wav")
	succes_sound, err = mix.LoadWAV(fullPathName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load Sound: %s\n", err)
		panic(err)
	}
	defer succes_sound.Free()
	mix.Volume(-1, 10)

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	//renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		//return 2
		panic(err)
	}
	defer renderer.Destroy()

	var rect sdl.Rect
	//var rects []sdl.Rect

	//--
	tetrominosBag = make([]int32, 14)
	tetrominosBag[0] = 1
	tetrominosBag[1] = 2
	tetrominosBag[2] = 3
	tetrominosBag[3] = 4
	tetrominosBag[4] = 5
	tetrominosBag[5] = 6
	tetrominosBag[6] = 7
	tetrominosBag[7] = 1
	tetrominosBag[8] = 2
	tetrominosBag[9] = 3
	tetrominosBag[10] = 4
	tetrominosBag[11] = 5
	tetrominosBag[12] = 6
	tetrominosBag[13] = 7
	idtetrominosBag = 14

	cellSize = int32(WIN_WIDTH / (NB_COLUMNS + 7))

	InitTetrominos()
	myRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	game = GameNew()
	curTetromino = nil
	nextTetromino = ShapeNew(int32(myRand.Intn(7)+1), (NB_COLUMNS+3)*cellSize, 10*cellSize)

	game.LoadHighScores("HighScores.txt")

	game.curMode = STANDBY
	processEvents = ProcessEventsStandBy

	startH := time.Now()
	startV := startH
	startR := startH

	running := true
	for running {

		//-- Draw Background
		renderer.SetDrawColor(48, 48, 255, 255)
		renderer.Clear()

		rect = sdl.Rect{X: int32(LEFT), Y: int32(TOP), W: int32(cellSize * NB_COLUMNS), H: int32(cellSize * NB_ROWS)}
		renderer.SetDrawColor(10, 10, 100, 255)
		renderer.FillRect(&rect)

		//-- Process current mode Events
		running = processEvents(renderer)

		if game.fQuitGame {
			break
		}

		if !running {
			//--
			id := game.IsHightScore(game.curScore)

			//-- Manage Game Over and User Escape
			if id >= 0 {
				//--
				game.InsertHightScore(id, game.userName, game.curScore)
				game.curMode = HIGHSCORES
				processEvents = ProcessEventsHightScores
				game.InitGame()
				curTetromino = nil
				nextTetromino = ShapeNew(int32(myRand.Intn(7)+1), (NB_COLUMNS+3)*cellSize, 10*cellSize)
				running = true
			} else {
				//--
				game.InitGame()
				curTetromino = nil
				nextTetromino = ShapeNew(int32(myRand.Intn(7)+1), (NB_COLUMNS+3)*cellSize, 10*cellSize)
				game.curMode = STANDBY
				processEvents = ProcessEventsStandBy
				running = true
			}

		}

		//-- Game Mode Update States
		if game.curMode == PLAY {

			if curTetromino != nil && !game.fPause {

				elapsedV := time.Since(startV)
				elapsedR := time.Since(startR)

				if game.nbCompledLines > 0 {
					if elapsedV.Milliseconds() > 250 {
						startV = time.Now()
						game.nbCompledLines--
						game.EraseFirstCompletedLine()
						succes_sound.Play(-1, 0)
					}

				} else if game.horizontalMove != 0 {
					elapsed := time.Since(startH)
					if elapsed.Milliseconds() > 20 {
						startH = time.Now()

						for iOffSet := 0; iOffSet < int(4); iOffSet++ {

							backupX := curTetromino.x
							curTetromino.x += game.horizontalMove

							if game.horizontalMove < 0 {
								if IsOutLeftBoardLimit(curTetromino) {
									curTetromino.x = backupX
									game.horizontalMove = 0
									break
								} else {
									if HitGround(curTetromino, game.board) {
										curTetromino.x = backupX
										game.horizontalMove = 0
										break
									}
								}
							} else if game.horizontalMove > 0 {
								if IsOutRightBoardLimit(curTetromino) {
									curTetromino.x = backupX
									game.horizontalMove = 0
									break
								} else {
									if HitGround(curTetromino, game.board) {
										curTetromino.x = backupX
										game.horizontalMove = 0
										break
									}
								}

							}

							if game.horizontalMove != 0 {
								if game.horizontalStartColumn != curTetromino.Column() {
									curTetromino.x = backupX
									game.horizontalMove = 0
									startH = time.Now()
									break
								}
							}

						}
					}

				} else if game.fDrop {

					if elapsedV.Milliseconds() > 10 {
						startV = time.Now()
						for iOffSet := 0; iOffSet < 6; iOffSet++ {
							//-- Move down to check
							curTetromino.y++
							if HitGround(curTetromino, game.board) {
								curTetromino.y--
								FreezeCurTetramino1()
								NewTetromino()
								game.fDrop = false
							} else if IsOutBottomLimit(curTetromino) {
								curTetromino.y--
								FreezeCurTetramino1()
								NewTetromino()
								game.fDrop = false
							}
							if game.fDrop {
								if game.velX != 0 {
									elapsed := time.Since(startH)

									if elapsed.Milliseconds() > 20 {

										backupX := curTetromino.x
										curTetromino.x += game.velX

										if game.velX < 0 {
											if IsOutLeftBoardLimit(curTetromino) {
												curTetromino.x = backupX
											} else {
												if HitGround(curTetromino, game.board) {
													curTetromino.x = backupX
												} else {
													startH = time.Now()
													game.horizontalMove = game.velX
													game.horizontalStartColumn = curTetromino.Column()
													break
												}
											}
										} else if game.velX > 0 {
											if IsOutRightBoardLimit(curTetromino) {
												curTetromino.x = backupX
											} else {
												if HitGround(curTetromino, game.board) {
													curTetromino.x = backupX
												} else {
													startH = time.Now()
													game.horizontalMove = game.velX
													game.horizontalStartColumn = curTetromino.Column()
													break
												}
											}

										}
									}

								}
							}
						}
					}

				} else {

					var limitElapse int64 = 25
					if game.fFastDown {
						limitElapse = 10
					}
					if elapsedV.Milliseconds() > limitElapse {
						startV = time.Now()

						for iOffSet := 0; iOffSet < 3; iOffSet++ {
							//-- Move down to check
							curTetromino.y++
							fMove := true
							if HitGround(curTetromino, game.board) {
								curTetromino.y--
								FreezeCurTetramino1()
								NewTetromino()
								fMove = false

							} else if IsOutBottomLimit(curTetromino) {
								curTetromino.y--
								FreezeCurTetramino1()
								NewTetromino()
								fMove = false
							}
							if fMove {
								if game.velX != 0 {
									elapsed := time.Since(startH)
									if elapsed.Milliseconds() > 15 {
										startH = time.Now()

										backupX := curTetromino.x
										curTetromino.x += game.velX

										if game.velX < 0 {
											if IsOutLeftBoardLimit(curTetromino) {
												curTetromino.x = backupX
											} else {
												if HitGround(curTetromino, game.board) {
													curTetromino.x = backupX
												} else {
													game.horizontalMove = game.velX
													game.horizontalStartColumn = curTetromino.Column()
													break
												}
											}
										} else if game.velX > 0 {
											if IsOutRightBoardLimit(curTetromino) {
												curTetromino.x = backupX
											} else {
												if HitGround(curTetromino, game.board) {
													curTetromino.x = backupX
												} else {
													game.horizontalMove = game.velX
													game.horizontalStartColumn = curTetromino.Column()
													break
												}
											}

										}
									}
								}
							}
						}

					}
				}

				//-- Check Game Over
				if game.IsGameOver() {

					//--
					id := game.IsHightScore(game.curScore)

					if id >= 0 {
						//--
						game.InsertHightScore(id, game.userName, game.curScore)
						game.curMode = HIGHSCORES
						processEvents = ProcessEventsHightScores
						game.InitGame()
						curTetromino = nil
					} else {
						//--
						game.curMode = GAMEOVER
						processEvents = ProcessEventsGameOver
						game.InitGame()
						curTetromino = nil

					}

				}

				if elapsedR.Milliseconds() > 500 {
					startR = time.Now()
					nextTetromino.RotateRight()

				}

			}

		}

		// rects = []sdl.Rect{{500, 300, 100, 100}, {200, 300, 200, 200}}
		// renderer.SetDrawColor(255, 0, 255, 255)
		// renderer.FillRects(rects)

		//------------------------------------------------------------
		//-- Draw Game

		//--
		game.DrawBoard(renderer)
		//--
		game.DrawScore(renderer)

		//--
		if curTetromino != nil {
			curTetromino.Draw(renderer)
		}
		if nextTetromino != nil {
			nextTetromino.Draw(renderer)
		}

		if game.curMode == STANDBY {
			game.DrawStandBy(renderer)

		} else if game.curMode == GAMEOVER {
			game.DrawGameOver(renderer)

		} else if game.curMode == HIGHSCORES {
			elapsedV := time.Since(startV)
			if elapsedV.Milliseconds() > 200 {
				startV = time.Now()
				game.iColorHighScore++
			}
			game.DrawHightScores(renderer)

		}

		//--
		renderer.Present()

		//sdl.Delay(1)

	}

}
