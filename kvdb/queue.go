package kvdb

type Pair struct {
	K string
	V string
}

var Q = make(chan Pair)

func KvLoop() {
	for pair := range Q {
		Set(pair.K, pair.V)
	}
}
