package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AndresLince/go-rest-websockets/models"
	"github.com/AndresLince/go-rest-websockets/repository"
	"github.com/AndresLince/go-rest-websockets/server"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type InsertPostRequest struct {
	PostContent string `json:"post_content"`
}
type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest = InsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			post := models.Post{
				Id:          id.String(),
				PostContent: postRequest.PostContent,
				UserId:      claims.UserId,
			}
			err = repository.InsertPost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostResponse{
				Id:          post.Id,
				PostContent: post.PostContent,
			})
		} else {
			fmt.Println("Entra4")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}
