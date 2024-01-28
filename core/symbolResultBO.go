package core

type SymbolResultBO struct {
	CountLong      float64
	CountShort     int
	WinLongCount   int
	LossLongCount  int
	WinShortCount  int
	LossShortCount int
	ProfitLong     float64
	ProfitShort    float64
	LossLong       float64
	LossShort      float64
	InstantWin     int
	InstantLoss    int

	DropDown      float64
	MaxDropDown   float64
	TotalLongFee  float64
	TotalShortFee float64
}

func GetSymbolResultBO() SymbolResultBO {
	bo := SymbolResultBO{
		CountLong:      0,
		CountShort:     0,
		WinLongCount:   0,
		LossLongCount:  0,
		WinShortCount:  0,
		LossShortCount: 0,
		ProfitLong:     0,
		ProfitShort:    0,
		LossLong:       0,
		LossShort:      0,
		InstantWin:     0,
		InstantLoss:    0,
		DropDown:       0,
		MaxDropDown:    0,
		TotalLongFee:   0,
		TotalShortFee:  0,
	}

	return bo
}

func CalculateOrdersResult(closedOrders []*OrderBO) SymbolResultBO {
	bo := GetSymbolResultBO()

	for _, v := range closedOrders {
		if v.GetDirection() == ORDER_LONG {
			bo.CountLong++
			if v.GetFinalProfit() > 0 {
				bo.WinLongCount++

				if v.IsInstantFillAndExit() {
					bo.InstantWin++
				}

				if bo.DropDown < bo.MaxDropDown {
					bo.MaxDropDown = bo.DropDown
				}
				bo.DropDown = 0
				bo.ProfitLong += v.GetFinalProfit() - v.GetFee()

			} else if v.GetFinalProfit() < 0 {
				bo.LossLongCount++
				bo.DropDown += v.GetFinalProfit() - v.GetFee()

				if v.IsInstantFillAndExit() {
					bo.InstantLoss++
				}
				bo.LossLong += v.GetFinalProfit() - v.GetFee()
			}
			bo.TotalLongFee += v.GetFee()

		} else if v.GetDirection() == ORDER_SHORT {
			bo.CountShort++
			if v.GetFinalProfit() > 0 {
				bo.WinShortCount++

				if v.IsInstantFillAndExit() {
					bo.InstantWin++
				}

				if bo.DropDown < bo.MaxDropDown {
					bo.MaxDropDown = bo.DropDown
				}
				bo.DropDown = 0
				bo.ProfitShort += v.GetFinalProfit() - v.GetFee()

			} else if v.GetFinalProfit() < 0 {
				bo.LossShortCount++
				bo.DropDown += v.GetFinalProfit() - v.GetFee()

				if v.IsInstantFillAndExit() {
					bo.InstantLoss++
				}
				bo.LossShort += v.GetFinalProfit() - v.GetFee()
			}
			bo.TotalShortFee += v.GetFee()
		}
	}

	return bo
}

func GetPfRf(profitTotal float64, lossToal float64, netProfit float64, maxDropDown float64) (float64, float64) {
	var profitFactor float64
	var recoveryFactor float64

	if lossToal == 0 {
		profitFactor = profitTotal
	} else {
		profitFactor = profitTotal / -lossToal
	}

	if maxDropDown == 0 {
		recoveryFactor = netProfit / 1
	} else {
		recoveryFactor = netProfit / -maxDropDown
	}

	return profitFactor, recoveryFactor
}
