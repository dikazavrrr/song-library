package limiter

import (
	"fmt"
	"net/http"
	"slices"
	"song-library/pkg/logger"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"golang.org/x/time/rate"
)

const (
	defaultRps              = 10
	defaultBurst            = 20
	defaultTTL              = time.Minute * 5
	defaultCleanupFrequency = time.Minute * 5
	defaultPeriod           = time.Second
)

const (
	localNetPrefix = "192.168.90"
	k8sPrefix      = "10."
	localhost      = "::1"
	ingressIP      = "31.171.171.71"
	notAllowedFmt  = "ratelimited for ip: %s, fwdfor: %s, xorigfwdfor: %s, clientIP: %s"
)

type (
	Limiter interface {
		Stop()
		visitor(string) *rate.Limiter
		whiteListed(string) bool
	}

	limiter struct {
		storage map[string]*record
		opts    *limiterOptions
		stop    chan struct{}
		limit   rate.Limit
		sync.RWMutex
	}

	record struct {
		lastSeen time.Time
		limiter  *rate.Limiter
	}

	limiterOptions struct {
		ttl           time.Duration
		customPeriod  bool
		period        time.Duration
		burst         int
		requests      int
		cleanupFreq   time.Duration
		allowedPrefix []string
		allowedIPs    []string
	}

	option func(*limiterOptions)
)

func Limit(l Limiter) func(http.Handler) http.Handler {
	const msg = "Too many requests"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			if !l.visitor(ip).Allow() {
				http.Error(w, msg, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GinLimit(l Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.GetHeader("x-original-forwarded-for")

		if l.whiteListed(ip) {
			c.Next()
		}

		if !l.visitor(ip).Allow() {
			logger.Error(fmt.Sprintf(notAllowedFmt, ip, c.GetString("x-forwarded-for"), ip, c.ClientIP()))
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}

// visitor lloks up entry in storage and returns *rate.Limiter, updating lastSeen field. Doesnt check if string is empty, so will return same updated limiter for all empty ip visitors.
func (lim *limiter) visitor(ip string) *rate.Limiter {
	lim.RLock()
	v, e := lim.storage[ip]
	lim.RUnlock()

	if !e {
		l := rate.NewLimiter(lim.limit, lim.opts.burst)

		lim.Lock()
		lim.storage[ip] = &record{
			lastSeen: time.Now(),
			limiter:  l,
		}
		lim.Unlock()

		return l
	}

	if v != nil {
		lim.Lock()
		v.lastSeen = time.Now()
		lim.Unlock()
	}

	return v.limiter
}

// Returns new instance of ratelimiter. If no opts are provided uses default settings. RPS = 10, BURST = 20, ttl and cleanup 5 minutes.
func New(opts ...option) Limiter {
	o := defautlOptions()

	for _, opt := range opts {
		opt(o)
	}

	lim := &limiter{
		storage: make(map[string]*record),
		opts:    o,
		stop:    make(chan struct{}),
		limit:   rate.Limit(float64(o.requests) / o.period.Seconds()),
	}

	go lim.scheduleCleanup()

	return lim
}

func defautlOptions() *limiterOptions {
	return &limiterOptions{
		ttl:           defaultTTL,
		requests:      defaultRps,
		burst:         defaultBurst,
		period:        defaultPeriod,
		cleanupFreq:   defaultCleanupFrequency,
		allowedPrefix: []string{k8sPrefix, localNetPrefix},
		allowedIPs:    []string{localhost, ingressIP},
	}
}

// RpsWithBurst sets custom rps and burst values. If burst is zero, all events will be blocked, unless rps is inf. If burst < rps requests will be limited by burst.
func RpsWithBurst(rps, burst int) option {
	if rps < 0 {
		rps = defaultRps
	}

	if burst < 0 {
		burst = defaultRps
	}

	return func(opts *limiterOptions) {
		opts.requests = rps
		opts.burst = burst
	}
}

// Rps sets allowed rps, burst will be disabled
func Rps(rps int) option {
	if rps < 0 {
		rps = defaultRps
	}

	return func(opts *limiterOptions) {
		opts.requests = rps
		opts.burst = rps
	}
}

func Burst(burst int) option {
	if burst < 0 {
		burst = defaultBurst
	}

	return func(opts *limiterOptions) {
		opts.burst = burst
	}
}

// Period sets allowed period when rps is smaller than 1. For example 1 request per 5 seconds. In most cases set burst to 1.
func Period(requests int, period time.Duration) option {
	if period < 0 {
		period = defaultPeriod
	}
	if requests < 0 {
		requests = 0
	}

	return func(opts *limiterOptions) {
		opts.period = period
		opts.customPeriod = true
		opts.requests = requests
	}
}

// CleanupFrequency sets how often to cleanup storage.
func CleanupFrequency(cf time.Duration) option {
	if cf <= 0 {
		cf = defaultCleanupFrequency
	}

	return func(opts *limiterOptions) {
		opts.cleanupFreq = cf
	}
}

// RecordTTL sets lifetime of every record, before it is expired.
func RecordTTL(ttl time.Duration) option {
	if ttl < 0 {
		ttl = defaultTTL
	}

	return func(opts *limiterOptions) {
		opts.ttl = ttl
	}
}

func AllowedIPs(ip ...string) option {
	return func(opts *limiterOptions) {
		opts.allowedIPs = append(opts.allowedIPs, ip...)
	}
}

func AllowedPrefixes(prefix ...string) option {
	return func(opts *limiterOptions) {
		opts.allowedPrefix = append(opts.allowedPrefix, prefix...)
	}
}

func (lim *limiter) scheduleCleanup() {
	ti := time.NewTicker(lim.opts.cleanupFreq)
	defer ti.Stop()

	for {
		select {
		case <-ti.C:
			lim.cleanup()
		case <-lim.stop:
			return
		}
	}
}

// cleanup cleans up storage map when called based on ttl option and last seen time
func (lim *limiter) cleanup() {
	exp := make([]string, len(lim.storage)>>1)

	lim.RLock()
	for k, v := range lim.storage {
		if v == nil {
			exp = append(exp, k)
		}

		if v == nil || time.Since(v.lastSeen) >= lim.opts.ttl {
			exp = append(exp, k)
		}
	}
	lim.RUnlock()

	lim.Lock()
	for _, k := range exp {
		delete(lim.storage, k)
	}
	lim.Unlock()
}

func (lim *limiter) whiteListed(ip string) bool {
	return slices.Contains(lim.opts.allowedIPs, ip) || lim.hasWhitelistedPrefix(ip)
}

func (lim *limiter) hasWhitelistedPrefix(ip string) bool {
	for _, v := range lim.opts.allowedPrefix {
		if strings.HasPrefix(ip, v) {
			return true
		}
	}

	return false
}

// Stop stops cleanup in limiter
func (lim *limiter) Stop() {
	close(lim.stop)
}
