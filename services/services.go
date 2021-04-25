package services

import (
	"fmt"
	"log"

	operations "github.com/ivanaghemo/GoML/operations"
	structs "github.com/ivanaghemo/GoML/structs"
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
