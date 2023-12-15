package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Broker   string
	Port     int
	Topic    string
	Username string
	Password string
}

func main() {
	// Чтение конфигурации из файла
	config := loadConfig("config.json")

	// Подключение к MQTT брокеру
	client := MQTT.NewClient(MQTT.NewClientOptions().
		SetClientID(fmt.Sprintf("client-%s", config.Topic)).
		AddBroker(config.Broker).
		SetUsername(config.Username).
		SetPassword(config.Password))
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Подписываемся на топик
	err = client.Subscribe(config.Topic, 0, func(c MQTT.Client, msg MQTT.Message) {
		// Если сообщение содержит значение "sleep", то отправляем компьютер в сон
		if string(msg.Payload()) == "sleep" {
			fmt.Println("Отправляем компьютер в сон")
			cmd := exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0")
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig(filename string) *Config {
	// Открытие файла
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Чтение конфигурации из файла
	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &config
}
