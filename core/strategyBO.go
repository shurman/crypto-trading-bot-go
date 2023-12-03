// StrategyBO
package core

type StrategyBO struct {
	info       *StrategyBase
	nextKline  chan Kline
	doneAction chan bool
}

type StrategyBase struct {
	name    string
	symbol  string
	execute func(*Kline, *StrategyBO)
}

func ConstructStrategyBase(_name string, _execute func(*Kline, *StrategyBO)) *StrategyBase {
	return &StrategyBase{
		name:    _name,
		execute: _execute,
	}
}

func InitialStrategy(_info StrategyBase, symbol string) *StrategyBO {
	_info.symbol = symbol
	newStrategy := &StrategyBO{
		info:       &_info,
		nextKline:  make(chan Kline),
		doneAction: make(chan bool),
	}
	go func() {
		for {
			nextKline := <-newStrategy.nextKline
			_info.execute(&nextKline, newStrategy)
			newStrategy.doneAction <- true
		}
	}()

	return newStrategy
}

func (bo *StrategyBO) GetSymbol() string {
	return bo.info.symbol
}

func (bo *StrategyBO) GetChanNextKline() chan Kline {
	return bo.nextKline
}

func (bo *StrategyBO) GetChanDoneAction() chan bool {
	return bo.doneAction
}

func (bo *StrategyBO) ToStandardId(id string) string {
	return bo.info.name + "-" + id
}
