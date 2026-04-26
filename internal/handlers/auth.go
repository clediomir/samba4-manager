package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"samba4-manager/internal/auth"
	"samba4-manager/internal/models"
)

// AuthLoginGET handles showing the login page
func (app *AppContext) AuthLoginGET(c *echo.Context) error {
	return c.Render(http.StatusOK, "auth/login", nil)
}

// AuthLoginPOST handles processing the login form
func (app *AppContext) AuthLoginPOST(c *echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Try to authenticate via LDAP
	u, err := auth.AuthenticateUser(app.Config, username, password)
	if err != nil {
		return c.Render(http.StatusBadRequest, "auth/login", map[string]interface{}{
			"Error": "Invalid credentials or LDAP error",
		})
	}

	// Determine admin status: if AdminGroup is empty, all users are admins (backwards compat).
	isAdmin := isAdminUser(u.MemberOf, app.Config.RBAC.AdminGroup)

	// Create session
	session, err := app.SessionMgr.CreateSession(u.SAMAccountName, c.RealIP(), c.Request().UserAgent(), isAdmin)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "auth/login", map[string]interface{}{
			"Error": "Failed to create session",
		})
	}

	app.SessionMgr.SetCookie(c.Response(), session.Token)

	return c.Redirect(http.StatusFound, "/dashboard")
}

// AuthLogout handles ending the session
func (app *AppContext) AuthLogout(c *echo.Context) error {
	if sessionRaw := c.Get("session"); sessionRaw != nil {
		if s, ok := sessionRaw.(*models.Session); ok {
			_ = app.SessionMgr.DeleteSession(s.Token)
		}
	}
	app.SessionMgr.ClearCookie(c.Response())
	return c.Redirect(http.StatusFound, "/auth/login")
}

// isAdminUser returns true if the user belongs to the configured AdminGroup.
// When adminGroup is empty, all authenticated users are treated as admins.
// Supports both simple group name and DN format (e.g. CN=Domain Admins,CN=Users,DC=domain,DC=tld).
func isAdminUser(memberOf []string, adminGroup string) bool {
	if adminGroup == "" {
		return true
	}
	adminGroupTrimmed := strings.TrimSpace(adminGroup)
	for _, g := range memberOf {
		gTrimmed := strings.TrimSpace(g)
		// Check exact match
		if strings.EqualFold(gTrimmed, adminGroupTrimmed) {
			return true
		}
		// Check DN format: CN=Domain Admins,CN=Users,...
		if strings.HasPrefix(strings.ToLower(gTrimmed), "cn="+strings.ToLower(adminGroupTrimmed)+",") {
			return true
		}
	}
	return false
}
