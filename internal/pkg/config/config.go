package config

import (
	"flag"
	"net/url"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

type Config struct {
	PORT                string
	DSP_CONNECTION_URLS []*string
	LOG_LVL             string
	SSP_TIMEOUT         time.Duration
	DSP_TIMEOUT         time.Duration
}

func InitConfig() (*Config, error) {
	port := flag.String("p", "8080", "server port")
	dspUrls := flag.String("d", "", "dsp url")
	logLVL := flag.String("l", "debug", "log lvl")
	dspTimeout := flag.Duration("dt", 200, "dsp request millisecond timeout")
	sspTimeout := flag.Duration("st", 250, "ssp request millisecond timeout")

	flag.Parse()

	if ok := govalidator.IsPort(*port); !ok {
		return nil, errors.Errorf("%s port flag variable is not convertible to int", port)
	}

	dspSlice := strings.Split(*dspUrls, ",")
	dspConnUrls := make([]*string, 0, len(dspSlice))

	for _, dsp := range dspSlice {
		if ok := govalidator.IsDialString(dsp); !ok {
			return nil, errors.Errorf("%s dsp flag variable is not convertible to slice net addr", dsp)
		}
		dspUrl := url.URL{
			Scheme: "http",
			Host:   dsp,
			Path:   "bid_request",
		}
		u := dspUrl.String()
		dspConnUrls = append(dspConnUrls, &u)
	}

	return &Config{
		PORT:                *port,
		DSP_CONNECTION_URLS: dspConnUrls,
		LOG_LVL:             *logLVL,
		SSP_TIMEOUT:         *sspTimeout * time.Millisecond,
		DSP_TIMEOUT:         *dspTimeout * time.Millisecond,
	}, nil
}
