package UserController

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rnschulenburg/gowrite-api-go/App/Repositories/UserRepository"
	"github.com/rnschulenburg/gowrite-api-go/routers/auth"
	"net/http"
	"strconv"
)

func FetchUserProjects(w http.ResponseWriter, r *http.Request) {

	authUserId, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userId := authUserId

	userId64, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil && userId64 != 0 {
		userId = int(userId64)
	} else {
		userId = authUserId
	}

	projects, err := UserRepository.FetchProjects(authUserId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
