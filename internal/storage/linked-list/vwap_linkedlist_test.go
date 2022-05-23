package linked_list_test

import (
	"container/list"
	"github.com/reactivejson/vwap-engine/internal/storage"
	"github.com/reactivejson/vwap-engine/internal/storage/linked-list"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

var points = map[string]storage.Point{
	"1": storage.NewPoint(1, 1, "TradingPair1"),
	"2": storage.NewPoint(2, 2, "TradingPair2"),
	"3": storage.NewPoint(3, 3, "TradingPair1"),
}

func TestLinkedListPush_withLimits_ShouldSucceed(t *testing.T) {
	t.Parallel()

	vwapQueue, err := linked_list.NewVwapLinkedList(2)
	require.NoError(t, err)

	vwapQueue.Push(points["1"])

	require.Equal(t, 1, int(vwapQueue.Size()))

	l := vwapQueue.GetDataPoints().(list.List)

	require.Equal(t, points["1"], l.Back().Value.(storage.Point))

	vwapQueue.Push(points["2"])
	l = vwapQueue.GetDataPoints().(list.List)
	require.Equal(t, 2, int(vwapQueue.Size()))
	require.Equal(t, points["2"], l.Back().Value.(storage.Point))

	vwapQueue.Push(points["3"])
	l = vwapQueue.GetDataPoints().(list.List)
	require.Equal(t, 2, int(vwapQueue.Size()))
	require.Equal(t, points["3"], l.Back().Value.(storage.Point))
}
func TestVwapLinkedList_Size(t *testing.T) {
	t.Parallel()

	vwapQueue, err := linked_list.NewVwapLinkedList(2)
	require.NoError(t, err)

	vwapQueue.Push(points["1"])
	require.Equal(t, 1, int(vwapQueue.Size()))

	vwapQueue.Push(points["2"])
	require.Equal(t, 2, int(vwapQueue.Size()))
}

func TestLinkedList_ConcurrencyMangnt(t *testing.T) {
	t.Parallel()

	vwapQueue, err := linked_list.NewVwapLinkedList(3)
	require.NoError(t, err)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		vwapQueue.Push(points["1"])
		wg.Done()
	}()

	go func() {
		vwapQueue.Push(points["2"])
		wg.Done()
	}()

	go func() {
		vwapQueue.Push(points["3"])
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, int(vwapQueue.Size()))
}

func TestVwapLinkedList_GetVwap_ShouldCompute_AndSucceed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Name     string
		Data     []storage.Point
		Expected map[string]float64
		Limit    uint
	}{
		{
			Name:  "Test Non existing Data",
			Limit: 3,
			Data: []storage.Point{
				storage.NewPoint(1, 1, "TradingPair1"),
				storage.NewPoint(3, 3, "TradingPair1"),
			},
			Expected: map[string]float64{
				"TradingPair1": 2.5,
				"TradingPair2": 0,
			},
		},
		{
			Name: "4 DataPoints and limited to 3",
			Data: []storage.Point{
				storage.NewPoint(1, 1, "TradingPair1"),
				storage.NewPoint(2, 2, "TradingPair2"),
				storage.NewPoint(3, 3, "TradingPair1"),
				storage.NewPoint(4, 4, "TradingPair2"),
			},
			Limit: 3,
			Expected: map[string]float64{
				"TradingPair1": 3,
				"TradingPair2": 3.3333333333333333,
			},
		},
		{
			Name: "2 DataPoints and limited to 3",
			Data: []storage.Point{
				storage.NewPoint(1, 1.1, "TradingPair1"),
				storage.NewPoint(2, 2.2, "TradingPair1"),
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

			vwapQueue, err := linked_list.NewVwapLinkedList(tt.Limit)
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
