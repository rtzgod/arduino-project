package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Data struct {
	data string
}

var (
	mu   sync.Mutex
	data Data
)

func RunServer() {
	http.HandleFunc("/data", handlePostData)  // Приём данных от Arduino
	http.HandleFunc("/roblox", handleGetData) // Отдача данных для Roblox

	fmt.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

// Обработчик для получения данных от Arduino (POST)
func handlePostData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST запросы", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newData Data
	if err := json.Unmarshal(body, &newData); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	data = newData
	mu.Unlock()

	fmt.Fprintf(w, "Данные получены: %+v\n", newData)
}

// Обработчик для отдачи данных в Roblox (GET)
func handleGetData(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
