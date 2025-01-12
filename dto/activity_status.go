package dto

type ActivityStatus int

const ActivePubs = 5000
const LowPubs = 20

const (
	Paused ActivityStatus = iota
	Low
	Active
)

func (s ActivityStatus) String() string {
	return [...]string{"Paused", "Low", "Active"}[s]
}
