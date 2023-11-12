// doubleTopBottom
package strategy

import (
	"log"
)

func doubleTopBottom(notify chan bool) {
	for {
		<-notify
		log.Println("doubleTopBottom")
	}
}
