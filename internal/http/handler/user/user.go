package user

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Amannigam1820/student-api-go/internal/storage"
	"github.com/Amannigam1820/student-api-go/internal/types"
	"github.com/Amannigam1820/student-api-go/internal/utils/response"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credetials types.User

		if err := json.NewDecoder(r.Body).Decode(&credetials); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Invalid Request")))
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
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})

	}
}
