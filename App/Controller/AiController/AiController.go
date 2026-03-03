package AiController

import (
	"github.com/rnschulenburg/gowrite-api-go/App/Services/AiService"
	"github.com/rnschulenburg/gowrite-api-go/Core/Http"
	"net/http"
)

// Ask
// return a tradingPlan
func Ask(w http.ResponseWriter, r *http.Request) {

	resp, err := AiService.GoChat(r)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	Http.JsonResponse(w, err, "success", "AiController.Ask", resp)
}
