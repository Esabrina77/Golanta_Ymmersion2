package routeur

import (
	"BlogYmmersion/controller"
	inittemplate "BlogYmmersion/templates"
	"fmt"
	"log"
	"net/http"
)

func InitServe() {
	FileServer := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", FileServer))
	http.HandleFunc("/connexion", controller.ConnexionHandler)
	http.HandleFunc("/inscription", controller.InscriptionHandler)
	http.HandleFunc("/home", controller.HomeHandler)
	http.HandleFunc("/category", controller.CategoryHandler)
	http.HandleFunc("/comments", controller.CommentsHandler)
	http.HandleFunc("/treatmentI", controller.TreatInscriptionHandler)
	http.HandleFunc("/treatmentC", controller.TreatConnexionHandler)
	http.HandleFunc("/submitComments", controller.SubmitCommentHandler)
	http.HandleFunc("/myList", controller.MyListHandler)
	http.HandleFunc("/ajouter-film", controller.AjouterFilmHandler)
	http.HandleFunc("/confirmation", controller.ConfirmationHandler)
	http.HandleFunc("/form", controller.FormHandler)
	http.HandleFunc("/404", controller.NotFoundHandler)
	http.HandleFunc("/search", controller.SearchHandler)
	http.HandleFunc("/notFound", controller.RessourceNotFoundHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		inittemplate.Temp.ExecuteTemplate(w, "404", nil)
	})
	if err := http.ListenAndServe(controller.Port, nil); err != nil {

		fmt.Printf("ERREUR LORS DE L'INITIATION DES ROUTES %v \n", err)

		log.Fatal(err)

	}
}
