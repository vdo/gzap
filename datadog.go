package gzap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/negroni"
	"go.uber.org/zap/zapcore"
)

// General purpose ResponseWriter that should be usable by major frameworks,
// or easily adapted to them
type ResponseWriter interface {
	http.ResponseWriter
	// Status returns the status code of the response or 0 if the response has
	// not been written
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int

}

// General purpose middleware can be used by multiple frameworks
func DatadogRequestLoggerMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	startTime := time.Now()
	fields := []zapcore.Field{
		String("http.method", r.Method),
		String("http.url", r.URL.String()),
	}

	// We have myriad approaches to extract the 'true' client ip. We should try everything.
	ip := r.Header.Get("Cf-Connecting-Ip")
	if ip == "" {
		ip = r.Header.Get("True-Client-Ip")
		if ip == "" {
			ip = r.Header.Get("X-Forwarded-For")
			if ip != "" {
				ips := strings.Split(ip, ",")
				ip = strings.TrimSpace(ips[0])
			} else {
				ip = r.RemoteAddr
			}
		}
	}

	if ip != "" {
		fields = append(fields, String("network.client.ip", ip))
	}

	userAgent := r.UserAgent()
	if userAgent != "" {
		fields = append(fields, String("http.useragent", userAgent))
	}

	referrer := r.Header.Get("Referer")
	if referrer != "" {
		fields = append(fields, String("http.referer", referrer))
	}

	requestId := r.Header.Get("X-Request-Id")
	if requestId == "" {
		requestId = r.Header.Get("X-Amzn-Trace-Id")
	}
	if requestId != "" {
		fields = append(fields, String("http.request_id", requestId))
	}

	next(rw, r)

	res := rw.(ResponseWriter)

	statusCode := res.Status()
	fields = append(fields, Int("http.status_code", statusCode))
	fields = append(fields, Int64("network.bytes_read", r.ContentLength))

	if res.Written() == true {
		fields = append(fields, Int("network.bytes_written", res.Size()))
	}


	var responseLogger LevedLogger = logger.Info
	if statusCode >= 400 && statusCode < 499 {
		responseLogger = logger.Warn
	} else if statusCode >= 500 && statusCode < 599 {
		responseLogger = logger.Error
	}

	duration := time.Since(startTime).Nanoseconds()
	fields = append(fields, Int64("duration", duration))

	summary := fmt.Sprintf(
		"%d %s %s (%s) %.2fms",
		statusCode,
		r.Method,
		r.URL.RequestURI(),
		ip,
		float64(duration)/float64(1e6),
	)

	responseLogger(summary, fields...)
}

// DatadogRequestLoggerHandler stores request-response information in attributes
// that are recognized by datadog without extra transformations.
func DatadogRequestLoggerHandler(handler http.Handler) http.Handler {
	n := negroni.New()
	n.Use(negroni.HandlerFunc(DatadogRequestLoggerMiddleware))
	n.UseHandler(handler)
	return n
}
