package strategy

type OrderDirection int

const (
	LONG  OrderDirection = 1
	SHORT OrderDirection = -1
)

func (dir OrderDirection) toString() string {
	if dir == LONG {
		return "LONG"
	}
	return "SHORT"
}
