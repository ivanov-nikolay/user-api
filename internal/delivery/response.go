package delivery

import (
	"net/http"

	"go.uber.org/zap"
)

func writeResponse(logger *zap.SugaredLogger, w http.ResponseWriter, dataJSON []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(statusCode)
	_, err := w.Write(dataJSON)
	if err != nil {
		logger.Errorf("error in writing response body: %s", err)
	}
}
