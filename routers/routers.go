package routers

import (
	"github.com/gorilla/mux"
	"github.com/rnschulenburg/gowrite-api-go/App/Controller/AiController"
	"github.com/rnschulenburg/gowrite-api-go/App/Controller/AuthController"
	"github.com/rnschulenburg/gowrite-api-go/App/Controller/EmptyController"
	"github.com/rnschulenburg/gowrite-api-go/App/Controller/ExportDocumentController"
	"github.com/rnschulenburg/gowrite-api-go/App/Controller/ProjectController"
	"github.com/rnschulenburg/gowrite-api-go/App/Controller/UserController"
	"github.com/rnschulenburg/gowrite-api-go/routers/auth"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/login", Login).Methods("POST")
	router.HandleFunc("/api/v1/auth/refresh", AuthController.Refresh).Methods("POST")
	router.HandleFunc("/api/v1/auth/logout", AuthController.Logout).Methods("POST")

	auth.Handler("/api/v1/user-projects", router, UserController.FetchUserProjects, "GET", "")
	auth.Handler("/api/v1/user-projects/{id:[0-9]+}", router, UserController.FetchUserProjects, "GET", "")

	auth.Handler("/api/v1/export-document/{projectName}", router, EmptyController.Get, "GET", "")
	auth.Handler("/api/v1/project/{projectName}", router, ProjectController.Get, "GET", "")
	auth.Handler("/api/v1/project/{projectName}", router, ProjectController.Save, "POST", "")

	auth.Handler("/api/v1/ask", router, AiController.Ask, "POST", "")
	auth.Handler("/api/v1/export-document/{projectName}", router, ExportDocumentController.Handle, "POST", "")

	auth.Handler("/api/v1/import-document", router, ExportDocumentController.ImportWord, "POST", "")

	return router
}
