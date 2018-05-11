package model

import (
	"fmt"
)

//import "fmt"

func findCell(search *Cell, list []*Cell) bool {
	for _, c := range list {
		if search == c {
			return true
		}
	}
	return false
}

func (m *Mapa) check(checks []*posicao, value byte) []*Cell {
	cels := make([]*Cell, 0, len(checks))
	for _, pos := range checks {
		c := m.getCell(pos)
		if c.CheckPossible(value) {
			cels = append(cels, c)
		}
	}
	return cels
}

func (m *Mapa) checkCelula(c *Cell) bool {
	if c.Get() != 0 {
		return false
	}

	vals := c.Possibles()
	if len(vals) == 1 {
		m.set(c, vals[0])
		return true
	}

	change := false
	indices := make([]*indice, 0, 9)
	rows, cols, area := c.impact()
	//PrintImpact(c, rows, cols, area)

	var count byte = byte(len(vals))
	for i := byte(0); i < count; i++ {
		//fmt.Println("\nChecando o valor: ", vals[i])

		ind := newIndice(vals[i])

		rowCells := m.check(rows, vals[i])
		//fmt.Println("Linhas encontradas: ", stringCells(rowCells))
		if len(rowCells) == 0 {
			// fmt.Println("Celula ", c.ref, "=", c, vals[i], " único na linha.")
			m.set(c, vals[i])
			return true
		} else {
			for _, cel := range rowCells {
				ind.addRow(cel)
			}
		}

		colCells := m.check(cols, vals[i])
		//fmt.Println("Colunas encontradas: ", stringCells(colCells))
		if len(colCells) == 0 {
			// fmt.Println("Celula ", c.ref, "=", c, vals[i], " único na coluna.")
			m.set(c, vals[i])
			return true
		} else {
			for _, cel := range colCells {
				ind.addCol(cel)
			}
		}

		areaCells := m.check(area, vals[i])
		//fmt.Println("Areas encontradas: ", stringCells(areaCells))
		if len(areaCells) == 0 {
			// fmt.Println("Celula ", c.ref, "=", c, vals[i], " único na área.")
			m.set(c, vals[i])
			return true
		} else {
			if checkOnlyInArea(vals[i], c, areaCells, rowCells, colCells) {
				change = true

				// refaz a checagem de área, pois podem/devem ter mudado
				areaCells = m.check(area, vals[i])
			}

			for _, cel := range areaCells {
				ind.addArea(cel)
			}
		}

		if ind.count > 0 {
			indices = append(indices, ind)
		}
	}

	if m.checkPairIndex(c, indices) {
		// m.Print()
		return true
	}

	return change
}

func checkOnlyInArea(value byte, c *Cell, areaCells, rowCells, colCells []*Cell) bool {
	// Iniciando como se todas estivessem na mesma cobertura
	change, sameRow, sameCol := false, true, true
	for _, areaCell := range areaCells {
		row, col := areaCell.atPos(c.ref)
		if !row {
			sameRow = false
		}
		if !col {
			sameCol = false
		}

	}

	// fmt.Println("Rows: ", stringCells(rowCells))
	// fmt.Println("Cols: ", stringCells(colCells))
	// fmt.Println("SameRow ", sameRow, "; SameCol ", sameCol)
	// Pause()

	if sameRow && len(rowCells) > len(areaCells) {
		// fmt.Println("Valor ", value, " apenas na linha de mesma área: ", stringCell(c), " e ", stringCells(areaCells))
		for x := byte(0); x < byte(len(rowCells)); x++ {
			rC := rowCells[x]
			if !findCell(rC, areaCells) {
				// fmt.Print(stringCell(rC), "; ")
				if rC.Cannot(value) {
					change = true
				}
			}
		}
		// fmt.Println()
	}

	if sameCol && len(colCells) > len(areaCells) {
		// fmt.Print("Valor ", value, " apenas na coluna de mesma área: ", stringCell(c), " e ", stringCells(areaCells))
		// fmt.Println(". Desmarcando de coluna em: ")
		for x := byte(0); x < byte(len(colCells)); x++ {
			cC := colCells[x]
			if !findCell(cC, areaCells) {
				// fmt.Print(stringCell(cC), "; ")
				if cC.Cannot(value) {
					change = true
				}
			}
		}
		fmt.Println()
	}

	return change
}

func (m *Mapa) checkExposedPairIndex(impacted []*posicao, v1, v2 byte, cels1, cels2 []*Cell) (change bool) {
	for _, c1 := range cels1 {
		for _, c2 := range cels2 {
			if c1 == c2 && len(c1.Possibles()) == 2 && c1.Possible(v1) && c1.Possible(v2) {
				// fmt.Printf("Par nu %s identificado para [%d %d]!\n", stringCell(c1), v1, v2)

				chng1 := m.cannot(impacted, v1, c1)
				chng2 := m.cannot(impacted, v2, c1)
				if chng1 || chng2 {
					change = true
				}
			}
		}
	}

	return change
}

func (m *Mapa) checkHidenPairIndex(c *Cell, v1, v2 byte, cels1, cels2 []*Cell) (change bool) {
	if len(cels1) == 1 && len(cels2) == 1 && cels1[0] == cels2[0] {
		cX := cels1[0]
		// fmt.Printf("Par escondido %s identificado para [%d %d]!\n", stringCell(cX), v1, v2)

		chng1 := c.cannotOthers(v1, v2)
		chng2 := cX.cannotOthers(v1, v2)
		if chng1 || chng2 {
			change = true
		}
	}

	return change
}

func (c *Cell) cannotOthers(others ...byte) (change bool) {
	cans := c.Possibles()
	for _, pos := range cans {
		if pos != others[0] && pos != others[1] {
			if c.Cannot(pos) {
				change = true
			}
		}
	}
	return change
}

func (m *Mapa) checkPairIndex(c *Cell, indices []*indice) (change bool) {
	rows, cols, area := c.impact()

	// fmt.Println("Pareando: ", stringCell(c), " com ", indices)
	// Pause()

	pairOnly := (len(indices) == 2)
	for i := 0; i < len(indices)-1; i++ {
		ind1 := indices[i] // varre do ZERO a LEN-2
		for j := i + 1; j < len(indices); j++ {
			ind2 := indices[j] // varre do I+1 a LEN-1

			if pairOnly {
				if m.checkExposedPairIndex(rows, ind1.value, ind2.value, ind1.rows, ind2.rows) {
					change = true
				}
				if m.checkExposedPairIndex(cols, ind1.value, ind2.value, ind1.cols, ind2.cols) {
					change = true
				}
				if m.checkExposedPairIndex(area, ind1.value, ind2.value, ind1.area, ind2.area) {
					change = true
				}
			} else {
				if m.checkHidenPairIndex(c, ind1.value, ind2.value, ind1.rows, ind2.rows) {
					change = true
				}
				if m.checkHidenPairIndex(c, ind1.value, ind2.value, ind1.cols, ind2.cols) {
					change = true
				}
				if m.checkHidenPairIndex(c, ind1.value, ind2.value, ind1.cols, ind2.cols) {
					change = true
				}
			}
		}
	}

	return change
}
