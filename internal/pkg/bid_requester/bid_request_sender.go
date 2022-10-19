package bid_requester

import (
	"net/http"
	"sync"
	"time"

	"github.com/Kodik77rus/resale-auction/internal/pkg/config"
	"github.com/Kodik77rus/resale-auction/internal/pkg/http_client"
	"github.com/Kodik77rus/resale-auction/internal/pkg/models"
	"github.com/Kodik77rus/resale-auction/internal/pkg/utils"
	"github.com/asaskevich/govalidator"
	"github.com/rs/zerolog/log"
)

type BidRequester struct {
	dspTimeout time.Duration
	httpClient *http_client.HttpClient
}

func InitBidRequester(
	config *config.Config,
	httpClient *http_client.HttpClient,
) *BidRequester {
	return &BidRequester{
		dspTimeout: config.DSP_TIMEOUT,
		httpClient: httpClient,
	}
}

func (b *BidRequester) Send(
	dsps []*string,
	bidRequest *models.DspBidRequest,
) ([]*models.DspBidRequestInfo, error) {
	body, err := utils.JsonMarshal(bidRequest)
	if err != nil {
		log.Error().
			Err(err).
			Interface("bid request", bidRequest).
			Msg("failed marshal bid request")
		return nil, err
	}

	dspCount := len(dsps)
	DspBidRequestInfo := make(chan *models.DspBidRequestInfo, dspCount)

	wg := sync.WaitGroup{}
	wg.Add(dspCount)

	for _, dsp := range dsps {
		go func(dsp *string) {
			defer wg.Done()

			var dspRespInfo models.DspBidRequestInfo

			dspRespInfo.DspInfo = *dsp
			dspRespInfo.DspBidRequest = bidRequest

			log.Info().Str("url", *dsp).Msg("start request to dsp")

			satusCode, respBody, err := b.httpClient.POST(
				dsp,
				body,
				b.dspTimeout,
			)
			if err != nil {
				log.Error().
					Err(err).
					Interface("bid request info", dspRespInfo).
					Msg("failed response to dsp")
				return
			}
			if satusCode != http.StatusOK {
				log.Warn().
					Int("status code", satusCode).
					Str("dsp", *dsp).
					Interface("bid request info", dspRespInfo).
					Msg("response to dsp not ok")
				return
			}

			var dspBidResponseDto models.DspBidResponse

			if ok := utils.IsValidJson(respBody); !ok {
				log.Error().
					Interface("bid request", dspRespInfo).
					Msg("invalid dsp response body")
				return
			}

			if err := utils.JsonUnmarshal(respBody, &dspBidResponseDto); err != nil {
				log.Error().
					Err(err).
					Str("dsp", *dsp).
					Interface("bid request", bidRequest).
					Msg("failed unmarshal bid response")
				return
			}

			dspRespInfo.DspBidResponse = &dspBidResponseDto

			ok, err := govalidator.ValidateStruct(dspBidResponseDto)
			if err != nil {
				log.Error().
					Err(err).
					Interface("dsp info", dspRespInfo).
					Msg("failed validate dsp response EMPTY_FIELD || WRONG_SCHEMA")
				return
			}
			if !ok {
				log.Warn().
					Interface("dsp info", dspRespInfo).
					Msg("invalid dsp response")
				return
			}
			if bidRequest.Id != dspBidResponseDto.Id {
				log.Warn().
					Interface("dsp info", dspRespInfo).
					Msg("invalid dsp response")
				return
			}

			DspBidRequestInfo <- &dspRespInfo
		}(dsp)
	}

	go func() {
		wg.Wait()
		close(DspBidRequestInfo)
	}()

	respSlice := make([]*models.DspBidRequestInfo, 0, dspCount)

	for resp := range DspBidRequestInfo {
		respSlice = append(respSlice, resp)
	}

	return respSlice, nil
}
