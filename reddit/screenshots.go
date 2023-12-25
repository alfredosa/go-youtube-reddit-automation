package reddit

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/playwright-community/playwright-go"
)

// TakeScreenShot takes a screenshot of a given URL and saves it with the provided ID.
func TakeScreenShot(url string, id string) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.SetViewportSize(1200, 1080)

	if _, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	_, err = page.Evaluate("window.scrollBy(0, 150)")
	if err != nil {
		log.Fatalf("could not scroll: %v", err)
	}

	screenshot_title := fmt.Sprintf("screenshots/%s.png", id)
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(screenshot_title),
	}); err != nil {
		log.Fatalf("could not create screenshot: %v", err)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}

	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}

	uniform, err := isUniformColor(screenshot_title)
	if err != nil {
		log.Fatalf("could not check if screenshot is uniform color: %v", err)
	}
	if uniform {
		log.Printf("Screenshot %s is uniform color, deleting", screenshot_title)
		if err = os.Remove(screenshot_title); err != nil {
			log.Fatalf("could not delete screenshot: %v", err)
		}
	}
}

func SetupPW() {
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}

}

func isUniformColor(imgPath string) (bool, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return false, err
	}

	bounds := img.Bounds()
	firstPixel := img.At(bounds.Min.X, bounds.Min.Y)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if img.At(x, y) != firstPixel {
				return false, nil
			}
		}
	}

	return true, nil
}
