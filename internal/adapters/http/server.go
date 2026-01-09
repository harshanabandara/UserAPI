package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "userapi/app/docs"
	"userapi/app/internal/core/ports"
)

type Server struct {
	UserService ports.UserService
	Router      *chi.Mux
	Validator   *validator.Validate
}

func initServer(server *Server) {
	server.Router.Get("/users", getAllUsers(server.UserService))
	server.Router.Get("/users/{userId}", getUser(server.UserService))
	server.Router.Post("/users", postUser(server.UserService, server.Validator))
	server.Router.Delete("/users/{userId}", deleteUser(server.UserService))
	server.Router.Patch("/users/{userId}", patchUser(server.UserService, server.Validator))
	// assign docs.
	server.Router.Get("/doc", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/doc/index.html", http.StatusMovedPermanently)
	})
	server.Router.Get("/doc/*", httpSwagger.WrapHandler)
}

func (server *Server) Start() error {
	initServer(server)
	err := http.ListenAndServe(":8080", server.Router)
	if err != nil {
		return err
	}
	return nil
}

// GetAllUsers godoc
//
//	@Summary		Get all users
//	@Description	Retrieves all users from the database.
//	@Tags users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array} UserResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/users [get]
func getAllUsers(service ports.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := service.GetAllUsers(r.Context())
		if err != nil {
			panic(err)
		}
		userDTOs := make([]UserResponse, len(users))
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

// GetUser godoc
//
//	@Summary		Get all users
//	@Description	Retrieves all users from the database.
//	@Tags users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array} UserResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{user_id} [get]
//	@Param user_id  path string true "User ID"
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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a user with first name, last name, and email, and other optional data.
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User payload"
// @Success 201 {object} UserResponse
// @Failure 400 {object} UserResponse
// @Router /users [post]
func postUser(service ports.UserService, validator *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read the request.
		user := CreateUserRequest{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		validationErr := validator.Struct(user)
		if validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(validationErr.Error()))
			return
		}
		createdUser, err := service.AddUser(r.Context(), user.getUser())
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

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update a user with first name, last name, and email.
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRequest true "User payload"
// @Success 201 {object} UserResponse
// @Failure 400 {object} UserResponse
// @Router /users/{user_id} [patch]
// @Param user_id  path string true "User ID"
func patchUser(service ports.UserService, validator *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		user := UserRequest{}

		err := json.NewDecoder(r.Body).Decode(&user)
		fmt.Println("user", user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		validationErr := validator.Struct(user)
		if validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(validationErr.Error()))
			return
		}
		fmt.Println("user", user.getUser())
		updateUser, err := service.UpdateUserByID(r.Context(), userID, user.getUser())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		userDTO := parseUserToUserDTO(updateUser)
		w.WriteHeader(http.StatusOK)
		blob, _ := json.Marshal(userDTO)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(blob)
		if err != nil {
			return
		}
	}
}

// DeleteUser godoc
// @Summary Delete an existing user
// @Description Delete a user by user id
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /users/{user_id} [delete]
// @Param user_id  path string true "User ID"
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
