package storage

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
}

//Point provides the data required to calculate a VWAP for a specific pair from coinbase.
// since this engine is dealing with financial operations, and for precise results maybe it's better to use libs
// like https://github.com/shopspring/decimal, but in this code we went with float64.
type Point interface {
	ComputePQ() float64
	//GetPrice returns Price of trading pair
	GetPrice() float64
	//GetQuantity returns Quantity of trading pair
	GetQuantity() float64
	//ProductID returns the TradingPair of coinbase product ID
	ProductID() string
}

//DataPoint provides the data required to calculate a VWAP for a specific pair from coinbase.
type DataPoint struct {
	//Price of trading pair
	Price float64
	//Quantity of trading pair
	Quantity float64
	//TradingPair is the coinbase product ID
	TradingPair string
}

func (d *DataPoint) ComputePQ() float64 {
	return d.Price * d.Quantity
}

func (d *DataPoint) GetPrice() float64 {
	return d.Price
}

func (d *DataPoint) GetQuantity() float64 {
	return d.Quantity
}

func (d *DataPoint) ProductID() string {
	return d.TradingPair
}

func NewPoint(p, q float64, t string) Point {
	return &DataPoint{
		Price:       p,
		Quantity:    q,
		TradingPair: t,
	}
}
