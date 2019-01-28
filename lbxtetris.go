package main

import (
	"time"
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	KEY_LEFT  uint = 65361
	KEY_UP    uint = 65362
	KEY_RIGHT uint = 65363
	KEY_DOWN  uint = 65364
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

type piece struct {
	board [BOARD_X_BLOCKS][BOARD_Y_BLOCKS] int
	typepiece int
}

type gameStatus struct {
	X int
	Y int
	speed int
	boardxsize int
	boardysize int
	unitSizeX  int
	unitSizeY  int
	flagEnd	   int
	board [BOARD_X_BLOCKS][BOARD_Y_BLOCKS] int
}

func (g *gameStatus) move( dir uint ) {
	switch dir {
		case KEY_LEFT:  if g.X > 0 { g.X-- }
		case KEY_UP:    if g.Y > 0 { g.Y-- }
		case KEY_RIGHT: if g.X < BOARD_X_BLOCKS -1 { g.X++ }
		case KEY_DOWN:  if g.Y < BOARD_Y_BLOCKS -1 { 
								g.Y++ 
						} else {
								g.flagEnd = 1
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
			fmt.Println ( "se acabo" )
			return
		} else {
			fmt.Println("Movimiento", g.X, g.Y)				
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

	// Event handlers
	board.Connect("draw", func(board *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(0, 0, 0)
		cr.Rectangle(float64(gs.X) * float64(gs.unitSizeX), float64(gs.Y) * float64(gs.unitSizeY), float64(gs.unitSizeX), float64(gs.unitSizeY) )
		cr.Fill()
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
