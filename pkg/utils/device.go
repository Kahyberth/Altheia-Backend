package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GetDeviceTypeFromUserAgent(userAgent string) string {
	userAgent = strings.ToLower(userAgent)

	if strings.Contains(userAgent, "mobile") || strings.Contains(userAgent, "android") {
		if strings.Contains(userAgent, "android") {
			return "Mobile App / Android"
		}
		return "Mobile / Unknown"
	}

	if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") {
		return "Mobile App / iOS"
	}

	if strings.Contains(userAgent, "chrome") {
		if strings.Contains(userAgent, "windows") {
			return "Chrome / Windows"
		} else if strings.Contains(userAgent, "mac") {
			return "Chrome / macOS"
		} else if strings.Contains(userAgent, "linux") {
			return "Chrome / Linux"
		}
		return "Chrome / Unknown"
	}

	if strings.Contains(userAgent, "firefox") {
		if strings.Contains(userAgent, "windows") {
			return "Firefox / Windows"
		} else if strings.Contains(userAgent, "mac") {
			return "Firefox / macOS"
		} else if strings.Contains(userAgent, "linux") {
			return "Firefox / Linux"
		}
		return "Firefox / Unknown"
	}

	if strings.Contains(userAgent, "safari") && !strings.Contains(userAgent, "chrome") {
		if strings.Contains(userAgent, "mac") {
			return "Safari / macOS"
		}
		return "Safari / Unknown"
	}

	if strings.Contains(userAgent, "edge") {
		return "Edge / Windows"
	}

	if strings.Contains(userAgent, "windows") {
		return "Browser / Windows"
	} else if strings.Contains(userAgent, "mac") {
		return "Browser / macOS"
	} else if strings.Contains(userAgent, "linux") {
		return "Browser / Linux"
	}

	return "Unknown Device"
}

func GetLocationFromIP(ipAddress string) string {
	if ipAddress == "127.0.0.1" || ipAddress == "::1" || strings.HasPrefix(ipAddress, "192.168.") || strings.HasPrefix(ipAddress, "10.") {
		return "Local Network"
	}

	return "Unknown Location"
}

func GetClientIP(xForwardedFor, xRealIP, remoteAddr string) string {

	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xRealIP != "" {
		return xRealIP
	}

	if strings.Contains(remoteAddr, ":") {
		parts := strings.Split(remoteAddr, ":")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	return remoteAddr
}

func GenerateNanoID() string {
	return uuid.New().String()
}
