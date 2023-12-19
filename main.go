package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"encoding/json"
	"os/signal"
	"syscall"

	"github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Broker   string
	Port     int
	Topic    string
	Username string
	Password string
	Debug bool
}


var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := string(msg.Topic())
	cmd := string(msg.Payload())
	
	log.Println(topic)
	log.Println(cmd)

	if cmd == "SLEEP"{
		log.Println("GO WINDOWS SLEEP")
		cmd := exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "Sleep")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}	
	}
}

var cH mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var dH mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Println(fmt.Sprintf("Disconnected: %s", err))
}

func main() {
	keepAlive := make(chan os.Signal)
    signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

	// Read config from file
	config := loadConfig("config.json")

	if (config.Debug){
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
	}

	opts := mqtt.NewClientOptions().AddBroker(config.Broker).SetClientID("gotrivial")
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetOnConnectHandler(cH)
	opts.SetConnectionLostHandler(dH)


	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe(config.Topic, 0, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	<-keepAlive
}

func loadConfig(filename string) *Config {
	// Open file 
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Read config from file
	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &config
}
