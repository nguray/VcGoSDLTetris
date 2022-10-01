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

	game := GameNew()

	game.LoadHighScores("HighScores.txt")

	game.curMode = STANDBY
	processEvents = game.ProcessEventsStandBy

	startH := time.Now()
	startV := startH
	startR := startH

	// game.curTetromino = game.nextTetromino
	// game.curTetromino.x = 5
	// game.curTetromino.y = 5 * cellSize

	// for c := 0; c < NB_COLUMNS; c++ {
	// 	game.board[15*NB_COLUMNS+c] = 1
	// }

	// for r := 0; r < NB_ROWS; r++ {
	// 	if r != 4 {
	// 		game.board[r*NB_COLUMNS+10] = 2
	// 	}
	// }

	// for r := 0; r < NB_ROWS; r++ {
	// 	game.board[r*NB_COLUMNS+1] = 2
	// }

	running := true
	for running {

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

			if id >= 0 {
				//--
				game.InsertHightScore(id, game.userName, game.curScore)
				game.curMode = HIGHSCORES
				processEvents = game.ProcessEventsHightScores
				game.InitGame()
				running = true
			} else {
				//--
				game.InitGame()
				game.curMode = STANDBY
				processEvents = game.ProcessEventsStandBy
				running = true
			}

		}

		if game.curMode == PLAY {

			if game.curTetromino != nil {

				//elapsed := time.Since(startH)
				elapsedV := time.Since(startV)
				elapsedR := time.Since(startR)

				if game.fDrop {
					if elapsedV.Milliseconds() > 10 {
						startV = time.Now()
						for iOffSet := 0; iOffSet < 6; iOffSet++ {
							//-- Move down to check
							game.curTetromino.y++
							idHit := game.curTetromino.HitGround1(renderer, game.board)
							if idHit >= 0 {
								game.curTetromino.y--
								game.FreezeCurTetramino1()
								game.NewTetromino()
								game.fDrop = false
							} else if game.curTetromino.OutBoardLimit1() {
								game.curTetromino.y--
								game.FreezeCurTetramino1()
								game.NewTetromino()
								game.fDrop = false
							}
							if game.fDrop {
								if game.velX != 0 {
									elapsed := time.Since(startH)
									if elapsed.Milliseconds() > 100 {
										game.curTetromino.x += game.velX
										idHit := game.curTetromino.HitGround1(renderer, game.board)
										if idHit >= 0 {
											game.curTetromino.x -= game.velX
										} else if game.curTetromino.OutBoardLimit1() {
											game.curTetromino.x -= game.velX
										} else {
											startH = time.Now()
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
							game.curTetromino.y++
							fMove := true
							idHit := game.curTetromino.HitGround1(renderer, game.board)
							if idHit >= 0 {
								//game.board[idHit] = 0
								game.curTetromino.y--
								game.FreezeCurTetramino1()
								game.NewTetromino()
								fMove = false

							} else if game.curTetromino.OutBoardLimit1() {
								game.curTetromino.y--
								game.FreezeCurTetramino1()
								game.NewTetromino()
								fMove = false
							}
							if fMove {
								if game.velX != 0 {
									elapsed := time.Since(startH)
									if elapsed.Milliseconds() > 100 {
										game.curTetromino.x += game.velX
										idHit := game.curTetromino.HitGround1(renderer, game.board)
										if idHit >= 0 {
											game.curTetromino.x -= game.velX
										} else if game.curTetromino.OutBoardLimit1() {
											game.curTetromino.x -= game.velX
										} else {
											startH = time.Now()
											break
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
								processEvents = game.ProcessEventsHightScores
								game.InitGame()
							} else {
								//--
								game.InitGame()
								game.curMode = STANDBY
								processEvents = game.ProcessEventsStandBy

							}

						}

					}
				}

				if elapsedR.Milliseconds() > 500 {
					startR = time.Now()
					game.nextTetromino.RotateRight()

				}

			}

		}

		// rects = []sdl.Rect{{500, 300, 100, 100}, {200, 300, 200, 200}}
		// renderer.SetDrawColor(255, 0, 255, 255)
		// renderer.FillRects(rects)

		//--
		game.DrawBoard(renderer)
		//--
		game.DrawScore(renderer)

		//--
		if game.curTetromino != nil {
			game.curTetromino.Draw(renderer)
		}
		if game.nextTetromino != nil {
			game.nextTetromino.Draw(renderer)
		}

		if game.curMode == STANDBY {
			game.DrawStandBy(renderer)
		} else if game.curMode == HIGHSCORES {
			game.DrawHightScores(renderer)
		}

		//--
		renderer.Present()

		//sdl.Delay(1)

	}

}
