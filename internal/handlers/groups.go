package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	goldap "github.com/go-ldap/ldap/v3"
	"samba4-manager/internal/models"
)

func (app *AppContext) GroupsListGET(c *echo.Context) error {
	groups, err := app.LDAPClient.GetAllGroups("")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "groups/list", map[string]interface{}{
		"Groups": groups,
	})
}

func (app *AppContext) GroupsFormGET(c *echo.Context) error {
	return c.Render(http.StatusOK, "groups/form", map[string]interface{}{
		"Group": nil,
	})
}

func (app *AppContext) GroupsEditGET(c *echo.Context) error {
	name := c.Param("name")

	filter := fmt.Sprintf("(sAMAccountName=%s)", goldap.EscapeFilter(name))
	groups, err := app.LDAPClient.GetAllGroups(filter)
	if err != nil {
		slog.Error("GroupsEditGET error", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if len(groups) == 0 {
		return c.String(http.StatusNotFound, "Group not found")
	}

	return c.Render(http.StatusOK, "groups/form", map[string]interface{}{
		"Group": groups[0],
	})
}

func (app *AppContext) GroupsSavePOST(c *echo.Context) error {
	groupName := c.FormValue("SAMAccountName")
	description := c.FormValue("Description")

	err := app.LDAPClient.CreateGroup(groupName, description)
	if err != nil {
		slog.Error("GroupsSavePOST error", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	username, _ := c.Get("username").(string)
	app.DB.Create(&models.AuditLog{
		AdminUser: username,
		Action:    "CREATE_GROUP",
		ObjectDN:  "CN=" + groupName,
		IPAddress: c.RealIP(),
	})

	return c.Redirect(http.StatusFound, "/groups")
}

func (app *AppContext) GroupsUpdatePOST(c *echo.Context) error {
	name := c.Param("name")
	description := c.FormValue("Description")

	filter := fmt.Sprintf("(sAMAccountName=%s)", goldap.EscapeFilter(name))
	groups, err := app.LDAPClient.GetAllGroups(filter)
	if err != nil || len(groups) == 0 {
		return c.String(http.StatusNotFound, "Group not found")
	}

	err = app.LDAPClient.UpdateGroup(groups[0].DN, description)
	if err != nil {
		slog.Error("GroupsUpdatePOST error", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	username, _ := c.Get("username").(string)
	app.DB.Create(&models.AuditLog{
		AdminUser: username,
		Action:    "UPDATE_GROUP",
		ObjectDN:  groups[0].DN,
		IPAddress: c.RealIP(),
	})

	return c.Redirect(http.StatusFound, "/groups")
}
