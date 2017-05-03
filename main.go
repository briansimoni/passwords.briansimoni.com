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
	"golang.org/x/net/html"
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

	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/secrets", secrets)
	http.HandleFunc("/delete-secret", deleteSecret)
	http.HandleFunc("/edit-secret", editSecret)
	http.HandleFunc("/addSecret", addSecret)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":" + port, context.ClearHandler(http.DefaultServeMux))
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
		http.Redirect(w, r, html.EscapeString("/?status=" + urlParam), http.StatusFound)
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
	requestedApplication := r.URL.Query().Get("app")
	userEmail := session.Values["userEmail"].(string)
	user := getUser(userEmail)

	//t := template.New("secrets template") // Create a template.
	t, err := template.ParseFiles("public/templates/secrets.html")  // Parse template file.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	if len(requestedApplication) != 0 {
		fmt.Println(requestedApplication)
		user.PageData = decryptUserSecret(user.Secrets, requestedApplication)
		err = t.Execute(w, user)  // merge.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		user.PageData = decryptUserApplications(user.Secrets)
		err = t.Execute(w, user)  // merge.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func deleteSecret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "simoni-session")
	if session.Values["authenticated"] != true {
		http.Error(w, "You are not authorized to view this page", http.StatusForbidden)
		return
	}

	email := session.Values["userEmail"].(string)
	user := getUser(email)
	app := r.URL.Query().Get("app")
	var appToDelete string
	for application, _ := range user.Secrets {
		if (decrypt(application) == app) {
			appToDelete = application
			break
		}
	}
	err := deleteSecretFromDB(appToDelete, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/secrets?status=deleted", http.StatusFound)
}

func editSecret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "simoni-session")
	if session.Values["authenticated"] != true {
		http.Error(w, "You are not authorized to view this page", http.StatusForbidden)
		return
	}
	email := session.Values["userEmail"].(string)
	user := getUser(email)

	r.ParseForm()
	application := r.Form.Get("application")
	password := r.Form.Get("password")

	fmt.Println("this is the application " + application)
	fmt.Println("this is the password " + password)

	var appToEdit string
	for app, _ := range user.Secrets {
		if decrypt(app) == application {
			appToEdit = app
			break
		}
	}
	err := updateApplicationPassword(encrypt(password), appToEdit, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/secrets?status=updated", http.StatusFound)

}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user := getUser(email)
	valid := compareHash(password, user.Hash)
	if !valid {
		urlParam := url.QueryEscape("Username or password is incorrect")
		http.Redirect(w, r, html.EscapeString("/?status=" + urlParam), http.StatusFound)
		return
	}


	session, _ := store.Get(r, "simoni-session")
	// Set some session values.
	session.Values["authenticated"] = true
	session.Values["userEmail"] = email
	session.Values["userId"] = user.Id

	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
	http.Redirect(w, r, "/secrets", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "simoni-session")
	delete(session.Values, "authenticated")
	delete(session.Values, "userEmail")
	delete(session.Values, "userId")
	session.Save(r, w)
	http.Redirect(w, r, "/?status=You have logged out successfully.", http.StatusFound)
}

func addSecret(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	application := r.Form.Get("application")
	password := r.Form.Get("password")

	session, _ := store.Get(r, "simoni-session")
	userId := session.Values["userId"].(string)
	user := getUser(session.Values["userEmail"].(string))
	user.PageData = decryptUserApplications(user.Secrets)

	// duplicate := false
	for app, _ := range user.PageData {
		if(app == application) {
			http.Redirect(w, r, "/secrets?status=failed", http.StatusFound)
		}
	}


	err := insertSecret(application, password, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/secrets?status=success", http.StatusFound)
}
