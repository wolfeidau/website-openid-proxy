package stores

// Sessions manage the storage of sessions
type Sessions interface{}

// Stores binds together all the stores into one structure
type Stores struct {
	Sessions Sessions
}
