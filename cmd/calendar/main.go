package main

import (
	"fmt"
	"github.com/Manabreaker/Calendar/configs"
	"github.com/Manabreaker/Calendar/internal/app/apiserver"
	"log"
)

const (
	configPath = "configs/config.yaml"
)

func main() {
	config, err := configs.NewConfig(configPath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Ошибка при обработке файла `%s`: %e", configPath, err))
	}
	serverConfig := &config.Server
	apiServer := apiserver.NewServer()
	log.Println("APIGateway server started successfully at " + serverConfig.Host + ":" + serverConfig.Port)
	if err := apiServer.Start(serverConfig); err != nil {
		log.Fatal(fmt.Sprintf("Ошибка при запуске сервера: %e", err))
	}
}
