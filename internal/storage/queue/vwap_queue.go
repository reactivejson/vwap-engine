package queue

import (
	"github.com/reactivejson/vwap-engine/internal/storage"
	"sync"
)

// VwapQueue represents a queue of DataPoints.
// Manipulation with ArrayList is slow because it internally uses an array. If any element is removed from the array, all the other elements are shifted in memory.
//Every time a new data point is added to the queue and saved for each trading pair, the VWAP computation is updated accordingly.
// For performance, and to avoid exponential complexity, the computation is cached for VWAP, CumulativeQuantity,
//and CumulativePriceQuantity for existing data points and updated with new entries.
type vwapQueue struct {
	mu sync.Mutex
	//DataPoints  is fast circular fifo data structure (aka., queue) with a specific limit.
	DataPoints              []storage.Point
	CumulativePriceQuantity map[string]float64 // Equation = Sum(Price*Quantity), Sum of Price * Quantity for each TradingPair for each window
	CumulativeQuantity      map[string]float64 // Equation = Sum(Quantity), Sum of Quantities for every TradingPair for each window
	VWAP                    map[string]float64 //Equation: VWAP = Sum(Price*Quantity) / Sum(Quantity) Volume Weighted Average Price is calculated for every TradingPair for each window
	// Limit sets the number of data points used to calculate the VWAP.
	Limit uint
}

//NewVwapQueue  creates a new VWAP queue and initializes all fields needed to make the VWAP Queue.
func NewVwapQueue(maxSize uint) (storage.Vwap, error) {
	return &vwapQueue{
		DataPoints:              []storage.Point{},
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
func (l *vwapQueue) GetDataPoints() any {
	return l.DataPoints
}

// GetVwap returns the VWAP for a  trading pair.
func (l *vwapQueue) GetVwap(tradingPair string) float64 {
	return l.VWAP[tradingPair]
}

// Push pushes an item onto the queue
//When Limit is reached, will delete  the first one.
func (l *vwapQueue) Push(d storage.Point) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Size() == l.Limit {
		l.remove()
	}

	l.computeVwap(d)
	l.DataPoints = append(l.DataPoints, d)

}

//computeVwap is used to compute the VWAP for a given trading pair.
func (l *vwapQueue) computeVwap(d storage.Point) {
	vw := d.ComputePQ()
	l.CumulativePriceQuantity[d.ProductID()] = l.CumulativePriceQuantity[d.ProductID()] + vw
	l.CumulativeQuantity[d.ProductID()] = l.CumulativeQuantity[d.ProductID()] + d.GetQuantity()
	//VWAP = Sum(Price*Quantity) / Sum(Quantity)
	l.VWAP[d.ProductID()] = l.CumulativePriceQuantity[d.ProductID()] / l.CumulativeQuantity[d.ProductID()]
}

// Remove removes 1st item from the queue.
func (l *vwapQueue) remove() {

	it := l.DataPoints[0]
	l.DataPoints[0] = nil

	// Subtract the values of 1st item from the VWAP computation..
	l.CumulativePriceQuantity[it.ProductID()] = l.CumulativePriceQuantity[it.ProductID()] - it.ComputePQ()
	l.CumulativeQuantity[it.ProductID()] = l.CumulativeQuantity[it.ProductID()] - it.GetQuantity()

	//VWAP = Sum(Price*Quantity) / Sum(Quantity)
	if l.CumulativeQuantity[it.ProductID()] != 0 {
		l.VWAP[it.ProductID()] = l.CumulativePriceQuantity[it.ProductID()] / l.CumulativeQuantity[it.ProductID()]
	}

	//removes 1st item from the queue
	l.DataPoints = l.DataPoints[1:]
}
