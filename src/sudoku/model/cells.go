package model

import (
	"errors"
	"fmt"
	"log"
)

var LogFatal = log.Fatal

func CellFactory(x, y, initValue byte) *Cell {
	c := Cell{}
	e := c.fill(x, y, initValue)
	if e != nil {
		LogFatal(e)
	}
	return &c
}

func validate(v byte) error {
	var min byte = 1
	var max byte = 9
	if v < min || v > max {
		txt := fmt.Sprint("Valor não permitido: ", v)
		return errors.New(txt)
	}
	return nil
}

type Cell struct {
	ref   *posicao
	Fixed bool
	value byte
	cans  int
}

func (c *Cell) Get() byte {
	return c.value
}

func (c *Cell) fill(x, y, v byte) error {
	row := x + 1
	col := y + 1
	ex := validate(row)
	ey := validate(col)
	if ex != nil || ey != nil {
		txt := "Indices devem estar entre 0 e 8. Verifique [ "
		if ex != nil {
			txt += fmt.Sprint("X=", x, " ")
		}
		if ey != nil {
			txt += fmt.Sprint("Y=", y, " ")
		}
		txt += "]"
		return errors.New(txt)
	}
	e := validate(v)
	if v == 0 || e == nil {
		c.ref = PosicaoFactory(x, y)
		c.value = v
		c.Fixed = (v != 0)
		c.cans = 0
		return nil
	}
	return e
}

func (c *Cell) validate(v byte) error {
	if c.Fixed {
		txt := fmt.Sprint("Célula Fixa ", stringCell(c), " não pode assumir ", v)
		return errors.New(txt)
	} else if v == 0 {
		return nil
	}
	return validate(v)
}

func (c *Cell) Set(v byte) error {
	e := c.validate(v)
	if e == nil {
		// fmt.Println("Marcando ", v, " para célula: ", stringCell(c))
		c.value = v
	}
	return e
}

func (c *Cell) RowArea() byte {
	return c.ref.rowArea()
}

func (c *Cell) ColArea() byte {
	return c.ref.colArea()
}

func (c *Cell) atPos(pos *posicao) (sRow, sCol bool) {
	sRow, sCol = pos.equalsPos(c.ref)
	return
}

func (c *Cell) samePos(pos *posicao) (bool, bool, bool) {
	return pos.samePos(c.ref)
}

func (c1 *Cell) sameCell(c2 *Cell) (bool, bool, bool) {
	return c2.samePos(c1.ref)
}

func (c *Cell) impact() ([]*posicao, []*posicao, []*posicao) {
	var rows = make([]*posicao, 0, 8)
	var cols = make([]*posicao, 0, 8)
	var area = make([]*posicao, 0, 8)
	rA := c.RowArea() * 3
	cA := c.ColArea() * 3
	for i := byte(0); i < 9; i++ {
		if i != c.ref.x {
			cols = append(cols, &posicao{i, c.ref.y})
		}
		if i != c.ref.y {
			rows = append(rows, &posicao{c.ref.x, i})
		}

		rI := i/3 + rA
		cI := i%3 + cA
		if rI != c.ref.x || cI != c.ref.y {
			area = append(area, &posicao{rI, cI})
		}
	}
	return rows, cols, area
}

func (c *Cell) CheckPossible(v byte) bool {
	return c.Get() == 0 && c.Possible(v)
}

func (c *Cell) Cannot(v byte) bool {
	return c.possibilite(v, false)
}

func (c *Cell) Can(v byte) bool {
	return c.possibilite(v, true)
}

func (c *Cell) possibilite(v byte, can bool) bool {
	if c.Fixed || c.Get() != 0 {
		return false
	}
	// fmt.Print("possibilite(")
	// fmt.Print(v)
	// fmt.Print(", ")
	// fmt.Print(can)
	// fmt.Print("): before=")
	// fmt.Print(c.cans)

	e := validate(v)
	if e != nil {
		LogFatal(e)
	}

	change := false
	var op byte = (v - 1)
	idx := 1 << op
	// fmt.Print("; idx=")
	// fmt.Print(idx)
	exist := (c.cans & idx) != 0
	if can != exist {
		change = true

		if can {
			c.cans |= idx
		} else {
			// fmt.Println("Removendo ", v, " de ", stringCell(c))
			if c.cans == idx {
				LogFatal("Desmarcando a última possibilidade. Falha em Regra!")
			}
			c.cans &^= idx
		}
	}

	// fmt.Print("; after=")
	// fmt.Print(c.cans)
	// fmt.Println(";")
	return change
}

func (c *Cell) Possible(v byte) bool {
	var idx byte = v - 1
	bitCheck := 1 << idx
	bit := c.cans & bitCheck
	return (bit != 0)
}

func (c *Cell) Possibles() []byte {
	// fmt.Print("possibles(): ")
	// fmt.Println(c.cans)

	count := 0
	p := make([]byte, 9)
	var i byte = 0
	op := 1
	for i < 9 {
		bit := c.cans & op
		// fmt.Print("** i=")
		// fmt.Print(i)
		// fmt.Print("; op=")
		// fmt.Print(op)
		// fmt.Print("; bit=")
		// fmt.Print(bit)
		// fmt.Print("; before=")
		// fmt.Print(count)
		// fmt.Print("; return=")
		// fmt.Print(p)

		if bit != 0 {
			p[count] = i + 1
			count++
		}
		// fmt.Print("; after=")
		// fmt.Println(count)
		op <<= 1
		i++
	}
	return p[:count]
}

func (c *Cell) String() string {
	if c.Get() != 0 {
		return fmt.Sprintf(" %d", c.Get())
	} else {
		p := c.Possibles()

		return fmt.Sprint(p)
	}
}
