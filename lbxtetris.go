package main

import (
	"time"
	"fmt"
	"math/rand"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	KEY_LEFT  uint = 65361
	KEY_UP    uint = 65362
	KEY_RIGHT uint = 65363
	KEY_DOWN  uint = 65364
	KEY_SPACE uint = 32
)

const (
	WINX_SIZE	= 600
	WINY_SIZE	= 800
)

const (
	INITIAL_SPEED_GAME	=	500
	SPEED_INCREMENT		=	20
)

const (
	BOARD_X_BLOCKS	=	10
	BOARD_Y_BLOCKS	=	20
) 

type gameStatus struct {
	speed int
	boardxsize int
	boardysize int
	unitSizeX  int
	unitSizeY  int
	flagEnd	   int
	board [BOARD_X_BLOCKS][BOARD_Y_BLOCKS] int
	piece [BOARD_X_BLOCKS][BOARD_Y_BLOCKS] int
}


func (g *gameStatus) newpiece( ) {
	
	 for x:=0 ; x < BOARD_X_BLOCKS ; x++ {
			for y:=0 ; y < BOARD_Y_BLOCKS ; y++ {
				g.piece[x][y] = 0
			}
	}
	
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
    


	switch (r.Intn(4)) {
		case 0:	for x:=3 ; x<8; x++ { g.piece[x][0] = 1 }
		case 1:	for x:=4 ; x<7; x++ { g.piece[x][0] = 2 }
				for x:=2 ; x<5; x++ { g.piece[x][1] = 2 }
		case 2: for x:=2 ; x<5; x++ { g.piece[x][0] = 3 }
				for x:=4 ; x<7; x++ { g.piece[x][1] = 3 }
		case 3: for x:=3 ; x<5; x++ { g.piece[x][0] = 4 }
				for x:=3 ; x<5; x++ { g.piece[x][1] = 4 }	
		default: fmt.Println("NOT DEFINED")
	}
	
} 


func (g *gameStatus) move( dir uint ) {
	var k int
	
	var aux[BOARD_X_BLOCKS][BOARD_Y_BLOCKS] int
	
	switch dir {
		case KEY_LEFT:  k = 0
						for y:=0 ; y < BOARD_Y_BLOCKS; y++ { k = k + g.piece[0][y] }
						if  k == 0  {
							for x:=1 ; x < BOARD_X_BLOCKS ; x++ {
								for y:=0 ; y < BOARD_Y_BLOCKS ; y++ {
									aux[x-1][y] = g.piece[x][y]
								}
							}				
							for x:=0 ; x < BOARD_X_BLOCKS ; x++ {
								for y:=0 ; y < BOARD_Y_BLOCKS ; y++ {
									if aux[x][y] != 0 && g.board[x][y] != 0 {
										return
									}
								}
							}
							g.piece = aux
						}
						
						
						
	
		case KEY_RIGHT: k = 0
						for y:=0 ; y < BOARD_Y_BLOCKS; y++ { k = k + g.piece[BOARD_X_BLOCKS-1][y] } 
						if k == 0 {
							for x:=BOARD_X_BLOCKS-2 ; x >= 0 ; x-- {
								for y:=0 ; y < BOARD_Y_BLOCKS ; y++ {
									aux[x+1][y] = g.piece[x][y]
								}
							}
							for x:=0 ; x < BOARD_X_BLOCKS ; x++ {
								for y:=0 ; y < BOARD_Y_BLOCKS ; y++ {
									if aux[x][y] != 0 && g.board[x][y] != 0 {
										return
									}
								}
							}
							g.piece = aux
						}
						
						
		case KEY_DOWN:  k = 0
						for x:=0 ; x < BOARD_X_BLOCKS; x++ { k = k + g.piece[x][BOARD_Y_BLOCKS-1] }
						if ( k == 0 ) { 
							for x:=0 ; x < BOARD_X_BLOCKS ; x++ {
								for y:=BOARD_Y_BLOCKS-2 ; y >=0 ; y-- {
									aux[x][y+1] = g.piece[x][y]
								}
							}
							for x:=0 ; x < BOARD_X_BLOCKS ; x++ {
								for y:=0 ; y < BOARD_Y_BLOCKS ; y++ {
									if aux[x][y] != 0 && g.board[x][y] != 0 {
										return
									}
								}
							}
							g.piece = aux			
						}
						
		case KEY_UP: 
		case KEY_SPACE:	g.move( KEY_DOWN )
						g.move( KEY_DOWN )
						g.move( KEY_DOWN )
						g.move( KEY_DOWN )
						g.move( KEY_DOWN )
						g.move( KEY_DOWN )

		
		default:	fmt.Println (" Tecla ",dir )  
	}
}

func (g *gameStatus) increaseSpeed ( ) {
	g.speed = g.speed - SPEED_INCREMENT
}

func (g *gameStatus) calculateUnitSize () {
	
	g.unitSizeX = g.boardxsize / BOARD_X_BLOCKS
	g.unitSizeY = g.boardysize / BOARD_Y_BLOCKS
	
}

func game( g * gameStatus, win *gtk.Window ) {
	
	for {
		time.Sleep( time.Duration(g.speed) * time.Millisecond)
		g.move ( KEY_DOWN )
		win.QueueDraw()
		if g.flagEnd == 1 {
			fmt.Println ( "the end" )
			return
		}
	}
}

func (g *gameStatus) merge (  ) {
	for x:=0 ; x < BOARD_X_BLOCKS ; x++ {
		for y:=0 ; y < BOARD_Y_BLOCKS ; y++ {
			g.board[x][y] =  g.board[x][y] + g.piece[x][y]
		}
	}
}

func (g *gameStatus) drawBoard ( cr *cairo.Context ) {
	
	cellp := 0
	
	for x, h := range g.board {
		for y, cell := range h {
			cellp = g.piece[x][y]
			if  cellp !=0 && cell !=0 {
				g.flagEnd = 1
				return
			}

			switch ( cell + cellp ) {
				case 0: cr.SetSourceRGB(0, 0, 0)
						cr.Rectangle(float64(x) * float64(g.unitSizeX), float64(y) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.008, float64(g.unitSizeY) - 0.005 )
						cr.Fill()
				case 1: cr.SetSourceRGB(0.7, 0, 0.1)
						cr.Rectangle(float64(x) * float64(g.unitSizeX), float64(y) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.008, float64(g.unitSizeY) - 0.005 )
						cr.Fill()
				case 2: cr.SetSourceRGB(0.3, 0.3, 0.7)
						cr.Rectangle(float64(x) * float64(g.unitSizeX), float64(y) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.008, float64(g.unitSizeY) - 0.005 )
						cr.Fill()
				case 3: cr.SetSourceRGB(0.5, 0.1, 0.6)
						cr.Rectangle(float64(x) * float64(g.unitSizeX), float64(y) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.008, float64(g.unitSizeY) - 0.005 )
						cr.Fill()
				case 4: cr.SetSourceRGB(0.2, 0.7, 0.9)
						cr.Rectangle(float64(x) * float64(g.unitSizeX), float64(y) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.008, float64(g.unitSizeY) - 0.005 )
						cr.Fill()
						
				
			}
		}
    
	}
	
	Loop:
	for x:=0 ; x < BOARD_X_BLOCKS ; x++ {
		for y:=0 ; y < BOARD_Y_BLOCKS ; y++ {
			if g.piece[x][y] != 0 && y == BOARD_Y_BLOCKS -1 {
				g.merge()
				g.newpiece()
				break Loop
			} else {
				if g.piece[x][y] != 0 && g.board[x][y+1] != 0 {
					g.merge()
					g.newpiece()
					break Loop
				}	
			}
		}
	}
	
}

func main() {
	gtk.Init(nil)


	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})


	board, _ := gtk.DrawingAreaNew()
	win.Add(board)

	// Set the default window size.
	win.SetDefaultSize(WINX_SIZE, WINY_SIZE)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	
	// Sizes
	gs := gameStatus{  speed: INITIAL_SPEED_GAME, boardxsize: board.GetAllocatedWidth( ), boardysize: board.GetAllocatedHeight( ) }
	gs.calculateUnitSize()
	gs.newpiece();

	// Event handlers
		
	board.Connect("draw", func(board *gtk.DrawingArea, cr *cairo.Context) {
		gs.drawBoard(cr);
	})
	win.Connect("key-press-event", func(win *gtk.Window, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		gs.move( keyEvent.KeyVal() )
		win.QueueDraw()
		
	})



	go game(&gs,win)

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run. 
	gtk.Main()
}
