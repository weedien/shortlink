package toolkit

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2/log"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetLinkCacheExpiration 计算短链接缓存有效期时间
func GetLinkCacheExpiration(validDate time.Time) time.Duration {
	if validDate.IsZero() {
		// 默认有效期为30天
		return time.Hour * 24 * 30
	}
	return validDate.Sub(time.Now())
}

// GetActualIp 获取用户真实IP
func GetActualIp(r *http.Request) string {
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = r.Header.Get("Proxy-Client-IP")
	}
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = r.Header.Get("WL-Proxy-Client-IP")
	}
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = r.Header.Get("HTTP_CLIENT_IP")
	}
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

// GetOs 获取用户访问操作系统
func GetOs(r *http.Request) string {
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	switch {
	case strings.Contains(userAgent, "windows"):
		return "Windows"
	case strings.Contains(userAgent, "mac"):
		return "Mac OS"
	case strings.Contains(userAgent, "linux"):
		return "Linux"
	case strings.Contains(userAgent, "android"):
		return "Android"
	case strings.Contains(userAgent, "iphone"), strings.Contains(userAgent, "ipad"):
		return "iOS"
	default:
		return "Unknown"
	}
}

// GetBrowser 获取用户访问浏览器
func GetBrowser(r *http.Request) string {
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	switch {
	case strings.Contains(userAgent, "edg"):
		return "Microsoft Edge"
	case strings.Contains(userAgent, "chrome"):
		return "Google Chrome"
	case strings.Contains(userAgent, "firefox"):
		return "Mozilla Firefox"
	case strings.Contains(userAgent, "safari"):
		return "Apple Safari"
	case strings.Contains(userAgent, "opera"):
		return "Opera"
	case strings.Contains(userAgent, "msie"), strings.Contains(userAgent, "trident"):
		return "Internet Explorer"
	default:
		return "Unknown"
	}
}

// GetDevice 获取用户访问设备
func GetDevice(r *http.Request) string {
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	if strings.Contains(userAgent, "mobile") {
		return "Mobile"
	}
	return "PC"
}

// GetNetwork 获取用户访问网络
func GetNetwork(r *http.Request) string {
	actualIp := GetActualIp(r)
	if strings.HasPrefix(actualIp, "192.168.") || strings.HasPrefix(actualIp, "10.") {
		return "WIFI"
	}
	return "Mobile"
}

// ExtractDomain 获取原始链接中的域名
func ExtractDomain(rawUrl string) string {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	domain := parsedUrl.Hostname()
	if strings.HasPrefix(domain, "www.") {
		domain = domain[4:]
	}
	return domain
}

const DefaultCacheValidTime = 86400000 // Example default cache valid time in milliseconds

func GetFaviconWithDefault(url string, defaultValue string) string {
	favicon, err := GetFavicon(url)
	if err != nil {
		return defaultValue
	}
	if favicon == "" {
		return defaultValue
	}
	return favicon
}

// GetFavicon 获取网站图标
func GetFavicon(rawUrl string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", fmt.Errorf("invalid Url: %v", err)
	}

	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		return "", fmt.Errorf("error while fetching Url: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("error while closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error while parsing document: %v", err)
		}
		faviconLink := doc.Find("link[rel~=(?i)^(shortcut )?icon]").First()
		if faviconLink != nil {
			return faviconLink.AttrOr("abs:href", ""), nil
		}
	}
	return "", nil
}
