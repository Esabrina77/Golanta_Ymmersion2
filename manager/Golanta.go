package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// structure de sauvegarde du login de chaque  user
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Pseudo   string `json:"pseudo"`
}

var ListUser []LoginUser

// Structure de sauvegarde des commentaires
type Comment struct {
	Pseudo      string `json:"pseudo"`
	NomPerso    string `json:"nom"`
	Commentaire string `json:"commentaire"`
}

// structure globale de chaque personnage
type Aventurier struct {
	ID        int      `json:"id"`
	Image     string   `json:"image"`
	Nom       string   `json:"nom"`
	Capacites []string `json:"capacites"`
}

const (
	CommentFile = "manager/comments.txt"
	DATA        = "DATA.json"
	DATAPerso   = "Personnage.Json"
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
	data, err := os.ReadFile("login.json")

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
func MarkLogin(email string, password string, pseudo string) {
	var newLogin = LoginUser{
		Email:    email,
		Password: password,
		Pseudo:   pseudo,
	}
	users := RetrieveUser()
	users = append(users, newLogin)

	//Convertir lelogin en JSON
	data, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	//Ecrire les données JSON dans le fichier
	err = os.WriteFile("login.json", data, 0666)
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

// CHARGER LES DONNEES DE DATA DES PERSONNAGES
func ChargerPersonnages() ([]Aventurier, error) {
	//lecture des données du fichier Json
	dataJSON, err := os.ReadFile(DATAPerso)
	if err != nil {
		return nil, err
	}

	//désérialiser les données json en une liste de commentaires
	var personnages []Aventurier
	err = json.Unmarshal(dataJSON, &personnages)

	if err != nil {
		return nil, err
	}
	fmt.Printf("list des PERSONNAGES : %#v\n", personnages)
	return personnages, nil
}

func SauvegarderPersonnages(personnages []Aventurier) error {
	data, err := json.MarshalIndent(personnages, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(DATAPerso, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
