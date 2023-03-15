package v3

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/emersion/go-ical"
)

func (h *Handler) Calendar(c *gin.Context) {
	cal := ical.NewCalendar()

	propProductID := ical.NewProp(ical.PropProductID)
	propProductID.Value = "gosd"
	cal.Props.Add(propProductID)

	propVersion := ical.NewProp(ical.PropVersion)
	propVersion.Value = "2.0"
	cal.Props.Add(propVersion)
	com := ical.NewComponent("VEVENT")

	propSummary := ical.NewProp(ical.PropSummary)
	propSummary.Value = "TTT"
	com.Props.Add(propSummary)

	DTStart := ical.NewProp(ical.PropDateTimeStart)
	DTStart.Value = "20230311T032125Z"
	com.Props.Add(DTStart)

	DTStamp := ical.NewProp(ical.PropDateTimeStamp)
	DTStamp.Value = "20230311T032125Z"
	com.Props.Add(DTStamp)

	Uid := ical.NewProp(ical.PropUID)
	Uid.Value = "47E591E3-945A-4057-9A9F-EC4F0809D9D3"
	com.Props.Add(Uid)

	cal.Children = append(cal.Children, com)

	//ev := ical.NewEvent()
	//ev.DateTimeStart(time.Now())
	cal.Children = append(cal.Children, )
	//c.Writer.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+"url.QueryEscape(blob.Name)")
	if err := ical.NewEncoder(c.Writer).Encode(cal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
