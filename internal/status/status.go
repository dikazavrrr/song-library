package status

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type APIStatusReader interface {
	ReadStatus() map[string]string
}

type Status struct {
	startUpTime time.Time
	memStat     runtime.MemStats
	httpClient  *http.Client
}

func NewStatus() *Status {
	return &Status{
		startUpTime: time.Now(),
		httpClient:  &http.Client{Timeout: 2 * time.Second},
	}
}

func (s *Status) ReadStatus() map[string]string {
	status := make(map[string]string)

	runtime.ReadMemStats(&s.memStat)

	status["utime"] = time.Since(s.startUpTime).String()
	status["memtot"] = fmt.Sprint(s.memStat.TotalAlloc/1048576) + "Mb"
	status["gc"] = fmt.Sprint(s.memStat.GCSys/1048576) + " Mb"
	status["sys"] = fmt.Sprint(s.memStat.Sys/1048576) + "Mb"

	return status
}
