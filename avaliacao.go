package main

import (
	"log"
	"net/http"
	"math/rand"
	"time"
	"sudoku/model"
	"sudoku/html"
))

func main() {
	rand.Seed(time.Now().Unix())

	http.HandleFunc("/play/", makeHandler(html.PlayHandler)
	http.HandleFunc("/game/", makeHandler(html.GameHandler)
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
	//mapa = resolver()

}

func resolver() *Mapa {
	input := [9][9]byte{
		{0, 3, 0, 0, 2, 0, 4, 7, 0},
		{0, 0, 0, 0, 0, 1, 0, 9, 0},
		{0, 0, 0, 7, 0, 0, 3, 0, 0},
		{0, 0, 5, 0, 8, 0, 1, 0, 0},
		{0, 9, 6, 0, 0, 5, 0, 0, 3},
		{0, 0, 3, 0, 7, 0, 5, 0, 0},
		{0, 0, 0, 1, 0, 0, 9, 0, 0},
		{0, 0, 0, 0, 0, 8, 0, 1, 0},
		{0, 8, 0, 0, 3, 0, 7, 2, 0},
	}
	return model.Avaliar("Aval", input)
}
