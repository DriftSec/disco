package core

import (
	"net/url"
	"regexp"
	"strings"
)

var linkFinderRegex = regexp.MustCompile(`(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')`)

func linkFinder(source string) ([]string, error) {
	var links []string
	// source = strings.ToLower(source)
	if len(source) > 1000000 {
		source = strings.ReplaceAll(source, ";", ";\r\n")
		source = strings.ReplaceAll(source, ",", ",\r\n")
	}
	source = decodeChars(source)

	match := linkFinderRegex.FindAllStringSubmatch(source, -1)
	for _, m := range match {
		matchGroup1 := filterNewLines(m[1])
		if matchGroup1 == "" {
			continue
		}
		links = append(links, matchGroup1)
	}
	links = Unique(links)
	return links, nil
}

func filterNewLines(s string) string {
	return regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(s), " ")
}

func decodeChars(s string) string {
	source, err := url.QueryUnescape(s)
	if err == nil {
		s = source
	}

	// In case json encoded chars
	replacer := strings.NewReplacer(
		`\u002f`, "/",
		`\u0026`, "&",
	)
	s = replacer.Replace(s)
	return s
}
