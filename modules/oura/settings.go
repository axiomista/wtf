package oura

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "Oura 💍"
)

// Settings defines the configuration properties for this module
type Settings struct {
	common *cfg.Common
	// Define your settings attributes here
	accessToken string // Oura personal access token (https://cloud.ouraring.com/docs/authentication#create-a-personal-access-token)
	myName      string // Name for your info header, optional
	days        int    // Days of Oura data to retrieve
}

// NewSettingsFromYAML creates a new settings instance from a YAML config block
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),
		// Configure your settings attributes here. See http://github.com/olebedev/config for type details
		accessToken: ymlConfig.UString("accessToken"), // Required
		myName:      ymlConfig.UString("myName", "My Oura"),
		days:        ymlConfig.UInt("days", 3),
	}

	return &settings
}
