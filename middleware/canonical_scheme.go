package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-gonic/gin"
)

func RedirectToServerAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverAddress := strings.TrimRight(system_setting.ServerAddress, "/")
		if serverAddress == "" {
			c.Next()
			return
		}

		base, err := url.Parse(serverAddress)
		if err != nil || base.Scheme != "https" || base.Host == "" {
			c.Next()
			return
		}

		host := forwardedHost(c)
		if !strings.EqualFold(host, base.Host) || forwardedProto(c) != "http" {
			c.Next()
			return
		}

		target := *base
		target.Path = c.Request.URL.Path
		target.RawQuery = c.Request.URL.RawQuery
		c.Redirect(http.StatusPermanentRedirect, target.String())
		c.Abort()
	}
}

func forwardedProto(c *gin.Context) string {
	if proto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto")); proto != "" {
		return strings.ToLower(strings.Split(proto, ",")[0])
	}
	if strings.Contains(strings.ToLower(c.GetHeader("CF-Visitor")), `"scheme":"http"`) {
		return "http"
	}
	if strings.Contains(strings.ToLower(c.GetHeader("CF-Visitor")), `"scheme":"https"`) {
		return "https"
	}
	if forwarded := strings.ToLower(c.GetHeader("Forwarded")); forwarded != "" {
		for _, part := range strings.Split(forwarded, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "proto=") {
				return strings.Trim(strings.TrimPrefix(part, "proto="), `"`)
			}
		}
	}
	if c.Request.TLS != nil {
		return "https"
	}
	return ""
}

func forwardedHost(c *gin.Context) string {
	if host := strings.TrimSpace(c.GetHeader("X-Forwarded-Host")); host != "" {
		return strings.Split(host, ",")[0]
	}
	return c.Request.Host
}
