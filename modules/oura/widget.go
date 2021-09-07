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
	baseURL           = "https://api.ouraring.com/v1"
	userInfoEndpoint  = "%s/userinfo"
	sleepEndpoint     = "%s/sleep?start=%s&end=%s"
	readinessEndpoint = "%s/readiness?start=%s&end=%s"
	activityEndpoint  = "%s/activity?start=%s&end=%s"
	bedtimeEndpoint   = "%s/bedtime?start=%s&end=%s"
)

var (
	accessToken string // Required config value
	dayOffset   int
	client      = http.DefaultClient
)

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, pages, settings.common),
		settings:   settings,
		pageTypes:  []string{"sleep", "readiness", "activity", "bedtime"},
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
		msg := fmt.Sprintf("%8sCould not get Oura data for %s", " ", widget.settings.myName)
		return fmt.Sprintf("\n\n%s", msg)
	}

	_, _, width, _ := widget.View.GetRect()
	userHeader := fmt.Sprintf("%s - [white]%d years[plum] - %5g kg\n", widget.settings.myName, userInfo.Age, userInfo.Weight)
	pageStar := fmt.Sprintf("[goldenrod]%s\n", widget.settings.common.PaginationMarker(len(widget.pageTypes), widget.idx, width-2))

	var pageHeader, pageData string
	switch pageName := widget.pageTypes[widget.idx]; pageName {
	case "sleep":
		pageHeader, pageData = widget.sleepPage()
	case "readiness":
		pageHeader, pageData = widget.readinessPage()
	case "activity":
		pageHeader, pageData = widget.activityPage()
	case "bedtime":
		pageHeader, pageData = widget.bedtimePage()
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
	h = "[orange]Sleep Summaries\n"
	sleepDays, err := getSleepSummaries()
	if err != nil {
		s = fmt.Sprintf("%8sCould not get Oura sleep data for %s", " ", widget.settings.myName)
		return h, fmt.Sprintf("\n\n\n%s", s)
	}
	for i := range sleepDays.SleepDays {
		s = fmt.Sprintf("%s%s", widget.formatSleep(sleepDays.SleepDays[i]), s)
	}
	return h, s
}

func (widget *Widget) readinessPage() (string, string) {
	var h, r string
	h = "[orange]Readiness Summaries\n"
	readinessDays, err := getReadinessSummaries()
	if err != nil {
		r = fmt.Sprintf("%8sCould not get Oura readiness data for %s", " ", widget.settings.myName)
		return h, fmt.Sprintf("\n\n\n%s", r)
	}
	for i := range readinessDays.ReadinessDays {
		r = fmt.Sprintf("%s%s", widget.formatReadiness(readinessDays.ReadinessDays[i]), r)
	}
	return h, r
}

func (widget *Widget) activityPage() (string, string) {
	var h, a string
	h = "[orange]Activity Summaries\n"
	activityDays, err := getActivitySummaries()
	if err != nil {
		a = fmt.Sprintf("%8sCould not get Oura activity data for %s", " ", widget.settings.myName)
		return h, fmt.Sprintf("\n\n\n%s", a)
	}
	for i := range activityDays.ActivityDays {
		a = fmt.Sprintf("%s%s", widget.formatActivity(activityDays.ActivityDays[i]), a)
	}
	return h, a
}

func (widget *Widget) bedtimePage() (string, string) {
	var h, b string
	b = "[orange]Ideal Bedtimes\n"
	bedtimes, err := getBedtimes()
	if err != nil {
		b = fmt.Sprintf("%8sCould not get Oura activity data for %s", " ", widget.settings.myName)
		return h, fmt.Sprintf("\n\n\n%s", b)
	}
	if bedtimes.IdealBedtimes[len(bedtimes.IdealBedtimes)-1].Status == "NOT_ENOUGH_DATA" {
		b = fmt.Sprintf("%8sNot enough sleep data for %s--\n%8sNo ideal bedtimes yet! ", " ", widget.settings.myName, " ")
		return h, fmt.Sprintf("\n\n\n%s", b)
	}
	for i := range bedtimes.IdealBedtimes {
		b = fmt.Sprintf("%s%s", widget.formatBedtime(bedtimes.IdealBedtimes[i]), b)
	}
	return h, b
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

func getSleepSummaries() (*SleepSummaries, error) {
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

func getReadinessSummaries() (*ReadinessSummaries, error) {
	var readinessResp *ReadinessSummaries
	readinessData, err := getData("readiness")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(readinessData, &readinessResp)
	if err != nil {
		return nil, err
	}
	return readinessResp, nil
}

func getActivitySummaries() (*ActivitySummaries, error) {
	var activityResp *ActivitySummaries
	activityData, err := getData("activity")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(activityData, &activityResp)
	if err != nil {
		return nil, err
	}
	return activityResp, nil
}

func getBedtimes() (*Bedtimes, error) {
	var bedtimesResp *Bedtimes
	bedtimesData, err := getData("bedtime")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bedtimesData, &bedtimesResp)
	if err != nil {
		return nil, err
	}
	return bedtimesResp, nil
}

func getData(endpointOption string) ([]byte, error) {
	// Construct endpoint for request
	var endpoint string
	switch switchEndpoint := endpointOption; switchEndpoint {
	case "sleep":
		endpoint = sleepEndpoint
	case "readiness":
		endpoint = readinessEndpoint
	case "activity":
		endpoint = activityEndpoint
	case "bedtime":
		endpoint = bedtimeEndpoint
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
	date := fmt.Sprintf("[grey]%s ", sleep.SummaryDate)
	score := fmt.Sprintf("[lightblue]score: %d\n", sleep.Score)
	asleep := fmt.Sprintf("[pink]asleep from %s to %s (%s)\n", sleep.BedtimeStart.Format("15:04 MST"), sleep.BedtimeEnd.Format("15:04 MST"), secondsToHuman(sleep.Duration))
	stages := fmt.Sprintf("[turquoise]rem: %s [skyblue]light: %s [chartreuse]deep: %s\n", secondsToHuman(sleep.Rem), secondsToHuman(sleep.Light), secondsToHuman(sleep.Deep))
	return fmt.Sprintf("%s%s%s%s \n", date, score, asleep, stages)
}

func (widget *Widget) formatReadiness(readiness Readiness) string {
	date := fmt.Sprintf("[grey]%s ", readiness.SummaryDate)
	score := fmt.Sprintf("[lightblue]score: %d\n", readiness.Score)
	hr := fmt.Sprintf("[pink]resting HR: %d bpm\n", readiness.ScoreRestingHr)
	recovery := fmt.Sprintf("[turquoise]recovery index: %d\n", readiness.ScoreRecoveryIndex)
	return fmt.Sprintf("%s%s%s%s \n", date, score, hr, recovery)
}

func (widget *Widget) formatActivity(activity Activity) string {
	date := fmt.Sprintf("[grey]%s ", activity.SummaryDate)
	score := fmt.Sprintf("[lightblue]score: %d\n", activity.Score)
	activeCals := fmt.Sprintf("[pink]active calories: %d kcal\n", activity.CalActive)
	steps := fmt.Sprintf("[turquoise]movement recorded: %d steps\n", activity.Steps)
	return fmt.Sprintf("%s%s%s%s \n", date, score, activeCals, steps)
}

func (widget *Widget) formatBedtime(bedtime Bedtime) string {
	date := fmt.Sprintf("[grey]%s ", bedtime.Date)
	status := fmt.Sprintf("[lightblue]status: %s\n", bedtime.Status)
	between := fmt.Sprintf("[pink]bedtime window is [chartreuse]between %d and %d \n", bedtime.BedtimeWindow.Start, bedtime.BedtimeWindow.End)
	return fmt.Sprintf("%s%s%s \n", date, status, between)
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
