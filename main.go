package main

import (
	"flag"
	"log"
	"strings"

	"github.com/driftsec/disco/core"
)

const defaultExtFilters = "png,apng,bmp,gif,ico,cur,jpg,jpeg,jfif,pjp,pjpeg,svg,tif,tiff,webp,xbm,3gp,aac,flac,mpg,mpeg,mp3,mp4,m4a,m4v,m4p,oga,ogg,ogv,mov,wav,webm,eot,woff,woff2,ttf,otf,css"
const defaultStatusFilters = "404"

type arrayHeaders []string

func (i *arrayHeaders) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayHeaders) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var urlin string
	var concur int
	var showB bool
	var extSrcs bool
	var statfils string
	var extfils string
	var headers arrayHeaders
	var cookies arrayHeaders
	var debug bool
	var blregex string

	flag.StringVar(&urlin, "u", "", "the target url")
	flag.IntVar(&concur, "t", 5, "max concurrent threads")
	flag.BoolVar(&showB, "b", false, "show the browser (for debugging)")
	flag.BoolVar(&extSrcs, "n", false, "dont fetch external sources")
	flag.StringVar(&statfils, "fc", defaultStatusFilters, "comma seperated list of status codes to filter, use \"+code1,code2\" to append to the default list")
	flag.StringVar(&extfils, "fe", defaultExtFilters, "comma seperated list of extensions to filter, use \"+ext1,ext2\" to append to the default list")
	flag.StringVar(&blregex, "fr", ".*?logout.*?|.*?login.*?", "URL blacklist regex, useful to avoid hitting logout etc. note: host is checked seperately")
	flag.Var(&headers, "H", "request header, use multiple times")
	flag.Var(&cookies, "C", "request cookie, use multiple times")
	flag.BoolVar(&debug, "v", false, "verbose/debug mode")

	flag.Parse()

	if urlin == "" {
		log.Fatal("-u is required")
	}

	var stats string
	if statfils != "" {
		if strings.HasPrefix(statfils, "+") {
			stats = defaultStatusFilters + "," + strings.TrimPrefix(statfils, "+")
		} else {
			stats = statfils
		}
	} else {
		stats = defaultStatusFilters
	}

	var exts string
	if extfils != "" {
		if strings.HasPrefix(extfils, "+") {
			exts = defaultExtFilters + "," + strings.TrimPrefix(extfils, "+")
		} else {
			exts = extfils
		}
	} else {
		exts = defaultStatusFilters
	}

	cfg := core.SessionConfig{
		ExtFilters:    exts,
		StatusFilters: stats,
		Concurrency:   concur,
		ShowBrowser:   showB,
		DoExtSources:  !extSrcs,
		Dbg:           debug,
		UrlFilter:     blregex,
	}
	cfg.Cookies = make(map[string]string)
	for _, c := range cookies {
		tmpc := strings.Split(c, "=")
		if len(tmpc) != 2 {
			log.Fatal("failed to parse cookie:", c)
		}
		cfg.Cookies[tmpc[0]] = tmpc[1]
	}
	cfg.Headers = make(map[string]string)
	for _, hl := range headers {
		tmp := strings.SplitN(hl, ":", 2)
		if len(tmp) != 2 {
			log.Fatal("failed to parse input header:", hl)
		}
		key := strings.Trim(tmp[0], " ")
		val := strings.Trim(tmp[1], " ")
		cfg.Headers[key] = val
	}

	session := core.New(&cfg)
	session.Start(urlin)

}

// TODO
// js link find needs testing and work. lots of junk urls/extra requests.
// need to deal with popups and alerts, they hang the process
