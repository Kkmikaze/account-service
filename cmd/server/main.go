package main

import (
	"account-service/cmd"
	"account-service/config"
)

func init() {
	config.InitConfig()
}

func main() {
	cmd.Execute()
}
