package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"samba4-manager/internal/models"
)

// SettingsData holds the settings data for the template
type SettingsData struct {
	SessionTimeoutMinutes int
	MaxLoginAttempts      int
	LockoutMinutes        int
	RequireTOTP           bool
	AdminGroup            string
	OperatorGroup         string
	HelpdeskGroup         string
	ReadonlyGroup         string
}

// getSettingValue retrieves a setting value from the DB or returns a default
func (app *AppContext) getSettingValue(key string, defaultValue string) string {
	var setting models.Setting
	if err := app.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		// Setting not found, create it with default value
		setting = models.Setting{Key: key, Value: defaultValue}
		app.DB.Create(&setting)
		return defaultValue
	}
	return setting.Value
}

// getIntSettingValue retrieves an integer setting value from the DB or returns a default
func (app *AppContext) getIntSettingValue(key string, defaultValue int) int {
	strValue := app.getSettingValue(key, fmt.Sprintf("%d", defaultValue))
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		// Invalid value in DB, reset to default
		app.DB.Model(&models.Setting{}).Where("key = ?", key).Update("value", fmt.Sprintf("%d", defaultValue))
		return defaultValue
	}
	return intValue
}

// getBoolSettingValue retrieves a boolean setting value from the DB or returns a default
func (app *AppContext) getBoolSettingValue(key string, defaultValue bool) bool {
	strValue := app.getSettingValue(key, fmt.Sprintf("%t", defaultValue))
	boolValue, err := strconv.ParseBool(strValue)
	if err != nil {
		// Invalid value in DB, reset to default
		app.DB.Model(&models.Setting{}).Where("key = ?", key).Update("value", fmt.Sprintf("%t", defaultValue))
		return defaultValue
	}
	return boolValue
}

// saveSettingValue saves a setting value to the DB
func (app *AppContext) saveSettingValue(key, value string) {
	var setting models.Setting
	if err := app.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		// Setting not found, create it
		setting = models.Setting{Key: key, Value: value}
		app.DB.Create(&setting)
	} else {
		// Update existing setting
		app.DB.Model(&setting).Update("value", value)
	}
}

func (app *AppContext) SettingsGET(c *echo.Context) error {
	// Load settings from DB or use defaults
	settingsData := SettingsData{
		SessionTimeoutMinutes: app.getIntSettingValue("session_timeout_minutes", 30),
		MaxLoginAttempts:      app.getIntSettingValue("max_login_attempts", 5),
		LockoutMinutes:        app.getIntSettingValue("lockout_minutes", 15),
		RequireTOTP:           app.getBoolSettingValue("require_totp", false),
		AdminGroup:            app.getSettingValue("admin_group", "Domain Admins"),
		OperatorGroup:         app.getSettingValue("operator_group", "SambaWebOperators"),
		HelpdeskGroup:         app.getSettingValue("helpdesk_group", "SambaWebHelpdesk"),
		ReadonlyGroup:         app.getSettingValue("readonly_group", "SambaWebReadOnly"),
	}

	return c.Render(http.StatusOK, "settings", map[string]interface{}{
		"Settings": settingsData,
		"Config":   app.Config,
	})
}

func (app *AppContext) SettingsPOST(c *echo.Context) error {
	// Parse settings update form
	sessionTimeoutStr := c.FormValue("session_timeout_minutes")
	maxLoginAttemptsStr := c.FormValue("max_login_attempts")
	lockoutMinutesStr := c.FormValue("lockout_minutes")
	requireTOTP := c.FormValue("require_totp") == "true"
	adminGroup := c.FormValue("admin_group")
	operatorGroup := c.FormValue("operator_group")
	helpdeskGroup := c.FormValue("helpdesk_group")
	readonlyGroup := c.FormValue("readonly_group")

	// Save settings to DB
	app.saveSettingValue("session_timeout_minutes", sessionTimeoutStr)
	app.saveSettingValue("max_login_attempts", maxLoginAttemptsStr)
	app.saveSettingValue("lockout_minutes", lockoutMinutesStr)
	app.saveSettingValue("require_totp", fmt.Sprintf("%t", requireTOTP))
	app.saveSettingValue("admin_group", adminGroup)
	app.saveSettingValue("operator_group", operatorGroup)
	app.saveSettingValue("helpdesk_group", helpdeskGroup)
	app.saveSettingValue("readonly_group", readonlyGroup)

	// Set flash message
	c.Set("flash", "Settings saved successfully!")

	return c.Redirect(http.StatusFound, "/settings")
}