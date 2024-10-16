package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "html/template"
	"bytes"
	"encoding/json"
	"os"
)

type FormData struct {
    AttackerName     string `form:"attacker_name"`
    AttackerCountry  string `form:"attacker_country"`
    DefenderName     string `form:"defender_name"`
    DefenderCountry  string `form:"defender_country"`
    LifeMode         string `form:"life_mode"`
    Simulates        int    `form:"simulates"`
	ServerNumber     int    `form:"server_number"`
}


func main() {
	var DEBUG = false 
    r := gin.Default()

    // Load HTML templates
    r.SetFuncMap(template.FuncMap{
        "safe": func(str string) template.HTML {
            return template.HTML(str)
        },
    })
    r.LoadHTMLGlob("templates/*")

	r.GET("static/style.css", func(c *gin.Context) {
		c.File("static/style.css")
	})

    // Serve the form
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "form.html", nil)
    })

    // Handle form submission
    r.POST("/generate", func(c *gin.Context) {
        var formData FormData
        if err := c.ShouldBind(&formData); err != nil {
            c.String(http.StatusBadRequest, "Error binding data: %s", err.Error())
            return
        }

        // Construct the JSON structure
        jsonData := map[string]interface{}{
            "attacker": map[string]interface{}{
                "country": formData.AttackerCountry,
                "name":    formData.AttackerName,
				"server":  formData.ServerNumber,
            },
            "defender": map[string]interface{}{
                "country": formData.DefenderCountry,
                "name":    formData.DefenderName,
				"server":  formData.ServerNumber,
            },
            "options": map[string]interface{}{
                "life-mode": formData.LifeMode,
                "simulates": "1000",
            },
        }

		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error marshalling JSON: %s", err.Error())
			return
		}

		if DEBUG  {
			resp, err := http.Post("http://localhost:8081/api.php", "application/json", bytes.NewBuffer(jsonBytes))
			if err != nil {
				c.String(http.StatusInternalServerError, "Error calling API: %s", err.Error())
				return
			}
		defer resp.Body.Close()

		var responseBody bytes.Buffer
		_, err = responseBody.ReadFrom(resp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading response body: %s", err.Error())
			return
		}

		c.HTML(http.StatusOK, "response.html", gin.H{
			"status": resp.Status,
			"body":   responseBody.String(),
		})
		} else {
		url := os.Getenv("BACKEND_URL")
		resp, err := http.Post(url+"/api.php", "application/json", bytes.NewBuffer(jsonBytes))
		if err != nil {
			c.String(http.StatusInternalServerError, "Error calling API: %s", err.Error())
			return
			}
		defer resp.Body.Close()

		var responseBody bytes.Buffer
		_, err = responseBody.ReadFrom(resp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading response body: %s", err.Error())
			return
		}

		c.HTML(http.StatusOK, "response.html", gin.H{
			"status": resp.Status,
			"body":   responseBody.String(),
		})
		}
    })

	r.LoadHTMLGlob("templates/*")

    r.Run(":8000") 
}

