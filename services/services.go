package services

import (
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
	operations "github.com/ivanaghemo/GoML/operations"
	structs "github.com/ivanaghemo/GoML/structs"
)

var (
	c        = cache.New(5*time.Minute, 10*time.Minute)
	cachekey = "cache/topsecret_split"
)

func TopSecretService(req structs.TopSecretRequest) (*structs.TopSecretResponse, error) {

	log.Println("[START] topSecretService")

	satelites, err := operations.IniSatelites(req)

	if err != nil {

		msg := fmt.Errorf("No se puede ejecutar la misión %s", err.Error())

		return nil, msg
	}

	distances, msgs := operations.ObtenerMensajeYDistancia(satelites)

	x, y := operations.GetLocation(distances...)

	if x == 0 && y == 0 {
		msg := fmt.Errorf("No se puede ubicar la nave o los satétiles no están en línea.")

		return nil, msg
	}

	message := operations.GetMessage(msgs...)

	if message == "" {
		msg := fmt.Errorf("Mensaje incompleto")

		return nil, msg
	}

	res := structs.TopSecretResponse{
		Position: structs.Point{
			X: x,
			Y: y,
		},
		Message: message,
	}

	log.Println("[FINAL] topSecretService")
	return &res, nil
}

func TopSecretSplitPostService(nombreSatelite string, req structs.TopSecretRequestSplit) (*structs.TopSecretSplitResponse, error) {

	satelite, err := operations.ObtenerUnSatelite(nombreSatelite)

	if err != nil {
		msg := fmt.Errorf("No se puede ejecutar la misión %s", err.Error())
		return nil, msg
	}

	satelite.Distance = req.Distance
	satelite.Message = req.Message

	operations.MofificarSatelite(satelite)

	res := structs.TopSecretSplitResponse{
		Message: fmt.Sprintf("Ok"),
	}

	c.Delete(cachekey)

	return &res, nil
}

func TopSecretSplitGetService() (*structs.TopSecretResponse, error) {
	//cache
	if x, found := c.Get(cachekey); found {
		ress := x.(*structs.TopSecretResponse)
		log.Println("Cache...")
		return ress, nil
	}
	satelites := operations.GetSatelitesOnlineSplit()

	if len(satelites) < 3 {
		msg := fmt.Errorf("No se ha podido conectar con los satélites")
		return nil, msg
	}

	distances, msgs := operations.ObtenerMensajeYDistancia(satelites)

	x, y := operations.GetLocationSplit(distances...)

	if x == 0 && y == 0 {
		msg := fmt.Errorf("No se puede ubicar la nave o los satétiles no están en línea.")
		return nil, msg
	}

	message := operations.GetMessage(msgs...)

	if message == "" {
		msg := fmt.Errorf("Mensaje incompleto")

		return nil, msg
	}

	res := structs.TopSecretResponse{
		Position: structs.Point{
			X: x,
			Y: y,
		},
		Message: message,
	}

	c.Set(cachekey, &res, cache.DefaultExpiration)

	return &res, nil
}
