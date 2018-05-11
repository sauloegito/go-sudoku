package model

import "fmt"

type indice struct {
	value byte
	count byte
	rows  []*Cell
	cols  []*Cell
	area  []*Cell
}

func newIndice(valor byte) *indice {
	ind := indice{value: valor}
	ind.rows = make([]*Cell, 0, 9)
	ind.cols = make([]*Cell, 0, 9)
	ind.area = make([]*Cell, 0, 9)
	return &ind
}

func (ind *indice) addRow(c *Cell) bool {
	if ind.checkCell(c) {
		if ind.addCell(ind.rows, c) {
			ind.rows = append(ind.rows, c)
			return true
		}
	}
	return false
}

func (ind *indice) addCol(c *Cell) bool {
	if ind.checkCell(c) {
		if ind.addCell(ind.cols, c) {
			ind.cols = append(ind.cols, c)
			return true
		}
	}
	return false
}

func (ind *indice) addArea(c *Cell) bool {
	if ind.checkCell(c) {
		if ind.addCell(ind.area, c) {
			ind.area = append(ind.area, c)
			return true
		}
	}
	return false
}

func (ind *indice) checkCell(c *Cell) bool {
	return c.CheckPossible(ind.value)
}

func (ind *indice) addCell(cels []*Cell, c *Cell) bool {
	for x := 0; x < len(cels); x++ {
		if cels[x] == c {
			return false
		}
	}
	ind.count++
	return true
}

func (ind *indice) String() string {
	txt := "\n" + fmt.Sprint(ind.value)
	txt += "\n\trows = " + stringCells(ind.rows)
	txt += "\n\tcols = " + stringCells(ind.cols)
	txt += "\n\tarea = " + stringCells(ind.area) + "\n"

	return txt
}
