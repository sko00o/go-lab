package main

import (
	"fmt"

	"github.com/go-rod/rod"
)

func main() {
	browser := rod.New().
		// SlowMotion(1 * time.Second).
		MustConnect().
		NoDefaultDevice()
	page := browser.
		MustPage("https://www.wikipedia.org/")

	// page = page.MustWindowFullscreen()

	page.MustElement("#searchInput").MustInput("earth")
	page.MustElement("#search-form > fieldset > button").MustClick()

	ele := page.MustElement("#mw-content-text > div.mw-parser-output > div:nth-child(3)")
	fmt.Println(ele.MustText())
}
