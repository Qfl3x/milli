package fileops

import (
	"milli/buffer"
)

type Operation = buffer.Operation
type Buffer = buffer.Buffer
type Cursor = buffer.Cursor

type Write struct {
	Index int
	Str string
}

func (*Write) Type() (string) {
	return "WRITE"
}
func (w *Write) Name() (string) {
	return "WRITE: " + string(w.Index) + ";" + w.Str
}

func simulateCursorForward(c Cursor, str string) (Cursor)  {
	x, y := c.X, c.Y
	for _, c := range str {
		if c == '\n' {
			x = 0
			y += 1
		} else {
			x += 1
		}
	}
	return Cursor{X:x, Y:y}
}


func (w *Write) Apply(b *Buffer) (error)  {
	b.Buf.Insert(w.Index, w.Str)
	err, cur := b.Ind2Cursor(w.Index)
	if err != nil {
		return err
	}
	b.Cur = &cur
	return nil
}

func (w *Write) Undo(b *Buffer) (error)  {
	b.Buf.Delete(w.Index, len(w.Str) -1)
	err, cur := b.Ind2Cursor(w.Index)
	if err != nil {
		return err
	}
	b.Cur = &cur
	return nil
}

func (w *Write) Add(s string, op string) (string) {
	w.Str += s
	return w.Str
}

type Delete struct {
	Index int
	Str string
}

func (*Delete) Type() (string) {
	return "DELETE"
}
func (d *Delete) Name() (string) {
	return "DELETE: " + string(d.Index) + ";" + string(len(d.Str))
}

func (d *Delete) Apply(b *Buffer) (error)  {
	length := len(d.Str)
	b.Buf.Delete(d.Index, length - 1)
	err, cur := b.Ind2Cursor(d.Index)
	if err != nil {
		return err
	}
	b.Cur = &cur
	return nil
}

func (d *Delete) Undo(b *Buffer) (error) {
	b.Buf.Insert(d.Index, d.Str)
	err, cur := b.Ind2Cursor(d.Index)
	if err != nil {
		return err
	}
	b.Cur = &cur
	return nil
}

func (d *Delete) Add(s string, op string) (string) {
	if op == "f" {
		d.Str += s
	} else if op == "b" {
		d.Str = s + d.Str
		d.Index -= 1
	}
	return ""
}
