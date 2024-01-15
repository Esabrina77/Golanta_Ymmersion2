package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// structure globale de chaque catégorie
type Category struct {
	Image string `json:"image"`
	ID    string `json:"id"`
	Name  string `json:"name"`

	Films []Film `json:"films"`
}

// structure globale de l'objet film
type Film struct {
	ID          string `json:"id"`
	Titre       string `json:"titre"`
	Auteur      string `json:"auteur"`
	Synopsis    string `json:"synopsis"`
	Image       string `json:"image"`
	Note        string `json:"note"`
	RealeseDate string `json:"realese_date"`
}
type DataCategory struct {
	Categories []Category `json:"categories"`
}

type Comment struct {
	Email       string `json:"email"`
	NomFilm     string `json:"nom_film"`
	Commentaire string `json:"commentaire"`
}
type FilmData struct {
	Films []Film `json:"films"`
}

// structure de sauvegarde du login de chaque  user
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ListUser []LoginUser

const (
	CommentFile = "manager/comments.txt"
	DATA        = "DATA.json"
	DATA2       = "AddArticle.Json"
)

func PrintColorResult(color string, message string) {
	colorCode := ""
	switch color {
	case "red":
		colorCode = "\033[31m"
	case "green":
		colorCode = "\033[32m"
	case "yellow":
		colorCode = "\033[33m"
	case "blue":
		colorCode = "\033[34m"
	case "purple":
		colorCode = "\033[35m"

	default: //REMETTRE LA COULEUR INITIALE (blanc)
		colorCode = "\033[0m"
	}
	fmt.Printf("%s%s\033[0m", colorCode, message)
}

func RetrieveUser() []LoginUser {
	data, err := os.ReadFile("manager/Login.txt")

	if err != nil {
		fmt.Printf("Erreur lors de la lecture du fichier:%v", err)
		return nil
	}
	var Users []LoginUser
	err = json.Unmarshal(data, &Users)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("list des users : %#v\n", Users)
	return Users
}

// Marquer ( enregistrer) les nouveaux users dans le fichiers De login
func MarkLogin(email string, password string) {
	var newLogin = LoginUser{
		Email:    email,
		Password: password,
	}
	users := RetrieveUser()
	users = append(users, newLogin)

	//Convertir lelogin en JSON
	data, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}
	//Ecrire les données JSON dans le fichier
	err = os.WriteFile("manager/Login.txt", data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("list des users : %#v\n", users)

}

// Enregistrer les commentaires
func SaveComment(newComments []Comment) error {
	// Charger les commentaires déjà enregistrés
	comments, err := LoadComments()
	if err != nil {
		return err
	}

	// Ajouter les nouveaux commentaires à la liste existante
	comments = append(comments, newComments...)

	// Convertir les commentaires en JSON
	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	// Ecrire les données JSON dans le fichier
	err = os.WriteFile(CommentFile, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

// Charger les commentaires à partir d'un fichier json
func LoadComments() ([]Comment, error) {
	//Vérifier si le fichier json

	_, err := os.Stat(CommentFile)
	if os.IsNotExist(err) {
		return []Comment{}, nil
	} else if err != nil {
		fmt.Printf("Erreur lors de la verification du fichier : %#v\n", err)
		return nil, err
	}

	//lecture des données du fichier Json
	dataJSON, err := os.ReadFile(CommentFile)
	if err != nil {
		log.Fatal(err)
	}
	//désérialiser les données json en une liste de commentaires
	var comments []Comment
	err = json.Unmarshal(dataJSON, &comments)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("list des commentaires : %#v\n", comments)
	return comments, err
}

// FONCTIONNALITE :RECHERCHE D UN FILM
// charger les données des films à partir de DATA
func LoadCategories() ([]Category, error) {
	file, err := os.Open(DATA)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var dataObj DataCategory
	err = json.Unmarshal(data, &dataObj)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return dataObj.Categories, nil
}

// RECHERCHER LE FILM DANS LES CATEGORIES
// pour une recherche plus étendue j'ai utilisé:
// -->> Knuth-Morris-Pratt (KMP). Cet algorithme permet de rechercher efficacement des motifs
// dans une chaîne de caractères.
func SearchFilm(Categories []Category, query string) []Film {
	var results []Film

	for _, category := range Categories {
		for _, film := range category.Films {
			if strings.Contains(strings.ToLower(film.Titre), strings.ToLower(query)) ||
				kmpSearch(strings.ToLower(film.Titre), strings.ToLower(query)) {
				results = append(results, film)
			}
		}
	}

	return results
}
func kmpSearch(text, pattern string) bool {
	lps := computeLPS(pattern)
	i, j := 0, 0

	for i < len(text) {
		if pattern[j] == text[i] {
			i++
			j++
		}

		if j == len(pattern) {
			return true
		} else if i < len(text) && pattern[j] != text[i] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}

	return false
}

func computeLPS(pattern string) []int {
	lps := make([]int, len(pattern))
	length := 0
	i := 1

	for i < len(pattern) {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}

	return lps
}

// CHARGER LES DONNEES DE DATA DES FILMS

func LoadFilmData() (FilmData, error) {
	var filmData FilmData

	file, err := os.Open(DATA2)
	if err != nil {
		if os.IsNotExist(err) {
			return filmData, nil
		}
		return filmData, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return filmData, err
	}

	err = json.Unmarshal(data, &filmData)
	if err != nil {
		return filmData, err
	}

	return filmData, nil
}

func SaveFilmData(filmData FilmData) error {
	data, err := json.MarshalIndent(filmData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(DATA2, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
