package buffer

import (
	"fmt"
	"os"
	"milli/rope"
)

type Operation interface {
	Name() string
	Type() string
	Add(s string, op string) string
	Apply(b *Buffer) error
	Undo(b *Buffer) error
}

type Cursor struct {
	X int
	Y int
}

type Buffer struct {
	Buf  rope.Rope
	Path string
	Ops []Operation
	OpsPointer int
	Cur *Cursor
}

func (b *Buffer) Write() error {
	err, s := b.Buf.String()
	if err != nil {
		return err
	}
	err = os.WriteFile(b.Path, []byte(s), 0666)
	return err
}


func InitializeBufferFile(path string) *Buffer {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Err")
	}
	return &Buffer{rope.NewRope(string(content)), path, make([]Operation, 256), 0, &Cursor{0,0}}
}

func split(buf string) ([]string) {
	var ret []string
	var cur string
	for i := 0; i < len(buf); i++ {
		if buf[i] == '\n' {
			ret = append(ret, cur)
			cur = ""
		} else {
			cur = cur + string(buf[i])
		}
	}
	return ret
}
func (b *Buffer) Split() (error, []string)  {
	err, s := b.Buf.String()
	if err != nil {
		return err, nil
	}
	return nil, split(s)
}


// func endOfFileCursor(buf []byte) (int, int)  {
// 	var x int
// 	var y int
// 	y = len(split(buf)) - 1
// 	x = len(split(buf)[y])
// 	return x, y
// }


func (b *Buffer) Insert(index int, x byte) () {
	b.Buf.Insert(index, string(x))
}


func (b *Buffer) Delete(index int)  {
	b.Buf.Delete(index, 0)
}

func (b *Buffer) ApplyLastOp()  {
	op := b.Ops[len(b.Ops)-1]
	op.Apply(b)
}


func (b *Buffer) MoveCursor(x int, y int) (error)  {
	err, lines := b.Split()
	if err != nil {
		return err
	}
	nlines := len(lines)
	cur := b.Cur
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

func (b *Buffer) MoveCursorBack() (error)  {
	err, lines := b.Split()
	if err != nil {
		return err
	}
	cur := b.Cur
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

func (b *Buffer) Cursor2ind() (error, int) {
	ind := 0
	err, splitLines := b.Split()
	if err != nil {
		return err, 0
	}
	cur := b.Cur
	for i := 0; i < cur.Y; i++ {
		ind += len(splitLines[i]) + 1
	}
	ind += cur.X
	return nil, ind
}

func (b *Buffer) Ind2Cursor(ind int) (error, Cursor) {
	if ind >= b.Buf.Length() {
		return fmt.Errorf("Index out of range"), Cursor{}
	}
	x, y := 0, 0
	for i := 0; i < ind; i++ {
		err, rune := b.Buf.Getindex(i)
		if err != nil {
			return err, Cursor{}
		}
		if rune == '\n' {
			x = 0
			y += 1
		} else {
			x += 1
		}
	}
	return nil, Cursor{x, y}
}
