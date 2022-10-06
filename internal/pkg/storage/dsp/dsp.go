package dsp

import (
	"net/http"
	"net/url"

	"github.com/Kodik77rus/resale-auction/internal/pkg/config"
	"github.com/Kodik77rus/resale-auction/internal/pkg/models"
)

type DspStorage struct {
	Dsps []*models.DspConfig
}

func InitDspStorage(
	config *config.Config,
) (*DspStorage, error) {
	dsps := make([]*models.DspConfig, len(config.DSP_CONNECTION_URLS))

	header := make(http.Header)
	header.Add("Content-Type", "application/json")

	for _, dspUrl := range config.DSP_CONNECTION_URLS {
		u, err := url.Parse(dspUrl)
		if err != nil {
			return nil, err
		}

		u.Path = "bid_request"

		dsps = append(dsps, &models.DspConfig{
			Endpoint:       u,
			RequestHeaders: header,
		})
	}

	return &DspStorage{
		Dsps: dsps,
	}, nil
}
