package main

import (
    "github.com/gin-gonic/gin"
    "fmt"
    "github.com/gocolly/colly/v2"
)

func main() {
    r := gin.Default()

    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204) // No content for OPTIONS requests
            return
        }
        c.Next()
    })

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, Gin!",
        })
    })

    // New route to scrape Amazon product info by ASIN
    r.GET("/amazon-product", func(c *gin.Context) {
        asin := c.Query("ASIN")
        // Scrape Amazon product info
        productInfo, err := scrapeAmazonProductInfo(asin)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        c.JSON(200, productInfo)
    })

    r.Run(":8080") // Run the web server on port 8080
}

func scrapeAmazonProductInfo(asin string) (map[string]string, error) {
    c := colly.NewCollector()

    // Define the URL of the Amazon product page
    url := "https://www.amazon.com/dp/" + asin

    productInfo := make(map[string]string)

    // Define the selectors to extract product information
    c.OnHTML("#titleSection", func(e *colly.HTMLElement) {
        productInfo["title"] = e.Text
    })

    c.OnHTML("#productDescription", func(e *colly.HTMLElement) {
        productInfo["description"] = e.Text
    })

    // Visit the Amazon product page and scrape the information
    err := c.Visit(url)
    if err != nil {
        return nil, err
    }

    fmt.Println("Description" + productInfo["description"])
    fmt.Println("Title" + productInfo["title"])    

    return productInfo, nil
}
