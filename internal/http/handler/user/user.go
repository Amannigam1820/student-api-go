package user

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Amannigam1820/student-api-go/internal/storage"
	"github.com/Amannigam1820/student-api-go/internal/types"
	"github.com/Amannigam1820/student-api-go/internal/utils/response"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwt_key = []byte("student_api_go")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func RegisterUser(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credetials types.User

		if err := json.NewDecoder(r.Body).Decode(&credetials); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request")))
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credetials.Password), bcrypt.DefaultCost)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request")))
			return
		}
		lastId, err := storage.RegisterUser(credetials.Username, string(hashedPassword))
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}

		slog.Info("User created SuccessFully", slog.String("UserId", fmt.Sprint(lastId)))
		response.WriteJson(w, http.StatusCreated, map[string]interface{}{"id": lastId, "message": "User Register Successfully"})

	}
}

func Login(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credential types.User

		if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request")))
			return
		}

		slog.Info("Received login request")

		user, err := storage.GetUserByUsername(credential.Username)
		if err != nil {
			slog.Error(err.Error())
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request")))
			return
		}

		// slog.Info("step2")

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password)); err != nil {

			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid username or password")))
			return

		}
		expirationTime := time.Now().Add(24 * time.Hour)
		Claims := &Claims{
			Username: user.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
		tokenString, err := token.SignedString(jwt_key)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("internal servr error")))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		slog.Info("User logged in successfully")
		response.WriteJson(w, http.StatusOK, map[string]interface{}{

			"message": "User logged in successfully",
			"token":   tokenString,
		})
	}
}
