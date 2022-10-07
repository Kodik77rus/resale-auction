package auction

import (
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/Kodik77rus/resale-auction/internal/pkg/bid_requester"
	"github.com/Kodik77rus/resale-auction/internal/pkg/config"
	"github.com/Kodik77rus/resale-auction/internal/pkg/models"
	"github.com/Kodik77rus/resale-auction/internal/pkg/storage/dsp"
	"github.com/Kodik77rus/resale-auction/internal/pkg/utils"
	"github.com/asaskevich/govalidator"
	"github.com/rs/zerolog/log"
)

type Auction struct{}

func InitAuction(
	config *config.Config,
	bidRequester *bid_requester.BidRequester,
	dspStorage *dsp.DspStorage,
	mu *http.ServeMux,
) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Error().
				Int("request status code", http.StatusMethodNotAllowed).
				Msg("http method not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().
				Err(err).
				Int("request status code", http.StatusInternalServerError).
				Msg("read body err")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if ok := utils.IsValidJson(body); !ok {
			log.Error().
				Int("request status code", http.StatusBadRequest).
				Msg("invalid request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var sspRequestDto models.SspRequest

		if err := utils.JsonUnmarshal(body, &sspRequestDto); err != nil {
			log.Error().
				Err(err).
				Int("request status code", http.StatusInternalServerError).
				Msg("can't unmarshal request body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ok, err := govalidator.ValidateStruct(sspRequestDto)
		if err != nil {
			log.Error().
				Err(err).
				Int("request status code", http.StatusBadRequest).
				Msg("can't validate ssp request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !ok {
			log.Error().
				Int("request status code", http.StatusBadRequest).
				Msg("invalid ssp request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		auctionLotsMap := make(map[uint][]*models.AuctionBid, len(sspRequestDto.Tiles))

		for _, sspTile := range sspRequestDto.Tiles {
			auctionLotsMap[sspTile.Id] = make([]*models.AuctionBid, 0, len(dspStorage.Dsps))
		}

		var dspBidRequstDto models.DspBidRequest

		dspBidRequstDto.Id = sspRequestDto.Id
		dspBidRequstDto.Context = sspRequestDto.Context

		dspImps := make([]models.RequestDspImp, 0, len(sspRequestDto.Tiles))

		for _, imp := range sspRequestDto.Tiles {
			dspImps = append(
				dspImps,
				models.RequestDspImp{
					Id:       imp.Id,
					Minwidth: imp.Width,
					Minheight: uint(
						math.Floor(
							float64(imp.Width) * float64(imp.Ratio),
						),
					),
				},
			)
		}

		dspsResponsesInfo, err := bidRequester.Send(dspStorage.Dsps, dspBidRequstDto)
		if err != nil {
			log.Error().
				Err(err).
				Int("request status code", http.StatusNoContent).
				Msg("failed dsps requests")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		dspResponsesCount := len(dspsResponsesInfo)

		if dspResponsesCount == 0 {
			log.Warn().
				Int("request status code", http.StatusNoContent).
				Msg("no responses from dsps")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		validDspsResps := make([]*models.DspBidRequestInfo, 0, dspResponsesCount)

		for _, dspResp := range dspsResponsesInfo {
			ok, err := govalidator.ValidateStruct(dspResp.DspBidResponse)
			if err != nil {
				log.Error().
					Err(err).
					Interface("dsp info", *dspResp).
					Msg("failed validate dsp response")
				continue
			}
			if !ok {
				log.Warn().
					Interface("dsp info", *dspResp).
					Msg("invalid dsp response")
				continue
			}
			if dspResp.DspBidResponse.Id != sspRequestDto.Id {
				log.Warn().
					Interface("dsp response id not equal ssp request id", *dspResp).
					Msg("invalid dsp response")
				continue
			}
			validDspsResps = append(validDspsResps, dspResp)
		}

		calculateAuctionParams(validDspsResps, auctionLotsMap)
		calculateWiners(auctionLotsMap)

		log.Info().Interface("auction result", auctionLotsMap)

		var sspResponseDto models.SspResponse

		sspResponseDto.Id = sspRequestDto.Id
		sspResponseDto.Imp = make([]models.SspImp, 0, len(auctionLotsMap))

		for _, sspTiles := range sspRequestDto.Tiles {
			winerImp, _ := auctionLotsMap[sspTiles.Id]
			sspResponseDto.Imp = append(sspResponseDto.Imp, models.SspImp{
				Id:     winerImp[0].Imp.Id,
				Width:  winerImp[0].Imp.Width,
				Height: winerImp[0].Imp.Height,
				Title:  winerImp[0].Imp.Title,
				Url:    winerImp[0].Imp.Url,
			})
		}

		resp, err := utils.JsonMarshal(sspResponseDto)
		if err != nil {
			log.Error().Err(err).Msg("failed marshal auction result")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}

	logMiddleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Info().Msgf("%s %s from %v", r.Method, r.URL.Path, r.RemoteAddr)
			start := time.Now()
			next.ServeHTTP(w, r)
			responseTime := time.Since(start)
			if responseTime.Milliseconds() >= config.SSP_TIMEOUT.Microseconds() {
				log.Warn().Msgf("%s %s from %v duration: %v to long", r.Method, r.URL.Path, r.RemoteAddr, responseTime)
			}
			log.Info().Msgf("%s %s from %v duration: %v", r.Method, r.URL.Path, r.RemoteAddr, responseTime)
		}
		return http.HandlerFunc(fn)
	}

	mu.Handle("/placements/request", logMiddleware(http.HandlerFunc(handler)))
}

func calculateAuctionParams(
	dspsResp []*models.DspBidRequestInfo,
	sspLots map[uint][]*models.AuctionBid,
) {
	for _, dsp := range dspsResp {
		for _, dspImp := range dsp.DspBidResponse.Imp {
			bids, ok := sspLots[dspImp.Id]
			if !ok {
				log.Warn().
					Interface("dsp", dsp.DspInfo).
					Interface("imp", dspImp).
					Msg("unknown dsp imp id")
				continue
			}
			sspLots[dspImp.Id] = append(bids, &models.AuctionBid{
				Dsp: dsp.DspInfo,
				Imp: dspImp,
			})
		}
	}
}

func calculateWiners(sspLots map[uint][]*models.AuctionBid) {
	for key, val := range sspLots {
		sort.SliceStable(val, func(i, j int) bool {
			return val[i].Imp.Price > val[j].Imp.Price
		})
		sspLots[key] = val
	}
}
