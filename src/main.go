package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	hydra  *hydra
}

func main() {
	CSRF := csrf.Protect(
		[]byte(os.Getenv("SYSTEM_SECRET")),
		csrf.Secure(false),
	)

	router := mux.NewRouter()
	router.Use(commonMiddleware)

	hydra := &hydra{
		client: new(http.Client),
	}

	router.HandleFunc("/", handleLayout(map[string]interface{}{
		"Title":         "Homepage",
		"HydraLoginURL": hydra.generateAuthenticationEndpoint(hydra.getOAuthClient()),
	}, "views/layout.html", "views/index.html"))
	router.HandleFunc("/callback", handleCallback)
	router.HandleFunc("/error", handleCallback)
	router.HandleFunc("/login", HandleLoginGet).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", CSRF(router)))
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Error or callback")
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		next.ServeHTTP(w, r)
	})
}

func handleLayout(context map[string]interface{}, filenames ...string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles(filenames...)
		log.Println(context)
		t.ExecuteTemplate(w, "layout", context)
	}
}

func HandleLoginGet(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/layout.html", "views/login.html")

	context := map[string]interface{}{
		"Title":          "Login",
		"Challenge":      r.URL.Query().Get("login_challenge"),
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	log.Println(context)
	t.ExecuteTemplate(w, "layout", context)
}
