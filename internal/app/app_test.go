package app

import (
	"github.com/reactivejson/vwap-engine/api/models"
	"github.com/reactivejson/vwap-engine/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

func TestParseData_ShouldSucceed(t *testing.T) {
	t.Parallel()

	data := &models.CoinbaseResponse{
		Price:     "1",
		ProductID: "TradingPair1",
		Size:      "1",
	}

	dataPoint, err := parseData(data)
	require.NoError(t, err)
	require.Equal(t, storage.DataPoint{
		Price:       1,
		Quantity:    1,
		TradingPair: "TradingPair1"},
		dataPoint)
}

func TestParseData_ShouldFail(t *testing.T) {
	t.Parallel()

	data := &models.CoinbaseResponse{
		Price:     "fail",
		ProductID: "TradingPair1",
		Size:      "1",
	}

	_, err := parseData(data)
	require.Error(t, err)

}
