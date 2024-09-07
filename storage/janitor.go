package storage

import "time"

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(fs *FileStorage) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			fs.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func startJanitor(fs *FileStorage, interval time.Duration) {
	j := &janitor{
		Interval: interval,
		stop:     make(chan bool),
	}
	go j.Run(fs)
}

func stopJanitor(fs *FileStorage) {
	fs.janitor.stop <- true
}
