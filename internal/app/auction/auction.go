package auction

import (
	"io/ioutil"
	"math"
	"net/http"
	"sort"

	"github.com/Kodik77rus/resale-auction/internal/pkg/bid_requester"
	"github.com/Kodik77rus/resale-auction/internal/pkg/config"
	"github.com/Kodik77rus/resale-auction/internal/pkg/models"
	"github.com/Kodik77rus/resale-auction/internal/pkg/utils"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Auction struct{}

func InitAuction(
	config *config.Config,
	bidRequester *bid_requester.BidRequester,
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
				Msg("invalid request body EMPTY_FIELD || WRONG_SCHEMA")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var sspRequestDto models.SspRequest

		if err := utils.JsonUnmarshal(body, &sspRequestDto); err != nil {
			log.Error().
				Err(err).
				Int("request status code", http.StatusInternalServerError).
				Msg("can't unmarshal request body EMPTY_FIELD || WRONG_SCHEMA")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ok, err := govalidator.ValidateStruct(sspRequestDto)
		if err != nil {
			log.Error().
				Err(err).
				Int("request status code", http.StatusBadRequest).
				Interface("ssp request", sspRequestDto).
				Msg("failed validate ssp request EMPTY_FIELD || WRONG_SCHEMA")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !ok {
			log.Error().
				Int("request status code", http.StatusBadRequest).
				Interface("ssp request", sspRequestDto).
				Msg("invalid ssp request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// log.Info().Interface("ssp request", sspRequestDto).Msg("get ssp request")

		sspTilesLen := len(sspRequestDto.Tiles)
		auctionLotsMap := make(map[uint][]*models.AuctionBid, sspTilesLen)

		for _, sspTile := range sspRequestDto.Tiles {
			auctionLotsMap[sspTile.Id] = make([]*models.AuctionBid, 0, len(config.DSP_CONNECTION_URLS))
		}

		var dspBidRequstDto models.DspBidRequest

		dspBidRequstDto.Id = sspRequestDto.Id
		dspBidRequstDto.Context = sspRequestDto.Context
		dspBidRequstDto.Imp = make([]*models.RequestDspImp, 0, sspTilesLen)

		for _, sspImp := range sspRequestDto.Tiles {
			dspBidRequstDto.Imp = append(
				dspBidRequstDto.Imp,
				&models.RequestDspImp{
					Id:       sspImp.Id,
					Minwidth: sspImp.Width,
					Minheight: uint(
						math.Floor(
							float64(sspImp.Width) * float64(sspImp.Ratio),
						),
					),
				},
			)
		}

		dspsResponsesInfo, err := bidRequester.Send(config.DSP_CONNECTION_URLS, &dspBidRequstDto)
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

		log.Info().Interface("dsp resps", dspsResponsesInfo)

		if err := calculateAuctionParams(dspsResponsesInfo, auctionLotsMap); err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		calculateWiners(auctionLotsMap)

		// log.Info().Interface("auction result", auctionLotsMap).Msg("auction result")

		var sspResponseDto models.SspResponse

		sspResponseDto.Id = sspRequestDto.Id
		sspResponseDto.Imp = make([]models.SspImp, 0, len(auctionLotsMap))

		for _, sspTiles := range sspRequestDto.Tiles {
			winerImp, ok := auctionLotsMap[sspTiles.Id]
			if !ok {
				log.Warn().
					Int("request status code", http.StatusNoContent)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			s := models.SspImp{
				Id:     winerImp[0].Imp.Id,
				Width:  winerImp[0].Imp.Width,
				Height: winerImp[0].Imp.Height,
				Title:  winerImp[0].Imp.Title,
				Url:    winerImp[0].Imp.Url,
			}
			log.Info().Interface("winner", s)
			sspResponseDto.Imp = append(
				sspResponseDto.Imp,
				s,
			)
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

	// logMiddleware := func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(
	// 		func(w http.ResponseWriter, r *http.Request) {
	// 			log.Info().Msgf(
	// 				"%s %s from %v",
	// 				r.Method,
	// 				r.URL.Path,
	// 				r.RemoteAddr,
	// 			)
	// 			start := time.Now()
	// 			next.ServeHTTP(w, r)
	// 			log.Info().Msgf(
	// 				"%s %s from %v duration: %v",
	// 				r.Method,
	// 				r.URL.Path,
	// 				r.RemoteAddr,
	// 				time.Since(start),
	// 			)
	// 		})
	// }

	// timeOutMiddleware := func(next http.Handler, duration time.Duration) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		ctx, cancel := context.WithTimeout(r.Context(), duration)
	// 		defer cancel()

	// 		r = r.WithContext(ctx)

	// 		processDone := make(chan struct{})
	// 		go func() {
	// 			next.ServeHTTP(w, r)
	// 			processDone <- struct{}{}
	// 		}()

	// 		select {
	// 		case <-ctx.Done():
	// 			log.Warn().Msg("ssp timeout!")
	// 			w.WriteHeader(http.StatusNoContent)
	// 		case <-processDone:
	// 		}
	// 	})
	// }

	mu.Handle(
		"/placements/request",
		// timeOutMiddleware(
		// logMiddleware(
		http.HandlerFunc(handler),
		// ),
		// config.SSP_TIMEOUT,
	)
	// )
}

func calculateAuctionParams(
	dspsResp []*models.DspBidRequestInfo,
	sspLots map[uint][]*models.AuctionBid,
) error {
	for _, dsp := range dspsResp {
		for _, dspImp := range dsp.DspBidResponse.Imp {
			bids, ok := sspLots[dspImp.Id]
			if !ok {
				log.Warn().
					Interface("dsp", dsp.DspInfo).
					Interface("imp", dspImp).
					Msg("unknown dsp imp id")
				return errors.New("unknown dsp imp id")
			}
			sspLots[dspImp.Id] = append(bids, &models.AuctionBid{
				Dsp: dsp.DspInfo,
				Imp: dspImp,
			})
		}
	}
	return nil
}

func calculateWiners(sspLots map[uint][]*models.AuctionBid) {
	log.Info().Interface("row", sspLots)
	for key, val := range sspLots {
		sort.SliceStable(val, func(i, j int) bool {
			return val[i].Imp.Price > val[j].Imp.Price
		})
		sspLots[key] = val
	}
}
