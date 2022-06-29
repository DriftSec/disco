package core

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/desertbit/readline"
)

var Pause bool

type Session struct {
	rqMap         *RequestMap
	showBrowser   bool
	wg            *sync.WaitGroup
	cdpCtx        *context.Context
	cdpCtxCancel  context.CancelFunc
	extFilters    []string
	statusFilters []int
	reqhist       []string
	rephist       []string
	baseUrl       string
	hostname      string
	concurrency   int
	doExtSources  bool
	reqHeaders    map[string]interface{}
	reqCookies    map[string]string
	limiter       chan bool
	Dbg           bool
	urlFilter     string
	// pause         bool
}

type SessionConfig struct {
	Headers       map[string]string
	Cookies       map[string]string
	ExtFilters    string
	StatusFilters string
	Concurrency   int
	ShowBrowser   bool
	DoExtSources  bool
	Dbg           bool
	UrlFilter     string
}

func New(cfg *SessionConfig) *Session {
	spn := Spin("initializing browser")

	ret := &Session{
		rqMap:         newRqm(),
		wg:            new(sync.WaitGroup),
		extFilters:    parseExtFilters(cfg.ExtFilters),
		statusFilters: parseStatusFilters(cfg.StatusFilters),
		concurrency:   cfg.Concurrency,
		showBrowser:   cfg.ShowBrowser,
		doExtSources:  cfg.DoExtSources,
		limiter:       make(chan bool, cfg.Concurrency),
		reqCookies:    cfg.Cookies,
		Dbg:           cfg.Dbg,
		urlFilter:     cfg.UrlFilter,
	}
	if ret.Dbg {
		DebugOn = true
	}
	ret.reqHeaders = make(map[string]interface{})
	for k, v := range cfg.Headers {
		ret.reqHeaders[k] = v
	}

	var ctx context.Context
	var cancel context.CancelFunc

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", !cfg.ShowBrowser),
		chromedp.Flag("ignore-certificate-errors", true),
		// chromedp.UserDataDir(dir),
	)

	ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = chromedp.NewContext(ctx)

	// create a timeout
	// taskCtx, cancel = context.WithTimeout(taskCtx, 100*time.Second)
	// defer cancel()

	ret.cdpCtx = &ctx
	ret.cdpCtxCancel = cancel
	spn.Stop()
	session2Grumble = ret
	initInteractive()
	return ret
}

func (s *Session) Start(targetUrl string) {
	go func() {
		tty, err := os.Open("/dev/tty")
		rl, err := readline.New("> ")
		rl.Config.Stdout = os.Stderr
		if err != nil {
			log.Fatal("failed to open tty")
		}
		defer tty.Close()
		inreader := bufio.NewScanner(tty)
		inreader.Split(bufio.ScanLines)
		for inreader.Scan() {
			instr := string(inreader.Bytes())
			args := strings.Split(strings.TrimSpace(instr), " ")
			os.Args = []string{}
			if len(args) == 1 && args[0] == "" {
				Pause = true
				initInteractive()
				err := app.RunWithReadline(rl)
				if err != nil {
					fmt.Println(err)
				}
				Pause = false

			}
		}
	}()

	u, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatal("failed to parse the Url")
	}
	s.hostname = u.Hostname()
	s.baseUrl = NormalizeHost(u.Scheme + "://" + u.Hostname() + ":" + u.Port())

	if err := chromedp.Run(*s.cdpCtx); err != nil {
		log.Fatal("failed to initialize the browser:", err)
	}

	spn := Spin("fetching external sources")
	// get extern sources
	if s.doExtSources && s.hostname != "127.0.0.1" {
		rawList := OtherSources(s.hostname, false)
		res := FilterExtensions(rawList, s.extFilters)
		res = NormalizeUrls(res, s.baseUrl)
		res = Unique(res)
		for _, seed := range res {
			s.doVisit(seed)
		}
	}
	spn.Stop()

	s.doVisit(targetUrl)
	s.wg.Wait()
	s.cdpCtxCancel()
}

func (s *Session) sameSite(turl string) bool {
	u, err := url.Parse(turl)
	if err != nil {
		return false
	}
	return u.Hostname() == s.hostname
}

func (s *Session) scopeOk(turl string) bool {
	// TODO scope check, baseurl or regex domain, regex to filter endpoints like logout,del user etc. need flags

	// works starts here.
	rex, err := regexp.Compile(s.urlFilter)
	if err != nil {
		log.Fatal(s.urlFilter, "is not a valid RegEx expression")
	}

	if rex.MatchString(turl) {
		return false
	}
	if !s.sameSite(turl) {
		return false
	}
	return true
}

func (s *Session) doVisit(turl string) {
	if Pause {
		for {
			time.Sleep(time.Millisecond * 300)
			if !Pause {
				break
			}
		}
	}
	if turl == "" {
		return
	}

	if !s.extensionOk(turl) {
		return
	}
	turl = NormalizeHost(turl)
	if ListContains(s.reqhist, turl) {
		return
	}
	s.reqhist = append(s.reqhist, turl)

	s.wg.Add(1)

	go s.chromeTab(turl)
}

func (s *Session) shouldReport(u string, sc int) bool {
	if Pause {
		for {
			time.Sleep(time.Millisecond * 300)
			if !Pause {
				break
			}
		}
	}
	u = NormalizeHost(u)

	if GetHostname(u) != s.hostname {
		return false
	}
	if !s.sameSite(u) {
		return false
	}
	if ListContains(s.rephist, u) {
		// we have already reported this one
		return false
	}

	if !s.extensionOk(u) {
		// its a filtered extension
		return false
	}

	if !s.statusCodeOk(sc) {
		return false
	}
	if s.Dbg {
		fmt.Println(u, sc)
	} else {
		fmt.Println(u)
	}
	s.rephist = append(s.rephist, u)
	return true

}

// extensionOk returns true if the extension is not filtered
func (s *Session) extensionOk(curl string) bool {
	ext := GetUrlExt(curl)
	return !ListContains(s.extFilters, ext)
}

// statusCodeOk returns true if the status code is not filtered
func (s *Session) statusCodeOk(sc int) bool {
	for _, s := range s.statusFilters {
		if s == sc {
			return false
		}
	}
	return true
}
