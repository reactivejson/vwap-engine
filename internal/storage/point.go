package storage

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

//Point provides the data required to calculate a VWAP for a specific pair from coinbase.
// since this engine is dealing with financial operations, and for precise results maybe it's better to use libs
// like https://github.com/shopspring/decimal, but in this code we went with float64.
type Point interface {
	//ComputePQ returns Price * Quantity
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
