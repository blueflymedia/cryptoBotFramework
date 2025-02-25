package bbgo

import (
	"context"

	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/bbgo/pkg/types"
)

// RoiTakeProfit force takes the profit by the given ROI percentage.
type RoiTakeProfit struct {
	Percentage fixedpoint.Value `json:"percentage"`

	session       *ExchangeSession
	orderExecutor *GeneralOrderExecutor
}

func (s *RoiTakeProfit) Bind(session *ExchangeSession, orderExecutor *GeneralOrderExecutor) {
	s.session = session
	s.orderExecutor = orderExecutor

	position := orderExecutor.Position()
	session.MarketDataStream.OnKLineClosed(func(kline types.KLine) {
		if kline.Symbol != position.Symbol || kline.Interval != types.Interval1m {
			return
		}

		closePrice := kline.Close
		if position.IsClosed() || position.IsDust(closePrice) {
			return
		}

		roi := position.ROI(closePrice)
		if roi.Compare(s.Percentage) > 0 {
			// stop loss
			Notify("[RoiTakeProfit] %s take profit is triggered by ROI %s/%s, price: %f", position.Symbol, roi.Percentage(), s.Percentage.Percentage(), kline.Close.Float64())
			_ = orderExecutor.ClosePosition(context.Background(), fixedpoint.One, "roiTakeProfit")
			return
		}
	})
}
