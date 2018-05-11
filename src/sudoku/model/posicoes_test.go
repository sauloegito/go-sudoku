package model

import (
	"testing"
)

func checkPosicao(p *posicao, x, y byte, t *testing.T) (iX, iY byte) {
	if p.x != x || p.y != y {
		t.Error("Instância incorreta: ", p)
	}

	iX = x / 3
	if p.rowArea() != iX {
		t.Error("Identificação de linha/area incorreta: ", p.rowArea(), iX)
	}
	iY = y / 3
	if p.colArea() != iY {
		t.Error("Identificação de coluna/area incorreta: ", p.colArea(), iY)
	}

	return
}

func TestPosicoes(t *testing.T) {
	factory := PosicaoFactory
	var x, y byte = 0, 4
	p1 := factory(x, y)
	iX, iY := checkPosicao(p1, x, y, t)

	var a, b byte = 2, 4
	p2 := factory(a, b)
	iA, iB := checkPosicao(p2, a, b, t)

	sameRow, sameCol, sameArea := p1.samePos(p2)
	if (a != x && sameRow) || (a == x && !sameRow) {
		t.Error("Mesma linha incorreta: ", p1, ",x:", x, "; ", p2, ",x:", a, "; result:", sameRow)
	}

	if (b != y && sameCol) || (b == y && !sameCol) {
		t.Error("Mesma coluna incorreta: ", p1, ",y:", y, "; ", p2, ",x:", b, "; result:", sameCol)
	}

	p2_ := factory(a, b)
	if p2 != p2_ {
		t.Error("Problema com a fábrica de posições")
	}

	if !(sameArea && iX == iA && iY == iB) {
		t.Error("Mesma area incorreta: ", p1, " area:{", iX, ",", iY, "}; ", p2, " area:{", iA, ",", iB, "}; result:", sameArea)
	}

}
