package config

import (
	"flag"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

type Config struct {
	PORT                string
	DSP_CONNECTION_URLS []string
}

func InitConfig() (*Config, error) {
	port := flag.String("p", "8080", "server port")
	dspUrls := flag.String("d", "", "dsp url")

	flag.Parse()

	if ok := govalidator.IsPort(*port); !ok {
		return nil, errors.Errorf("%s port flag variable is not int convertible to int", port)
	}

	dspSlice := strings.Split(*dspUrls, ",")

	for _, dsp := range dspSlice {
		if ok := govalidator.IsDialString(dsp); !ok {
			return nil, errors.Errorf("%s dsp flag variable is not convertible to net addr", dsp)
		}
	}

	return &Config{
		PORT:                *port,
		DSP_CONNECTION_URLS: dspSlice,
	}, nil
}
