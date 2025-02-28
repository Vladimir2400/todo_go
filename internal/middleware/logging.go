package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Логируем начало запроса
		startTime := time.Now()
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		//передаем управление след запросу
		next.ServeHTTP(w, r)

		//Логируем завершение запроса
		duration := time.Since(startTime)
		log.Printf("Response: %s %s (%s)", r.Method, r.URL.Path, duration)
	})
}
