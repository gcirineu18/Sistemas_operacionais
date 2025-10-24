package main

import (
	"log"
	"time"

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
	Begin int `json:begin`
	Duration int `json:duration`
	Priority int `json:priority`
}



func main(){

	r:= gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
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
			return
		}

		if body.Quantum <= 0 {
			c.JSON(400, gin.H{"error": "Quantum deve ser maior que 0"})
			return
		}

		for i, p := range body.Input {
			if p.Duration <= 0 {
				c.JSON(400, gin.H{
					"error":   "Duração inválida em um dos processos",
					"process": i + 1,
				})
				return
			}
			if p.Priority < 0 {
				c.JSON(400, gin.H{
					"error":   "Prioridade inválida em um dos processos",
					"process": i + 1,
				})
				return
			}
		}


		log.Printf("Algoritmo: %s, Quantum: %d, Aging: %d", body.Alg, body.Quantum, body.Aging)

		tempoMedioVida, tempoMedioEspera,trocasContexto ,diagramaTempo:= processScheduler(body)
		
		c.JSON(200, gin.H{
			"tempoMedioVida":tempoMedioVida, 
			"tempoMedioEspera":tempoMedioEspera,
			"trocasContexto": trocasContexto,
			"diagramaTempo": diagramaTempo,
		})
	
	})

	r.Run(":8080")
}
