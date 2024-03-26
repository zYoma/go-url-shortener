package libs

import (
	"compress/gzip"
	"io"
	"net/http"
)

// compressWriter является обёрткой над http.ResponseWriter, добавляющей поддержку
// сжатия ответов сервера в формате gzip. Это позволяет автоматически сжимать
// отправляемые данные, если статус ответа меньше 300.
type compressWriter struct {
	w  http.ResponseWriter // Исходный writer для отправки ответов.
	zw *gzip.Writer        // Writer для сжатия данных в формате gzip.
}

// NewCompressWriter создаёт и возвращает новый экземпляр compressWriter.
//
// w: http.ResponseWriter, к которому будет применяться сжатие.
//
// Возвращает инициализированный экземпляр compressWriter.
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header возвращает заголовки HTTP ответа.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write сжимает и записывает данные в исходный http.ResponseWriter.
//
// p: данные для записи.
//
// Возвращает количество байт, записанных в исходный writer, и возможную ошибку.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader устанавливает HTTP статус ответа и добавляет заголовок
// "Content-Encoding: gzip", если статус меньше 300.
//
// statusCode: HTTP статус для установки в ответе.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и освобождает все ресурсы. Должен быть вызван
// после завершения записи данных.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser для чтения сжатых gzip данных.
// Позволяет декомпрессировать данные, принимаемые от клиента, прозрачно для сервера.
type compressReader struct {
	r  io.ReadCloser // Исходный reader для чтения сжатых данных.
	zr *gzip.Reader  // Reader для декомпрессии данных.
}

// NewCompressReader создаёт и возвращает новый экземпляр compressReader для чтения
// данных сжатых в формате gzip.
//
// r: io.ReadCloser с сжатыми данными для чтения.
//
// Возвращает инициализированный экземпляр compressReader и возможную ошибку при создании.
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read декомпрессирует и читает сжатые данные из исходного потока.
//
// p: буфер для записи декомпрессированных данных.
//
// Возвращает количество байт, записанных в буфер, и возможную ошибку при чтении.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает и освобождает все ресурсы, связанные с чтением сжатых данных.
// Должен быть вызван после завершения чтения.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
