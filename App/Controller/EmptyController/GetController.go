package EmptyController

import (
	"github.com/rnschulenburg/gowrite-api-go/Core/Http"
	"net/http"
)

// Get
// return empty
func Get(w http.ResponseWriter, r *http.Request) {
	Http.JsonResponse(w, nil, "success", "EmptyController.Get", nil)
}
