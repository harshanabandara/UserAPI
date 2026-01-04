package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"userapi/app/internal/core/domain"
	"userapi/app/internal/core/ports"
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
		users, err := service.GetAllUsers(r.Context())
		if err != nil {
			panic(err)
		}
		userDTOs := make([]UserDTO, len(users))
		for i, user := range users {
			userDTOs[i] = parseUserToUserDTO(user)
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
		user, err := userService.GetUserById(r.Context(), userID)
		userDTO := parseUserToUserDTO(user)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		blob, _ := json.Marshal(userDTO)
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
		createdUser, err := service.AddUser(r.Context(), user)
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
		parsedUser := parseUserToUserDTO(createdUser)
		blob, err := json.Marshal(parsedUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(blob)
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
		updateUser, err := service.UpdateUserByID(r.Context(), userID, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		blob, _ := json.Marshal(updateUser)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(blob)
		if err != nil {
			return
		}
	}
}

func deleteUser(service ports.UserService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		err := service.DeleteUserByID(r.Context(), userID)
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
