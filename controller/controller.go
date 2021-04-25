package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	services "github.com/ivanaghemo/GoML/services"
	structs "github.com/ivanaghemo/GoML/structs"
)

func PostTopSecrets(c *gin.Context) {
	var msgSecret structs.TopSecretRequest

	//Obtengo json con distancia y mensajes
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//fmt.Fprintf(w, "No se pueden obtener los mensajes")
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "No se pueden obtener los mensaje.",
		})
		return
	}

	err2 := json.Unmarshal(reqBody, &msgSecret)
	if err2 != nil {
		//fmt.Fprintf(w, "El mensaje es ilegible.")
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "El mensaje es ilegible.",
		})
		return

	}

	//Invoco servicio que obtiene el punto de coordenadas y mensaje
	res, err := services.TopSecretService(msgSecret)

	if err != nil {
		log.Println("[topSecretService] No se pueden obtener los datos." + err.Error())
		//fmt.Fprintf(w, "No se puede obtener datos: "+err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "No se puede obtener datos: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
	return

}

func PostTopSecretSplit(c *gin.Context) {

	var (
		secretRequest  = structs.TopSecretRequestSplit{}
		nombreSatelite = c.Param("satellite_name")
	)

	if err := c.Bind(&secretRequest); err != nil {
		log.Printf("El mensaje TopSecretRequestSplit no se puede leer: %s", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "El mensaje TopSecretRequestSplit no se puede leer: " + err.Error(),
		})
		return
	}

	log.Println("[START] PostTopSecretSplit con..", secretRequest)

	res, err := services.TopSecretSplitPostService(nombreSatelite, secretRequest)

	if err != nil {
		log.Println("[PostTopSecretSplit] No se pueden obtener los datos." + err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "No se puede obtener datos: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
	return

}

func GetTopSecretSplit(c *gin.Context) {

	log.Println("[START] GetTopSecretSplit..")

	res, err := services.TopSecretSplitGetService()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "No se puede obtener datos: " + err.Error(),
		})
	}

	c.JSON(http.StatusOK, res)

}
