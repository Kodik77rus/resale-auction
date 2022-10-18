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
) *DspStorage {
	dsps := make([]*models.DspConfig, 0, len(config.DSP_CONNECTION_URLS))

	header := make(http.Header)
	header.Add("Content-Type", "application/json")

	for _, dspUrl := range config.DSP_CONNECTION_URLS {
		dsps = append(dsps, &models.DspConfig{
			Endpoint: &url.URL{
				Host: dspUrl,
				Path: "bid_request",
			},
			RequestHeaders: header,
		})
	}

	return &DspStorage{
		Dsps: dsps,
	}
}
