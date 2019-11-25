package utils

import "time"

func ExecuteEvery(f func(), d time.Duration) {
	wait := time.NewTicker(d).C
	for range wait {
		f()
	}
}
