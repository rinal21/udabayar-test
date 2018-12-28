package helper

import "fmt"

type UnitHelper struct{}

func (u *UnitHelper) WeightFormat(weight int) string {
	if weight >= 1000 {
		return fmt.Sprintf("%.2f kg", float64(weight)/1000)
	}

	return fmt.Sprintf("%d gram", weight)
}

func NewUnitHelper() *UnitHelper {
	return &UnitHelper{}
}
