package main

import (
	"ebpf-firewall/dbLayer"
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexedwards/argon2id"
)

type accountSignUp struct {
	UserName  string `validate:"required,min=5,max=64"`
	Passwd    string `validate:"required,min=10,max=64"`
	FirstName string `validate:"required,min=2,max=64"`
	LastName  string `validate:"required,min=2,max=64"`
	Email     string `validate:"required,email,min=8,max=64"`
}

func (app *application) postAccountSignUp(w http.ResponseWriter, r *http.Request) {
	var acc accountSignUp
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		log.Printf("<ERROR>\t\t[(Sign-up)json decode failed]\n%s\n\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := app.validate.Struct(acc); err != nil {
		log.Printf("<WARNING>\t\t[(Sign-up)json fields failed to match struct requirements]\n%s\n\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := argon2id.CreateHash(acc.Passwd, argon2id.DefaultParams)
	if err != nil {
		log.Printf("<ERROR>\t\t[(Sign-up)argon2id failed to create hash]\n%s\n\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	arg := dbLayer.CreateAccountParams{
		Username:   acc.UserName,
		Passwdhash: hash,
		Powerlevel: 0,
		Firstname:  acc.FirstName,
		Lastname:   acc.LastName,
		Email:      acc.Email,
	}
	if err := app.queries.CreateAccount(app.ctx, arg); err != nil {
		log.Printf("<WARNING>\t\t[(Sign-up)failed to create account]\n%s\n\n", err)
		w.WriteHeader(http.StatusConflict)
		return
	} else {
		log.Printf("<INFO>\t\t[(Sign-up)succesfully created user]\nuser name: %s\n\n", arg.Username)
		return
	}
}

type accountSignIn struct {
	UserName string `validate:"required,min=5,max=64"`
	Passwd   string `validate:"required,min=10,max=64"`
}

func (app *application) postAccountSignIn(w http.ResponseWriter, r *http.Request) {
	var acc accountSignIn
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		log.Printf("<ERROR>\t\t[(Sign-in)json decode failed]\n%s\n\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := app.validate.Struct(acc); err != nil {
		log.Printf("<WARNING>\t\t[(Sign-in)json fields failed to match struct requirements]\n%s\n\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accStruct, err := app.queries.RetrieveAccount(app.ctx, acc.UserName)
	if err != nil {
		log.Printf("<WARNING>\t\t[(Sign-in)failed to retrieve account]\n%s\n\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(acc.Passwd, accStruct.Passwdhash)
	if err != nil {
		log.Printf("<ERROR>\t\t[(Sign-in)argon2id failed to create hash]\n%s\n\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if match {
		if err := app.queries.DeleteBearerToken(app.ctx, acc.UserName); err != nil {
			log.Printf("<ERROR>\t\t[(Sign-in)failed to invalidate bearer token]\n%s\n\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Printf("<INFO>\t\t[(Sign-in)succesfully invalidated bearer tokens]\nuser name: %s\n\n", acc.UserName)
		}
		if err := app.generateBearerToken(w, accStruct); err != nil {
			log.Printf("<ERROR>\t\t[(Sign-in)failed to generate bearer token]\n%s\n\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("<INFO>\t\t[(Sign-in)succesfull sign-in]\nuser name: %s\n\n", acc.UserName)
		return
	} else {
		log.Printf("<WARNING>\t\t[(Sign-in)password does not match]\nuser name: %s\n\n", acc.UserName)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
