package main

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"task/configs"
	"task/internal/floodcontrol"
	"time"
)

func main() {
	yamlFile, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		panic(err)
	}

	var conf configs.Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}

	floodControl := floodcontrol.NewMemcachedFloodControl(conf)

	userID := int64(123)
	for i := 0; i < 10; i++ {
		ctx := context.WithValue(context.Background(), "ttl", int64(conf.WindowSize.Seconds()))
		_, err := floodControl.Check(ctx, userID)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Message")
		}

		time.Sleep(1 * time.Second)
	}

}
