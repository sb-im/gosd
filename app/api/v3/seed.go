package v3

import (
	"sb.im/gosd/app/model"
)

func (h *Handler) InitSeed() {

	// init seed:
	// - team id: 1
	// - user id: 1, belong to team 1
	// - session id: 1, belong to user 1 && team id 1
	const (
		tid = 1
		uid = 1
		sid = 1
	)

	team := &model.Team{
		Name: "default",
	}
	if h.teamIsExist(tid) {
		h.orm.First(team, tid)
	} else {
		h.orm.Create(team)
	}

	user := &model.User{
		Username: "demo",
		Team:     team,
		Teams:    []model.Team{*team},
	}
	if h.userIsExist(uid) {
		h.orm.First(user, uid)
	} else {
		h.orm.Create(user)
	}

	session := &model.Session{
		Team: *team,
		User: *user,
		IP:   "root",
	}
	if h.sessionIsExist(sid) {
		h.orm.First(session, sid)
	} else {
		h.orm.Create(session)
	}
}

func (h Handler) sessionIsExist(id uint) bool {
	var count int64
	h.orm.Find(&model.Session{}, id).Count(&count)
	if count > 0 {
		return true
	}
	return false
}
