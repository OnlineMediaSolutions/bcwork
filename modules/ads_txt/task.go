package adstxt

const (
	// actions
	AdsTxtTaskActionAdd    = "ADD"
	AdsTxtTaskActionRemove = "REMOVE"
	// types
	AdsTxtTaskTypePublisher     = "PUBLISHER"
	AdsTxtTaskTypeDemandPartner = "DEMAND_PARTNER"
)

type AdsTxtTask struct {
	Action string
	Type   string
	Old    interface{}
	New    interface{}
}

// TODO: needs publisher management
func processAdsTxtTaskTypePublisher(task AdsTxtTask) error {
	return nil
}

func processAdsTxtTaskTypeDemandPartner(task AdsTxtTask) error {
	return nil
}
