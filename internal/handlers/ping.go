package handlers

import (
	"net/http"
)

func (h *HandlerService) Ping(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	err := h.provider.Ping(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
