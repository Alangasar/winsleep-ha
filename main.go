package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"encoding/json"
	"flag"

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
	// topic := string(msg.Topic())
	cmd := string(msg.Payload())
	
	if cmd == "SLEEP"{
		log.Println("GO WINDOWS SLEEP")
		cmd := exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}	
	}
}


func main() {
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

	num := flag.Int("num", 1, "The number of messages to publish or subscribe (default 1)")

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe(config.Topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	receiveCount := 0
	choke := make(chan [2]string)
	
	for receiveCount < *num {
		incoming := <-choke
		log.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		receiveCount++
	}
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
