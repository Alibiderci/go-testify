package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	// здесь нужно создать запрос к сервису
	req := httptest.NewRequest("GET", "cafe?count=10&city=moscow", nil)
	reqWithWrongCity := httptest.NewRequest("GET", "cafe?count=10&city=karaganda", nil) // здесь нужно создать запрос к сервису

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// здесь нужно добавить необходимые проверки
	// тест на корректный статус код и наличие тела ответа
	expectedStatus := 200

	require.NotEmpty(t, responseRecorder.Body.String())
	require.Equal(t, expectedStatus, responseRecorder.Code)

	// тест на корректное количество кафе
	assert.Len(t, responseRecorder.Body.String(), totalCount)

	// проверяем город
	wrongCityRR := httptest.NewRecorder()
	wrongCityHandler := http.HandlerFunc(mainHandle)
	wrongCityHandler.ServeHTTP(wrongCityRR, reqWithWrongCity)
	require.Equal(t, 400, wrongCityRR.Code)
	require.Equal(t, "wrong city value", wrongCityRR.Body.String())

}
