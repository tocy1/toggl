package handlers

import (
	"crypto/subtle"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/tocy1/toggl/config"
	"github.com/tocy1/toggl/db"
)

type APIHandler struct {
	AdminUsername string
	AdminPassword string
	Logger        lager.Logger
	db            *db.MariaDBDataStore
}

func NewApiHandler(sc config.Server, logger lager.Logger, db *db.MariaDBDataStore) APIHandler {
	return APIHandler{
		sc.AdminUsername,
		sc.AdminPassword,
		logger.Session("handler"),
		db,
	}
}

func (h *APIHandler) authorized(r *http.Request) bool {
	username, password, isOk := r.BasicAuth()
	if isOk {
		return subtle.ConstantTimeCompare([]byte(h.AdminUsername), []byte(username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(h.AdminPassword), []byte(password)) == 1
	}
	return false
}
