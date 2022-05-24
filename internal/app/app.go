package app

import (
	"context"
	"fmt"
	"github.com/reactivejson/vwap-engine/api/models"
	"github.com/reactivejson/vwap-engine/internal/storage"
	"strconv"
	"time"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

//Run the App context, subscribe to the ws, and initiate the storage and calculation for the trading pairs.
//It is resilient tolerant. It will gracefully shut down and can receive an interrupt signal and safely to close the connexion.
func (s *Context) Run(ctx context.Context) (err error) {
	receiver := make(chan *models.CoinbaseResponse)

	err = s.wsReceiver.Subscribe(s.cfg.TradingPairs)
	if err != nil {
		return fmt.Errorf("failed to subscribe err: %w", err)
	}

	s.wsReceiver.Read(ctx, receiver)

	for response := range receiver {
		//Skip non valid responses
		if response.Price == "" {
			continue
		}

		if dataPoint, err := parseData(response); err != nil {
			return err
		} else {
			s.queue.Push(dataPoint)
		}

		// Log VWAPs of trading pairs to stdout.
		fmt.Println(time.Now().Format(time.UnixDate))
		fmt.Println("VWAPs:", s.queue)
	}
	return
}

//Convert JSON response (models.CoinbaseResponse) to a storage.Point
func parseData(response *models.CoinbaseResponse) (storage.Point, error) {
	price, err := strconv.ParseFloat(response.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing price %s: %w", response.Price, err)
	}

	quantity, err := strconv.ParseFloat(response.Size, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing quantity %s: %w", response.Size, err)
	}

	return storage.NewPoint(
		price,
		quantity,
		response.ProductID,
	), nil

}
