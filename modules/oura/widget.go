package oura

import (
	"fmt"
	"math"
	"strconv"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget
	ouraData  *OuraData
	pageTypes []string
	idx       int
	settings  *Settings
	client    *Client
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, pages, settings.common),
		settings:   settings,
		pageTypes:  []string{"sleep", "readiness", "activity", "bedtime"},
		ouraData:   &OuraData{},
		client:     NewClient(settings),
	}
	widget.initializeKeyboardControls()
	widget.View.SetScrollable(true)
	widget.View.SetWrap(true)
	widget.View.SetWordWrap(true)
	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {
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
	var pageHeader, pageData string
	err := widget.getAllData()
	if err != nil {
		msg := fmt.Sprintf("%8sCould not get Oura data for %s\n%8sCheck personal access token in config!\n%8s%s", " ", widget.settings.myName, " ", " ", err.Error())
		return fmt.Sprintf("\n\n%s", msg)
	}

	_, _, width, _ := widget.View.GetRect()
	userHeader := fmt.Sprintf("%s - [white]%d years[plum] - %5g kg\n", widget.settings.myName, widget.ouraData.UserInfo.Age, widget.ouraData.UserInfo.Weight)
	pageStar := fmt.Sprintf("[goldenrod]%s\n", widget.settings.common.PaginationMarker(len(widget.pageTypes), widget.idx, width-2))

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

func (widget *Widget) getAllData() error {
	// UserInfo
	userInfo, err := widget.client.getUserInfo()
	if err != nil {
		return err
	}
	widget.ouraData.UserInfo = userInfo

	// SleepDays
	sleepDays, err := widget.client.getSleeps()
	if err != nil {
		return err
	}
	widget.ouraData.SleepDays = sleepDays

	// ReadinessDays
	readinessDays, err := widget.client.getReadiness()
	if err != nil {
		return err
	}
	widget.ouraData.ReadinessDays = readinessDays

	// ActivityDays
	activityDays, err := widget.client.getActivity()
	if err != nil {
		return err
	}
	widget.ouraData.ActivityDays = activityDays

	// IdealBedtimes
	idealBedtimes, err := widget.client.getBedtimes()
	if err != nil {
		return err
	}
	widget.ouraData.IdealBedtimes = idealBedtimes

	return nil
}

// Construct display header & body for each page type

func (widget *Widget) sleepPage() (string, string) {
	var h, s string
	h = "[orange]Sleep Summaries\n"
	for i := range widget.ouraData.SleepDays {
		s = fmt.Sprintf("%s%s", widget.formatSleep(widget.ouraData.SleepDays[i]), s)
	}
	if len(widget.ouraData.SleepDays) <= 0 {
		s = fmt.Sprintf("%8sNo Oura sleep data returned for %s", " ", widget.settings.myName)
	}
	return h, s
}

func (widget *Widget) readinessPage() (string, string) {
	var h, r string
	h = "[orange]Readiness Summaries\n"
	for i := range widget.ouraData.ReadinessDays {
		r = fmt.Sprintf("%s%s", widget.formatReadiness(widget.ouraData.ReadinessDays[i]), r)
	}
	if len(widget.ouraData.ReadinessDays) <= 0 {
		r = fmt.Sprintf("%8sNo Oura readiness data returned for %s", " ", widget.settings.myName)
	}
	return h, r
}

func (widget *Widget) activityPage() (string, string) {
	var h, a string
	h = "[orange]Activity Summaries\n"
	for i := range widget.ouraData.ActivityDays {
		a = fmt.Sprintf("%s%s", widget.formatActivity(widget.ouraData.ActivityDays[i]), a)
	}
	if len(widget.ouraData.ActivityDays) <= 0 {
		a = fmt.Sprintf("%8sNo Oura activity data returned for %s", " ", widget.settings.myName)
	}
	return h, a
}

func (widget *Widget) bedtimePage() (string, string) {
	var h, b string
	h = "[orange]Ideal Bedtimes\n"
	if widget.ouraData.IdealBedtimes[len(widget.ouraData.IdealBedtimes)-1].Status == "NOT_ENOUGH_DATA" {
		b = fmt.Sprintf("%8sNot enough sleep data for %s--\n%8sNo ideal bedtimes yet! ", " ", widget.settings.myName, " ")
		return h, fmt.Sprintf("\n\n\n%s", b)
	}
	for i := range widget.ouraData.IdealBedtimes {
		b = fmt.Sprintf("%s%s", widget.formatBedtime(widget.ouraData.IdealBedtimes[i]), b)
	}
	if len(widget.ouraData.IdealBedtimes) <= 0 {
		b = fmt.Sprintf("%8sNo Oura bedtime data returned for %s", " ", widget.settings.myName)
	}
	return h, b
}

// Format date lines for module display

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

// Helpers

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
