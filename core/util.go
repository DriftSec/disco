package core

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/driftsec/wappalyzer"
	"github.com/gocolly/colly"
	"github.com/tj/go-spin"
)

type spinner struct {
	Spn     *spin.Spinner
	Msg     string
	stopsig bool
}

func Spin(msg string) *spinner {
	spn := spinner{
		Spn: spin.New(),
	}
	go func() {
		for {
			if spn.stopsig {
				os.Stderr.WriteString("\033[2K\r") //\033[2K\r
				break
			}
			os.Stderr.WriteString(fmt.Sprintf("\r  \033[92m%s\033[m %s ", msg, spn.Spn.Next()))
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return &spn
}

func (spn *spinner) Stop() {
	spn.stopsig = true
}

func parseExtFilters(in string) []string {
	t := strings.ReplaceAll(in, " ", "")
	l := strings.Split(t, ",")
	var ret []string
	for _, e := range l {
		e = strings.Trim(e, ".")
		ret = append(ret, e)
	}
	return ret
}

func parseStatusFilters(in string) []int {
	t := strings.ReplaceAll(in, " ", "")
	l := strings.Split(t, ",")
	var ret []int
	for _, e := range l {
		eint, err := strconv.Atoi(e)
		if err != nil {
			log.Fatal(e, "is not a valid status code")
		}
		ret = append(ret, eint)
	}
	ret = append(ret, 0)
	return ret
}

func Unique(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func GetUrlExt(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	pos := strings.LastIndex(u.Path, ".")
	if pos == -1 {
		return ""
	}
	return u.Path[pos+1 : len(u.Path)]
}

func ListContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func FilterExtensions(urls []string, filters []string) []string {
	var ret []string

	for _, u := range urls {
		uc := strings.Split(u, "?")[0]
		uext := GetUrlExt(uc)
		if !ListContains(filters, uext) {
			ret = append(ret, u)
		}

	}
	return ret
}

func NormalizeHost(curl string) string {
	u, err := url.Parse(curl)
	if err != nil {
		return "ERROR HANDLE THIS"
	}
	sch := u.Scheme
	port := u.Port()
	hst := u.Hostname()

	if port == "" {
		switch sch {
		case "http":
			port = "80"
		case "https":
			port = "443"
		default:
			return "ERROR HANDLE THIS"
		}
	}
	var xtra string
	if u.Fragment != "" {
		xtra = "#" + u.Fragment
	}
	return sch + "://" + hst + ":" + port + u.RequestURI() + xtra

}

func NormalizeUrls(ulist []string, baseurl string) []string {
	var ret []string
	baseurl = strings.TrimSuffix(baseurl, "/")
	baseurl = strings.TrimSuffix(baseurl, "?")
	for _, ulurl := range ulist {
		u, err := url.Parse(ulurl)
		if err != nil {
			continue
		}
		newu := baseurl + u.RequestURI()
		newu = strings.TrimSuffix(newu, "/")
		newu = strings.TrimSuffix(newu, "?")
		ret = append(ret, newu)
	}
	return ret
}

func GetHostname(curl string) string {
	u, err := url.Parse(curl)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func CollyHeaders2Wapp(rh *colly.Response) wappalyzer.MapStrOrArray {
	headers := make(wappalyzer.MapStrOrArray)
	for k, v := range *rh.Headers {
		lowerCaseKey := strings.ToLower(k)
		headers[lowerCaseKey] = v
	}
	return headers
}

func CollyCookies2Wapp(rh *colly.Response) wappalyzer.MapStr {
	headers := CollyHeaders2Wapp(rh)
	cookies := make(wappalyzer.MapStr)
	for _, cookie := range headers["set-cookie"] {
		keyValues := strings.Split(cookie, ";")
		for _, keyValueString := range keyValues {
			keyValueSlice := strings.Split(keyValueString, "=")
			if len(keyValueSlice) > 1 {
				key, value := keyValueSlice[0], keyValueSlice[1]
				cookies[key] = value
			}
		}
	}
	return cookies
}

func isMimeType(data string) bool {
	mimes := []string{"application/x-www-form-urlencoded", "text/xml", "application/epub+zip", "application/gzip", "application/java-archive", "application/json", "application/ld+json", "application/msword", "application/octet-stream", "application/ogg", "application/pdf", "application/rtf", "application/vnd.amazon.ebook", "application/vnd.apple.installer+xml", "application/vnd.mozilla.xul+xml", "application/vnd.ms-excel", "application/vnd.ms-fontobject", "application/vnd.ms-powerpoint", "application/vnd.oasis.opendocument.presentation", "application/vnd.oasis.opendocument.spreadsheet", "application/vnd.oasis.opendocument.text", "application/vnd.openxmlformats-officedocument.presentationml.presentation", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.rar", "application/vnd.visio", "application/x-7z-compressed", "application/x-abiword", "application/x-bzip", "application/x-bzip2", "application/x-csh", "application/x-freearc", "application/xhtml+xml", "application/x-httpd-php", "application/xmlifnotreadablefromcasualusers", "application/x-sh", "application/x-shockwave-flash", "application/x-tar", "application/zip", "audio/aac", "audio/midiaudio/x-midi", "audio/mpeg", "audio/ogg", "audio/opus", "audio/wav", "audio/webm", "font/otf", "font/ttf", "font/woff", "font/woff2", "image/bmp", "image/gif", "image/jpeg", "image/png", "image/svg+xml", "image/tiff", "image/vnd.microsoft.icon", "image/webp", "text/calendar", "text/css", "text/csv", "text/html", "text/javascript", "text/plain", "text/xmlifreadablefromcasualusers", "video/3gpp", "video/3gpp2", "video/mp2t", "video/mpeg", "video/ogg", "video/webm", "video/x-msvideo"}
	return ListContains(mimes, data)
}

func (s *Session) absoluteUrl(newPath string, curUrl string) string {
	u, err := url.Parse(curUrl)
	if err != nil {
		return ""
	}
	dir := path.Dir(u.Path)
	u.Path = dir

	curUrl = u.String() + "/"
	if isMimeType(newPath) {
		return ""
	}
	if len(newPath) < 4 {
		return ""
	}
	if newPath[:2] == "//" {
		newPath = "https:" + newPath
	}

	if newPath[:4] != "http" {
		if newPath[:1] == "/" {
			newPath = s.baseUrl + strings.TrimPrefix(newPath, "/")

		} else {
			newPath = curUrl + newPath
		}
	}

	return newPath
}

