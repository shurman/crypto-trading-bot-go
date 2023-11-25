// StrategyBO
package core

type StrategyBO struct {
	name       string
	nextKline  chan Kline
	doneAction chan bool
	//Execute   func(StrategyBO)
}

func ConstructStrategy(_name string, _execute func(*Kline, *StrategyBO)) *StrategyBO {
	newStrategy := &StrategyBO{
		name:       _name,
		nextKline:  make(chan Kline),
		doneAction: make(chan bool),
		//Execute:   _execute,
	}
	go func() {
		for {
			nextKline := <-newStrategy.nextKline
			_execute(&nextKline, newStrategy)
			newStrategy.doneAction <- true
		}
	}()

	return newStrategy
}

// func (bo *StrategyBO) GetName() string {
// 	return bo.name
// }

func (bo *StrategyBO) GetChanNextKline() chan Kline {
	return bo.nextKline
}

func (bo *StrategyBO) GetChanDoneAction() chan bool {
	return bo.doneAction
}

func (bo *StrategyBO) ToStandardId(id string) string {
	return bo.name + "-" + id
}
