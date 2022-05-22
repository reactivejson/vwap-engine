package storage

import (
	"fmt"
	"sync"
)

// VwapQueue represents a queue of DataPoints.
//Every time a new data point is added to the queue and saved for each trading pair, the VWAP computation is updated accordingly.
// For performance, and to avoid exponential complexity, the computation is cached for VWAP, CumulativeQuantity,
//and CumulativePriceQuantity for existing data points and updated with new entries.
type vwapQueue struct {
	mu sync.Mutex
	//DataPoints  is fast circular fifo data structure (aka., queue) with a specific limit.
	DataPoints              []DataPoint
	CumulativePriceQuantity map[string]float64 // Equation = Sum(Price*Quantity), Sum of Price * Quantity for each TradingPair for each window
	CumulativeQuantity      map[string]float64 // Equation = Sum(Quantity), Sum of Quantities for every TradingPair for each window
	VWAP                    map[string]float64 //Equation: VWAP = Sum(Price*Quantity) / Sum(Quantity) Volume Weighted Average Price is calculated for every TradingPair for each window
	// Limit sets the number of data points used to calculate the VWAP.
	Limit uint
}

//NewVwapQueue  creates a new VWAP queue and initializes all fields needed to make the VWAP Queue.
func NewVwapQueue(dataPoint []DataPoint, maxSize uint) (Vwap, error) {

	if len(dataPoint) > int(maxSize) {
		return &vwapQueue{}, fmt.Errorf("initial datapoints exceeds maxSize")
	}

	return &vwapQueue{
		DataPoints:              dataPoint,
		Limit:                   maxSize,
		CumulativePriceQuantity: make(map[string]float64),
		CumulativeQuantity:      make(map[string]float64),
		VWAP:                    make(map[string]float64),
	}, nil
}

// Size returns the length of the queue.
func (l *vwapQueue) Size() uint {
	return uint(len(l.DataPoints))
}

// Get returns the data point item at given index.
func (l *vwapQueue) Get(i int) DataPoint {
	return l.DataPoints[i]
}

// GetVwap returns the VWAP for a  trading pair.
func (l *vwapQueue) GetVwap(tradingPair string) float64 {
	return l.VWAP[tradingPair]
}

// Push pushes an item onto the queue
//When Limit is reached, will delete  the first one.
func (l *vwapQueue) Push(d DataPoint) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Size() == l.Limit {
		l.remove()
	}

	l.computeVwap(d)
	l.DataPoints = append(l.DataPoints, d)

}

//computeVwap is used to compute the VWAP for a given trading pair.
func (l *vwapQueue) computeVwap(d DataPoint) {
	vw := d.Price * d.Quantity
	l.CumulativePriceQuantity[d.TradingPair] = l.CumulativePriceQuantity[d.TradingPair] + vw
	l.CumulativeQuantity[d.TradingPair] = l.CumulativeQuantity[d.TradingPair] + d.Quantity
	//VWAP = Sum(Price*Quantity) / Sum(Quantity)
	l.VWAP[d.TradingPair] = l.CumulativePriceQuantity[d.TradingPair] / l.CumulativeQuantity[d.TradingPair]
}

// Remove removes 1st item from the queue.
func (l *vwapQueue) remove() {

	it := l.DataPoints[0]
	l.DataPoints[0] = DataPoint{}

	// Subtract the values of 1st item from the VWAP computation..
	l.CumulativePriceQuantity[it.TradingPair] = l.CumulativePriceQuantity[it.TradingPair] - it.Price*it.Quantity
	l.CumulativeQuantity[it.TradingPair] = l.CumulativeQuantity[it.TradingPair] - it.Quantity

	//VWAP = Sum(Price*Quantity) / Sum(Quantity)
	if l.CumulativeQuantity[it.TradingPair] != 0 {
		l.VWAP[it.TradingPair] = l.CumulativePriceQuantity[it.TradingPair] / l.CumulativeQuantity[it.TradingPair]
	}

	//removes 1st item from the queue
	l.DataPoints = l.DataPoints[1:]
}
