package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"os"
	"strconv"
)

// HttpLogMiddleware provides a simple AuthBasic implementation. On failure, a 401 HTTP response
//is returned. On success, the wrapped middleware is called, and the userId is made available as
// request.Env["REMOTE_USER"].(string)
type HttpLogMiddleware struct {

	// Realm name to display to the user. Required.
	Realm string

	// // Callback function that should perform the authentication of the user based on userId and
	// // password. Must return true on success, false on failure. Required.
	// Authenticator func(userId string, password string) bool

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userId string, request *rest.Request) bool
}

// MiddlewareFunc makes HttpLogMiddleware implement the Middleware interface.
func (mw *HttpLogMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {

	if mw.Realm == "" {
		log.Fatal("Realm is required")
	}
	fmt.Println("fdadas")
	// if mw.Authenticator == nil {
	// 	log.Fatal("Authenticator is required")
	// }

	if mw.Authorizator == nil {
		mw.Authorizator = func(userId string, request *rest.Request) bool {
			return true
		}
	}

	return func(writer rest.ResponseWriter, request *rest.Request) {

		// authHeader := request.Header.Get("Authorization")
		// if authHeader == "" {
		// 	mw.unauthorized(writer)
		// 	return
		// }

		// providedUserId, providedPassword, err := mw.decodeBasicAuthHeader(authHeader)

		// if err != nil {
		// 	Error(writer, "Invalid authentication", http.StatusBadRequest)
		// 	return
		// }

		// if !mw.Authenticator(providedUserId, providedPassword) {
		// 	mw.unauthorized(writer)
		// 	return
		// }

		fmt.Println("fdadfas")
		if !mw.Authorizator("1", request) {
			// mw.unauthorized(writer)
			return
		}
		// request.Env["REMOTE_USER"] = providedUserId

		handler(writer, request)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("CMD Port Path ")
		return
	}
	portStr := os.Args[1]
	pathStr := os.Args[2]
	port, err := strconv.Atoi(portStr)
	if (err != nil) || (port == 0) {
		fmt.Println("Port Input error")
		return
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&HttpLogMiddleware{
		Realm: "test zone",
		Authorizator: func(userId string, request *rest.Request) bool {

			fmt.Println("fdadas")
			fmt.Println(request.Host)
			return true
		},
	})
	router, err := rest.MakeRouter(
		rest.Get("/test", GetTest),
		rest.Post("/test", PostTest),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	http.Handle("/", http.FileServer(http.Dir(pathStr)))

	log.Fatal(http.ListenAndServe(":"+portStr, nil))

}

type Config struct {
	Name  string
	Value string
}

func GetTest(w rest.ResponseWriter, r *rest.Request) {
	var config Config
	config.Name = "Get"
	config.Value = "Get return test"
	w.WriteJson(config)
}

func PostTest(w rest.ResponseWriter, r *rest.Request) {
	var config Config
	config.Name = "Post"
	config.Value = "Post return test"
	w.WriteJson(config)

}
