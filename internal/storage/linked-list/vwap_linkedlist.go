package linked_list

import (
	"container/list"
	"github.com/reactivejson/vwap-engine/internal/storage"
	"strings"

	"sync"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

// vwapLinkedList represents a doubly linked list as a queue of DataPoints.
// Manipulation with Linked List is faster than Array List because it uses a doubly linked list, so no bit shifting is required in memory.
//Every time a new data point is added to the queue and saved for each trading pair, the VWAP computation is updated accordingly.
// For performance, and to avoid exponential complexity, the computation is cached for VWAP, CumulativeQuantity,
//and CumulativePriceQuantity for existing data points and updated with new entries.
type vwapLinkedList struct {
	mu sync.Mutex
	//The arrays allocated in memory are never returned. Therefor A dynamic doubly Linked list structure, is better to be used for a long-living queue.
	//DataPoints  is fast circular fifo data structure (aka., Linked list queue) with a specific limit.
	DataPoints              *list.List         //doubly linked list as a queue
	CumulativePriceQuantity map[string]float64 // Equation = Sum(Price*Quantity), Sum of Price * Quantity for each TradingPair for each window
	CumulativeQuantity      map[string]float64 // Equation = Sum(Quantity), Sum of Quantities for every TradingPair for each window
	VWAP                    map[string]float64 //Equation: VWAP = Sum(Price*Quantity) / Sum(Quantity) Volume Weighted Average Price is calculated for every TradingPair for each window
	// Limit sets the number of data points used to calculate the VWAP.
	Limit uint
}

//NewVwapLinkedList  creates a new VWAP queue and initializes all fields needed to make the VWAP Queue.
func NewVwapLinkedList(maxSize uint) (storage.Vwap, error) {
	return &vwapLinkedList{
		DataPoints:              list.New(),
		Limit:                   maxSize,
		CumulativePriceQuantity: make(map[string]float64),
		CumulativeQuantity:      make(map[string]float64),
		VWAP:                    make(map[string]float64),
	}, nil
}

// Size returns the length of the queue.
func (l *vwapLinkedList) Size() uint {
	return uint(l.DataPoints.Len())
}

func (l *vwapLinkedList) GetDataPoints() any {
	return *l.DataPoints
}

// GetVwap returns the VWAP for a  trading pair.
func (l *vwapLinkedList) GetVwap(tradingPair string) float64 {
	return l.VWAP[tradingPair]
}

func (l *vwapLinkedList) GetVwaps() map[string]float64 {
	return l.VWAP
}

// Push pushes an item onto the queue
//When Limit is reached, will delete  the first one.
func (l *vwapLinkedList) Push(d storage.Point) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Size() == l.Limit {
		l.remove()
	}

	l.computeVwap(d)
	l.DataPoints.PushBack(d)
}

//computeVwap is used to compute the VWAP for a given trading pair.
func (l *vwapLinkedList) computeVwap(d storage.Point) {
	vw := d.ComputePQ()
	l.CumulativePriceQuantity[d.ProductID()] = l.CumulativePriceQuantity[d.ProductID()] + vw
	l.CumulativeQuantity[d.ProductID()] = l.CumulativeQuantity[d.ProductID()] + d.GetQuantity()
	//VWAP = Sum(Price*Quantity) / Sum(Quantity)
	l.VWAP[d.ProductID()] = l.CumulativePriceQuantity[d.ProductID()] / l.CumulativeQuantity[d.ProductID()]
}

// Remove removes 1st item from the queue.
func (l *vwapLinkedList) remove() {

	it := l.DataPoints.Front()
	// Subtract the values of 1st item from the VWAP computation.
	l.CumulativePriceQuantity[it.Value.(storage.Point).ProductID()] = l.CumulativePriceQuantity[it.Value.(storage.Point).ProductID()] - it.Value.(storage.Point).ComputePQ()
	l.CumulativeQuantity[it.Value.(storage.Point).ProductID()] = l.CumulativeQuantity[it.Value.(storage.Point).ProductID()] - it.Value.(storage.Point).GetQuantity()

	//VWAP = Sum(Price*Quantity) / Sum(Quantity)
	if l.CumulativeQuantity[it.Value.(storage.Point).ProductID()] != 0 {
		l.VWAP[it.Value.(storage.Point).ProductID()] = l.CumulativePriceQuantity[it.Value.(storage.Point).ProductID()] / l.CumulativeQuantity[it.Value.(storage.Point).ProductID()]
	}

	//removes 1st item from the queue
	l.DataPoints.Remove(it)
}

func (l *vwapLinkedList) String() string {
	return strings.Join(storage.Format(l.VWAP), " | ")
}
