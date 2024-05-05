package kvdb

type Pair struct {
	K string
	V string
}

var Q = make(chan Pair)

func KvLoop() {
	for {
		select {
		case pair := <-Q:
			Set(pair.K, pair.V)
		}
	}
}
