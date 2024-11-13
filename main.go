package main

import (
	"io"
	"log"
	"milli/buffer"
	"milli/bufferview"
	"milli/fileops"
	"os"

	"github.com/gdamore/tcell/v2"
)

type UI struct {
	Writer io.Writer
	Reader io.Reader
}

func NewUI() *UI {
	return &UI{os.Stdout, os.Stdin}
}



func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Clear screen
	s.Clear()
	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	binit := buffer.InitializeBufferFile("test.file")
	view := bufferview.BufferView{B:binit, TopLine: 0}
	b := view.B
	var cursX int
	var cursY int
	file, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	logger := log.New(file, "",  1)
	var lastOperation buffer.Operation
	cursShow := true
	for {
		s.Clear()
		cursX, cursY = view.B.Cur.X, view.B.Cur.Y
		bufferview.ShowBufferView(s, &view, defStyle, cursShow)
		cursShow = !cursShow
		s.Show()
		err, ind := view.Cursor2ind()
		if err != nil {
			panic(err)
		}
		ev := s.PollEvent()

		logger.Println("Last Operation", lastOperation)
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			key := ev.Key()
			logger.Println("Key pressed", tcell.KeyNames[key])
			logger.Println("BS", tcell.KeyBackspace)
			// logger.Println("Cursor position :", cursX, cursY)
			switch key {
			case tcell.KeyEscape:
				return
			case tcell.KeyUp:
				cursY -= 1
				view.B.MoveCursor(cursX, cursY)
				lastOperation = nil
			case tcell.KeyDown:
				cursY += 1
				view.B.MoveCursor(cursX, cursY)
				lastOperation = nil
			case tcell.KeyLeft:
				cursX -= 1
				view.B.MoveCursor(cursX, cursY)
				lastOperation = nil
			case tcell.KeyRight:
				cursX += 1
				view.B.MoveCursor(cursX, cursY)
				lastOperation = nil
			case tcell.KeyBackspace2:
				err, rune := b.Buf.Getindex(ind - 1)
				if err != nil {
					panic(err)
				}
				b.Delete(ind - 1)
				logger.Println("Rune Deleted :", string(rune))
				if lastOperation != nil && lastOperation.Type() == "DELETE" {
					lastOperation.Add(string(rune), "b")
				} else {
					op := &fileops.Delete{Index: ind - 1, Str: string(rune)}
					b.Ops[b.OpsPointer] = op
					b.OpsPointer += 1
					lastOperation = op
				}
				view.MoveCursorBack()
			case tcell.KeyDelete:
				err, rune := b.Buf.Getindex(ind )
				if err != nil {
					panic(err)
				}
				b.Delete(ind)
				if lastOperation != nil && lastOperation.Type() == "DELETE" {
					lastOperation.Add(string(rune), "f")
				} else {
					op := &fileops.Delete{Index: ind, Str: string(rune)}
					b.Ops[b.OpsPointer] = op
					b.OpsPointer += 1
					lastOperation = op
				}
			case tcell.KeyEnter:
				rune := byte('\n')
				if lastOperation != nil && lastOperation.Type() == "WRITE" {
					lastOperation.Add(string(rune), "")
				} else {
					op := &fileops.Write{Index: ind, Str: string(rune)}
					b.Ops[b.OpsPointer] = op
					b.OpsPointer += 1
					lastOperation = op
				}
				b.Insert(ind, rune)
				cursX = 0
				cursY += 1
				view.MoveCursor(cursX, cursY)
			case tcell.KeyRune:
				rune := ev.Rune()
				b.Insert(ind, byte(rune))
				if lastOperation != nil && lastOperation.Type() == "WRITE" {
					lastOperation.Add(string(rune), "")
				} else {
					op := &fileops.Write{Index: ind, Str: string(rune)}
					b.Ops[b.OpsPointer] = op
					b.OpsPointer += 1
					lastOperation = op
				}
				cursX += 1
				view.MoveCursor(cursX, cursY)
			case tcell.KeyCtrlZ:
				if b.OpsPointer > 0{
					last := b.Ops[b.OpsPointer - 1]
					last.Undo(b)
					b.OpsPointer -= 1
					lastOperation = nil
					s.Sync()
				}
			}
		}
	}
}

