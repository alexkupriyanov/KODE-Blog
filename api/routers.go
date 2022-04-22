package routes

import (
	"KODE-Blog/api/controllers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World!")
	_ = r
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"Delete message",
		strings.ToUpper("Delete"),
		"/message/{id}",
		controllers.DeleteMessage,
	},

	Route{
		"Get message details",
		strings.ToUpper("Get"),
		"/message/{id}",
		controllers.GetMessageDetails,
	},

	Route{
		"Like message",
		strings.ToUpper("Post"),
		"/message/{id}/like",
		controllers.Like,
	},

	Route{
		"Get page 1 message list",
		strings.ToUpper("Get"),
		"/message/list",
		controllers.GetMessageList,
	},

	Route{
		"Get page {page} message list",
		strings.ToUpper("Get"),
		"/message/list/{page}",
		controllers.GetMessageList,
	},

	Route{
		"Crete message",
		strings.ToUpper("Post"),
		"/message",
		controllers.CreateMessage,
	},

	Route{
		"Login",
		strings.ToUpper("Post"),
		"/user/login",
		controllers.Login,
	},

	Route{
		"Logout",
		strings.ToUpper("Post"),
		"/user/logout",
		controllers.Logout,
	},

	Route{
		"Create user",
		strings.ToUpper("Post"),
		"/user",
		controllers.CreateUser,
	},

	Route{
		"Download",
		strings.ToUpper("Get"),
		"/files/{file}",
		controllers.Download,
	},
	Route{
		"Healthcheck",
		strings.ToUpper("Get"),
		"/health",
		func (w http.ResponseWriter, r *http.Request) {
			_ = w
			_ = r
		    return
		},
	},
}
