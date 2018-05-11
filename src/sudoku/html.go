package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"sudoku/model"
	"math/rand"
	"io/ioutil"
	"regexp"
)

var templateSudoku = template.New("sudoku")
var templateFuncs = template.FuncMap{"cellDeal": CellDealWith}
var sudokuTemplate = template.Must(templateSudoku.Funcs(templateFuncs).ParseFiles("tmplt/sudoku.html"))

var possiblesTemplate = template.Must(template.ParseFiles("tmplt/possibles.html"))

func renderSudoku(w http.ResponseWriter, m *model.Mapa) error {
	return sudokuTemplate.ExecuteTemplate(w, "sudoku.html", m)
}

func renderPossibles(w io.Writer, c model.Cell) error {
	return possiblesTemplate.ExecuteTemplate(w, "possibles.html", makePossibles(c.Possibles()))
}

func registerError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|play)/([1-4]\\d{0,3})$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		log.Println("Handling:", m)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func PlayHandler(w http.ResponseWriter, r *http.Request, id string) {
	id = id[:1] // Só o primeiro caracter
	arrMapa, err := loadMapas(id)
	if err != nil {
		id += "001"
	} else {
		idx := rand.Intn(len(arrMapa))
		id = arrMapa[idx].id
	}

	http.Redirect(w, r, "/game/"+id, http.StatusFound)
}

func GameHandler(w http.ResponseWriter, r *http.Request, id string) {
	m, err := loadMapa(id)
	if err != nil {
		http.Redirect(w, r, "/new/"+id, http.StatusFound)
		return
	}
	renderSudoku(w, m)
}

func loadMapa(id string) (*model.Mapa, error) {
	filename := "maps/" + id[:1] + "/" + id[1:] + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func valuePos(c model.Cell, value byte) string {
	if c.Possible(value) {
		return fmt.Sprint(value)
	}
	return " "
}

func makePossibles(pos []byte) (retorno [3][3]string) {
	var i int = 0
	for x := byte(0); x < 9; x++ {
		a, b := x/3, x%3
		retorno[a][b] = " "
		if i < len(pos) {
			y := x + 1
			if pos[i] == y {
				retorno[a][b] = fmt.Sprint(pos[i])
			} else if pos[i] > y {
				continue
			}
			i++ // quando pos[i] <= y
		}
	}
	return retorno
}

func CellDealWith(args ...interface{}) string {
	ok := false
	var c model.Cell
	if len(args) == 1 {
		c, ok = args[0].(model.Cell)
	}
	if !ok {
		for x := 0; x < len(args); x++ {
			fmt.Printf("%T %v\n", args[x], args[x])
		}
		return "Not Cell: " + fmt.Sprint(len(args)) + fmt.Sprint(args...)
	}

	var classes string

	x := c.RowArea() + c.ColArea()
	if x%2 == 0 {
		classes = "realce "
	}

	if c.Get() != 0 {
		classes += "valued "
		if c.Fixed {
			classes += "fixed "
		}
		return fmt.Sprint("<div class=\"", classes, "\">", c.Get(), "</div>")
	}

	buf := new(bytes.Buffer)
	err := renderPossibles(buf, c)
	if err != nil {
		model.LogFatal("Fail: ", err, "\nCelula: ", c)
	}

	return fmt.Sprint("<div class=\"", classes, "\">", buf, "</div>")
}
