package gologin

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

type User interface {
	GetId() string // get the user id
}

type GoLogin struct {
	CookieStore             *sessions.CookieStore
	LoadUser                func(userid string) User
	PermissionDeniedHandler http.HandlerFunc
}

// Log the user in
//
// This stores data in a session cookie. The request context's "currentuser" variable is set to the
// user object that is passed.
func (gl *GoLogin) LoginUser(user User, w http.ResponseWriter, r *http.Request) error {
	userId := user.GetId()
	session, err := gl.CookieStore.Get(r, SESSION_NAME)

	if err != nil {
		return err
	}

	session.Values[SV_USERID] = userId
	session.Values[SV_LOGGEDIN] = true
	session.Save(r, w)
	context.Set(r, CV_CURRENTUSER, user)

	return err
}

func (gl *GoLogin) LogoutUser(w http.ResponseWriter, r *http.Request) error {
	session, err := gl.CookieStore.Get(r, SESSION_NAME)

	if err != nil {
		return err
	}

	session.Values[SV_LOGGEDIN] = false
	session.Save(r, w)
	context.Set(r, CV_CURRENTUSER, nil)

	return err
}

// Check if the user for this request is logged in
func (gl *GoLogin) IsLoggedIn(r *http.Request) bool {
	session, err := gl.CookieStore.Get(r, "user")
	if err != nil {
		return false
	}

	loggedIn, ok := session.Values[SV_LOGGEDIN].(bool)
	if !ok {
		return false
	} else {
		return loggedIn
	}
}

// HTTP middleware function
//
// This sets the request context's "currentuser" variable to the user object returned by LoadUser().
func (gl *GoLogin) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// load the user object, if any
	session, _ := gl.CookieStore.Get(r, SESSION_NAME)

	if uid, ok := session.Values[SV_USERID].(string); ok {
		user := gl.LoadUser(uid)
		if user != nil {
			context.Set(r, CV_CURRENTUSER, user)
		}
	}

	next(w, r)
}
