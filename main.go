package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"strconv"

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
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(int64(params.Height), int64(params.Height)),
		// 打开页面
		chromedp.Navigate(params.URL),
		chromedp.FullScreenshot(&buf, params.Quality),
	)
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
