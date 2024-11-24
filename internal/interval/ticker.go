package interval

import "time"

func OnInterval(duration time.Duration, handler func()) *time.Ticker {
	ticker := time.NewTicker(duration)

	go func() {
		for range ticker.C {
			handler()
		}
	}()

	return ticker
}
