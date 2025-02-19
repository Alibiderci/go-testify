package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	// здесь нужно создать запрос к сервису
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	cafeList := strings.Split(body, ",")

	// тест на корректное количество кафе
	assert.Len(t, cafeList, totalCount)
}

func TestMainHandlerForCorrectStatusAndBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// тест на корректный статус код и наличие тела ответа
	expectedStatus := 200

	require.NotEmpty(t, responseRecorder.Body.String())
	require.Equal(t, expectedStatus, responseRecorder.Code)
}

func TestMainHandlerWhenCityIsWrong(t *testing.T) {
	// запрос с заведомо неверным городом
	req := httptest.NewRequest("GET", "/cafe?count=4&city=karaganda", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// тест на корректный статус код и тела ответа
	require.Equal(t, 400, responseRecorder.Code)
	require.Equal(t, "wrong city value", responseRecorder.Body.String())
}
