package handlers

import (
	"net/http"
)

// Ping проверяет доступность и работоспособность хранилища данных,
// используя провайдер хранилища. Метод предназначен для использования в качестве
// простого эндпоинта проверки состояния сервиса (health check).
//
// В случае успешной проверки возвращает HTTP-статус 200 (OK) и тело "OK",
// указывая на то, что сервис и его зависимости функционируют нормально.
// Если проверка не удалась, клиенту возвращается HTTP-статус 500 (Internal Server Error),
// сигнализируя о возникших проблемах с доступностью или работоспособностью хранилища данных.
//
// Параметры:
//
//	w http.ResponseWriter: интерфейс для отправки HTTP ответов.
//	req *http.Request: структура, представляющая HTTP запрос.
func (h *HandlerService) Ping(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	err := h.provider.Ping(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
