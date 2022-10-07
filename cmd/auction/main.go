package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Kodik77rus/resale-auction/internal/app/auction"
	"github.com/Kodik77rus/resale-auction/internal/pkg/bid_requester"
	"github.com/Kodik77rus/resale-auction/internal/pkg/config"
	"github.com/Kodik77rus/resale-auction/internal/pkg/http_client"
	"github.com/Kodik77rus/resale-auction/internal/pkg/storage/dsp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
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

	initZeroLogger(cnfg)

	httpClient := http_client.InitHttpClient()

	bidRequester := bid_requester.InitBidRequester(cnfg, httpClient)

	dsps := dsp.InitDspStorage(cnfg)

	mux := &http.ServeMux{}

	auction.InitAuction(cnfg, bidRequester, dsps, mux)

	if err := http.ListenAndServe(
		net.JoinHostPort("", cnfg.PORT),
		mux,
	); err != nil {
		return err
	}

	return nil
}

func initZeroLogger(cnfg *config.Config) {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.999"

	lvl, err := zerolog.ParseLevel(cnfg.LOG_LVL)
	if err != nil {
		log.Println("init zero logger : error parse config level", "err: ", err)
		lvl = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(lvl)
}
