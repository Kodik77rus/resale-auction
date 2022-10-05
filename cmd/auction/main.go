package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Kodik77rus/resale-auction/internal/pkg/config"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Println("main : shutting down", "err: ", err)
		os.Exit(1)
	}
}

func run() error {
	cnfg, err := config.InitConfig()
	if err != nil {
		return errors.Errorf("failed parse config err: %v", err)
	}

	fmt.Print(cnfg)
	return nil
}
