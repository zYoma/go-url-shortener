package handlers

import (
	"net"
	"net/http"

	"github.com/go-chi/render"
	"github.com/zYoma/go-url-shortener/internal/models"
)

// GetStats обрабатывает запрос /api/internal/stats.
// Проверяет, принадлежит ли IP-адрес клиента доверенной подсети,
// и возвращает статистику службы в формате JSON.
// Если IP-адрес клиента не принадлежит доверенной подсети,
// возвращает статус 403 Forbidden.
// Если происходит ошибка при получении статистики из базы данных,
// возвращает статус 500 Internal Server Error.
func (h *HandlerService) GetStats(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if h.cfg.TrustedSubnet == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Получение IP-адреса клиента из заголовка X-Real-IP
	clientIP := req.Header.Get("X-Real-IP")
	if clientIP == "" {
		// Если заголовок X-Real-IP отсутствует, получаем IP-адрес из RemoteAddr
		var err error
		clientIP, _, err = net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// Проверка, принадлежит ли IP-адрес клиента доверенной подсети
	_, trustedIPNet, err := net.ParseCIDR(h.cfg.TrustedSubnet)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	clientAddr := net.ParseIP(clientIP)
	if !trustedIPNet.Contains(clientAddr) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	response, err := h.provider.GetServiceStats(ctx)
	if err != nil {
		render.JSON(w, req, models.Error("failed get stats from db"))
		return
	}

	render.JSON(w, req, response)
}
