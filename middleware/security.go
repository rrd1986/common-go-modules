package middleware

import "net/http"

//OptionsContentTypeHeader represents the X-Content-Type-Options Header param
const OptionsContentTypeHeader = "X-Content-Type-Options"

//XSSProtectionHeader represents the X-XXS-Protection Header param
const XSSProtectionHeader = "X-XSS-Protection"

// CacheControlHeader represents the Cache Control Header param
const CacheControlHeader = "Cache-Control"

//PragmaHeader represents the Pragma Header param
const PragmaHeader = "Pragma"

// cspHeader represents the CSP Header param
const CSPHeader = "Content-Security-Policy"

// hstsHeader represents the HSTS Header param
const HSTSHeader = "Strict-Transport-Security"

//Middleware to handle the security headers
func SecurityWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(OptionsContentTypeHeader, "nosniff")
		w.Header().Set(XSSProtectionHeader, "1; mode=block")
		w.Header().Set(CacheControlHeader, "no-store")
		w.Header().Set(PragmaHeader, "no-cache")
		w.Header().Set(CSPHeader, "frame-ancestors 'self'")
		w.Header().Set(HSTSHeader, "max-age=31536000")
		next.ServeHTTP(w, r)
	})
}
