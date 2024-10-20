package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "html/template"
	"bytes"
	"encoding/json"
	"os"
	"log"
	"fmt"
	"strings"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

type FormData struct {
    AttackerName     string `form:"attacker_name"`
    AttackerCountry  string `form:"attacker_country"`
	AttackerServer   int    `form:"attacker_server"`
    DefenderName     string `form:"defender_name"`
    DefenderCountry  string `form:"defender_country"`
	DefenderServer   int    `form:"defender_server"`
    LifeMode         string `form:"life_mode"`
    Simulates        int    `form:"simulates"`
}

type Response struct {
	WinChance float64 `json:"win-chance"`
	LoseChance float64 `json:"lose-chance"`
	DrawChance float64 `json:"draw-chance"`
	Details Details `json:"details"`
}

type Details struct {
	Fights int `json:"fights"`
	Wins int `json:"wins"`
	Loses int `json:"loses"`
	Draws int `json:"draws"`
}


func main() {
    r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

    // Load HTML templates
    r.SetFuncMap(template.FuncMap{
        "safe": func(str string) template.HTML {
            return template.HTML(str)
        },
    })
    r.LoadHTMLGlob("templates/*")

	r.Static("/assets", "./assets") 

    // Serve the form
    r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		attackerName := session.Get("attacker_name")
		defenderName := session.Get("defender_name")
		attackerServer := session.Get("attacker_server")
		defenderServer := session.Get("defender_server")
        c.HTML(http.StatusOK, "form.html", gin.H{
			"AttackerName": attackerName,
			"DefenderName": defenderName,
			"AttackerServer": attackerServer,
			"DefenderServer": defenderServer,
		})
    })

	r.GET("/turma", func(c *gin.Context) {
		session := sessions.Default(c)
		attackerName := session.Get("attacker_name")
		defenderName := session.Get("defender_name")
		attackerServer := session.Get("attacker_server")
		defenderServer := session.Get("defender_server")
		c.HTML(http.StatusOK, "turma.html", gin.H{
			"AttackerName": attackerName,
			"DefenderName": defenderName,
			"AttackerServer": attackerServer,
			"DefenderServer": defenderServer,
		})
	})
    // Handle form submission
    r.POST("/generate", func(c *gin.Context) {
		session := sessions.Default(c)

		response := simulateBattle("arena", c, session)
		parsedResponse := parseResponse(response)
		c.HTML(http.StatusOK, "response.html", gin.H{
			"winChance": parsedResponse.WinChance,
			"loseChance": parsedResponse.LoseChance,
			"drawChance": parsedResponse.DrawChance,
		})
    })

	r.POST("/generate-turma", func(c *gin.Context) {
		session := sessions.Default(c)

		response := simulateBattle("turma", c, session)
		parsedResponse := parseResponse(response)
		c.HTML(http.StatusOK, "turma-response.html", gin.H{
			"winChance": parsedResponse.WinChance,
			"loseChance": parsedResponse.LoseChance,
			"drawChance": parsedResponse.DrawChance,
		})
    })

	r.LoadHTMLGlob("templates/*")

    r.Run(":8000") 
}


func parseResponse(response string) *Response {
	
	trimmedResponse := strings.TrimPrefix(response, "null")
	log.Println("trimmedResponse:", trimmedResponse)

	var apiResponse Response
	err := json.Unmarshal([]byte(trimmedResponse), &apiResponse)
	if err != nil {
		log.Println("Error unmarshaling json:", err)
	}

	fmt.Printf("Win Chance: %.2f", apiResponse.WinChance)
	fmt.Printf("Lose Chance: %.2f", apiResponse.LoseChance)
	fmt.Printf("Draw Chance: %.2f", apiResponse.DrawChance)

	return &apiResponse
}

func simulateBattle(mode string, c *gin.Context, s sessions.Session) (response string) {
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
				"server":  formData.AttackerServer,
            },
            "defender": map[string]interface{}{
                "country": formData.DefenderCountry,
                "name":    formData.DefenderName,
				"server":  formData.DefenderServer,
            },
            "options": map[string]interface{}{
                "life-mode": formData.LifeMode,
                "simulates": "1000",
            },
        }
		
		s.Set("attacker_name", formData.AttackerName)
		s.Set("defender_name", formData.DefenderName)
		s.Set("attacker_server", formData.AttackerServer)
		s.Set("defender_server", formData.DefenderServer)
		s.Save()

		log.Println("Request:")
		log.Println(formData.AttackerName)
		log.Println(formData.AttackerServer)
		log.Println(formData.DefenderName)
		log.Println(formData.DefenderServer)

		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error marshalling JSON: %s", err.Error())
			return
		}
		
		url := "http://localhost:8080"
	    envUrl, exists := os.LookupEnv("BACKEND_URL")
	    if exists {
	        url = envUrl
	    }
		path := "/" + mode + ".php"
		resp, err := http.Post(url+path, "application/json", bytes.NewBuffer(jsonBytes))
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
		response = responseBody.String()
		log.Println("Full response: ", response)

		return response
	}
