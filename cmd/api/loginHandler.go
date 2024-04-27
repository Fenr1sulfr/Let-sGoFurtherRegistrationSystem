package main

import (
	"net/http"
	"time"

	"sulfurAuth.net/internal/db"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
	token := db.TokenModel{
		DB: app.DB,
	}
	database := db.UserModel{
		DB: app.DB,
	}
	var userReq UserRequest

	err := app.readJson(w, r, &userReq)
	if err != nil {
		app.logger.Fatal("Bad request")
		return
	}

	data := &db.User{
		Email:    userReq.Email,
		Password: userReq.Password,
	}
	app.logger.Println("Processssing interactions with db...")
	err = database.Insert(data)
	if err != nil {
		app.logger.Fatal(err)
	}
	tokens, err := token.New(int64(data.ID), 3*24*time.Hour, db.ScopeActivation)
	if err != nil {
		app.logger.Fatal(err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": data}, nil)
	if err != nil {
		app.logger.Fatal(err)
	}
	dataToken := map[string]any{
		"activationToken": tokens.Plaintext,
		"userID":          data.ID,
	}
	err = app.mailer.Send(data.Email, "welcome.tmpl", dataToken)
	if err != nil {
		app.logger.Fatal(err)
	}
	app.logger.Println("Success!")
}
