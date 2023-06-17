// Simple API for the challenge located at https://github.com/fetch-rewards/receipt-processor-challenge
// First time using Go!


package main
import (
	"fmt"

	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"math"
	"strings"
	"regexp"
	"time"
	"strconv"
) 

type receipt struct {
	ID string `json:"id"`
	Retailer string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total string `json:"total"` 
	Items []item `json:"items`
}

type item struct {
	ShortDescription string `json:"shortDescription"`
	Price string `json:"price"`
}



func main() {
	router := gin.Default()
	router.POST("/receipts/process", postReceipt)
	router.GET("/receipts", getReceipts)
	router.GET("/receipts/:id/points", getPoints)

	router.Run("localhost:8080")
}

var receipts = []receipt {}


// generates ID for receipt and adds it to memory
func postReceipt(c *gin.Context) {
	var newReceipt receipt

	if err := c.BindJSON(&newReceipt); err != nil {
		fmt.Println(err)
	return
	}
	newReceipt.ID = uuid.New().String()
	
	receipts = append(receipts, newReceipt)
	c.IndentedJSON(http.StatusCreated, gin.H{"id": newReceipt.ID})
}


// get all receipts
func getReceipts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, receipts)
}


// gets the number of points the receipt with the given ID, if it exists
func getPoints(c *gin.Context) {
	id := c.Param("id")

	for _, receipt := range receipts {
		if receipt.ID == id {
			points := calculatePoints(receipt)
			c.IndentedJSON(http.StatusOK, gin.H{"points": points})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Receipt not found!"})
}


// calcuates the number of points for the given receipt.
func calculatePoints(receipt receipt) string {
	points := 0

	// One point for every alphanumeric character in the retailer name. (probably a better way to do this?)
	for _, char := range receipt.Retailer {
		if regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(string(char)) {
			points += 1
		}
	}
	
	if total, err := strconv.ParseFloat(receipt.Total, 32); err == nil {
		// 50 points if the total is a round dollar amount with no cents.
		if math.Mod(total, 1) == 0 {
			points += 50
		}

		// 25 points if the total is a multiple of 0.25.
		if math.Mod(total, 0.25) == 0 {
			points += 25
		}
	}

	// 5 points for every two items on the receipt.	
	points += (len(receipt.Items) / 2) * 5


	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	for _, item := range receipt.Items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if math.Mod(float64(len(trimmed)), 3) == 0 {
			if price, err := strconv.ParseFloat(item.Price, 32); err == nil {
				price *= 0.2
				points += int(math.Ceil(price))
			}
		}
	}

	// 6 points if the day in the purchase date is odd.
	if date, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if math.Mod(float64(date.Day()), 2) == 1 {
			points += 6
		}
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		hour, minute, _ := purchaseTime.Clock()
		if (hour == 14 && minute > 0) || (hour == 15) {
			points += 10
		}
		
	}

	return strconv.FormatInt(int64(points), 10)

}