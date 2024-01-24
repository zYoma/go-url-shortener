package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func (h *HandlerService) GetUserURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Printf("GET WITH USER: %s", userID)
	response, err := h.provider.GetUserURLs(ctx, h.cfg.BaseShortURL, userID)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	render.JSON(w, req, response)
}
