package storage

import "fmt"

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

// Vwap represents a queue of DataPoints and their VWAPs.
//DataPoints  is fast circular fifo data structure (aka., queue) with a specific limit.
type Vwap interface {
	// Push pushes an item onto the queue and calculates the new VWAP.
	//When Limit is reached, will delete  the first one.
	Push(d Point)

	// Size returns the length of the dara points queue.
	Size() uint

	// GetDataPoints returns the data point items.
	GetDataPoints() any

	// GetVwap returns the VWAP for a  trading pair.
	GetVwap(tradingPair string) float64

	//GetVwaps returns the VWAPs for trading pairs
	GetVwaps() map[string]float64
}

func Format[M ~map[K]V, K comparable, V any](m M) (result []string) {
	result = make([]string, 0, len(m))
	for k, v := range m {
		result = append(result, fmt.Sprintf("%v (%v)", k, v))
	}
	return
}
