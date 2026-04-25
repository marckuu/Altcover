package tools

import (
	"book-cover-service/dto"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	errorResponse := dto.NewErrorResponse(err.Error())
	if convertErr := json.NewEncoder(w).Encode(errorResponse); convertErr != nil {
		fmt.Println("Ошибка при записи ответа с информацией об ошибке")
	}
}
