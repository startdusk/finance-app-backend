package v1

// ActCreated is an act indicates that create action was finished.
type ActCreated struct {
	Created bool `json:"created"`
}

// ActDeleted is an act indicates that deleted action was finished.
type ActDeleted struct {
	Deleted bool `json:"deleted"`
}

// ActUpdated is an act indicates that updated action was finished.
type ActUpdated struct {
	Updated bool `json:"updated"`
}
