package handlers

import (
	"gorm.io/gorm"
	"samba4-manager/internal/auth"
	"samba4-manager/internal/config"
	"samba4-manager/internal/ldap"
)

// AppContext holds dependencies for the handlers
type AppContext struct {
	Config     *config.Config
	DB         *gorm.DB
	LDAPClient *ldap.Client
	SessionMgr *auth.SessionManager
}
