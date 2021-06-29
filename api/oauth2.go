package api

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/server"
)

func (h *handler) oauthInit() {
	h.oauth.SetAllowGetAccessRequest(true)
	h.oauth.SetClientInfoHandler(server.ClientFormHandler)
	h.oauth.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	h.oauth.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	h.oauth.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		log.Info(username,": " , password)

		user, err := h.store.UserByUsername(username)
		if err != nil {
			log.Error(err)
			return
		}

		// TODO: This UserByUsername sql.ErrNoRows not return error
		// Maybe this need Change
		if user == nil {
			err = errors.New("Not user found")
			return
		}

		if err = h.store.CheckPassword(username, password); err != nil {
			log.Error(err)
			return
		}

		userID = strconv.FormatInt(user.ID, 10)

		return
	})
}

func (h *handler) authorizeHandler(w http.ResponseWriter, r *http.Request) {
	err := h.oauth.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *handler) oAuthHandler(w http.ResponseWriter, r *http.Request) {
	h.oauth.HandleTokenRequest(w, r)
}
