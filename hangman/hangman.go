package hangman

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

type Variable struct {
	Mot       string
	Motcrypte string
	Faute     int
	Fin       int
	Entrer    string
}

func lecturelist() []string { // fonction qui met la liste de mot dans un tableau
	var tableau []string
	file, err := os.Open("hangman/assets/mot/word.txt")
	if err != nil {
		os.Exit(1) // si pas de liste alors
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		scanner := scanner.Text()
		tableau = append(tableau, scanner)
	}
	return tableau
}

func rdm(max int) int { // fonction qui fait un aléatoire simple
	rand.Seed(time.Now().UnixNano())
	nombre := rand.Intn(max)
	return nombre
}

func Motrdm() string { // fonction qui choisit un mot aléatoire
	var Variable *Variable = new(Variable)
	liste := lecturelist()
	nombreligne := len(liste) - 1
	lignemot := rdm(nombreligne)
	Variable.Mot = liste[lignemot]
	return Variable.Mot
}
func rdmblacklist(max int, blacklist []int) int { // fonction qui fait un aléatoire avec une blacklist
	exclue := map[int]bool{}
	for _, index := range blacklist {
		exclue[index] = true
	}
	for {
		nombre := rand.Intn(max)
		if !exclue[nombre] {
			return nombre
		}
	}
}
func premiermot(mot string) []int { // fonction qui gère le calcul de lettre à révéler, leurs positions
	var blacklist []int
	n := len(mot)/2 - 1
	var positionlettre []int
	for i := 0; i < n; i++ {
		a := rdmblacklist(len(mot), blacklist[:])
		blacklist = append(blacklist, a)
		positionlettre = append(positionlettre, a)
	}
	return positionlettre
}

func Printmot(mot string) string { // fonction qui print les mots cryptés en ajoutant les lettres de base + celles découvertes
	var motcrypte string
	positionlettre := premiermot(mot)
	runemot := []rune(mot)
	for i := 0; i < len(mot); i++ {
		motcrypte += "_"
	}
	runecrypte := []rune(motcrypte)
	for j := 0; j < len(positionlettre); j++ {
		for k := 0; k < len(mot); k++ {
			if k == positionlettre[j] {
				runecrypte[k] = runemot[k]
			}
		}
	}
	motcrypte = string(runecrypte)
	return motcrypte
}

func Equal(tableau1, tableau2 []rune) bool { // fonction qui compare 2 tableaux de rune
	if len(tableau1) != len(tableau2) {
		return false
	}
	for i, v := range tableau1 {
		if v != tableau2[i] {
			return false
		}
	}
	return true
}
