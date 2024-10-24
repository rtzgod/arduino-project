package main

import (
	"bufio"
	"github.com/rtzgod/arduino-project/internal/config"
	"github.com/rtzgod/arduino-project/internal/httpclient"
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {
	// Настройка порта (замени /dev/ttyACM0 на свой порт)
	cfg := config.MustLoad()

	serialCfg := &serial.Config{
		Name:        cfg.Serial.Name,        // Укажи правильный порт (например, /dev/ttyUSB0)
		Baud:        cfg.Serial.Baud,        // Скорость Serial
		ReadTimeout: cfg.Serial.ReadTimeout, // Таймаут чтения
	}
	port, err := serial.OpenPort(serialCfg)
	if err != nil {
		log.Fatal("port")
	}
	defer port.Close()

	client := httpclient.New("http://tty.kz:8080/data")

	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		line := scanner.Text()
		client.Post(line)
		time.Sleep(cfg.Serial.ReadTimeout)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Ошибка чтения: %v", err)
	}
}
