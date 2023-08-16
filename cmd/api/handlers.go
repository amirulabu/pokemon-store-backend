package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/flow"
	"github.com/amirulabu/pokemon-store-backend/internal/password"
	"github.com/amirulabu/pokemon-store-backend/internal/pokemon"
	"github.com/amirulabu/pokemon-store-backend/internal/request"
	"github.com/amirulabu/pokemon-store-backend/internal/response"
	"github.com/amirulabu/pokemon-store-backend/internal/utils"
	"github.com/amirulabu/pokemon-store-backend/internal/validator"

	"github.com/pascaldekloe/jwt"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func getURLQueryParamInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	return utils.GetInt(value, defaultValue)
}

func (app *application) getPokemons(w http.ResponseWriter, r *http.Request) {
	limit := getURLQueryParamInt(r, "limit", 20)
	offset := getURLQueryParamInt(r, "offset", 0)

	data, err := pokemon.GetPokemons(offset, limit, app.config.baseURL)
	if err != nil {
		app.serverError(w, r, err)
	}

	resErr := response.JSON(w, http.StatusOK, data)
	if resErr != nil {
		app.serverError(w, r, resErr)
	}
}

func (app *application) getPokemonByNameOrId(w http.ResponseWriter, r *http.Request) {
	name := flow.Param(r.Context(), "nameOrId")

	data, err := pokemon.GetSinglePokemon(name)
	if err != nil {
		app.serverError(w, r, err)
	}

	resErr := response.JSON(w, http.StatusOK, data)
	if resErr != nil {
		app.serverError(w, r, resErr)
	}
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	existingUser, err := app.db.GetUserByEmail(input.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(existingUser == nil, "Email", "Email is already in use")
	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "Email", "Must be a valid email address")

	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "Password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "Password", "Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.db.InsertUser(input.Email, hashedPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.db.GetUserByEmail(input.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(user != nil, "Email", "Email address could not be found")

	if user != nil {
		passwordMatches, err := password.Matches(input.Password, user.HashedPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	var claims jwt.Claims
	claims.Subject = strconv.Itoa(user.ID)

	expiry := time.Now().Add(24 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]string{
		"AuthenticationToken":       string(jwtBytes),
		"AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CurrentPassword string `json:"CurrentPassword"`
		NewPassword     string `json:"NewPassword"`
		Validator       validator.Validator
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	currentUser := contextGetAuthenticatedUser(r)

	input.Validator.CheckField(input.CurrentPassword != "", "CurrentPassword", "Current Password is required")
	input.Validator.CheckField(input.NewPassword != "", "NewPassword", "New Password is required")
	input.Validator.CheckField(len(input.NewPassword) >= 8, "NewPassword", "New Password is too short")
	input.Validator.CheckField(len(input.NewPassword) <= 72, "NewPassword", "New Password is too long")
	input.Validator.CheckField(validator.NotIn(input.NewPassword, password.CommonPasswords...), "NewPassword", "New Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	passwordMatches, err := password.Matches(input.CurrentPassword, currentUser.HashedPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(passwordMatches, "CurrentPassword", "Current Password is incorrect")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.NewPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.db.UpdateUserHashedPassword(currentUser.ID, hashedPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}

func (app *application) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.db.GetAllUsers()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, users)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) changePasswordById(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId      int    `json:"UserId"`
		NewPassword string `json:"NewPassword"`
		Validator   validator.Validator
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.UserId != 0, "UserId", "UserId is required")
	input.Validator.CheckField(input.NewPassword != "", "NewPassword", "New Password is required")
	input.Validator.CheckField(len(input.NewPassword) >= 8, "NewPassword", "New Password is too short")
	input.Validator.CheckField(len(input.NewPassword) <= 72, "NewPassword", "New Password is too long")
	input.Validator.CheckField(validator.NotIn(input.NewPassword, password.CommonPasswords...), "NewPassword", "New Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.NewPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.db.UpdateUserHashedPassword(input.UserId, hashedPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
