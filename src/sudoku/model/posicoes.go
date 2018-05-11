package model

var posicoesIniciadas bool
var posicoes [9][9]*posicao

type posicao struct {
	x byte
	y byte
}

func (p1 *posicao) equalsPos(p2 *posicao) (bool, bool) {
	return p1.x == p2.x, p1.y == p2.y
}

func (p *posicao) rowArea() byte {
	return p.x / 3
}

func (p *posicao) colArea() byte {
	return p.y / 3
}

func (p1 *posicao) samePos(p2 *posicao) (bool, bool, bool) {
	sameArea := p1.rowArea() == p2.rowArea() && p1.colArea() == p2.colArea()
	sameRow, sameCol := p1.equalsPos(p2)
	return sameRow, sameCol, sameArea
}

func initPosicoes() {
	for i := byte(0); i < 9; i++ {
		for j := byte(0); j < 9; j++ {
			posicoes[i][j] = &posicao{i, j}
		}
	}
	posicoesIniciadas = true
}

func PosicaoFactory(x, y byte) *posicao {
	if !posicoesIniciadas {
		initPosicoes()
	}
	return posicoes[x][y]
}
