package v3

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"sb.im/gosd/app/model"

	"github.com/emersion/go-ical"
	"github.com/gin-gonic/gin"
)

func canlendarJobToEvent(name string, job *model.Job) *ical.Component {
	toiCalTime := func(t time.Time) string {
		// https://www.rfc-editor.org/rfc/rfc5545#section-3.3.5
		layout := "20060102T150405Z"
		return t.UTC().Format(layout)
	}

	ev := ical.NewComponent("VEVENT")
	propSummary := ical.NewProp(ical.PropSummary)
	propSummary.Value = fmt.Sprintf("%s#%d", name, job.Index)
	ev.Props.Add(propSummary)

	DTStart := ical.NewProp(ical.PropDateTimeStart)
	DTStart.Value = toiCalTime(job.StartedAt)
	ev.Props.Add(DTStart)

	DTStamp := ical.NewProp(ical.PropDateTimeStamp)
	DTStamp.Value = toiCalTime(job.StartedAt)
	ev.Props.Add(DTStamp)

	DTEnd := ical.NewProp(ical.PropDateTimeEnd)
	DTEnd.Value = toiCalTime(job.StartedAt.Add(time.Duration(math.Abs(float64(job.Duration)) * float64(time.Second))))
	ev.Props.Add(DTEnd)

	Uid := ical.NewProp(ical.PropUID)
	Uid.Value = fmt.Sprintf("%d@sblab.xyz", job.ID)
	ev.Props.Add(Uid)
	return ev
}

// @Summary calendar
// @Schemes Calendar
// @Description iCalendar Protocol
// @Tags calendar
// @Produce text/calendar
// @Success 200
// @Router /calendar.ics [get]
func (h *Handler) Calendar(c *gin.Context) {
	var tasks []model.Task
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	if err := h.orm.WithContext(c).Offset((page-1)*size).Limit(size).Preload("Jobs").Find(&tasks, "team_id = ?", h.getCurrent(c).TeamID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cal := ical.NewCalendar()

	propProductID := ical.NewProp(ical.PropProductID)
	propProductID.Value = "gosd"
	cal.Props.Add(propProductID)

	propVersion := ical.NewProp(ical.PropVersion)
	propVersion.Value = "2.0"
	cal.Props.Add(propVersion)

	for _, task := range tasks {
		for _, job := range task.Jobs {
			ev := canlendarJobToEvent(task.Name, &job)
			cal.Children = append(cal.Children, ev)
		}
	}

	// Content-Type: text/calendar; charset=UTF-8
	c.Writer.Header().Set("Content-Type", "text/calendar; charset=UTF-8")
	if err := ical.NewEncoder(c.Writer).Encode(cal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
