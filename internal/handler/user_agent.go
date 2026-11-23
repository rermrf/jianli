package handler

import "strings"

type userAgentMetadata struct {
	Device, OS, OSVersion, Browser, BrowserVersion string
}

func deriveUserAgentMetadata(userAgent string) userAgentMetadata {
	ua := strings.ToLower(userAgent)
	m := userAgentMetadata{Device: "Unknown", OS: "Unknown", Browser: "Unknown"}
	// Match browser-specific tokens before their Chrome/Safari compatibility tokens.
	for _, browser := range []struct{ name, token string }{
		{"WeChat", "micromessenger/"}, {"Edge", "edgios/"}, {"Edge", "edga/"}, {"Edge", "edg/"}, {"Edge", "edge/"},
		{"Opera", "opr/"}, {"Opera", "opios/"}, {"Samsung Internet", "samsungbrowser/"},
		{"Chrome", "crios/"}, {"Chrome", "chrome/"}, {"Firefox", "fxios/"}, {"Firefox", "firefox/"},
	} {
		if strings.Contains(ua, browser.token) {
			m.Browser, m.BrowserVersion = browser.name, userAgentVersion(ua, browser.token)
			break
		}
	}
	if m.Browser == "Unknown" && strings.Contains(ua, "safari/") {
		m.Browser, m.BrowserVersion = "Safari", userAgentVersion(ua, "version/")
	}

	ipadDesktop := strings.Contains(ua, "macintosh") && strings.Contains(ua, "mobile/")
	switch {
	case strings.Contains(ua, "windows"):
		m.OS = "Windows"
		// Windows 10 and 11 share the NT 10.0 user-agent token.
		m.OSVersion = map[string]string{"10.0": "10/11", "6.3": "8.1", "6.2": "8", "6.1": "7", "6.0": "Vista", "5.1": "XP"}[userAgentVersion(ua, "windows nt ")]
	case strings.Contains(ua, "android"):
		m.OS, m.OSVersion = "Android", userAgentVersion(ua, "android ")
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod") || ipadDesktop:
		m.OS = "iOS"
		m.OSVersion = userAgentVersion(ua, " os ")
	case strings.Contains(ua, "mac os x") || strings.Contains(ua, "macintosh"):
		m.OS, m.OSVersion = "macOS", userAgentVersion(ua, "mac os x ")
	case strings.Contains(ua, "cros "):
		m.OS = "ChromeOS"
	case strings.Contains(ua, "linux"):
		m.OS = "Linux"
	}

	switch {
	case strings.Contains(ua, "bot") || strings.Contains(ua, "spider") || strings.Contains(ua, "crawler"):
		m.Device = "Bot"
	case strings.Contains(ua, "ipad") || strings.Contains(ua, "tablet") || ipadDesktop || (m.OS == "Android" && !strings.Contains(ua, "mobile")):
		m.Device = "Tablet"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipod") || strings.Contains(ua, "mobile"):
		m.Device = "Mobile"
	case m.OS != "Unknown":
		m.Device = "Desktop"
	}
	return m
}

func userAgentVersion(ua, token string) string {
	index := strings.Index(ua, token)
	if index < 0 {
		return ""
	}
	value := ua[index+len(token):]
	end := 0
	for end < len(value) && (value[end] >= '0' && value[end] <= '9' || value[end] == '.' || value[end] == '_') {
		end++
	}
	return strings.ReplaceAll(value[:end], "_", ".")
}
