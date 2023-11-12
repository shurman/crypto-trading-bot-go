// doubleTopBottom
package strategy

import (
	"crypto-trading-bot-go/core"
	"log"
)

func doubleTopBottom(notify chan bool) {
	for {
		<-notify
		log.Println("doubleTopBottom")
		log.Println(len(core.KLineSlice))
	}
}
