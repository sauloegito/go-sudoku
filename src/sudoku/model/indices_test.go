package model

import (
	"fmt"
	"testing"
)

func TestIndices(t *testing.T) {
	var x byte = 4
	ind := newIndice(x)

	if ind.value != x || ind.count > 0 {
		t.Error("Instância de Indice incorreta", ind)
	}

	var a, b byte = 2, 3
	fixedCel := CellFactory(b, a, x)
	if ind.addArea(fixedCel) {
		t.Error("Inclusão de celula em área não permitida se valor já definido. ", ind)
	}

	if ind.addRow(fixedCel) {
		t.Error("Inclusão de celula em linha não permitida se valor já definido. ", ind)
	}

	if ind.addCol(fixedCel) {
		t.Error("Inclusão de celula em coluna não permitida se valor já definido. ", ind)
	}

	varCel1 := CellFactory(a, b, 0)
	if ind.addArea(varCel1) {
		t.Error("Inclusão de celula em área não permitida se possibilidade não marcada. ", ind)
	}

	varCel1.Can(x)
	if !ind.addArea(varCel1) {
		t.Error("Inclusão de celula em área deveria ser permitida. ", ind)
	}

	if ind.addArea(varCel1) {
		t.Error("Inclusão de mesma celula mais de uma vez não é permitida. ", ind)
	}

	varCel2 := CellFactory(b, a, 0)
	varCel2.Can(x)
	if !ind.addArea(varCel2) {
		t.Error("Inclusão de segunda celula em área deveria ser permitida. ", ind)
	}

	if !ind.addRow(varCel1) {
		t.Error("Inclusão de celula em linha deveria ser permitida. ", ind)
	}

	if !ind.addCol(varCel2) {
		t.Error("Inclusão de celula em coluna deveria ser permitida. ", ind)
	}

	if ind.count != 4 {
		t.Error("Contagem de índices inconsistente. Esperado: 4 Obtido:", ind.count, "Indice:", ind)
	}

	fmt.Println("Indice montado:", ind)
}
