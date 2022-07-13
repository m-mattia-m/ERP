package customers

import (
	"erp/db"
	"erp/users"
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Main(r *gin.RouterGroup) {

	db.CreateCustomersTable()

	r.POST("/createCustomer", users.BasicAuth, createCustomer)
	r.POST("/editCustomer/:id", users.BasicAuth, editCustomer)
	r.POST("/editCustomer", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /editCustomer/id")
	})
	r.GET("/deleteCustomer/:id", users.BasicAuth, deleteCustomer)
	r.GET("/deleteCustomer", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /deleteCustomer/id")
	})
	r.GET("/getCustomer/:id", users.BasicAuth, getCustomer)
	r.GET("/getCustomer", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getCustomer/id")
	})
	r.GET("/getCustomers", users.BasicAuth, getCustomers)
}

func createCustomer(c *gin.Context) {
	currentUserID := c.Request.Header.Get("userid")

	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	street := c.PostForm("street")
	streetNr := c.PostForm("streetNr")
	postcode := c.PostForm("postcode")
	city := c.PostForm("city")
	email := c.PostForm("email")
	telefon := c.PostForm("telefon")
	currentCustomer, msg := newCustomer(firstname, lastname, street, streetNr, postcode, city, email, telefon, currentUserID)
	if msg == "" {
		saveCustomerOnDB(currentCustomer)
		c.JSON(200, currentCustomer)
	} else {
		c.JSON(400, gin.H{"error": msg})
	}
}

func editCustomer(c *gin.Context) {
	currentUserID := c.Request.Header.Get("userid")
	id := c.Param("id")
	var customer Customer = getCustomersFromDBById(id)

	if customer.Id == id {
		customer.Firstname = c.PostForm("firstname")
		customer.Lastname = c.PostForm("lastname")
		customer.Street = c.PostForm("street")
		customer.StreetNr = c.PostForm("streetNr")
		customer.Postcode = c.PostForm("postcode")
		customer.City = c.PostForm("city")
		customer.Email = c.PostForm("email")
		customer.Telefon = c.PostForm("telefon")
		customer.CreatedBy = currentUserID
		updateCustomerOnDB(customer)
		c.JSON(200, customer)
	} else {
		c.JSON(400, gin.H{"error": "No customer found with the id: " + id})
	}
}

func getCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer Customer = getCustomersFromDBById(id)
	if customer.Id == id {
		c.JSON(200, customer)
	} else {
		c.JSON(400, gin.H{"error": "No customer found with the id: " + id})
	}
}

func getCustomers(c *gin.Context) {
	var customers []Customer = getCustomersFromDB()
	if customers != nil {
		c.JSON(200, customers)
	} else {
		c.JSON(400, gin.H{"error": "No users found"})
	}
}

func deleteCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer Customer = getCustomersFromDBById(id)

	if customer.Id == id {
		deleteCustomerFromDB(customer)
		c.JSON(200, gin.H{"message": "delete customer with the id: " + id})
	} else {
		c.JSON(400, gin.H{"error": "No customer found with the id: " + id})
	}
}

func newCustomer(firstname string, lastname string, street string, streetNr string, postcode string, city string, email string, telefon string, currentUserID string) (Customer, string) {
	var customers []Customer = getCustomersFromDB()
	i := sort.Search(len(customers), func(i int) bool {
		return customers[i].Firstname == firstname && customers[i].Lastname == lastname && customers[i].Street == street && customers[i].StreetNr == streetNr && customers[i].Postcode == postcode && customers[i].City == city
	})
	if i < len(customers) && (customers[i].Firstname == firstname && customers[i].Lastname == lastname && customers[i].Street == street && customers[i].StreetNr == streetNr && customers[i].Postcode == postcode && customers[i].City == city) {
		customer := Customer{}
		return customer, "Customer already exists"
	}
	i = sort.Search(len(customers), func(i int) bool { return email <= customers[i].Email })
	if i < len(customers) && customers[i].Email == email {
		customer := Customer{}
		return customer, "email already exists"
	}
	i = sort.Search(len(customers), func(i int) bool { return telefon <= customers[i].Telefon })
	if i < len(customers) && customers[i].Telefon == telefon {
		customer := Customer{}
		return customer, "telefon already exists"
	}

	var currentCustomer = new(Customer)
	currentCustomer.Id = uuid.New().String()
	currentCustomer.Firstname = firstname
	currentCustomer.Lastname = lastname
	currentCustomer.Street = street
	currentCustomer.StreetNr = streetNr
	currentCustomer.Postcode = postcode
	currentCustomer.City = city
	currentCustomer.Email = email
	currentCustomer.Telefon = telefon
	currentCustomer.CreatedBy = currentUserID

	return *currentCustomer, ""
}

func getCustomersFromDB() []Customer {
	var customers []Customer
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Firstname`, `Lastname`, `Street`, `StreetNr`, `Postcode`, `City`, `Email`, `Telefon`, `CreatedBy` FROM `customer`;")
	if err != nil {
		fmt.Println("[DB]: Can't Select customers from DB \t-->\t" + err.Error())
	}
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.Id, &customer.Firstname, &customer.Lastname, &customer.Street, &customer.StreetNr, &customer.Postcode, &customer.City, &customer.Email, &customer.Telefon, &customer.CreatedBy)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to Customer-Object (CustomerId: %v) \t-->\t"+err.Error(), customer.Id)
		}
		customers = append(customers, customer)
	}
	rows.Close()
	return customers
}

func saveCustomerOnDB(customer Customer) {
	var query string = "INSERT INTO `customer` (`Id`, `Firstname`, `Lastname`, `Street`, `StreetNr`, `Postcode`, `City`, `Email`, `Telefon`, `CreatedBy`) VALUES ('" + customer.Id + "', '" + customer.Firstname + "', '" + customer.Lastname + "', '" + customer.Street + "', '" + customer.StreetNr + "', '" + customer.Postcode + "', '" + customer.City + "', '" + customer.Email + "', '" + customer.Telefon + "', '" + customer.CreatedBy + "');"
	db.RunSqlQueryWithoutReturn(query)
}

func updateCustomerOnDB(customer Customer) {
	var query string = "UPDATE `customer` SET `Firstname` = '" + customer.Firstname + "', `Lastname` = '" + customer.Lastname + "', `Street` = '" + customer.Street + "', `StreetNr` = '" + customer.StreetNr + "', `Postcode` = '" + customer.Postcode + "', `City` = '" + customer.City + "', `Email` = '" + customer.Email + "', `Telefon` = '" + customer.Telefon + "', `CreatedBy` = '" + customer.CreatedBy + "' WHERE `Id` = '" + customer.Id + "';"
	db.RunSqlQueryWithoutReturn(query)
}

func deleteCustomerFromDB(customer Customer) {
	var query string = "DELETE FROM `customer` WHERE `Id` = '" + customer.Id + "';"
	db.RunSqlQueryWithoutReturn(query)
}

func getCustomersFromDBById(id string) Customer {
	var customers []Customer
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Firstname`, `Lastname`, `Street`, `StreetNr`, `Postcode`, `City`, `Email`, `Telefon`, `CreatedBy` FROM `customer` WHERE Id='" + id + "';")
	if err != nil {
		fmt.Println("[DB]: Can't Select customer by Id from DB \t-->\t" + err.Error())
	}

	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.Id, &customer.Firstname, &customer.Lastname, &customer.Street, &customer.StreetNr, &customer.Postcode, &customer.City, &customer.Email, &customer.Telefon, &customer.CreatedBy)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to User-Object (userId: %v) \t-->\t"+err.Error(), customer.Id)
		}
		customers = append(customers, customer)
	}
	rows.Close()
	return customers[0]
}

type Customer struct {
	Id        string
	Firstname string
	Lastname  string
	Street    string
	StreetNr  string
	Postcode  string
	City      string
	Email     string
	Telefon   string
	CreatedBy string
}
