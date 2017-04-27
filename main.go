package main

import (
	"net/http"
	"os"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
	"errors"
	"net/url"
	"html/template"
)

var masterKey string
var store *sessions.CookieStore

func main() {

	key, err := getMasterKey()
	if err != nil {
		panic(errors.New("Unable to read master key"))
	}
	masterKey = key
	store = sessions.NewCookieStore([]byte(masterKey))


	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	fmt.Println("Listening on port " + port)

	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/secrets", secrets)
	http.HandleFunc("/addSecret", addSecret)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":" + port, context.ClearHandler(http.DefaultServeMux))
}


// hello world, the web server
func HelloServer(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "simoni-session")
	if session.Values["test"] != nil {
		var value string
		value = session.Values["test"].(string)
		fmt.Fprint(w, "this is a database test: " + value)
	} else {
		fmt.Fprint(w, "there is nothing in the session")
	}
}


func signup(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	fmt.Println("password from the form: " + password)
	hash := generateHash(password)
	err := insertUser(email, hash)
	if err != nil {
		fmt.Println(err.Error())
		urlParam := url.QueryEscape(err.Error())
		http.Redirect(w, r, "/?status=" + urlParam, http.StatusFound)
		return
	}

	session, _ := store.Get(r, "simoni-session")
	// Set some session values.
	session.Values["authenticated"] = true
	session.Values["userEmail"] = email

	//TODO: just return the user id from the insert function
	user := getUser(email)
	session.Values["userId"] = user.Id

	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
	http.Redirect(w, r, "/secrets", http.StatusFound)
}

func secrets(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "simoni-session")
	if session.Values["authenticated"] != true {
		http.Error(w, "You are not authorized to view this page", http.StatusForbidden)
		return
	}
	userEmail := session.Values["userEmail"].(string)
	user := getUser(userEmail)

	//t := template.New("secrets template") // Create a template.
	t, err := template.ParseFiles("public/templates/secrets.html")  // Parse template file.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Secrets = decryptUserSecrets(user.Secrets)
	err = t.Execute(w, user)  // merge.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user := getUser(email)
	valid := compareHash(password, user.Hash)
	if !valid {
		urlParam := url.QueryEscape("Username or password is incorrect")
		http.Redirect(w, r, "/?status=" + urlParam, http.StatusFound)
		return
	}


	session, _ := store.Get(r, "simoni-session")
	// Set some session values.
	session.Values["authenticated"] = true
	session.Values["authenticated"] = true
	session.Values["userEmail"] = email
	session.Values["userId"] = user.Id

	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
	http.Redirect(w, r, "/secrets", http.StatusFound)
}

func addSecret(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	application := r.Form.Get("application")
	password := r.Form.Get("password")

	session, _ := store.Get(r, "simoni-session")
	userId := session.Values["userId"].(string)

	err := insertSecret(application, password, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/secrets?status=success", http.StatusFound)
}
