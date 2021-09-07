package oura

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget
	pageTypes []string
	idx       int
	settings  *Settings
}

const (
	baseURL          = "https://api.ouraring.com/v1"
	sleepEndpoint    = "%s/sleep?start=%s&end=%s"
	userInfoEndpoint = "%s/userinfo"
)

var (
	accessToken = "" // Required config value
	dayOffset   = -3
	client      = http.DefaultClient
)

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, pages, settings.common),
		settings:   settings,
		// pageTypes: []string{"sleep", "activity", "readiness", "bedtime"},
		pageTypes: []string{"sleep", "balderdash"},
	}
	accessToken = widget.settings.accessToken
	dayOffset = widget.settings.days * -1

	widget.initializeKeyboardControls()
	widget.View.SetScrollable(true)
	widget.View.SetWrap(true)
	widget.View.SetWordWrap(true)
	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {
	accessToken = widget.settings.accessToken
	dayOffset = widget.settings.days * -1
	// The last call should always be to the display function
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) display() {
	widget.Redraw(func() (string, string, bool) {
		return widget.CommonSettings().Title, widget.content(), false
	})
}

func (widget *Widget) content() string {
	userInfo, err := getUserInfo()
	if err != nil {
		msg := fmt.Sprintf("Could not get Oura data for %s", widget.settings.myName)
		return fmt.Sprintf("\n\n%s", utils.CenterText(msg, 16))
	}

	_, _, width, _ := widget.View.GetRect()
	userHeader := fmt.Sprintf("%s - [white]%d years[plum] - %5g kg\n", widget.settings.myName, userInfo.Age, userInfo.Weight)
	pageStar := fmt.Sprintf("[goldenrod]%s\n", widget.settings.common.PaginationMarker(len(widget.pageTypes), widget.idx, width-2))

	var pageHeader, pageData string
	switch pageName := widget.pageTypes[widget.idx]; pageName {
	case "sleep":
		pageHeader, pageData = widget.sleepPage()
	// case "activity":
	// case "readiness":
	// case "bedtime":
	default:
		msg := fmt.Sprintf("[blue]Could not find specified page: \"%s\"\n", pageName)
		notFoundPage := fmt.Sprintf("%s%s%s", userHeader, pageStar, msg)
		return notFoundPage
	}

	fullPage := fmt.Sprintf("%s%s%s%s", userHeader, pageStar, pageHeader, pageData)
	return fullPage
}

func (widget *Widget) sleepPage() (string, string) {
	var h, s string
	h = "[orange]Sleep Data\n"
	sleepDays, err := getSleepDays()
	if err != nil {
		s = fmt.Sprintf("Could not get Oura data for %s", widget.settings.myName)
		return h, fmt.Sprintf("\n\n\n%s", utils.CenterText(s, 16))
	}
	for i := range sleepDays.SleepDays {
		s += widget.formatSleep(sleepDays.SleepDays[i])
	}
	return h, s
}

func getUserInfo() (*UserInfo, error) {
	var userInfoResp *UserInfo
	userData, err := getData("user")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(userData, &userInfoResp)
	if err != nil {
		return nil, err
	}
	return userInfoResp, err
}

func getSleepDays() (*SleepSummaries, error) {
	var sleepResp *SleepSummaries
	sleepData, err := getData("sleep")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(sleepData, &sleepResp)
	if err != nil {
		return nil, err
	}
	return sleepResp, nil
}

func getData(endpointOption string) ([]byte, error) {
	// Construct endpoint for request
	var endpoint string
	switch switchEndpoint := endpointOption; switchEndpoint {
	case "sleep":
		endpoint = sleepEndpoint
	// case "activity":
	// 	endpoint = activityEndpoint
	// case "readiness":
	// 	endpoint = readinessEndpoint
	// case "bedtime":
	// 	endpoint = bedtimeEndpoint
	case "user":
		endpoint = userInfoEndpoint
	default:
		return nil, fmt.Errorf("\"%s\" is not a valid endpoint option", endpointOption)
	}
	if endpointOption == "user" {
		endpoint = fmt.Sprintf(endpoint, baseURL)
	} else {
		start := time.Now().AddDate(0, 0, dayOffset).Format("2006-01-02")
		end := time.Now().AddDate(0, 0, 0).Format("2006-01-02") // Today
		endpoint = fmt.Sprintf(endpoint, baseURL, start, end)
	}

	// Build request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Send request and ensure success
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, err
	}
	// Reads response body to []byte
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (widget *Widget) formatSleep(sleep Sleep) string {
	asleep := fmt.Sprintf("asleep from %s to %s (%s)", sleep.BedtimeStart.Format("15:04 MST"), sleep.BedtimeEnd.Format("15:04 MST"), secondsToHuman(sleep.Duration))
	stages := fmt.Sprintf("[turquoise]rem: %s [skyblue]light: %s [chartreuse]deep: %s", secondsToHuman(sleep.Rem), secondsToHuman(sleep.Light), secondsToHuman(sleep.Deep))
	return fmt.Sprintf("[grey]%s [lightblue]score: %d\n[pink]%s\n%s\n\n", sleep.SummaryDate, sleep.Score, asleep, stages)
}

func secondsToHuman(input int) string {
	hours := math.Floor(float64(input) / 60 / 60)
	input = input % (60 * 60)
	minutes := math.Floor(float64(input) / 60)
	if hours > 0 {
		return fmt.Sprintf("%s %s", plural(int(hours), "hr"), plural(int(minutes), "min"))
	}
	if minutes > 0 {
		return plural(int(minutes), "min")
	}
	return plural(int(input), "sec")
}

func plural(count int, singular string) string {
	if count <= 1 {
		return strconv.Itoa(count) + " " + singular
	}
	return strconv.Itoa(count) + " " + singular + "s"

}
