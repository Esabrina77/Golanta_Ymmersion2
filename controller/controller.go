package controller

import (
	"Golanta/manager"
	inittemplate "Golanta/templates"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"

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

// Sécurisation des routes/gestions des erreurs de chargement de pages
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	inittemplate.Temp.ExecuteTemplate(w, "404", nil)
}

func RessourceNotFoundHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "notFound", nil)
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

func TreatInscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session
	//recupérer les données du formulaire d'enregistrement
	email := r.FormValue("email")
	password := r.FormValue("password")
	pseudo := r.FormValue("pseudo")

	//Enregistrer le nouvel Utilisateur
	users := manager.RetrieveUser()
	var login bool

	for _, user := range users {
		if user.Email == email && user.Password == password && user.Pseudo == pseudo {
			//verifier si le login est déjà enregistré
			login = true
		}
	}
	if login {

		http.Redirect(w, r, "/connexion?error=already_registred", http.StatusFound)
	} else {
		//IL S AGIT D'UNE PREMIERE CONNEXION !
		//rediriger vers la page dc'acceuil & enregistrer le login
		manager.MarkLogin(email, password, pseudo)

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
	pseudo := r.FormValue("pseudo")

	fmt.Println("l' email:", email)
	fmt.Println("le password:", password)
	users := manager.RetrieveUser()
	var login bool

	for _, user := range users {
		if user.Email == email && user.Password == password && user.Pseudo == pseudo {
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

		session.Values["pseudo"] = pseudo
		fmt.Println("PSEUDO RECU", pseudo)
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
	inittemplate.Temp.ExecuteTemplate(w, "home", nil)
}

func AjouterPersoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == "POST" {
		nom := r.FormValue("nom")
		capacites := r.Form["capacite"]

		aventurier := manager.Aventurier{
			Nom:       nom,
			Capacites: capacites,
		}
		personnages, err := manager.ChargerPersonnages()
		if err != nil {
			fmt.Printf("ERREUR: %#v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		//Ajout d'un nouveau personnage
		personnages = append(personnages, aventurier)
		//sauvegarde du personnages
		err = manager.SauvegarderPersonnages(personnages)
		if err != nil {
			fmt.Printf("ERREUR: %#v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
	}
}

func MyListHandler(w http.ResponseWriter, r *http.Request) {
	personnages, err := manager.ChargerPersonnages()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Passez les données des films à votre template HTML
	perso := struct {
		Personnages []manager.Aventurier
	}{
		Personnages: personnages,
	}
	inittemplate.Temp.ExecuteTemplate(w, "myList", perso)
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
	pseudo, err := GetEmailSession(r)
	if err != nil {
		//Rediriger l'utilisateur vers la page de connexion
		http.Redirect(w, r, "/connexion", http.StatusFound)
		return
	}

	comment := manager.Comment{
		Pseudo:      pseudo,
		NomFilm:     nomFilm,
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
	pseudo, ok := session.Values["pseudo"].(string)

	if !ok {
		return "", errors.New("utilisateur non authentifié")
	}
	return pseudo, nil
}
