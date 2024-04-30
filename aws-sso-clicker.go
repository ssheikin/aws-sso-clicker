package main

import (
	"bufio"
	"context"
	"errors"
	"github.com/fatih/color"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/theckman/yacspin"
	"os"
	"regexp"
	"time"
)

// Time before MFA step times out
const MFA_TIMEOUT = 20

var cfg = yacspin.Config{
	Frequency:         100 * time.Millisecond,
	CharSet:           yacspin.CharSets[59],
	Suffix:            "AWS SSO Signing in: ",
	SuffixAutoColon:   false,
	Message:           "",
	StopCharacter:     "✓",
	StopFailCharacter: "✗",
	StopMessage:       "Logged in successfully",
	StopFailMessage:   "Log in failed",
	StopColors:        []string{"fgGreen"},
}
var spinner, _ = yacspin.New(cfg)

func main() {
	spinner.Start()

	// get sso url from stdin
	url := getUrl()
	// start aws sso login
	ssoLogin(url)

	spinner.Stop()
	time.Sleep(1 * time.Second)
}

func ssoLogin(url string) {
	spinner.Message(color.MagentaString("init page in the default browser \n"))
	spinner.Pause()

	u := launcher.NewUserMode().Bin("/Applications/Chromium.app/Contents/MacOS/Chromium").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	err := rod.Try(func() {
		page := browser.MustPage(url).MustWindowMaximize()
		spinner.Unpause()
		spinner.Message("click 1/2")
		page.Timeout(MFA_TIMEOUT*time.Second).MustElementR("button", "Confirm and continue").MustWaitEnabled().MustClick()
		page.MustWaitNavigation()
		spinner.Message("click 2/2")
		page.Timeout(MFA_TIMEOUT*time.Second).MustElementR("button", "Allow access").MustWaitEnabled().MustClick()
		page.MustWaitNavigation()
		page.MustClose()
	})

	if errors.Is(err, context.DeadlineExceeded) {
		panic("Timed out waiting for MFA")
	} else if err != nil {
		panic(err.Error())
	}
}

// returns sso url from stdin.
func getUrl() string {
	spinner.Message("reading url from stdin")

	scanner := bufio.NewScanner(os.Stdin)
	url := ""
	for url == "" {
		scanner.Scan()
		t := scanner.Text()
		r, _ := regexp.Compile("^https.*user_code=([A-Z]{4}-?){2}")

		if r.MatchString(t) {
			url = t
		}
	}

	return url
}

// print error message and exit
func panic(errorMsg string) {
	red := color.New(color.FgRed).SprintFunc()
	spinner.StopFailMessage(red("Login failed error - " + errorMsg))
	spinner.StopFail()
	os.Exit(1)
}

// print error message
func error(errorMsg string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	spinner.Message("Warn: " + yellow(errorMsg))
}
