package router

import (
	"log"
	"net/http"
	"strings"
)

// RouteRegistrar antarmuka abstraksi untuk mendaftarkan rute HTTP
type RouteRegistrar interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Handle(pattern string, handler http.Handler)
}

// RouteInfo menyimpan metadata rute yang terdaftar
type RouteInfo struct {
	Method string
	Path   string
}

// Router pembungkus http.ServeMux dengan kapabilitas inspeksi rute dan mode debug
type Router struct {
	*http.ServeMux
	routes []RouteInfo
}

// New membuat instance Router baru
func New() *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
		routes:   make([]RouteInfo, 0),
	}
}

// HandleFunc mendaftarkan function handler dan mencatat rute untuk inspeksi debug
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.record(pattern)
	r.ServeMux.HandleFunc(pattern, handler)
}

// Handle mendaftarkan http.Handler dan mencatat rute untuk inspeksi debug
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.record(pattern)
	r.ServeMux.Handle(pattern, handler)
}

func (r *Router) record(pattern string) {
	trimmed := strings.TrimSpace(pattern)
	parts := strings.SplitN(trimmed, " ", 2)
	if len(parts) == 2 {
		r.routes = append(r.routes, RouteInfo{
			Method: parts[0],
			Path:   parts[1],
		})
	} else {
		r.routes = append(r.routes, RouteInfo{
			Method: "ALL",
			Path:   parts[0],
		})
	}
}

// Routes mengembalikan daftar semua rute yang telah didaftarkan
func (r *Router) Routes() []RouteInfo {
	return r.routes
}

// TotalRoutes mengembalikan jumlah rute yang terdaftar
func (r *Router) TotalRoutes() int {
	return len(r.routes)
}

// PrintBanner menampilkan ringkasan rute, status port, dan informasi aplikasi saat booting
func (r *Router) PrintBanner(port, env string, debugMode bool) {
	log.Println("========================================================================")
	log.Printf("🚀  MODULAR MONOLITH WEB API (Golang net/http)")
	log.Printf("📡  Port Server : :%s", port)
	log.Printf("🔒  Environment : %s", env)
	log.Printf("🛠️   Debug Mode  : %t", debugMode)
	log.Printf("📋  Total Routes: %d Rute Terdaftar", len(r.routes))
	log.Println("------------------------------------------------------------------------")
	log.Println("  DAFTAR ROUTER ENDPOINT:")
	for _, route := range r.routes {
		color := "\033[0;32m" // default green (GET)
		switch route.Method {
		case "POST":
			color = "\033[0;34m" // blue
		case "PUT":
			color = "\033[1;33m" // yellow
		case "DELETE":
			color = "\033[0;31m" // red
		}
		reset := "\033[0m"
		log.Printf("  %s[%-6s]%s %s", color, route.Method, reset, route.Path)
	}
	log.Println("========================================================================")
	if debugMode {
		log.Printf("💡 Debug Mode AKTIF: Semua Incoming (-->) & Outgoing (<--) request akan dicatat.")
	}
	log.Println("")
}
