// Package middleware holds Gin middleware shared across routes.
package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that reflects allowed origins and handles
// preflight requests. Pass "*" to allow any origin. Origins ending with ":*"
// match any port on the same scheme and hostname, e.g. http://localhost:*.
func CORS(allowed []string) gin.HandlerFunc {
	allow := make(map[string]bool, len(allowed))
	var wildcardPorts []originWildcardPort
	for _, o := range allowed {
		if parsed, ok := parseOriginWildcardPort(o); ok {
			wildcardPorts = append(wildcardPorts, parsed)
			continue
		}
		allow[o] = true
	}
	wildcard := allow["*"]

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && (wildcard || allow[origin] || matchOriginWildcardPort(origin, wildcardPorts)) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type originWildcardPort struct {
	scheme string
	host   string
}

func parseOriginWildcardPort(origin string) (originWildcardPort, bool) {
	if !strings.HasSuffix(origin, ":*") {
		return originWildcardPort{}, false
	}

	u, err := url.Parse(strings.TrimSuffix(origin, ":*"))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return originWildcardPort{}, false
	}

	return originWildcardPort{scheme: u.Scheme, host: u.Hostname()}, true
}

func matchOriginWildcardPort(origin string, allowed []originWildcardPort) bool {
	u, err := url.Parse(origin)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	for _, item := range allowed {
		if u.Scheme == item.scheme && u.Hostname() == item.host {
			return true
		}
	}
	return false
}
