/*
 * Copyright (c) 2019 Cesar Lombao <cesar.lombao@gmail.com>
 *
 * This code can be fouind at https://github.com/lombao/golbxtetris
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */


package main

import (
	"time"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"os"
	"os/user"
	"path/filepath"
	"bufio"
	"sync"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	KEY_LEFT  	uint = 65361
	KEY_UP    	uint = 65362
	KEY_RIGHT 	uint = 65363
	KEY_DOWN  	uint = 65364
	KEY_SPACE 	uint = 32
	KEY_PAUSE 	uint = 112
	KEY_NEWGAME uint = 110
)

const (
	WINX_SIZE	= 480
	WINY_SIZE	= 880
)

const (
	INITIAL_SPEED_GAME	=	500
	SPEED_DECREMENT		=	30
)

const (
	BOARD_X_BLOCKS	=	10
	BOARD_Y_BLOCKS	=	20
) 



////////////////////////////////////////////////////
////////////////////////////////////////////////////

const (
	MAX_NUM_HOF_ENTRIES = 5
	FILE_HOF	=	".lbxgames/tetris.dat"
)

type hallOfFame struct {
	name [MAX_NUM_HOF_ENTRIES] string
	points [MAX_NUM_HOF_ENTRIES] int
}


func (h *hallOfFame) readRecords () {
	
	usr, _ := user.Current()
	homedir := usr.HomeDir
	file := filepath.Join(homedir, FILE_HOF)
	
	for a:=0; a<MAX_NUM_HOF_ENTRIES; a++ {
			h.name[a] = "_"
			h.points[a] = 0 
	}
	
	if _, err := os.Stat(file); err == nil {
		// path/to/whatever exists
		file, _ := os.Open(file)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		a:=0
		for scanner.Scan() {
			l := scanner.Text()
			words := strings.Fields(l)
			h.name[a] = words[0]
			h.points[a],_ = strconv.Atoi(words[1])
			a++
		}	

	} 	
}

func (h *hallOfFame) writeRecords () {
	
	usr, _ := user.Current()
	homedir := usr.HomeDir
	file := filepath.Join(homedir, FILE_HOF)
	 
	f, _ := os.Create(file)
    defer f.Close() 
    
    for a:=0 ; a<MAX_NUM_HOF_ENTRIES; a++ {
		p := strconv.Itoa(h.points[a])
		k := h.name[a] + " " + p + "\n"
		f.WriteString(k)
	}
	f.Sync()
}



////////////////////////////////////////////////////
////////////////////////////////////////////////////

type gameStatus struct {
	speed int
	boardxsize int
	boardysize int
	nextpxsize int
	nextpysize int
	unitSizeX  int
	unitSizeY  int
	nextpunitSizeX  int
	nextpunitSizeY  int
	flagEnd	   int
	flagPause  int
	flagDataW  int
	board [BOARD_X_BLOCKS+8][BOARD_Y_BLOCKS+8] int
	piece [4][4] int
	nextpiece [4][4] int
	posX	int
	posY	int
	points 	int
	maxpoints int
	mux	sync.Mutex
}


func (g *gameStatus) firstpiece( ) {
	
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
    p := (r.Intn(19) + 1) 
    
	switch (p) {
		case 1:	g.piece[1][2]  = p; g.piece[2][2]  = p; g.piece[2][1]  = p; g.piece[3][1]  = p; 
		case 2: g.piece[1][0]  = p; g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; 
		
		case 3: g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[1][2]  = p; g.piece[2][2]  = p; 
		
		case 4: g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[3][2]  = p; 
	   	case 5: g.piece[2][0]  = p; g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[1][2]  = p; 
	   	
		case 6: g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[1][2]  = p; g.piece[3][2]  = p; 
	    case 7: g.piece[2][0]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[3][1]  = p; 
		case 8: g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[3][1]  = p; g.piece[2][2]  = p; 
		case 9: g.piece[2][0]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[1][1]  = p; 
		
		case 10: g.piece[0][1]  = p; g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[3][1]  = p; 
		case 11: g.piece[2][0]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[2][3]  = p; 
		
		case 12: g.piece[1][2]  = p; g.piece[2][2]  = p; g.piece[3][2]  = p; g.piece[3][1]  = p; 
		case 13: g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[2][3]  = p; 
		case 14: g.piece[1][2]  = p; g.piece[2][2]  = p; g.piece[3][2]  = p; g.piece[1][3]  = p; 
		case 15: g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[2][3]  = p; g.piece[3][3]  = p; 
		
		case 16: g.piece[1][1]  = p; g.piece[1][2]  = p; g.piece[2][2]  = p; g.piece[3][2]  = p; 
		case 17: g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[2][3]  = p; g.piece[3][1]  = p; 
		case 18: g.piece[1][2]  = p; g.piece[2][2]  = p; g.piece[3][2]  = p; g.piece[3][3]  = p; 
		case 19: g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[2][3]  = p; g.piece[1][3]  = p; 
		
		
		default: fmt.Println("NOT DEFINED")
	}
	
	g.nextpiece = g.piece
	
} 


func (g *gameStatus) newpiece( ) {
	
	
	g.piece = g.nextpiece
	
	for x:=0 ; x < 4 ; x++ {
			for y:=0 ; y < 4 ; y++ {
				g.nextpiece[x][y] = 0
			}
	}
	
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
    p := (r.Intn(19) + 1) 
    
	switch (p) {
		case 1:	g.nextpiece[1][2]  = p; g.nextpiece[2][2]  = p; g.nextpiece[2][1]  = p; g.nextpiece[3][1]  = p; 
		case 2: g.nextpiece[1][0]  = p; g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; 
		
		case 3: g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[1][2]  = p; g.nextpiece[2][2]  = p; 
		
		case 4: g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; g.nextpiece[3][2]  = p; 
	   	case 5: g.nextpiece[2][0]  = p; g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[1][2]  = p; 
	   	
		case 6: g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; g.nextpiece[1][2]  = p; g.nextpiece[3][2]  = p; 
	    case 7: g.nextpiece[2][0]  = p; g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; g.nextpiece[3][1]  = p; 
		case 8: g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[3][1]  = p; g.nextpiece[2][2]  = p; 
		case 9: g.nextpiece[2][0]  = p; g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; g.nextpiece[1][1]  = p; 
		
		case 10: g.nextpiece[0][1]  = p; g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[3][1]  = p; 
		case 11: g.nextpiece[2][0]  = p; g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; g.nextpiece[2][3]  = p; 
		
		case 12: g.nextpiece[1][2]  = p; g.nextpiece[2][2]  = p; g.nextpiece[3][2]  = p; g.nextpiece[3][1]  = p; 
		case 13: g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; g.nextpiece[2][3]  = p; 
		case 14: g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[3][1]  = p; g.nextpiece[1][2]  = p; 
		case 15: g.nextpiece[1][1]  = p; g.nextpiece[1][2]  = p; g.nextpiece[1][3]  = p; g.nextpiece[2][3]  = p; 
		
		case 16: g.nextpiece[1][1]  = p; g.nextpiece[1][2]  = p; g.nextpiece[2][2]  = p; g.nextpiece[3][2]  = p; 
		case 17: g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; g.nextpiece[2][3]  = p; g.nextpiece[3][1]  = p; 
		case 18: g.nextpiece[1][1]  = p; g.nextpiece[2][1]  = p; g.nextpiece[3][1]  = p; g.nextpiece[3][2]  = p; 
		case 19: g.nextpiece[2][1]  = p; g.nextpiece[2][2]  = p; g.nextpiece[2][3]  = p; g.nextpiece[1][3]  = p; 
		
		default: fmt.Println("NOT DEFINED")
	}
	
	g.posY = 0
	g.posX = ( (BOARD_X_BLOCKS+8) / 2 ) - 4
	
	g.points++
	
} 

// returns 1 if move is no valid/collision, if ok returns 0
func (g *gameStatus) checkMove ( ) int {
	
	for  x:=g.posX ; x < g.posX + 4 ; x++  {
		for y := g.posY ; y < g.posY + 4 ; y++ {
			if g.board[x][y] != 0 && g.piece[x-g.posX][y-g.posY] != 0 {
				if g.posY <= 4 { g.flagEnd = 1 }
				return 1
			}
			if g.piece[x-g.posX][y-g.posY] != 0 {
				if x < 4 { return 1 }
				if x >= BOARD_X_BLOCKS+4 { return 1 }
				if y >= BOARD_Y_BLOCKS+4 { return 1 }
			}
			
		}
	}
	return 0
}


func (g *gameStatus) merge (  ) {
	
	for  x:=g.posX ; x < g.posX + 4 ; x++  {
		for y := g.posY ; y < g.posY + 4 ; y++ {
			g.board[x][y] = g.board[x][y] + g.piece[x-g.posX][y-g.posY]
		}
	}

	
	k := 1
	for y := BOARD_Y_BLOCKS + 4 ; y >= 4 ; y-- {
		k = 1
		for x:= 4 ;x < BOARD_X_BLOCKS + 4; x++ {
			k = k * g.board[x][y]
		}
		if k != 0 {
			g.points += 10
			g.speed = g.speed - SPEED_DECREMENT
			fmt.Println("Speed: ",g.speed)
			for hy := y ; hy >= 4 ; hy -- {
				for hx := 4; hx < BOARD_X_BLOCKS + 4; hx++ {
					g.board[hx][hy] = g.board[hx][hy-1]
				}
			}
			y++
		}
		
	}
}

func (g *gameStatus) move( dir uint ) {
	
	g.mux.Lock()
	defer g.mux.Unlock()
	
	var aux1[4][4] int
	var aux2[4][4] int
	
	if g.flagEnd == 1 && dir != KEY_NEWGAME { return }
	if g.flagPause == 1 && dir != KEY_PAUSE { return } 
	
	
	switch dir {
		case KEY_LEFT: 	g.posX = g.posX - 1
						if g.checkMove() == 1 { g.posX++ }
							
	
		case KEY_RIGHT: g.posX = g.posX +1 
						if g.checkMove() == 1 { g.posX-- }
																		
		case KEY_DOWN:  g.posY++
						if g.checkMove() == 1 { 
							if  g.flagEnd == 1 {
								g.posY--
								fmt.Println ( "G A M E  O V E R")
								return
							}
							g.posY--
							g.merge()
							g.newpiece()
						}
						
		case KEY_UP: 	aux1 = g.piece
						for x:=0 ; x < 4 ; x++ {
							for y:=0 ; y < 4 ; y++ {
								aux2[3-y][x] = g.piece[x][y]
							}
						}
						g.piece = aux2
						if g.checkMove() == 1 {
							g.piece = aux1
						}
							
		case KEY_SPACE:	g.posY++
						for g.checkMove() == 0   { g.posY++ }
						g.posY--
						g.merge()
						g.newpiece()
						
		case KEY_PAUSE:	if g.flagPause == 1 { 
							g.flagPause = 0
							fmt.Println("Resume")
						} else { 
							g.flagPause  = 1 
							fmt.Println("Pause")
						}
						
		case KEY_NEWGAME: 
						if g.flagEnd == 1 {
								g.flagEnd = 0 
								g.points = 0 
								g.flagDataW = 0
								for x:=0 ; x < BOARD_X_BLOCKS + 8; x++ {
									for y:=0 ; y< BOARD_Y_BLOCKS + 8 ; y++ {
										g.board[x][y] = 0
									}
								}
								g.newpiece()
						}
						
		default:	fmt.Println (" Key Press ",dir )  
	}
}

func (g *gameStatus) drawPoints ( cr *cairo.Context ) {

	g.mux.Lock()
	defer g.mux.Unlock()
	
		cr.SelectFontFace( "DejaVu Sans", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)	
		cr.MoveTo( 2, 40 )
		cr.SetFontSize(20)
		cr.SetSourceRGB(0.6,0.4,0.3)
		cr.ShowText( "POINTS " )
		cr.SetSourceRGB(0.8,0.2,0)
		cr.ShowText( strconv.Itoa(g.points) )
		cr.ShowText( "   " )
		cr.SetSourceRGB(0.6,0.4,0.3)
		cr.ShowText( "RECORD " )
		cr.SetSourceRGB(0.8,0.2,0)
		cr.ShowText( strconv.Itoa(g.maxpoints) )
		



}



func (g *gameStatus) drawNextPiece ( cr *cairo.Context ) {

	g.mux.Lock()
	defer g.mux.Unlock()
		
	for x:= 0 ; x < 4  ; x++ {
		for y:=0 ; y < 4 ; y++ {
			cell := g.nextpiece[x][y]
			switch ( cell ) {
				case 0:
						cr.SetSourceRGB(0, 0, 0)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX), float64(g.nextpunitSizeY)  )
						cr.Fill()
						
						
				case 1,2: 
						cr.SetSourceRGB(0.3, 0.3, 0.7)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX)- 1, float64(g.nextpunitSizeY) - 1 )
						cr.Fill()
						
				case 3: cr.SetSourceRGB(0.4, 0.7, 0.2)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX)- 1, float64(g.nextpunitSizeY) - 1 )
						cr.Fill()
						
						
				case 4, 5: 
						cr.SetSourceRGB(0.2, 0.7, 0.9)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX)- 1, float64(g.nextpunitSizeY) - 1 )
						cr.Fill()
						
				case 6,7,8, 9:  
						cr.SetSourceRGB(0.5, 0.3, 0.8)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX)- 1, float64(g.nextpunitSizeY) - 1 )
						cr.Fill()	
						
				case 10,11: 
						cr.SetSourceRGB(0.7, 0.1, 0.1)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX)- 1, float64(g.nextpunitSizeY) - 1 )
						cr.Fill()	
						
				case 12,13,14,15:
						cr.SetSourceRGB(0.3, 0.9, 0.4)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX)- 1, float64(g.nextpunitSizeY) - 1 )
						cr.Fill()	
						
				case 16,17,18,19:
						cr.SetSourceRGB(0.6, 0.3, 0.2)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX)- 1, float64(g.nextpunitSizeY) - 1 )
						cr.Fill()	
						
				default:cr.SetSourceRGB(0.5, 0.5, 0.5)
						cr.Rectangle(float64(x) * float64(g.nextpunitSizeX), float64(y) * float64(g.nextpunitSizeY), float64(g.nextpunitSizeX)- 1, float64(g.nextpunitSizeY) - 1 )
						cr.Fill() 			
			}
		}
	}
}

func (g *gameStatus) drawBoard ( cr *cairo.Context ) {
	
	g.mux.Lock()
	defer g.mux.Unlock()
	
	for x:= 4 ; x < BOARD_X_BLOCKS  + 4; x++ {
		for y:=4 ; y < BOARD_Y_BLOCKS + 4; y++ {
			cell := g.board[x][y]
			
			if x >= g.posX && x < g.posX+4 && y >= g.posY && y < g.posY+4 {
				   cell = cell + g.piece[x-g.posX][y-g.posY]
			}

			switch ( cell ) {
				case 0:
						cr.SetSourceRGB(255, 0, 0)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX)  , float64(y-4) * float64(g.unitSizeY) , float64(g.unitSizeX) , float64(g.unitSizeY)  )
						cr.Fill()
						cr.SetSourceRGB(0, 0, 0)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX) - 0.3 , float64(y-4) * float64(g.unitSizeY) , float64(g.unitSizeX) , float64(g.unitSizeY)  )
						cr.Fill()
						cr.SetSourceRGB(255, 1, 1)
				
						
				case 1,2: 
						cr.SetSourceRGB(0.3, 0.3, 0.7)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX), float64(y-4) * float64(g.unitSizeY), float64(g.unitSizeX)- 1, float64(g.unitSizeY) - 1 )
						cr.Fill()
						
				case 3: cr.SetSourceRGB(0.4, 0.7, 0.2)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX), float64(y-4) * float64(g.unitSizeY), float64(g.unitSizeX)- 1, float64(g.unitSizeY) - 1 )
						cr.Fill()
						
						
				case 4, 5: 
						cr.SetSourceRGB(0.2, 0.7, 0.9)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX), float64(y-4) * float64(g.unitSizeY), float64(g.unitSizeX)- 1, float64(g.unitSizeY) - 1 )
						cr.Fill()
						
				case 6,7,8, 9:  
						cr.SetSourceRGB(0.5, 0.3, 0.8)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX), float64(y-4) * float64(g.unitSizeY), float64(g.unitSizeX)- 1, float64(g.unitSizeY) - 1 )
						cr.Fill()	
						
				case 10,11: 
						cr.SetSourceRGB(0.7, 0.1, 0.1)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX), float64(y-4) * float64(g.unitSizeY), float64(g.unitSizeX)- 1, float64(g.unitSizeY) - 1 )
						cr.Fill()	
						
				case 12,13,14,15:
						cr.SetSourceRGB(0.3, 0.9, 0.4)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX), float64(y-4) * float64(g.unitSizeY), float64(g.unitSizeX)- 1, float64(g.unitSizeY) - 1 )
						cr.Fill()	
						
				case 16,17,18,19:
						cr.SetSourceRGB(0.6, 0.3, 0.2)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX), float64(y-4) * float64(g.unitSizeY), float64(g.unitSizeX)- 1, float64(g.unitSizeY) - 1 )
						cr.Fill()	
						
				default:cr.SetSourceRGB(0.5, 0.5, 0.5)
						cr.Rectangle(float64(x-4) * float64(g.unitSizeX), float64(y-4) * float64(g.unitSizeY), float64(g.unitSizeX)- 1, float64(g.unitSizeY) - 1 )
						cr.Fill() 			
			}
		}
    
	}
	
	if ( g.flagPause == 1) { 
		cr.SetSourceRGB(255, 2, 2)
		cr.SelectFontFace( "Courier", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
		
		cr.MoveTo( float64(3) * float64(g.unitSizeX), float64(10) * float64(g.unitSizeY) )
		cr.SetFontSize(42)
		cr.ShowText("PAUSE")
		cr.MoveTo( float64(1) * float64(g.unitSizeX), float64(11) * float64(g.unitSizeY) )
		cr.SetFontSize(24)
		cr.ShowText("Press P again to continue")
	}
	
	if ( g.flagEnd == 1) { 
		cr.SetSourceRGB(1, 0, 0)
		cr.SelectFontFace( "Courier", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
		
		cr.MoveTo( float64(1) * float64(g.unitSizeX), float64(10) * float64(g.unitSizeY) )
		cr.SetFontSize(42)
		cr.ShowText("G A M E  O V E R")
		cr.MoveTo( float64(1) * float64(g.unitSizeX), float64(11) * float64(g.unitSizeY) )
		cr.SetFontSize(24)
		cr.ShowText("Press N to play again")
	}
	
	
}

func (g *gameStatus) calculateUnitSize () {
	
	g.unitSizeX = g.boardxsize / BOARD_X_BLOCKS
	g.unitSizeY = g.boardysize / BOARD_Y_BLOCKS
	
	g.nextpunitSizeX = g.nextpxsize / 4
	g.nextpunitSizeY = g.nextpysize / 4
	
}

func game( g * gameStatus, win *gtk.Window , hof * hallOfFame) {
	
	for {
		if  g.flagEnd == 1  {
			g.maxpoints = hof.points[0]
			if g.flagDataW == 0 {
				L:
				for k :=0 ; k < MAX_NUM_HOF_ENTRIES ; k++ {
					if g.points > hof.points[k] {
						for z := MAX_NUM_HOF_ENTRIES - 1; z > k ; z-- {
							hof.name[z] = hof.name[z-1]
							hof.points[z] = hof.points[z-1]
						}
						hof.name[k] = "_"
						hof.points[k] = g.points
						hof.writeRecords()
						g.flagDataW = 1
						break L
					}
				}
			} 	
		} 
		
		time.Sleep( time.Duration(g.speed) * time.Millisecond)
		
		if g.flagEnd == 0 {
			g.move ( KEY_DOWN )
			win.QueueDraw()
		}
		
	}
}

func main() {

	hof := hallOfFame { }
	hof.readRecords()


	gtk.Init(nil)


	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("GO LBX TETRIS")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Set the default window size.
	win.SetDefaultSize(WINX_SIZE, WINY_SIZE)

	// Create a new grid widget to arrange child widgets
	grid, _ := gtk.GridNew()
	grid.SetRowSpacing( 10 )
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	
	leftbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL,50)
	
	board, _ := gtk.DrawingAreaNew()
	nextp, _ := gtk.DrawingAreaNew()
    nextp.SetSizeRequest(70,70)
    points,_ := gtk.DrawingAreaNew()
    points.SetSizeRequest(300,70)
    
	
	leftbox.Add(nextp)
	leftbox.Add(points)

	
	grid.Add(leftbox)
	grid.Add(board)
			
	board.SetHExpand(true)
	board.SetVExpand(true)

	
	

	// Recursively show all widgets contained in this window.
	win.Add(grid)
	win.ShowAll()

	
	// Sizes
	gs := gameStatus{  	speed: INITIAL_SPEED_GAME, 
						boardxsize: board.GetAllocatedWidth( ), 
						boardysize: board.GetAllocatedHeight( ), 
						nextpxsize: nextp.GetAllocatedWidth( ), 
						nextpysize: nextp.GetAllocatedHeight( ),
					}
	gs.calculateUnitSize()
	gs.firstpiece()
	gs.newpiece()
	gs.maxpoints = hof.points[0]

	// Event handlers
		
	board.Connect("draw", func(board *gtk.DrawingArea, cr *cairo.Context) {
		gs.drawBoard(cr);
	})
	nextp.Connect("draw", func(nextp *gtk.DrawingArea, cr *cairo.Context) {
		gs.drawNextPiece(cr);
	})
	points.Connect("draw", func(points *gtk.DrawingArea, cr *cairo.Context) {
		gs.drawPoints(cr);
	})
	win.Connect("key-press-event", func(win *gtk.Window, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		gs.move( keyEvent.KeyVal() )
		
		win.QueueDraw()	
	})



	go game(&gs,win,&hof)

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run. 
	gtk.Main()
}
