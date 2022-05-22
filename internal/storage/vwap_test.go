package storage_test

import (
	"github.com/reactivejson/vwap-engine/internal/storage"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var dps = map[string]storage.DataPoint{
	"1": {1, 1, "TradingPair1"},
	"2": {2, 2, "TradingPair2"},
	"3": {3, 3, "TradingPair1"},
}

func TestPush_withLimits_ShouldSucceed(t *testing.T) {
	t.Parallel()

	vwapQueue, err := storage.NewVwapQueue([]storage.DataPoint{}, 2)
	require.NoError(t, err)

	vwapQueue.Push(dps["1"])
	require.Equal(t, 1, int(vwapQueue.Size()))
	require.Equal(t, dps["1"], vwapQueue.Get(0))

	vwapQueue.Push(dps["2"])
	require.Equal(t, 2, int(vwapQueue.Size()))
	require.Equal(t, dps["2"], vwapQueue.Get(1))

	vwapQueue.Push(dps["3"])
	require.Equal(t, 2, int(vwapQueue.Size()))
	require.Equal(t, dps["3"], vwapQueue.Get(1))
}
func TestVwapQueue_Size(t *testing.T) {
	t.Parallel()

	vwapQueue, err := storage.NewVwapQueue([]storage.DataPoint{}, 2)
	require.NoError(t, err)

	vwapQueue.Push(dps["1"])
	require.Equal(t, 1, int(vwapQueue.Size()))

	vwapQueue.Push(dps["2"])
	require.Equal(t, 2, int(vwapQueue.Size()))
}

func TestConcurrencyMangnt(t *testing.T) {
	t.Parallel()

	vwapQueue, err := storage.NewVwapQueue([]storage.DataPoint{}, 3)
	require.NoError(t, err)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		vwapQueue.Push(dps["1"])
		wg.Done()
	}()

	go func() {
		vwapQueue.Push(dps["2"])
		wg.Done()
	}()

	go func() {
		vwapQueue.Push(dps["3"])
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, int(vwapQueue.Size()))
}

func TestVwapQueue_GetVwap_ShouldCompute_AndSucceed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Name     string
		Data     []storage.DataPoint
		Expected map[string]float64
		Limit    uint
	}{
		{
			Name:  "Test Non existing Data",
			Limit: 3,
			Data: []storage.DataPoint{
				{1, 1, "TradingPair1"},
				{3, 3, "TradingPair1"},
			},
			Expected: map[string]float64{
				"TradingPair1": 2.5,
				"TradingPair2": 0,
			},
		},
		{
			Name: "4 DataPoints and limited to 3",
			Data: []storage.DataPoint{
				{1, 1, "TradingPair1"},
				{2, 2, "TradingPair2"},
				{3, 3, "TradingPair1"},
				{4, 4, "TradingPair2"},
			},
			Limit: 3,
			Expected: map[string]float64{
				"TradingPair1": 3,
				"TradingPair2": 3.3333333333333333,
			},
		},
		{
			Name: "2 DataPoints and limited to 3",
			Data: []storage.DataPoint{
				{1, 1.1, "TradingPair1"},
				{2, 2.2, "TradingPair1"},
			},
			Expected: map[string]float64{
				"TradingPair1": 1.6666666666666665,
			},
			Limit: 3,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			vwapQueue, err := storage.NewVwapQueue([]storage.DataPoint{}, tt.Limit)
			require.NoError(t, err)

			for _, d := range tt.Data {
				vwapQueue.Push(d)
			}

			for k := range tt.Expected {
				require.Equal(t, tt.Expected[k], vwapQueue.GetVwap(k))
			}
		})
	}
}
