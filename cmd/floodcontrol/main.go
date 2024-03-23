package main

import (
	"context"
	"fmt"
	"task/configs"
	"task/internal/floodcontrol"
	"time"
)

func main() {
	conf, err := configs.LoadConfig("configs/config.yaml")
	if err != nil {
		panic(err)
	}

	floodControl := floodcontrol.NewMemcachedFloodControl(conf)
	userID := int64(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 10; i++ {
		_, err := floodControl.Check(ctx, userID)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Message from user %d\n", userID)
		}

		time.Sleep(1 * time.Second)
	}
}
