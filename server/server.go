package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/aarnaud/http-mitigation/config"
	"time"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"github.com/aarnaud/http-mitigation/db"
	"github.com/google/uuid"
	"github.com/dgrijalva/jwt-go"
)

// TODO: Use configuration for secret
var mySigningKey = []byte("secret")

func Start() {
	router := mux.NewRouter()

	// Route for the challenge
	router.HandleFunc("/__protection/{challenge}", func(w http.ResponseWriter, r *http.Request) {
		challengeHandler(w, r)
	})
	// All route catch by the defaultHandler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defaultHandler(w, r)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Config.HTTPPort), router))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// Get request cookie
	cookie, _ := r.Cookie(config.Config.CookieName)
	domain := r.Header.Get("X-Original-Host")

	// If token in cookie is valid
	if cookie != nil {
		// Get token from cookie
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		if err != nil {
			log.Error(err)
		}
		//Check the token
		if token.Valid {
			log.Debugf("Token valid")
			w.WriteHeader(200)
			return
		}
	}

	_, _, allowed := db.Limiter.AllowMinute(domain, config.Config.Threshold1)
	if allowed {
		w.WriteHeader(200)
		return
	}

	// anything else ask to resolve the challenge
	getChallenge(w, r)
}

func challengeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// If challenge request, check validity
	url := db.Client.Get(fmt.Sprintf("challenge:%s", vars["challenge"])).Val()
	if url != "" {
		expires := time.Now().Add(time.Minute)

		// Generate JWT token
		claims := make(jwt.MapClaims)
		claims["exp"] = expires.Unix()
		claims["iat"] = time.Now().Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(mySigningKey)

		if err != nil {
			log.Error(err)
			w.WriteHeader(500)
		}
		// Test cookie value...
		cookie := &http.Cookie{
			Name:    config.Config.CookieName,
			Path:    "/",
			Value:   tokenString,
			Expires: expires,
		}
		http.SetCookie(w, cookie)
		log.Debugf("Challenge resolved")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(403)
}


func getChallenge(w http.ResponseWriter, r *http.Request){
	var url string
	uri := r.Header.Get("X-Original-URI")
	query := r.Header.Get("X-Original-Query")
	token, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
	// Client has 2 second to resolve the challenge
	if query != "" {
		url = fmt.Sprintf("%s%s", uri, query)
	} else {
		url = uri
	}
	db.Client.Set(fmt.Sprintf("challenge:%s", token.String()), url, 2*time.Second)
	w.Header().Add("X-challenge", fmt.Sprintf("/__protection/%s", token.String()))
	w.WriteHeader(200)
}