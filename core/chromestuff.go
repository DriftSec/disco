package core

import (
	"context"
	"strings"
	"sync"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/inspector"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type RequestMap struct {
	Mx *sync.Mutex
	m  map[network.RequestID]*network.Request
}

func newRqm() *RequestMap {
	return &RequestMap{
		Mx: &sync.Mutex{},
		m:  make(map[network.RequestID]*network.Request),
	}
}
func (r *RequestMap) add(rid network.RequestID, rp *network.Request) {
	r.Mx.Lock()
	defer r.Mx.Unlock()
	r.m[rid] = rp
}

func (r *RequestMap) remove(rid network.RequestID) {
	r.Mx.Lock()
	defer r.Mx.Unlock()
	delete(r.m, rid)
}

func (r *RequestMap) get(rid network.RequestID) *network.Request {
	r.Mx.Lock()
	defer r.Mx.Unlock()

	if val, ok := r.m[rid]; ok {
		return val
	}
	return nil
}

func (s *Session) chromeTab(turl string) {
	if turl == "" {
		return
	}
	s.limiter <- true
	ctx, cancel := chromedp.NewContext(*s.cdpCtx)
	defer func() {
		cancel()
		s.wg.Done()
		<-s.limiter
	}()

	var jslinks []string
	var forms []string
	s.createListenTargets(ctx, cancel)
	err := chromedp.Run(ctx, chromedp.Tasks{
		network.Enable(),
		fetch.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(s.reqHeaders)),
		s.setcookies(),
		chromedp.Navigate(turl),
		// chromedp.WaitVisible(`body`, chromedp.BySearch),
		chromedp.Evaluate(GetLinks, &jslinks),
		chromedp.Evaluate(FireEvents, nil),
		chromedp.Evaluate(GetForms, &forms),
	})
	cancel()

	if err != nil {
		if !strings.Contains(err.Error(), "ERR_CONNECTION_ABORTED") {
			Eprint(err, turl)
		}
		return
	}

	for _, l := range forms {
		s.doVisit(l)
	}
	for _, l := range jslinks {
		s.doVisit(l)
	}

}

func (s *Session) createListenTargets(tabCtx context.Context, tabCtxCancel context.CancelFunc) {
	chromedp.ListenTarget(tabCtx, func(ev interface{}) {
		switch ev := ev.(type) {
		// cancel tabs if they crash
		case *inspector.EventTargetCrashed:
			tabCtxCancel()
			// close alert boxes
		case *page.EventJavascriptDialogOpening:
			go func() {
				if err := chromedp.Run(tabCtx,
					page.HandleJavaScriptDialog(false),
				); err != nil {
					tabCtxCancel()
				}
			}()
			// intercept requests for the scope check
		case *fetch.EventRequestPaused:
			go func() {
				c := chromedp.FromContext(tabCtx)
				ctx := cdp.WithExecutor(tabCtx, c.Target)
				if !s.scopeOk(ev.Request.URL) {
					Dprint("scope fail, skipping:", ev.Request.URL)
					fetch.FailRequest(ev.RequestID, network.ErrorReasonConnectionAborted).Do(ctx)
				} else {
					fetch.ContinueRequest(ev.RequestID).Do(ctx)
				}
			}()
			// catch responses and determine if we should report
		case *network.EventResponseReceived:
			request := s.rqMap.get(ev.RequestID)
			if request != nil {
				if request.HasPostData {
					s.shouldReport(request.URL+"?"+request.PostData+"&isPOST", int(ev.Response.Status))
				} else {
					s.shouldReport(request.URL, int(ev.Response.Status))
				}
			} else {
				s.shouldReport(ev.Response.URL, int(ev.Response.Status))
			}

			go func() {
				c := chromedp.FromContext(tabCtx)
				rbp := network.GetResponseBody(ev.RequestID)
				body, err := rbp.Do(cdp.WithExecutor(tabCtx, c.Target))
				if err != nil {
					return
				}
				if err == nil {
					lnks, _ := linkFinder(string(body))
					for _, l := range lnks {
						nu := s.absoluteUrl(l, ev.Response.URL)
						if nu == "" {
							continue
						}

						s.doVisit(nu)
					}

				}
			}()
			// catch requests, add req data to map so we can reference it later
		case *network.EventRequestWillBeSent:
			s.rqMap.add(ev.RequestID, ev.Request)
			// if s.Dbg {
			// 	fmt.Println("WILL-SEND:", ev.Request.URL)
			// }

		}
	})
}

// setcookies returns a task to navigate to a host with the passed cookies set
// on the network request.
func (s *Session) setcookies() chromedp.Tasks {
	if len(s.reqCookies) < 1 {
		return nil
	}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			// expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			for k, v := range s.reqCookies {
				err := network.SetCookie(k, v).WithPath("/").
					// WithExpires(&expr).
					WithDomain(s.hostname).
					// WithHTTPOnly(true).
					Do(ctx)
				if err != nil {
					Eprint(err)
					return err
				}
			}
			return nil
		}),
	}
}
