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
	board [BOARD_X_BLOCKS+4][BOARD_Y_BLOCKS+4] int
	piece [4][4] int
	posX	int
	posY	int
}


func (g *gameStatus) newpiece( ) {
	
	 for x:=0 ; x < 4 ; x++ {
			for y:=0 ; y < 4 ; y++ {
				g.piece[x][y] = 0
			}
	}
	
	s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
    p := (r.Intn(8) + 1) 
    
	switch (p) {
		case 1:	g.piece[1][2]  = p; g.piece[2][2]  = p; g.piece[2][1]  = p; g.piece[3][1]  = p; 
		case 2: g.piece[1][0]  = p; g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; 
		case 3: g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[1][2]  = p; g.piece[2][2]  = p; 
		case 4: g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[3][2]  = p; 
	   	case 5: g.piece[2][0]  = p; g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[1][2]  = p; 
		case 6: g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[1][2]  = p; g.piece[3][2]  = p; 
	    case 7: g.piece[2][0]  = p; g.piece[2][1]  = p; g.piece[2][2]  = p; g.piece[3][1]  = p; 
		case 8: g.piece[1][1]  = p; g.piece[2][1]  = p; g.piece[3][1]  = p; g.piece[2][2]  = p; 
		default: fmt.Println("NOT DEFINED")
	}
	
	g.posY = 0
	g.posX = ( (BOARD_X_BLOCKS+4) / 2 ) - 2
	fmt.Println ( "Created new piece ",p,g.posX,g.posY )
} 

// returns 1 if move is no valid/collision, if ok returns 0
func (g *gameStatus) checkMove ( ) int {
	
	for  x:=g.posX ; x < g.posX + 4 ; x++  {
		for y := g.posY ; y < g.posY + 4 ; y++ {
			if g.board[x][y] != 0 && g.piece[x-g.posX][y-g.posY] != 0 {
				return 1
			}
			if g.piece[x-g.posX][y-g.posY] != 0 {
				if x < 2 { return 1 }
				if x >= BOARD_X_BLOCKS+2 { return 1 }
				if y >= BOARD_Y_BLOCKS+2 { return 1 }
			}
			
		}
	}
	fmt.Println ( "Returing checkmove 0")
	return 0
}


func (g *gameStatus) merge (  ) {
	
	for  x:=g.posX ; x < g.posX + 4 ; x++  {
		for y := g.posY ; y < g.posY + 4 ; y++ {
			g.board[x][y] = g.board[x][y] + g.piece[x-g.posX][y-g.posY]
		}
	}
}

func (g *gameStatus) move( dir uint ) {
	
	var aux1[4][4] int
	var aux2[4][4] int
	
	switch dir {
		case KEY_LEFT: 	g.posX = g.posX - 1
						if g.checkMove() == 1 { g.posX++ }
							
	
		case KEY_RIGHT: g.posX = g.posX +1 
						if g.checkMove() == 1 { g.posX-- }
																		
		case KEY_DOWN:  g.posY++
						if g.checkMove() == 1 { 
							fmt.Println ( "we reach the end") 
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

		default:	fmt.Println (" Tecla ",dir )  
	}
}



func (g *gameStatus) drawBoard ( cr *cairo.Context ) {
	
	for x:= 2 ; x < BOARD_X_BLOCKS  + 2; x++ {
		for y:=2 ; y < BOARD_Y_BLOCKS + 2; y++ {
			cell := g.board[x][y]
			
			if x >= g.posX && x < g.posX+4 && y >= g.posY && y < g.posY+4 {
				   cell = cell + g.piece[x-g.posX][y-g.posY]
			}

			switch ( cell ) {
				case 0: cr.SetSourceRGB(0, 0, 0)
						cr.Rectangle(float64(x-2) * float64(g.unitSizeX), float64(y-2) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.018, float64(g.unitSizeY) - 0.015 )
						cr.Fill()
				case 1: cr.SetSourceRGB(0.7, 0, 0.1)
						cr.Rectangle(float64(x-2) * float64(g.unitSizeX), float64(y-2) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.018, float64(g.unitSizeY) - 0.015 )
						cr.Fill()
				case 2: cr.SetSourceRGB(0.3, 0.3, 0.7)
						cr.Rectangle(float64(x-2) * float64(g.unitSizeX), float64(y-2) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.018, float64(g.unitSizeY) - 0.015 )
						cr.Fill()
				case 3: cr.SetSourceRGB(0.5, 0.1, 0.6)
						cr.Rectangle(float64(x-2) * float64(g.unitSizeX), float64(y-2) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.018, float64(g.unitSizeY) - 0.015 )
						cr.Fill()
				case 4: cr.SetSourceRGB(0.2, 0.7, 0.9)
						cr.Rectangle(float64(x-2) * float64(g.unitSizeX), float64(y-2) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.018, float64(g.unitSizeY) - 0.015 )
						cr.Fill()
				default:cr.SetSourceRGB(0.5, 0.5, 0.5)
						cr.Rectangle(float64(x-2) * float64(g.unitSizeX), float64(y-2) * float64(g.unitSizeY), float64(g.unitSizeX)- 0.018, float64(g.unitSizeY) - 0.015 )
						cr.Fill() 			
			}
		}
    
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
