package http

import (
	"encoding/json"
	"errors"
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
			serverErr := fmt.Errorf("error getting all users: %w", err)
			fmt.Println(serverErr)
			http.Error(w, errors.New("could not retrieve the users").Error(), http.StatusInternalServerError)
			return
		}
		userDTOs := make([]UserResponse, len(users))
		for i, user := range users {
			userDTOs[i] = parseUserToUserDTO(user)
		}
		blob, err := json.Marshal(users)
		if err != nil {
			marshalErr := fmt.Errorf("error marshalling users: %w", err)
			fmt.Println(marshalErr)
			http.Error(w, errors.New("could not retrieve the users").Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(blob)
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
		if err != nil {
			serverErr := fmt.Errorf("could not retrieve the user: %w", err)
			fmt.Println(serverErr)
			notFoundErr := fmt.Errorf("could not retrieve the user %s", userID)
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(notFoundErr.Error()))
			return
		}
		userDTO := parseUserToUserDTO(user)
		blob, err := json.Marshal(userDTO)
		if err != nil {
			marshalErr := fmt.Errorf("error marshalling user: %w", err)
			fmt.Println(marshalErr)
			http.Error(w, errors.New("could not retrieve the user").Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(blob)
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
			decodeError := fmt.Errorf("could not decode the request body: %w", err)
			fmt.Println(decodeError)
			_, _ = w.Write([]byte(decodeError.Error()))
			return
		}
		validationErr := validator.Struct(user)
		if validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			validationErr2 := fmt.Errorf("could not validate the request: %w", validationErr)
			fmt.Println(validationErr2)
			_, _ = w.Write([]byte(validationErr2.Error()))
			return
		}
		createdUser, err := service.AddUser(r.Context(), user.getUser())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			userErr := fmt.Errorf("could not add the user: %w", err)
			fmt.Println(userErr)

			_, _ = w.Write([]byte(errors.New("could not add the user").Error()))
			return
		}
		// check the user id.
		if createdUser.UserID == "" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(errors.New("could not create user").Error()))
			return
		}
		parsedUser := parseUserToUserDTO(createdUser)
		blob, err := json.Marshal(parsedUser)
		if err != nil {
			parsingErr := fmt.Errorf("could not parse the user: %w", err)
			fmt.Println(parsingErr)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(errors.New("could not retrieve the user").Error()))
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
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			decodeError := fmt.Errorf("could not decode the request body: %w", err)
			fmt.Println(decodeError)
			_, _ = w.Write([]byte(decodeError.Error()))
			return
		}
		validationErr := validator.Struct(user)
		if validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(validationErr.Error()))
			return
		}
		updateUser, err := service.UpdateUserByID(r.Context(), userID, user.getUser())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			userErr := fmt.Errorf("could not update user: %w", err)
			fmt.Println(userErr)
			_, _ = w.Write([]byte(errors.New("could not update user").Error()))
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
			deleteErr := fmt.Errorf("could not delete user: %w", err)
			fmt.Println(deleteErr)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(errors.New("could not delete user").Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("deleted"))
		if err != nil {
			return
		}
	}
}
