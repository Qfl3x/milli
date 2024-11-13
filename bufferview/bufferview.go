package bufferview

import (
	"milli/buffer"
	"github.com/gdamore/tcell/v2"
)


func showLine(s tcell.Screen, line string, y int, defStyle tcell.Style)  {
	for i := 0; i < len(line); i++ {
		s.SetContent(i, y, rune(line[i]), nil, defStyle)
	}
}

type BufferView struct {
	B *buffer.Buffer
	TopLine int
	
}

func (b *BufferView) MoveCursor(x int, y int) (error)  {
	err, lines := b.B.Split()
	if err != nil {
		return err
	}
	nlines := len(lines)
	cur := b.B.Cur
	if y >= nlines - 1 {
		cur.Y = nlines - 1
	} else if y <= 0 {
		cur.Y = 0
	} else {
		cur.Y = y
	}
	linelength := len(lines[cur.Y])

	if x >= linelength {
		cur.X = linelength
	} else if x <= 0 {
		cur.X = 0
	} else {
		cur.X = x
	}
	return nil
}

func (b *BufferView) MoveCursorBack() (error)  {
	err, lines := b.B.Split()
	if err != nil {
		return err
	}
	cur := b.B.Cur
	if cur.X > 0 {
		cur.X -= 1
		return nil
	}
	if cur.Y > 0 {
		cur.Y -= 1
		cur.X = len(lines[cur.Y])
		return nil
	}
	return nil
}

func (b *BufferView) Cursor2ind() (error, int) {
	ind := 0
	err, splitLines := b.B.Split()
	if err != nil {
		return err, 0
	}
	cur := b.B.Cur
	for i := 0; i < cur.Y; i++ {
		ind += len(splitLines[i]) + 1
	}
	ind += cur.X
	return nil, ind
}
func ShowBufferView(s tcell.Screen, view *BufferView, defStyle tcell.Style, cursShow bool) (error) {
	b := view.B
	err, lines := b.Split()
	if err != nil {
		return err
	}
	cur := b.Cur
	_, screenLength := s.Size()
	if cur.Y < view.TopLine {
		view.TopLine = cur.Y
	}
	if cur.Y > view.TopLine + screenLength - 1 {
		view.TopLine = cur.Y - screenLength + 1
	}
	for i := view.TopLine; i < min(len(lines), view.TopLine + screenLength) ; i++ {
		showLine(s, lines[i], i - view.TopLine, defStyle)
	}
	if cursShow {
		s.ShowCursor(cur.X, cur.Y - view.TopLine)
	}
	return nil
}
