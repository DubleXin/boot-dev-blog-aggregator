package main

import (
	"fmt"
	"github.com/DubleXin/aggregator/internal/config"
)

func main() {
	fmt.Println("Reading cofig...")

	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error, ocured:\n\t", err.Error())
	}
	fmt.Println("Current config:")
	fmt.Println(cfg.DBUrl)
	fmt.Println(cfg.CurrentUser)
	user := "Rundom_user_nAme1234"
	fmt.Println("Writing", user, "to config...")
	err = cfg.SetUser(user)
	if err != nil {
		fmt.Println("Error, ocured:\n\t", err.Error())
	}

	fmt.Println("Reading cofig...")

	cfg, err = config.ReadConfig()
	if err != nil {
		fmt.Println("Error, ocured:\n\t", err.Error())
	}
	fmt.Println("Current config:")
	fmt.Println(cfg.DBUrl)
	fmt.Println(cfg.CurrentUser)
}
