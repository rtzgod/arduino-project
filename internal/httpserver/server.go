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
	Data string `json:"data"`
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
	// Настройка заголовков CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")              // Разрешить все домены
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Разрешить POST и OPTIONS
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Разрешить заголовок Content-Type

	if r.Method == http.MethodOptions {
		// Ответ на предварительный запрос
		w.WriteHeader(http.StatusNoContent)
		return
	}

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
	// Настройка заголовков CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить все домены
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Разрешить GET и OPTIONS
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Разрешить заголовок Content-Type

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
