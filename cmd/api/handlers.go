package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/flow"
	"github.com/amirulabu/pokemon-store-backend/internal/password"
	"github.com/amirulabu/pokemon-store-backend/internal/request"
	"github.com/amirulabu/pokemon-store-backend/internal/response"
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

func (app *application) getPokemons(w http.ResponseWriter, r *http.Request) {
	type Results struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	type PokemonList struct {
		Count    int         `json:"count"`
		Next     string      `json:"next"`
		Previous interface{} `json:"previous"`
		Results  []Results   `json:"results"`
	}

	res, err := http.Get("https://pokeapi.co/api/v2/pokemon?offset=20&limit=20")
	if err != nil {
		app.serverError(w, r, err)
	}

	var data PokemonList

	jsonErr := json.NewDecoder(res.Body).Decode(&data)
	if jsonErr != nil {
		app.serverError(w, r, jsonErr)
	}

	resErr := response.JSON(w, http.StatusOK, data)
	if resErr != nil {
		app.serverError(w, r, resErr)
	}
}

func (app *application) getPokemonByNameOrId(w http.ResponseWriter, r *http.Request) {
	type SinglePokemon struct {
		Abilities []struct {
			Ability struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
		BaseExperience int `json:"base_experience"`
		Forms          []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"forms"`
		GameIndices []struct {
			GameIndex int `json:"game_index"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"game_indices"`
		Height                 int    `json:"height"`
		HeldItems              []any  `json:"held_items"`
		ID                     int    `json:"id"`
		IsDefault              bool   `json:"is_default"`
		LocationAreaEncounters string `json:"location_area_encounters"`
		Moves                  []struct {
			Move struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move"`
			VersionGroupDetails []struct {
				LevelLearnedAt  int `json:"level_learned_at"`
				MoveLearnMethod struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"move_learn_method"`
				VersionGroup struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"version_group"`
			} `json:"version_group_details"`
		} `json:"moves"`
		Name      string `json:"name"`
		Order     int    `json:"order"`
		PastTypes []any  `json:"past_types"`
		Species   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"species"`
		Sprites struct {
			BackDefault      string `json:"back_default"`
			BackFemale       any    `json:"back_female"`
			BackShiny        string `json:"back_shiny"`
			BackShinyFemale  any    `json:"back_shiny_female"`
			FrontDefault     string `json:"front_default"`
			FrontFemale      any    `json:"front_female"`
			FrontShiny       string `json:"front_shiny"`
			FrontShinyFemale any    `json:"front_shiny_female"`
			Other            struct {
				DreamWorld struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"dream_world"`
				Home struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"home"`
				OfficialArtwork struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"official-artwork"`
			} `json:"other"`
			Versions struct {
				GenerationI struct {
					RedBlue struct {
						BackDefault      string `json:"back_default"`
						BackGray         string `json:"back_gray"`
						BackTransparent  string `json:"back_transparent"`
						FrontDefault     string `json:"front_default"`
						FrontGray        string `json:"front_gray"`
						FrontTransparent string `json:"front_transparent"`
					} `json:"red-blue"`
					Yellow struct {
						BackDefault      string `json:"back_default"`
						BackGray         string `json:"back_gray"`
						BackTransparent  string `json:"back_transparent"`
						FrontDefault     string `json:"front_default"`
						FrontGray        string `json:"front_gray"`
						FrontTransparent string `json:"front_transparent"`
					} `json:"yellow"`
				} `json:"generation-i"`
				GenerationIi struct {
					Crystal struct {
						BackDefault           string `json:"back_default"`
						BackShiny             string `json:"back_shiny"`
						BackShinyTransparent  string `json:"back_shiny_transparent"`
						BackTransparent       string `json:"back_transparent"`
						FrontDefault          string `json:"front_default"`
						FrontShiny            string `json:"front_shiny"`
						FrontShinyTransparent string `json:"front_shiny_transparent"`
						FrontTransparent      string `json:"front_transparent"`
					} `json:"crystal"`
					Gold struct {
						BackDefault      string `json:"back_default"`
						BackShiny        string `json:"back_shiny"`
						FrontDefault     string `json:"front_default"`
						FrontShiny       string `json:"front_shiny"`
						FrontTransparent string `json:"front_transparent"`
					} `json:"gold"`
					Silver struct {
						BackDefault      string `json:"back_default"`
						BackShiny        string `json:"back_shiny"`
						FrontDefault     string `json:"front_default"`
						FrontShiny       string `json:"front_shiny"`
						FrontTransparent string `json:"front_transparent"`
					} `json:"silver"`
				} `json:"generation-ii"`
				GenerationIii struct {
					Emerald struct {
						FrontDefault string `json:"front_default"`
						FrontShiny   string `json:"front_shiny"`
					} `json:"emerald"`
					FireredLeafgreen struct {
						BackDefault  string `json:"back_default"`
						BackShiny    string `json:"back_shiny"`
						FrontDefault string `json:"front_default"`
						FrontShiny   string `json:"front_shiny"`
					} `json:"firered-leafgreen"`
					RubySapphire struct {
						BackDefault  string `json:"back_default"`
						BackShiny    string `json:"back_shiny"`
						FrontDefault string `json:"front_default"`
						FrontShiny   string `json:"front_shiny"`
					} `json:"ruby-sapphire"`
				} `json:"generation-iii"`
				GenerationIv struct {
					DiamondPearl struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"diamond-pearl"`
					HeartgoldSoulsilver struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"heartgold-soulsilver"`
					Platinum struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"platinum"`
				} `json:"generation-iv"`
				GenerationV struct {
					BlackWhite struct {
						Animated struct {
							BackDefault      string `json:"back_default"`
							BackFemale       any    `json:"back_female"`
							BackShiny        string `json:"back_shiny"`
							BackShinyFemale  any    `json:"back_shiny_female"`
							FrontDefault     string `json:"front_default"`
							FrontFemale      any    `json:"front_female"`
							FrontShiny       string `json:"front_shiny"`
							FrontShinyFemale any    `json:"front_shiny_female"`
						} `json:"animated"`
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"black-white"`
				} `json:"generation-v"`
				GenerationVi struct {
					OmegarubyAlphasapphire struct {
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"omegaruby-alphasapphire"`
					XY struct {
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"x-y"`
				} `json:"generation-vi"`
				GenerationVii struct {
					Icons struct {
						FrontDefault string `json:"front_default"`
						FrontFemale  any    `json:"front_female"`
					} `json:"icons"`
					UltraSunUltraMoon struct {
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"ultra-sun-ultra-moon"`
				} `json:"generation-vii"`
				GenerationViii struct {
					Icons struct {
						FrontDefault string `json:"front_default"`
						FrontFemale  any    `json:"front_female"`
					} `json:"icons"`
				} `json:"generation-viii"`
			} `json:"versions"`
		} `json:"sprites"`
		Stats []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		} `json:"stats"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
		Weight int `json:"weight"`
	}
	name := flow.Param(r.Context(), "nameOrId")
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	res, err := http.Get(url)
	if err != nil {
		app.serverError(w, r, err)
	}

	var data SinglePokemon

	jsonErr := json.NewDecoder(res.Body).Decode(&data)
	if jsonErr != nil {
		app.serverError(w, r, jsonErr)
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

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "Email", "Must be a valid email address")
	input.Validator.CheckField(existingUser == nil, "Email", "Email is already in use")

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

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
