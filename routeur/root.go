package routeur

import (
	"Golanta/controller"
	inittemplate "Golanta/templates"
	"fmt"
	"log"
	"net/http"
)

func InitServe() {
	FileServer := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", FileServer))
	http.HandleFunc("/connexion", controller.ConnexionHandler)
	http.HandleFunc("/inscription", controller.InscriptionHandler)
	http.HandleFunc("/comments", controller.CommentsHandler)
	http.HandleFunc("/treatmentI", controller.TreatInscriptionHandler)
	http.HandleFunc("/treatmentC", controller.TreatConnexionHandler)
	http.HandleFunc("/submitComments", controller.SubmitCommentHandler)
	http.HandleFunc("/form", controller.FormHandler)
	http.HandleFunc("/404", controller.NotFoundHandler)
	http.HandleFunc("/notFound", controller.RessourceNotFoundHandler)
	http.HandleFunc("/story", controller.StoryHandler)
	http.HandleFunc("/myList", controller.MyListHandler)
	http.HandleFunc("/ajout_perso", controller.AjouterPersoHandler)
	http.HandleFunc("/confirmation", controller.ConfirmationHandler)
	http.HandleFunc("/confirmationSup", controller.ConfirmationSupHandler)
	http.HandleFunc("/supprimer", controller.SupprimerPersoHandler)
	http.HandleFunc("/modifier", controller.ModifierPersoHandler)
	http.HandleFunc("/search", controller.RecherchePersonnageHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		inittemplate.Temp.ExecuteTemplate(w, "404", nil)
	})
	if err := http.ListenAndServe(controller.Port, nil); err != nil {

		fmt.Printf("ERREUR LORS DE L'INITIATION DES ROUTES %v \n", err)

		log.Fatal(err)

	}
}
