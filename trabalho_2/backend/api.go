package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

)

type ContextBody struct{
	Alg string `json:"alg"`
	Quantum int `json:"quantum"`
	Aging int `json:"aging"`
	Input []Processes `json:"input"`
}

type Processes struct{
	Begin int  `json:"begin"`
	Duration int `json:"duration"`
	Priority int `json:"priority"`
}



func main(){
	godotenv.Load()
	frontendURL := os.Getenv("FRONTEND_URL")
	r:= gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{frontendURL},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Origin", "Content-Type", "OPTIONS"},
		ExposeHeaders: []string{"Content-Lenght"},
		AllowCredentials: false,
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/processes", func(c *gin.Context){
		var body ContextBody 

		// Desesserilizar JSON recebido no corpo da requisição
		if err := c.ShouldBindJSON(&body); err!= nil{
			c.JSON(400, gin.H{"error": err.Error()})
			log.Print("error: "  + err.Error())
			return
		}

		fmt.Println(body.Input)

		if body.Quantum <= 0 {
			c.JSON(400, gin.H{"error": "Quantum deve ser maior que 0"})
			return
		}

		for i, p := range body.Input {

			if p.Begin == 0 && p.Duration == 0 && p.Priority == 0{
				c.JSON(400, gin.H{
					"error": "Entrada inválida, por favor, tente novamente.",			
				})
				return
			}

			if  p.Duration <= 0 {
				c.JSON(400, gin.H{
					"error":  fmt.Sprintf("Duração inválida no processo %d", i+1),
				})
				return
			} else if p.Priority < 0 {
				c.JSON(400, gin.H{
					"error":  fmt.Sprintf("Prioridade inválida no processo %d", i+1),
				})
				return
			}  else if p.Begin < 0 {
				c.JSON(400, gin.H{
					"error":  fmt.Sprintf("Tempo de início inválido no processo %d", i+1),
				})
				return
			}	
		}


		log.Printf("Algoritmo: %s, Quantum: %d, Aging: %d", body.Alg, body.Quantum, body.Aging)

		tempoMedioVida, tempoMedioEspera,trocasContexto ,diagramaTempo, ordemProcessos:= processScheduler(body)
		
		c.JSON(200, gin.H{
			"tempoMedioVida":tempoMedioVida, 
			"tempoMedioEspera":tempoMedioEspera,
			"trocasContexto": trocasContexto,
			"diagramaTempo": diagramaTempo,
			"ordemProcessos": ordemProcessos,
		})
	
	})

	r.Run(":8081")
}
