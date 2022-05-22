package storage

//DataPoint provides the data required to calculate a VWAP for a specific pair from coinbase.
// since this engine is dealing with financial operations, and for precise results maybe it's better to use libs
// like https://github.com/shopspring/decimal, but in this code we went with float64.
type DataPoint struct {
	//Price of trading pair
	Price float64
	//Quantity of trading pair
	Quantity float64
	//TradingPair is the coinbase product ID
	TradingPair string
}

// Vwap represents a queue of DataPoints and their VWAPs.
//DataPoints  is fast circular fifo data structure (aka., queue) with a specific limit.
type Vwap interface {
	// Push pushes an item onto the queue and calculates the new VWAP.
	//When Limit is reached, will delete  the first one.
	Push(d DataPoint)

	// Size returns the length of the dara points queue.
	Size() uint

	// Get returns the data point item at given index.
	Get(index int) DataPoint

	// GetVwap returns the VWAP for a  trading pair.
	GetVwap(tradingPair string) float64
}
