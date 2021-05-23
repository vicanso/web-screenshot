package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
)

var devtoolsWsURL = os.Getenv("DEV_TOOLS_WS_URL")

type screenshotParams struct {
	Width   int
	Height  int
	URL     string
	Quality int
	Visible string
	Delay   time.Duration
	Header  http.Header
}

func captureScreenshot(ctx context.Context, params screenshotParams) ([]byte, error) {
	var buf []byte
	// 如果有配置dev tools ws url
	if devtoolsWsURL != "" {
		allocatorContext, cancel := chromedp.NewRemoteAllocator(ctx, devtoolsWsURL)
		defer cancel()
		ctx = allocatorContext
	}

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	actions := []chromedp.Action{
		chromedp.EmulateViewport(int64(params.Height), int64(params.Height)),
	}
	if len(params.Header) != 0 {
		header := make(map[string]interface{})
		for key, values := range params.Header {
			header[key] = strings.Join(values, ", ")
		}
		actions = append(
			actions,
			network.Enable(),
			network.SetExtraHTTPHeaders(header),
		)
	}

	// 打开页面
	actions = append(actions, chromedp.Navigate(params.URL))

	// 延时
	if params.Delay != 0 {
		actions = append(actions, chromedp.Sleep(params.Delay))
	}
	// 等待元素可访问
	if params.Visible != "" {
		actions = append(actions, chromedp.WaitVisible(params.Visible))
	}

	// 截屏
	actions = append(actions, chromedp.FullScreenshot(&buf, params.Quality))

	err := chromedp.Run(ctx, actions...)

	return buf, err
}

func captureScreenshotHandler(c *elton.Context) (err error) {
	url := c.QueryParam("url")
	if url == "" {
		err = errors.New("url can not be nil")
		return
	}
	params := screenshotParams{
		URL: url,
	}
	params.Width, _ = strconv.Atoi(c.QueryParam("width"))
	params.Height, _ = strconv.Atoi(c.QueryParam("height"))
	params.Quality, _ = strconv.Atoi(c.QueryParam("quality"))
	params.Delay, _ = time.ParseDuration(c.QueryParam("delay"))
	params.Visible = c.QueryParam("visible")
	if c.QueryParam("overrideHeader") != "" {
		params.Header = c.Request.Header
	}
	data, err := captureScreenshot(c.Context(), params)
	if err != nil {
		return
	}
	c.SetContentTypeByExt(".png")
	c.BodyBuffer = bytes.NewBuffer(data)
	return
}

func main() {
	e := elton.New()

	e.Use(middleware.NewDefaultError())

	e.GET("/ping", func(c *elton.Context) error {
		c.BodyBuffer = bytes.NewBufferString("pong")
		return nil
	})

	e.GET("/capture-screenshot", captureScreenshotHandler)

	log.Default().Println("web screenshot server is running")

	err := e.ListenAndServe(":7000")
	if err != nil {
		panic(err)
	}
}
