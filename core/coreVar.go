// coreVar
package core

var (
	NotifyNewKline = make(chan bool)
	klineSlice     []Kline
	klineLen       int = 0
)

type Kline struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

func GetKlineSliceLen() int {
	return klineLen
}

func appendKlineSlice(newKline Kline) {
	klineSlice = append(klineSlice, newKline)
	klineLen = len(klineSlice)
}

func GetLastKline(nth int) Kline {
	if klineLen-nth < 0 {
		panic("Index out of range")
	}
	return klineSlice[klineLen-nth]
}
