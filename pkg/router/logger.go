package router

import (
	"log"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += int64(n)
	return n, err
}

// RequestLogger middleware untuk mencatat lalu lintas incoming request dan outgoing response secara realtime
func RequestLogger(debugMode bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			clientIP := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				clientIP = forwarded
			}

			// 1. Catat INCOMING REQUEST
			if debugMode {
				log.Printf("--> [IN]  %-6s %s | Client IP: %s", r.Method, r.URL.Path, clientIP)
			}

			// 2. Bungkus response writer untuk menangkap status code dan byte size
			recorder := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // default jika WriteHeader tidak dipanggil secara eksplisit
			}

			// 3. Jalankan handler berikutnya
			next.ServeHTTP(recorder, r)

			// 4. Hitung durasi latency eksekusi
			latency := time.Since(start)

			// 5. Tentukan warna status HTTP
			statusColor := "\033[0;32m" // 2xx hijau
			if recorder.statusCode >= 300 && recorder.statusCode < 400 {
				statusColor = "\033[0;36m" // 3xx cyan
			} else if recorder.statusCode >= 400 && recorder.statusCode < 500 {
				statusColor = "\033[1;33m" // 4xx kuning
			} else if recorder.statusCode >= 500 {
				statusColor = "\033[0;31m" // 5xx merah
			}
			resetColor := "\033[0m"

			// 6. Catat OUTGOING RESPONSE
			log.Printf("<-- [OUT] %s%d %s%s | %-6s %s | Latency: %s | Size: %d bytes",
				statusColor,
				recorder.statusCode,
				http.StatusText(recorder.statusCode),
				resetColor,
				r.Method,
				r.URL.Path,
				formatLatency(latency),
				recorder.bytesWritten,
			)
		})
	}
}

func formatLatency(d time.Duration) string {
	if d < time.Millisecond {
		return d.String() // microsecond
	}
	return d.Round(time.Microsecond).String()
}
