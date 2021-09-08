package oura

import (
	"time"
)

//The structs I used for unmashalling Oura API responses are
// built on ones from the nicest Oura API package,
//**https://github.com/lildude/oura**, with gratitude.

// OuraData is a struct containing all data available for a user through the Oura API.
type OuraData struct {
	UserInfo      *UserInfo   `json:"user_info"`
	SleepDays     []Sleep     `json:"sleep"`
	ReadinessDays []Readiness `json:"readiness"`
	ActivityDays  []Activity  `json:"activity"`
	IdealBedtimes []Bedtime   `json:"ideal_bedtimes"`
}

// UserInfo is the information for the current user
type UserInfo struct {
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
	Gender string  `json:"gender"`
	Email  string  `json:"email"`
}

// SleepSummaries represents all sleep summaries for the period requested
type SleepSummaries struct {
	SleepDays []Sleep `json:"sleep"`
}

// ReadinessSummaries represents all readiness summaries for the period requested
type ReadinessSummaries struct {
	ReadinessDays []Readiness `json:"readiness"`
}

// ActivitySummaries represents all activity in the summary for the period requested
type ActivitySummaries struct {
	ActivityDays []Activity `json:"activity"`
}

// Bedtimes represents all ideal bedtimes for the period requested
type Bedtimes struct {
	IdealBedtimes []Bedtime `json:"ideal_bedtimes"`
}

// Sleep represents a single day's sleep summary
type Sleep struct {
	SummaryDate       string    `json:"summary_date"`
	PeriodID          int       `json:"period_id"`
	IsLongest         int       `json:"is_longest"`
	Timezone          int       `json:"timezone"`
	BedtimeStart      time.Time `json:"bedtime_start"`
	BedtimeEnd        time.Time `json:"bedtime_end"`
	Score             int       `json:"score"`
	ScoreTotal        int       `json:"score_total"`
	ScoreDisturbances int       `json:"score_disturbances"`
	ScoreEfficiency   int       `json:"score_efficiency"`
	ScoreLatency      int       `json:"score_latency"`
	ScoreRem          int       `json:"score_rem"`
	ScoreDeep         int       `json:"score_deep"`
	ScoreAlignment    int       `json:"score_alignment"`
	Total             int       `json:"total"`
	Duration          int       `json:"duration"`
	Awake             int       `json:"awake"`
	Light             int       `json:"light"`
	Rem               int       `json:"rem"`
	Deep              int       `json:"deep"`
	OnsetLatency      int       `json:"onset_latency"`
	Restless          int       `json:"restless"`
	Efficiency        int       `json:"efficiency"`
	MidpointTime      int       `json:"midpoint_time"`
	HrLowest          int       `json:"hr_lowest"`
	HrAverage         float32   `json:"hr_average"`
	Rmssd             int       `json:"rmssd"`
	BreathAverage     float32   `json:"breath_average"`
	TemperatureDelta  float32   `json:"temperature_delta"`
	Hypnogram5Min     string    `json:"hypnogram_5min"`
	Hr5min            []int     `json:"hr_5min"`
	Rmssd5min         []int     `json:"rmssd_5min"`
}

// Readiness represents a single day's readiness summary
type Readiness struct {
	SummaryDate          string `json:"summary_date"`
	PeriodID             int    `json:"period_id"`
	Score                int    `json:"score"`
	ScorePreviousNight   int    `json:"score_previous_night"`
	ScoreSleepBalance    int    `json:"score_sleep_balance"`
	ScorePreviousDay     int    `json:"score_previous_day"`
	ScoreActivityBalance int    `json:"score_activity_balance"`
	ScoreRestingHr       int    `json:"score_resting_hr"`
	ScoreHrvBalance      int    `json:"score_hrv_balance"`
	ScoreRecoveryIndex   int    `json:"score_recovery_index"`
	ScoreTemperature     int    `json:"score_temperature"`
	RestModeState        int    `json:"rest_mode_state"`
}

// Activity represents a single day's activity summary
type Activity struct {
	SummaryDate            string    `json:"summary_date"`
	DayStart               time.Time `json:"day_start"`
	DayEnd                 time.Time `json:"day_end"`
	Timezone               int       `json:"timezone"`
	Score                  int       `json:"score"`
	ScoreStayActive        int       `json:"score_stay_active"`
	ScoreMoveEveryHour     int       `json:"score_move_every_hour"`
	ScoreMeetDailyTargets  int       `json:"score_meet_daily_targets"`
	ScoreTrainingFrequency int       `json:"score_training_frequency"`
	ScoreTrainingVolume    int       `json:"score_training_volume"`
	ScoreRecoveryTime      int       `json:"score_recovery_time"`
	DailyMovement          int       `json:"daily_movement"`
	NonWear                int       `json:"non_wear"`
	Rest                   int       `json:"rest"`
	Inactive               int       `json:"inactive"`
	InactivityAlerts       int       `json:"inactivity_alerts"`
	Low                    int       `json:"low"`
	Medium                 int       `json:"medium"`
	High                   int       `json:"high"`
	Steps                  int       `json:"steps"`
	CalTotal               int       `json:"cal_total"`
	CalActive              int       `json:"cal_active"`
	MetMinInactive         int       `json:"met_min_inactive"`
	MetMinLow              int       `json:"met_min_low"`
	MetMinMediumPlus       int       `json:"met_min_medium_plus"`
	MetMinMedium           int       `json:"met_min_medium"`
	MetMinHigh             int       `json:"met_min_high"`
	AverageMet             float32   `json:"average_met"`
	Class5min              string    `json:"class_5min"`
	Met1min                []float32 `json:"met_1min"`
	RestModeState          int       `json:"rest_mode_state"`
}

// Bedtime represents a single day's bedtime recommendation
type Bedtime struct {
	Date          string `json:"date"`
	BedtimeWindow struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"bedtime_window"`
	Status string `json:"status"`
}
