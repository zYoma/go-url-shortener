package handlers

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"go.uber.org/zap"
)

// DeleteShortListURL обрабатывает HTTP-запросы для удаления списка коротких URL, принадлежащих пользователю.
// В теле запроса ожидается JSON-массив, содержащий строки с короткими URL для удаления. Метод аутентифицирует пользователя,
// декодирует тело запроса и, в случае успешной аутентификации и декодирования, ставит задачу на удаление URL в фоновом режиме.
// В ответ клиенту отправляется HTTP-статус 202 (Accepted), указывающий на то, что запрос принят к обработке,
// но процесс удаления может быть выполнен позже.
//
// Для обработки удаления используется фоновый механизм: задачи на удаление помещаются в канал, из которого они
// будут извлечены и обработаны в другой горутине. Это позволяет методу быстро отвечать клиенту и осуществлять
// фактическое удаление асинхронно.
//
// В случае ошибок аутентификации, декодирования тела запроса или если тело запроса оказывается пустым,
// клиенту возвращается соответствующий HTTP-статус ошибки и описание ошибки в формате JSON.
//
// Параметры:
//
//	w http.ResponseWriter: интерфейс для отправки HTTP ответов.
//	req *http.Request: структура, представляющая HTTP запрос.
func (h *HandlerService) DeleteShortListURL(w http.ResponseWriter, req *http.Request) {
	userID, err := getUserFromRequest(req.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var listURL []string

	err = render.DecodeJSON(req.Body, &listURL)

	if errors.Is(err, io.EOF) {
		logger.Log.Error("request body is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, req, models.Error("empty request"))
		return
	}
	if err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, req, models.Error("failed to decode request"))
		return
	}

	go func() {
		select {
		case h.delChan <- models.UserListURLForDelete{UserID: userID, URLS: listURL}:
			// Успешно отправлено
		case <-time.After(time.Second * 5): // Например, 5 секунд таймаут
			// Обработка таймаута, возможно запись в лог
			logger.Log.Error("timeout when sending to delChan")
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}
