package httpclient

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Структура для сериализации данных
type RequestData struct {
	Data string `json:"data"`
}

type Client struct {
	Url string
}

func New(url string) *Client {
	return &Client{
		Url: url,
	}
}

// Изменение метода Post для правильной сериализации
func (c *Client) Post(data string) {
	// Создаем объект RequestData с полем data
	requestData := RequestData{
		Data: data,
	}

	// Кодируем объект в JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		return
	}

	// Отправляем POST-запрос
	resp, err := http.Post(c.Url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Ошибка отправки:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: сервер вернул статус %d\n", resp.StatusCode)
	} else {
		log.Println("Данные успешно отправлены")
	}
}
