package controller

import (
	"BlogYmmersion/manager"
	inittemplate "BlogYmmersion/templates"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const Port = "localhost:8080"

var store = sessions.NewCookieStore([]byte(SecretKey()))

func SecretKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

func RessourceNotFoundHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "connexion", nil)
}
func ConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	inittemplate.Temp.ExecuteTemplate(w, "confirmation", nil)

}

func ConnexionHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "connexion", nil)
}

func FormHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "form", nil)
}
func InscriptionHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "inscription", nil)
}
func ErrorHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "error", nil)
}
func TreatInscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session
	//recupérer les données du formulaire d'enregistrement
	email := r.FormValue("email")
	password := r.FormValue("password")

	//Enregistrer le nouvel Utilisateur
	users := manager.RetrieveUser()
	var login bool

	for _, user := range users {
		if user.Email == email && user.Password == password {
			//verifier si le login est déjà enregistré
			login = true
		}
	}
	if login {

		http.Redirect(w, r, "/connexion?error=already_registred", http.StatusFound)
	} else {
		//IL S AGIT D'UNE PREMIERE CONNEXION !
		//rediriger vers la page dc'acceuil & enregistrer le login
		manager.MarkLogin(email, password)

		i := 0
		//Creer une nouvelle session & stocker l'email
		var err error
		session, err = store.Get(r, "session-name")
		for i > 1 {
			if err != nil {
				http.Error(w, "ERREUR DE SESSION_1", http.StatusInternalServerError)
				return
			}
		}

		session.Values["email"] = email
		fmt.Println("EMAIL RECU", email)
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERREUR DE SESSION_2", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home?success=Login_registred", http.StatusFound)
	}

}
func TreatConnexionHandler(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session
	//recupérer les données du formulaire de connexion
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println("l' email:", email)
	fmt.Println("le password:", password)
	users := manager.RetrieveUser()
	var login bool

	for _, user := range users {
		if user.Email == email && user.Password == password {
			//verifier si le login est correcte
			login = true
			break
		}
	}
	if login {
		i := 0
		//Creer une nouvelle session & stocker l'email
		var err error
		session, err = store.Get(r, "session-name")
		for i > 1 {
			if err != nil {
				http.Error(w, "ERREUR DE SESSION_1", http.StatusInternalServerError)
				return
			}
		}

		session.Values["email"] = email
		fmt.Println("EMAIL RECU", email)
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERREUR DE SESSION_2", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		//rediriger vers la page de connexion avec un message d'erreur
		http.Redirect(w, r, "/connexion?error=invalid_login_try_again", http.StatusFound)
	}

}
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	//ouverture et lecture des données json à partir du fichier
	dataJSON, err := os.ReadFile("DATA.json")

	//gestion de l'erreur à la lecture du fichier
	if err != nil {
		log.Fatal(err)
		manager.PrintColorResult("red", " ERREUR LORS DE LA LECTURE DU FICHIER")
		fmt.Printf("ALERTE: %#v", err)
		return
	}

	//Désérialiser les données JSON dans
	// la structure de données(Analyse des reponses json)
	var data manager.DataCategory
	err = json.Unmarshal(dataJSON, &data)

	//gestion de l'erreur à la lecture du fichier
	if err != nil {
		log.Fatal(err)
		manager.PrintColorResult("red", " ERREUR LORS DE LA LECTURE DU FICHIER")
		fmt.Printf("ALERTE: %#v", err)
		return
	}

	// Extraire les 8 premiers films de toutes les catégories
	for i := range data.Categories {
		if len(data.Categories[i].Films) > 2 {
			data.Categories[i].Films = data.Categories[i].Films[:2]
		}
	}

	// Vérifier si la tranche a au moins 11 éléments avant d'accéder au 11e élément
	if len(data.Categories) > 2 {
		fmt.Printf("ALERTE: %#v", data.Categories[2])
	}

	inittemplate.Temp.ExecuteTemplate(w, "home", data)
}

// Pouvoir effectuer le recherche d'un film
// gestionnaire de requete "get"
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	categories, err := manager.LoadCategories()
	if err != nil {
		http.Error(w, "ERREUR DE SESSION_search", http.StatusInternalServerError)
		return
	}
	results := manager.SearchFilm(categories, query)

	//Au cas où aucun resultat ne correspond à la recherche
	if len(results) == 0 {
		http.Redirect(w, r, "/notFound", http.StatusFound)
		return
	}

	inittemplate.Temp.ExecuteTemplate(w, "search", results)
}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {

	//Récuperer l'IP de la catégorie en paramètre de la requetes
	categoryID := r.URL.Query().Get("id")

	//lecture des données json à partir du fichier
	dataJSON, err := os.ReadFile("DATA.json")
	//gestion de l'erreur à la lecture du fichier
	if err != nil {
		log.Fatal(err)
		manager.PrintColorResult("red", " ERREUR LORS DE LA LECTURE DU FICHIER")
		fmt.Printf("ALERTE: %#v", err)
		return
	}

	//Désérialiser les données JSON dans
	// la structure de données(Analyse des reponses json)
	var data manager.DataCategory
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Trouver la catégorie correspondante à l'ID
	var category manager.Category
	for _, c := range data.Categories {
		if c.ID == categoryID {
			category = c
			break
		}
	}

	inittemplate.Temp.ExecuteTemplate(w, "category", category)
}

// Sécurisation des routes/gestions des erreurs de chargement de pages
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	inittemplate.Temp.ExecuteTemplate(w, "404", nil)
}
func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := manager.LoadComments()
	if err != nil {
		//Gerer l'erreur lors du chargement des commentaires
		log.Println("'erreur lors du chargement des commentaires CommentHandler", err)
		//Rediriger l'utilisateur vers la page d'erreurs
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	inittemplate.Temp.ExecuteTemplate(w, "comments", comments)
}

func SubmitCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentaire := r.FormValue("commentaire")
	nomFilm := r.FormValue("nom_film")
	userEmail, err := GetEmailSession(r)
	if err != nil {
		//Rediriger l'utilisateur vers la page de connexion
		http.Redirect(w, r, "/connexion", http.StatusFound)
		return
	}

	comment := manager.Comment{
		Email:   userEmail,
		NomFilm: nomFilm,

		Commentaire: commentaire,
	}

	// ENREGISTRER LE COMMENTAIRE mis à jour dans le fichier
	err = manager.SaveComment([]manager.Comment{comment})
	if err != nil {
		//Gerer l'erreur lors de la sauvegarde des commentaires
		log.Println("Erreur lors de la sauvegarde des commentaires :", err)
		//Rediriger l'utilisateur vers la page d'erreurs
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/comments", http.StatusSeeOther)
}

// Récuperer l'email à partir de la session ouverte
func GetEmailSession(r *http.Request) (string, error) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		return "", err
	}
	//Vérifier si le user est authentifié dans la session
	userEmail, ok := session.Values["email"].(string)

	if !ok {
		return "", errors.New("utilisateur non authentifié")
	}
	return userEmail, nil
}

func AjouterFilmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	titre := r.Form.Get("titre")
	auteur := r.Form.Get("auteur")
	synopsis := r.Form.Get("synopsis")

	film := manager.Film{Titre: titre, Auteur: auteur, Synopsis: synopsis}

	filmData, err := manager.LoadFilmData()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	filmData.Films = append(filmData.Films, film)

	err = manager.SaveFilmData(filmData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}

func MyListHandler(w http.ResponseWriter, r *http.Request) {
	filmData, err := manager.LoadFilmData()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Passez les données des films à votre template HTML
	films := filmData.Films
	inittemplate.Temp.ExecuteTemplate(w, "myList", films)
}
