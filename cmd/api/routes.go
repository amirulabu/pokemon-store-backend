package main

import (
	"net/http"

	"github.com/alexedwards/flow"
)

func (app *application) routes() http.Handler {
	mux := flow.New()

	mux.NotFound = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.Use(app.recoverPanic)
	mux.Use(app.authenticate)

	mux.HandleFunc("/", app.home, "GET")

	mux.HandleFunc("/another", app.another, "GET")

	mux.HandleFunc("/status", app.status, "GET")
	mux.HandleFunc("/users", app.createUser, "POST")
	mux.HandleFunc("/authentication-tokens", app.createAuthenticationToken, "POST")
	mux.HandleFunc("/pokemon/:nameOrId", app.getPokemonByNameOrId, "GET")
	mux.HandleFunc("/pokemon", app.getPokemons, "GET")

	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requireAuthenticatedUser)

		mux.HandleFunc("/protected", app.protected, "GET")
		mux.HandleFunc("/change-password", app.changePassword, "POST")

		mux.Group(func(mux *flow.Mux) {
			mux.Use(app.requireAdminUser)

			mux.HandleFunc("/admin/protected", app.protected, "GET")
			mux.HandleFunc("/admin/users", app.getAllUsers, "GET")
			mux.HandleFunc("/admin/change-user-password", app.changePasswordById, "POST")
		})
	})

	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requireBasicAuthentication)

		mux.HandleFunc("/basic-auth-protected", app.protected, "GET")
	})

	return mux
}
