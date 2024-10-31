package main

import (
	"main/config"
	"main/internal/cmd"

	"github.com/yanun0323/pkg/logs"
)

func main() {
	if err := config.Init("config", "./config", "../config", "../../config"); err != nil {
		logs.Fatalf("init config: %v", err)
	}

	cmd.Run()
}
