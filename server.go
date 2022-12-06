package main

import (
	"fmt"
	"hangman/hangman"
	"log"
	"net/http"
	"text/template"
)

type Variable struct {
	Mot       string
	Motcrypte string
	Faute     int
	Fin       int
	Entrer    string
	Boucle    int
}

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html", "./templates/header.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Infos(w http.ResponseWriter, r *http.Request, infos *Variable) {
	template, err := template.ParseFiles("./pages/hangman.html", "./templates/header.html", "./templates/variable.html", "./templates/forms.html", "./templates/dessin.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, infos)
}

func Reset(w http.ResponseWriter, r *http.Request, Variable *Variable) {
	Variable.Faute = -1
	Variable.Fin = 0
	Variable.Boucle = 0
	Variable.Entrer = ""
	http.Redirect(w, r, "/hangman", http.StatusSeeOther)
}

func Victoire(w http.ResponseWriter, r *http.Request, Variable *Variable) {
	template, err := template.ParseFiles("./pages/victoire.html", "./templates/header.html", "./templates/variable.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, Variable)
}

func Defaite(w http.ResponseWriter, r *http.Request, Variable *Variable) {
	template, err := template.ParseFiles("./pages/defaite.html", "./templates/header.html", "./templates/variable.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, Variable)
}

func Credit(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/credit.html", "./templates/header.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func User(w http.ResponseWriter, r *http.Request, Variable *Variable) {
	tmpl := template.Must(template.ParseFiles("./templates/forms.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	a := r.FormValue("entrer")
	if len(a) > 0 {
		Variable.Entrer = a
	}
	tmpl.Execute(w, struct{ Success bool }{true})

}

func hang(w http.ResponseWriter, r *http.Request, Variable *Variable) {
	User(w, r, Variable)
	if len(Variable.Entrer) > 0 {
		jeu(Variable)
	}
	if Variable.Fin == 1 {
		http.Redirect(w, r, "/defaite", http.StatusSeeOther)
	} else if Variable.Fin == 2 {
		http.Redirect(w, r, "/victoire", http.StatusSeeOther)
	}
}

func jeu(Variable *Variable) {
	main1(Variable)
}

func main1(Variable *Variable) {
	runemot := []rune(Variable.Mot)
	runelecture := []rune(Variable.Entrer)
	runecrypte := []rune(Variable.Motcrypte)
	boucle := 0
	double := 0
	fautemot := 0
	for boucle < 1 {
		if len(runelecture) > 1 {
			if hangman.Equal(runemot, runelecture) {
				runecrypte = runemot
			} else if double < 1 {
				fautemot += 1
			}
		} else if len(runelecture) < 2 {
			for i := 0; i < len(Variable.Mot); i++ {
				if runemot[i] == runelecture[0] {
					runecrypte[i] = runelecture[0]
					boucle += 1
					double += 1
				} else if double > 1 {
					double -= 1
					Variable.Faute -= double
				}
			}
		}
		if boucle != 1 {
			Variable.Faute += 1
			if fautemot == 1 {
				Variable.Faute += 1
			}
			boucle = 1
		}
		Variable.Motcrypte = string(runecrypte)

		if Variable.Fin == 1 {
			boucle = 1
			fmt.Println("Le mot était :", Variable.Mot)
		} else if Variable.Mot == Variable.Motcrypte {
			boucle = 1
			fmt.Println("Felicitation")
			Variable.Fin = 2
			Variable.Faute -= 1
		}
	}
	Variable.Fin = dessin(Variable)
}

func dessin(Variable *Variable) int { // fonction qui dessine les pendues
	if Variable.Faute >= 10 {
		Variable.Fin = 1
	}
	return Variable.Fin
}

func choixmot(Variable *Variable) {
	Variable.Mot = hangman.Motrdm()
	fmt.Println(Variable.Mot)
	Variable.Motcrypte = hangman.Printmot(Variable.Mot)
}
func main() {
	var Variable *Variable = new(Variable) // permet d'utiliser la structure
	Variable.Boucle = 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r)
	})
	http.HandleFunc("/hangman", func(w http.ResponseWriter, r *http.Request) {
		if Variable.Boucle < 1 {
			choixmot(Variable)
			Variable.Fin = 0
			Variable.Faute = 0
			Variable.Boucle = 1
		}
		hang(w, r, Variable)
		Infos(w, r, Variable)
	})

	http.HandleFunc("/credit", func(w http.ResponseWriter, r *http.Request) {
		Credit(w, r)
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		Reset(w, r, Variable)
	})
	http.HandleFunc("/victoire", func(w http.ResponseWriter, r *http.Request) {
		Victoire(w, r, Variable)
	})
	http.HandleFunc("/defaite", func(w http.ResponseWriter, r *http.Request) {
		Defaite(w, r, Variable)
	})
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fi := http.FileServer(http.Dir("./hangman/assets/"))
	http.Handle("/hangman/assets/", http.StripPrefix("/hangman/assets/", fi))
	http.ListenAndServe(":8000", nil)
}
