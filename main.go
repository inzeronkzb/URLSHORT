package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urls = make(map[string]string) // Cоздаем мапу, которая будет хранить сокращенный url как ключ и полный как значение

func main() {
	// Обработка маршрутов
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/short/", handleRedirect)
	// Запуск сервера
	fmt.Println("URL Shortener is running on :3030")
	http.ListenAndServe(":3030", nil)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	// Проверяем является ли метод HTTP запроса POST
	if r.Method == http.MethodPost { 
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	// HTML форма, если метод запроса не POST, чтобы пользователь мог ввести url для сокращения
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<form method="post" action="/shorten">
				<input type="url" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form>
		</body>
		</html>
	`)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Получение оригинального url
	originalURL := r.FormValue("url")
	if originalURL == "" {
  		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
  		return
	}

	// Генерация укороченного ключа и сохранение в мапу
	shortKey := generateShortKey()
	urls[shortKey] = originalURL

	// Здесь создается строка для укороченного URL, которая включает в себя хост и путь к ресурсам
	shortenedURL := fmt.Sprintf("http://localhost:3030/short/%s", shortKey)

	// Отправка ответа в виде HTML
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<p>Original URL: `, originalURL, `</p>
			<p>Shortened URL: <a href="`, shortenedURL, `">`, shortenedURL, `</a></p>
		</body>
		</html>
	`)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Извлечение укороченного ключа из URL
	shortKey := strings.TrimPrefix(r.URL.Path, "/short/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Поиск оригинального URL в мапе
	originalURL, found := urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	// Перенаправление на оригинальный URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())
	// Генерация укороченного ключа
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}