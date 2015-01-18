GoLogin
=======

GoLogin is a login middleware for Go web apps inspired by Flask login.

Example usage (main.go):

	package main

	import (
		"github.com/veegee/gologin"
		"github.com/gorilla/mux"
	)

	var CookieStore *sessions.CookieStore
	var LoginManager *gologin.GoLogin

	func LoadUser(userid string) gologin.User {
		user := &dbmodels.User{}

		// ... query the database and populate the user object ...

		return user
	}

	func PermissionDenied(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "error: not logged in")
	}

	func Index(w http.ResponseWriter, r *http.Request) {
		if !LoginManager.IsLoggedIn(r) {
			LoginManager.PermissionDeniedHandler(w, r)
			return
		}
	}

	func Login(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := context.Get(r, gologin.CV_CURRENTUSER).(*dbmodels.Entity)

		if currentUser != nil && ok {
			// already logged in
		}

		userid := ... look up user by email and verify password ...

		user := LoginManager.LoadUser(userid)

		if user != nil {
			LoginManager.LoginUser(user, w, r)
			fmt.Fprintf(w, "Login success")
		} else {
			http.Error(w, "Login failed", http.StatusForbidden)
		}
	}

	func main() {
		CookieStore = sessions.NewCookieStore([]byte("something-very-secret"))
		LoginManager = &gologin.GoLogin{CookieStore, LoadUser, PermissionDenied}

		router := mux.NewRouter()
		router.HandleFunc("/", Index)
		router.HandleFunc("/login", Login)

		n := negroni.New()
		n.Use(globals.LoginManager)
		n.UseHandler(router)
		n.Run("0.0.0.0:8002")
	}
