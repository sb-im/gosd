package daemon

import (
	"sb.im/gosd/app/model"

	"gorm.io/gorm"
)

func DatabaseMigrate(orm *gorm.DB) {
	orm.AutoMigrate(&model.Team{})
	orm.AutoMigrate(&model.User{})
	orm.AutoMigrate(&model.Session{})
	orm.AutoMigrate(&model.UserTeam{})

	orm.AutoMigrate(&model.Schedule{})
	orm.AutoMigrate(&model.Node{})
	orm.AutoMigrate(&model.Task{})
	orm.AutoMigrate(&model.Blob{})
	orm.AutoMigrate(&model.Job{})
	orm.AutoMigrate(&model.Profile{})
}

func DatabaseSeed(orm *gorm.DB) error {
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

	if err := orm.FirstOrCreate(team, tid).Error; err != nil {
		return err
	}

	user := &model.User{
		Username: "demo",
		Team:     team,
		Teams:    []model.Team{*team},
	}

	if err := orm.FirstOrCreate(user, uid).Error; err != nil {
		return err
	}

	session := &model.Session{
		Team: *team,
		User: *user,
		IP:   "root",
	}

	if err := orm.FirstOrCreate(session, sid).Error; err != nil {
		return err
	}

	// Task
	if err := orm.FirstOrCreate(&model.Task{
		Name:   "demo-task",
		TeamID: team.ID,
		NodeID: 1,
	}).Error; err != nil {
		return err
	}

	return nil
}
