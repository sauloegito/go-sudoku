package model

import (
	"fmt"
)

func stringCell(c *Cell) string {
	return fmt.Sprint(c.ref, "=", c)
}

func stringCells(cels []*Cell) string {
	txt := "["
	for x := 0; x < len(cels); x++ {
		if x > 0 {
			txt += " "
		}
		txt += stringCell(cels[x])
	}
	txt += "]"
	return txt
}

func stringPosicoes(arr []*posicao) string {
	txt := "["
	for i := 0; i < len(arr); i++ {
		if i > 0 {
			txt += " "
		}
		txt += fmt.Sprint(arr[i])
	}
	txt += "]"
	return txt
}

func PrintImpact(c *Cell, rows, cols, area []*posicao) {
	fmt.Println("Impactados por: ", stringCell(c))
	fmt.Println("\trows: ", stringPosicoes(rows))
	fmt.Println("\tcols: ", stringPosicoes(cols))
	fmt.Println("\tarea: ", stringPosicoes(area))
	//Pause()
}
