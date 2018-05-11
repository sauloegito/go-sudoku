package model

import (
	"fmt"
	"testing"
)

func TestFailIntanceCell(t *testing.T) {
	origLogFatal := LogFatal

	// after this test, replace the original fatal function
	defer func() { LogFatal = origLogFatal }()

	errors := []string{}
	LogFatal = func(args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprint(args))
		}
	}

	CellFactory(12, 2, 0)
	CellFactory(2, 32, 0)
	CellFactory(9, 9, 0)
	CellFactory(2, 5, 18)

	c := CellFactory(3, 4, 0)
	chng := c.Can(3)
	if !chng {
		errors = append(errors, "Mudança deve ser permitida: "+stringCell(c)+" "+fmt.Sprint(chng))
	}
	c.Cannot(3)

	fixed := CellFactory(3, 5, 2)
	err := fixed.Set(0)
	if err != nil {
		errors = append(errors, fmt.Sprint(err))
	}

	if len(errors) != 6 {
		t.Errorf("Experados 6 erros, encontrados %v.\n%v", len(errors), errors)
	}

	fmt.Println("Erros mapeados:")
	for i, err := range errors {
		fmt.Println("\t", i, "-", err)
	}
}

func TestInstanceCell(t *testing.T) {
	var x, y, v byte = 4, 0, 6
	c := CellFactory(x, y, v)
	var a, b byte = 3, 0
	p := PosicaoFactory(a, b)

	sameRow, sameCol, sameArea := c.samePos(p)
	if !(!sameRow && sameCol && sameArea) {
		t.Error("Cel:", stringCell(c), "; Pos:", p, " compare:", sameRow, sameCol, sameArea)
	}

	wait := CellFactory(a, b, 0)
	err := wait.Set(v)
	if err != nil {
		t.Error("Preenchimento deve ser permitido para ", stringCell(wait), "; ", err)
	}
	err = wait.Set(0)
	if err != nil {
		t.Error("Limpar preenchimetno deve ser permitido para ", stringCell(wait), "; ", err)
	}

	waitRow, waitCol := wait.atPos(p)
	if !waitRow || !waitCol {
		t.Error("Wait:", stringCell(wait), "; Pos:", p, " compare:", waitRow, waitCol)
	}

	sameRow, sameCol, sameArea = c.sameCell(wait)
	if !(!sameRow && sameCol && sameArea) {
		t.Error("Cel-1:", stringCell(c), "; Cel-2:", stringCell(wait), " compare:", sameRow, sameCol, sameArea)
	}

	if wait.Get() != 0 {
		t.Error("Celula inicializada: ", stringCell(wait))
	}

	var l, m, n byte = 2, 4, 7
	cans := []byte{l, m, n}
	for _, can := range cans {
		wait.Can(can)
	}
	retorno := wait.Possibles()
	erro := len(cans) != len(retorno)
	for x := 0; !erro && x < len(retorno); x++ {
		if cans[x] != retorno[x] {
			erro = true
			break
		}
	}
	if erro || !wait.Possible(l) || !wait.Possible(m) || !wait.Possible(n) {
		t.Error("Montagem de possibilidades com problema. Origem:", cans, "; Resultado:", retorno)
	}

	wait.Cannot(m)
	retorno = wait.Possibles()
	if len(retorno) != 2 || !wait.Possible(l) || wait.Possible(m) || !wait.Possible(n) {
		t.Error("Remoção de ", b, " com problema. Origem:", cans, "; Resultado:", retorno)
	}

	e := wait.Set(l)
	if e != nil {
		t.Error("Problema ao atribuir:", a, ". ", e)
	}
	if wait.Get() != l {
		t.Error("Atribuição com problema. Esperado:", a, "; Resultado:", wait.Get())
	}

	fmt.Println("Após SET ", stringCell(wait), " Cans:", wait.Possibles())

	wait.Cannot(l)
	wait.Can(m)

	fmt.Println("Após CNG ", stringCell(wait), " Cans:", wait.Possibles())

	if !wait.Possible(l) || wait.Possible(m) || !wait.Possible(n) {
		t.Error("Após Set, não pode mudar lista de possibilidades. ", stringCell(wait), " p:", wait.Possibles())
	}

	rows, cols, area := c.impact()
	PrintImpact(c, rows, cols, area)

	sizes := []int{len(rows), len(cols), len(area)}
	for _, size := range sizes {
		if size != 8 {
			t.Error("Lista de impactos de ", c.ref, " com tamanho inesperado: ", sizes, "\nLinhas: ", stringPosicoes(rows), "\nColunas: ", stringPosicoes(cols), "\nÁrea: ", stringPosicoes(area))
		}
	}

	for _, p := range rows {
		sameRow, _, _ := c.samePos(p)
		if !sameRow || c.ref.x != p.x {
			t.Error("Lista de impactos de ", c.ref, " com Linha equivocada: ", p)
		}
	}

	for _, p := range cols {
		_, sameCol, _ := c.samePos(p)
		if !sameCol || c.ref.y != p.y {
			t.Error("Lista de impactos de ", c.ref, " com Coluna equivocada: ", p)
		}
	}

	for _, p := range area {
		_, _, sameArea := c.samePos(p)
		if !sameArea {
			t.Error("Lista de impactos de ", c.ref, " com Área equivocada: ", p)
		}
	}

}

func TestErroAreaImpact(t *testing.T) {
	cel := CellFactory(3, 8, 0)
	rows, cols, area := cel.impact()
	PrintImpact(cel, rows, cols, area)

	if len(area) != 8 {
		t.Error("Problemas")
	}
	rA, cA := cel.RowArea()*3, cel.ColArea()*3

	var i byte = 0
	for x := byte(0); x < 9; x++ {
		r := x/3 + rA
		c := x%3 + cA
		if cel.ref.x == r && cel.ref.y == c {
			continue
		}

		if area[i].x != r || area[i].y != c {
			t.Error("Area inesperada {", r, c, "} para", stringCell(cel), i, area[i])
		}
		i++
	}

}
