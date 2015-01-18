package gologin

const VERSION = "0.1.0"

// request context variable names
const (
	CV_CURRENTUSER = "currentuser" // currently logged-in user (User)
)

// session cookie variable names
const (
	SESSION_NAME = "user"     // gorilla sessions session name
	SV_USERID    = "userid"   // stored user id (string)
	SV_LOGGEDIN  = "loggedin" // logged in status (bool)
)
