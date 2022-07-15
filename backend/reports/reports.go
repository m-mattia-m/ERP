package reports

import (
	"erp/db"
	"erp/users"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func Main(r *gin.RouterGroup) {
	db.CreateReportsTable()

	r.POST("/createReport", users.BasicAuth, createReport)
	r.POST("/editReport/:id", users.BasicAuth, editReport)
	r.POST("/editReport", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /editReport/id")
	})
	r.GET("/deleteReport/:id", users.BasicAuth, deleteReport)
	r.GET("/deleteReport", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /deleteReport/id")
	})
	r.GET("/getReport/:id", users.BasicAuth, getReport)
	r.GET("/getReport", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getReport/id")
	})
	r.GET("/getReports", users.BasicAuth, getReports)

}

func createReport(c *gin.Context) {

	var duration float32
	currentUserID := c.Request.Header.Get("userid")

	title := c.PostForm("title")
	description := c.PostForm("description")
	date, _ := time.Parse("2006-01-02 15:04:05", c.PostForm("date"))
	durationParse, _ := strconv.ParseFloat(c.PostForm("duration"), 32)
	duration = float32(durationParse)
	customerId := c.PostForm("customerId")
	createdBy := currentUserID

	fmt.Println("[DB]: Date: " + date.String())

	currentReport, msg := newReport(title, description, date, duration, customerId, createdBy)
	if msg == "" {
		saveReportOnDB(currentReport)
		c.JSON(200, currentReport)
	} else {
		c.JSON(400, gin.H{"error": msg})
	}
}

func editReport(c *gin.Context) {
	id := c.Param("id")

	var report Report = getReportsFromDBById(id)

	if report.Id == id {
		report.Title = c.PostForm("title")
		report.Description = c.PostForm("description")
		report.Date, _ = time.Parse("2006-01-02 15:04:05", c.PostForm("date"))
		fmt.Println("[DB]: Date: " + report.Date.String())
		durationParse, _ := strconv.ParseFloat(c.PostForm("duration"), 32)
		report.Duration = float32(durationParse)
		report.CustomerId = c.PostForm("customerId")
		report.CreatedBy = c.Request.Header.Get("userid")
		updateReportOnDB(report)
		c.JSON(200, report)
	} else {
		c.JSON(400, gin.H{"error": "No report found with the id: " + id})
	}
}

func deleteReport(c *gin.Context) {
	id := c.Param("id")
	var report Report = getReportsFromDBById(id)

	if report.Id == id {
		deleteReportFromDB(report)
		c.JSON(200, gin.H{"message": "delete report with the id: " + id})
	} else {
		c.JSON(400, gin.H{"error": "No report found with the id: " + id})
	}
}

func getReport(c *gin.Context) {
	id := c.Param("id")
	var report Report = getReportsFromDBById(id)

	if report.Id == id {
		c.JSON(200, report)
	} else {
		c.JSON(400, gin.H{"error": "No report found with the id: " + id})
	}
}

func getReports(c *gin.Context) {
	var reports []Report = getReportsFromDB()
	if reports != nil {
		c.JSON(200, reports)
	} else {
		c.JSON(400, gin.H{"error": "No report found"})
	}
}

func newReport(Title string, Description string, Date time.Time, Duration float32, CustomerId string, CreatedBy string) (Report, string) {
	var reports []Report = getReportsFromDB()
	i := sort.Search(len(reports), func(i int) bool {
		return reports[i].Title >= Title && reports[i].CustomerId >= CustomerId && reports[i].Date.Equal(Date)
	})
	if i < len(reports) && (reports[i].Title >= Title && reports[i].CustomerId >= CustomerId && reports[i].Date.Equal(Date)) {
		report := Report{}
		return report, "Report already exists"
	}

	var currentReport = new(Report)
	currentReport.Id = uuid.New().String()
	currentReport.Title = Title
	currentReport.Description = Description
	currentReport.Date = Date
	currentReport.Duration = Duration
	currentReport.CustomerId = CustomerId
	currentReport.CreatedBy = CreatedBy

	return *currentReport, ""
}

func getReportsFromDB() []Report {
	var reports []Report
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Title`, `Description`, `Date`, `Duration`, `CustomerId`, `CreatedBy` FROM `reports`;")
	if err != nil {
		fmt.Println("[DB]: Can't Select Reports from DB \t-->\t" + err.Error())
	}
	for rows.Next() {
		var report Report
		err := rows.Scan(&report.Id, &report.Title, &report.Description, &report.Date, &report.Duration, &report.CustomerId, &report.CreatedBy)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to Report-Object (reportId: %v) \t-->\t"+err.Error(), report.Id)
		}
		reports = append(reports, report)
	}
	rows.Close()
	return reports
}

func saveReportOnDB(report Report) {
	var query string = "INSERT INTO `reports` (`Id`, `Title`, `Description`, `Date`, `Duration`, `CustomerId`, `CreatedBy`) VALUES ('" + report.Id + "', '" + report.Title + "', '" + report.Description + "', '" + report.Date.Format("2006-01-02 15:04:05") + "', '" + fmt.Sprintf("%f", report.Duration) + "', '" + report.CustomerId + "', '" + report.CreatedBy + "');"
	db.RunSqlQueryWithoutReturn(query)
}

func updateReportOnDB(report Report) {
	var query string = "UPDATE `reports` SET `Title` = '" + report.Title + "', `Description` = '" + report.Description + "', `Date` = '" + report.Date.Format("2006-01-02 15:04:05") + "', `Duration` = '" + fmt.Sprintf("%f", report.Duration) + "', `CustomerId` = '" + report.CustomerId + "', `CreatedBy` = '" + report.CreatedBy + "' WHERE `Id` = '" + report.Id + "';"
	db.RunSqlQueryWithoutReturn(query)
}

func deleteReportFromDB(report Report) {
	var query string = "DELETE FROM `reports` WHERE `Id` = '" + report.Id + "';"
	db.RunSqlQueryWithoutReturn(query)
}

func getReportsFromDBById(id string) Report {
	var reports []Report
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Title`, `Description`, `Date`, `Duration`, `CustomerId`, `CreatedBy` FROM `reports` WHERE `Id` = '" + id + "';")
	if err != nil {
		fmt.Println("[DB]: Can't Select report by Id from DB \t-->\t" + err.Error())
	}

	for rows.Next() {
		var report Report
		err := rows.Scan(&report.Id, &report.Title, &report.Description, &report.Date, &report.Duration, &report.CustomerId, &report.CreatedBy)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to Report-Object (reportId: %v) \t-->\t"+err.Error(), report.Id)
		}
		reports = append(reports, report)
	}
	rows.Close()
	return reports[0]
}

type Report struct {
	Id               string
	Title            string
	Description      string
	Date             time.Time
	Duration         float32
	CustomerId       string
	CreatedBy        string
	CreatedTimestamp time.Time
}
