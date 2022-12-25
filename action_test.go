package chromedpundetected

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {
	testRun(t,
		3,
		NewConfig(
			WithTimeout(20*time.Second),
			WithHeadless(),
		),
		func(ctx context.Context) error {
			version := make(map[string]string)
			err := chromedp.Run(ctx,
				RunCommandWithRes("Browser.getVersion", nil, &version),
			)
			t.Log("Version:", version)
			return err
		},
	)
}

func TestBlockURLs(t *testing.T) {
	btn := `//button[@title="Akkoord"]`

	testRun(t,
		3,
		NewConfig(
			WithTimeout(20*time.Second),
			WithHeadless(),
		),
		func(ctx context.Context) error {
			if err := chromedp.Run(ctx,
				chromedp.Navigate("https://www.nu.nl/"),
				chromedp.WaitVisible(btn),
				chromedp.Click(btn),
			); err != nil {
				return err
			}

			if err := chromedp.Run(ctx,
				BlockURLs("*.nu.nl"),
				chromedp.Navigate("https://www.nu.nl/"),
				chromedp.WaitVisible(btn),
			); err != nil && !errors.Is(err, context.DeadlineExceeded) {
				return err
			}

			return nil
		},
	)
}

func TestCookiesExtract(t *testing.T) {
	btn := `//button[@title="Akkoord"]`

	testRun(t,
		3,
		NewConfig(
			WithTimeout(20*time.Second),
			WithHeadless(),
		),
		func(ctx context.Context) error {
			var cookies []Cookie
			if err := chromedp.Run(ctx,
				chromedp.Navigate("https://www.nu.nl/"),
				chromedp.WaitVisible(btn),
				chromedp.Click(btn),
				chromedp.Sleep(2*time.Second),
				SaveCookies(&cookies),
			); err != nil {
				return err
			}
			require.Greater(t, len(cookies), 0, "cookies len > 0")

			t.Log("Cookies:")
			for _, c := range cookies {
				t.Log(c.Name, "=", c.Value)
			}

			return nil
		},
	)
}
