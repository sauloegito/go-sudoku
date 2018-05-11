package model

import (
	"bufio"
	"fmt"
	"os"
)

func Pause() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func Avaliar(label string, matriz [9][9]byte) *Mapa {
	m := NewMapa(label, matriz)
	full, change, count := m.Resolve()
	fmt.Println(label, "Full: ", full, "; Change: ", change, "; Count: ", count)
	m.Print()
	return m
}

type Mapa struct {
	label string
	Cells  [9][9]*Cell
	values [9][]*Cell
}

func (m Mapa) Print() {
	for i := 0; i < 9; i++ {
		fmt.Println(m.Cells[i])
	}
	Pause()
}

func  NewMapa(label string, matriz [9][9]byte) (m *Mapa) {
	m = &Mapa{label: label}
	waiting := make([]*Cell, 0, 72)
	for i := byte(0); i < 9; i++ {
		m.values[i] = make([]*Cell, 0, 9)
	}

	for i := byte(0); i < 9; i++ {
		for j := byte(0); j < 9; j++ {
			value := matriz[i][j]
			c := CellFactory(i, j, value)
			m.Cells[i][j] = c
			if c.Get() != 0 {
				m.put(c)
			} else {
				waiting = append(waiting, c)
			}
		}
	}

	//fmt.Println("Entrada Inicial ", label)
	// m.Print()

	for x := 0; x < len(waiting); x++ {
		w := waiting[x]
	valores:
		for i := byte(0); i < 9; i++ {
			vI := m.values[i]
			for idx := byte(0); idx < byte(len(vI)); idx++ {
				sRow, sCol, sArea := w.sameCell(vI[idx])
				if sRow || sCol || sArea {
					continue valores
				}
			}
			w.Can(i + 1)
		}
	}
	return m
}

func (m *Mapa) put(c *Cell) {
	idx := c.Get() - 1
	m.values[idx] = append(m.values[idx], c)
}

func (m *Mapa) set(c *Cell, value byte) {
	e := c.Set(value)
	if e != nil {
		LogFatal(e)
	}
	m.put(c)

	rows, cols, area := c.impact()
	// fmt.Println("\tLinhas:", stringPosicoes(rows))
	m.cannot(rows, value, nil)
	// fmt.Println("\tColunas:", stringPosicoes(cols))
	m.cannot(cols, value, nil)
	// fmt.Println("\tÁrea:", stringPosicoes(area))
	m.cannot(area, value, nil)
	// fmt.Println("\tEncerrando SET")
}

func (m *Mapa) getCell(pos *posicao) *Cell {
	return m.Cells[pos.x][pos.y]
}

func (m *Mapa) cannot(checks []*posicao, value byte, except *Cell) bool {
	change := false
	for _, pos := range checks {
		c := m.getCell(pos)
		if c.Get() == 0 && c != except {
			if c.Cannot(value) {
				change = true
			}
		}
	}
	return change
}

func (m *Mapa) Resolve() (bool, bool, int) {
	full := false
	change := true
	count := 0
	for !full && change {
		// fmt.Println("Resolver Step ", count)
		count++
		change = false
		// m.Print()
		for row := byte(0); row < 9; row++ {
			for col := byte(0); col < 9; col++ {
				if m.checkCelula(m.Cells[row][col]) {
					change = true
					//m.Print()
				}
			}
		}

		// Verificação de conclusão
		full = true
		for i := byte(0); (i < 9) && full; i++ {
			count := byte(len(m.values[i]))
			full = full && (count == 9)
		}

	}
	return full, change, count
}
