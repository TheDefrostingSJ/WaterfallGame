package main;

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

// Fle Change test

type Page struct {
	Title string
	Body string
	Users []string
	RoundNumber int
}

var templates = template.Must(template.ParseFiles(
	"edit.html",
	"game.html",
	"view.html",
	"loginConfirm.html",
	))
var validPath = regexp.MustCompile("^/(edit|save|view|game)/([a-zA-Z0-9]+)$")

var users []string
var round int

func (p *Page) save() error {
	//filename := p.Title + ".txt"
	return nil
}

func loadPage(title string) (*Page, error) {
	return &Page{Title: title, Users: users, RoundNumber: round}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
			http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
			p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	name := r.FormValue("body")
	if name == "" {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	users = append(users, name)
	p := &Page{Title: title, Body: name, Users: users}
	err := p.save()
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	renderTemplate(w, "loginConfirm", p)
}

func gameHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
			p = &Page{Title: title}
	}
	round = round + 1
	renderTemplate(w, "game", p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
			m := validPath.FindStringSubmatch(r.URL.Path)
			if m == nil {
					http.NotFound(w, r)
					return
			}
			fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	// Log in

	fmt.Println("Please Login:")
	fmt.Println("\n\nhttp://localhost:8080/view/Waterfall")

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	// Round Start - Game Logic


	http.HandleFunc("/game/", makeHandler(gameHandler))



	log.Fatal(http.ListenAndServe(":8080", nil))
}
