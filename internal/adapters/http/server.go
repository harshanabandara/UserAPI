package http

import (
	"UserApi/internal/core/domain"
	"UserApi/internal/core/ports"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	UserService ports.UserService
	Router      *chi.Mux
}

func initServer(server *Server) {
	server.Router.Get("/users", getAllUsers(server.UserService))
	server.Router.Get("/users/{userId}", getUser(server.UserService))
	server.Router.Post("/users", postUser(server.UserService))
	server.Router.Delete("/users/{userId}", deleteUser(server.UserService))
	server.Router.Patch("/users/{userId}", patchUser(server.UserService))

}

func (server *Server) Start() error {
	initServer(server)
	err := http.ListenAndServe(":8080", server.Router)
	if err != nil {
		return err
	}
	return nil
}

func getAllUsers(service ports.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := service.GetAllUsers()
		if err != nil {
			panic(err)
		}
		blob, err := json.Marshal(users)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(blob)
		if err != nil {
			return
		}
	}
}

func getUser(userService ports.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		user, err := userService.GetUserById(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		blob, _ := json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(blob)
		if err != nil {
			//log something here.
			return
		}

	}
}

func postUser(service ports.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read the request.
		user := domain.User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		createdUser, err := service.AddUser(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		// check the user id.
		if createdUser.UserID == "" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Could not create user"))
			return
		}

		err = json.NewEncoder(w).Encode(createdUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
	}
}

func patchUser(service ports.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		user := domain.User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		updateUser, err := service.UpdateUserByID(userID, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		blob, _ := json.Marshal(updateUser)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(blob)
	}
}

func deleteUser(service ports.UserService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		err := service.DeleteUserByID(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("deleted"))
		if err != nil {
			return
		}
	}
}
