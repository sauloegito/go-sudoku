package html

import (
	"fmt"
	"sudoku/model"
	"testing"
)

func TestHtml(t *testing.T) {
	c := model.CellFactory(3, 4, 0)
	c.Can(1)
	c.Can(7)

	s := valuePos(*c, 3)
	if " " != s {
		t.Errorf("Deveria ser vazio: %q", s)
	}

	s = valuePos(*c, 7)
	if "7" != s {
		t.Errorf("Deveria ser sete: %q", s)
	}

	c.Can(5)
	c.Can(9)
	arr := makePossibles(c.Possibles())
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			fmt.Printf("%q", arr[x][y])
		}
		fmt.Println()
	}

	fmt.Println(CellDealWith(*c))

	c = model.CellFactory(3, 4, 7)
	fmt.Println(CellDealWith(*c))
}
