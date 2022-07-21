package timetracking

import (
	"erp/db"
	"erp/users"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Main(r *gin.RouterGroup) {

	db.CreateTimetrackingTable()

	r.GET("/event", users.BasicAuth, event)
	r.GET("/getMonthFromUser/:id", users.BasicAuth, getMonthFromUser)
	r.GET("/getMonthFromUser", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getMonthFromUser/id")
	})
	r.GET("/getAllFromUser/:id", users.BasicAuth, getAllFromUser)
	r.GET("/getAllFromUser", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getAllFromUser/id")
	})
	r.GET("/getLastEventFromUser/:id", users.BasicAuth, getLastEventFromUser)
	r.GET("/getLastEventFromUser", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getLastEventFromUser/id")
	})
}

func event(c *gin.Context) {
	event := c.Request.Header.Get("event")
	userId := c.Request.Header.Get("userid")
	date, _ := time.Parse("2006-01-02 15:04:05", c.PostForm("date"))

	fmt.Println("[DB]: Date: " + date.String())

	currentTime, msg := newTimetracking(event, userId)
	if msg == "" {
		saveTimeOnDB(currentTime)
		c.JSON(200, getLastTimeFromDbByUser(userId))
	} else {
		c.JSON(400, gin.H{"error": msg})
	}
}

func getMonthFromUser(c *gin.Context) {
	userId := c.Param("id")
	dateVal := c.Request.Header.Get("date")
	date, _ := time.Parse("2006-01", dateVal)

	var times []Timetracking = getTimeMonthFromDByUser(userId, date)
	if times != nil {
		c.JSON(200, times)
	} else {
		c.JSON(400, gin.H{"error": "No time found"})
	}
}

func getAllFromUser(c *gin.Context) {
	id := c.Param("id")
	var timetracking []Timetracking = getTimeFromDbByUser(id)
	if timetracking != nil {
		c.JSON(200, timetracking)
	} else {
		c.JSON(400, gin.H{"error": "No time found"})
	}
}

func getLastEventFromUser(c *gin.Context) {
	userid := c.Param("id")
	var time Timetracking = getLastTimeFromDbByUser(userid)

	if time.CreatedBy == userid {
		c.JSON(200, time)
	} else {
		c.JSON(400, gin.H{"error": "No time found"})
	}
}

func newTimetracking(event string, createdBy string) (Timetracking, string) {
	var currentTime = new(Timetracking)
	currentTime.Id = uuid.New().String()
	currentTime.Event = event
	currentTime.CreatedBy = createdBy

	return *currentTime, ""
}

func saveTimeOnDB(timetracking Timetracking) {
	var query string = "INSERT INTO `timetracking` (`Id`, `Event`, `CreatedBy`) VALUES ('" + timetracking.Id + "', '" + timetracking.Event + "', '" + timetracking.CreatedBy + "');"
	db.RunSqlQueryWithoutReturn(query)
}

func getTimeFromDbByUser(userId string) []Timetracking {
	var times []Timetracking
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Timestamp`, `Event`, `CreatedBy` FROM `timetracking` WHERE `CreatedBy` = '" + userId + "';")
	if err != nil {
		fmt.Println("[DB]: Can't Select Timetrackings from DB \t-->\t" + err.Error())
	}
	for rows.Next() {
		var timetracking Timetracking
		err := rows.Scan(&timetracking.Id, &timetracking.Timestamp, &timetracking.Event, &timetracking.CreatedBy)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to Timetracking-Object (timetrackingID: %v) \t-->\t"+err.Error(), timetracking.Id)
		}
		times = append(times, timetracking)
	}
	rows.Close()
	return times
}

func getLastTimeFromDbByUser(userId string) Timetracking {
	var times []Timetracking
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Timestamp`, `Event`, `CreatedBy` FROM `timetracking` WHERE `CreatedBy` = '" + userId + "' ORDER BY Timestamp DESC LIMIT 1;")
	if err != nil {
		fmt.Println("[DB]: Can't Select Timetrackings from DB \t-->\t" + err.Error())
	}
	for rows.Next() {
		var timetracking Timetracking
		err := rows.Scan(&timetracking.Id, &timetracking.Timestamp, &timetracking.Event, &timetracking.CreatedBy)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to Timetracking-Object (timetrackingID: %v) \t-->\t"+err.Error(), timetracking.Id)
		}
		times = append(times, timetracking)
	}
	rows.Close()
	return times[0]
}

func getTimeMonthFromDByUser(userId string, date time.Time) []Timetracking {
	var times []Timetracking
	month := strconv.Itoa(int(date.Month()))
	year := strconv.Itoa(date.Year())
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Timestamp`, `Event`, `CreatedBy` FROM `timetracking` WHERE `CreatedBy` = '" + userId + "' AND MONTH(Timestamp) = " + month + " AND YEAR(Timestamp) = " + year + ";")
	if err != nil {
		fmt.Println("[DB]: Can't Select Timetrackings from DB \t-->\t" + err.Error())
	}
	for rows.Next() {
		var timetracking Timetracking
		err := rows.Scan(&timetracking.Id, &timetracking.Timestamp, &timetracking.Event, &timetracking.CreatedBy)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to Timetracking-Object (timetrackingID: %v) \t-->\t"+err.Error(), timetracking.Id)
		}
		times = append(times, timetracking)
	}
	rows.Close()
	return times
}

type Timetracking struct {
	Id        string
	Timestamp time.Time
	Event     string
	CreatedBy string
}
