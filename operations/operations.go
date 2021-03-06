package operations

import (
	"fmt"
	"math"

	structs "github.com/ivanaghemo/GoML/structs"
)

var (
	kenobi          = structs.Point{X: -500, Y: -200}
	skywalker       = structs.Point{X: 100, Y: -100}
	sato            = structs.Point{X: 500, Y: 100}
	satelites       = iniResistencia()
	onlineSatelites = []structs.Satellite{}
)

// InicioResistencia: inicializamos los satétiles que nos permiten iniciar la trilateración
func iniResistencia() []structs.Satellite {
	return []structs.Satellite{
		nuevoSatelite("kenobi", 0, kenobi, nil),
		nuevoSatelite("skywalker", 0, skywalker, nil),
		nuevoSatelite("sato", 0, sato, nil),
	}
}

//Crea un satétile
func nuevoSatelite(n string, d float32, p structs.Point, m []string) structs.Satellite {
	return structs.Satellite{
		Name:     n,
		Distance: d,
		Location: p,
		Message:  m,
	}
}

func IniSatelites(req structs.TopSecretRequest) ([]structs.Satellite, error) {

	satellites := []structs.Satellite{}

	for _, s := range req.Satellites {

		sat, err := buscarSateliteXNombre(s.Name, satelites)

		if err != nil {
			return nil, err
		}

		sat.Distance = s.Distance
		sat.Message = s.Message
		satellites = append(satellites, sat)
	}

	onlineSatelites = satellites

	return satellites, nil
}

func buscarSateliteXNombre(n string, satss []structs.Satellite) (structs.Satellite, error) {

	for _, s := range satss {
		if s.Name == n {
			return s, nil
		}
	}

	err := fmt.Errorf("Satétile %s no existe.", n)

	return structs.Satellite{}, err
}

func ObtenerMensajeYDistancia(satss []structs.Satellite) ([]float32, [][]string) {

	var (
		distances = []float32{}
		msgs      = [][]string{}
	)

	for _, s := range satss {

		distances = append(distances, s.Distance)
		msgs = append(msgs, s.Message)
	}

	return distances, msgs
}

//func (s *Satellites) GetLocation(d1, d2, d3 float64) (coordinates Point) {
func GetLocation(distances ...float32) (x, y float32) {

	x, y, err := trilateracion(onlineSatelites)

	if err != nil {
		return 0, 0
	}
	return x, y
}

func normalize(point structs.Point) float32 {
	return float32(math.Pow(math.Pow(float64(point.X), 2)+math.Pow(float64(point.Y), 2), .5))
}

func GetMessage(messages ...[]string) (msg string) {

	if len(messages) < 1 {
		return ""
	}

	l := len(messages[0])

	for _, m := range messages {
		if len(m) != l {
			return ""
		}
	}

	m1 := messages[0]
	m2 := messages[1]
	m3 := messages[2]

	finalMessage := ""

	for i := 0; i < l; i++ {

		wordlist := []string{}

		wordlist = append(wordlist, m1[i])
		wordlist = append(wordlist, m2[i])
		wordlist = append(wordlist, m3[i])

		word := selectWord(wordlist)

		if word == "" {
			return ""
		}

		if finalMessage == "" {
			finalMessage = word
		} else {
			finalMessage = fmt.Sprintf("%s %s", finalMessage, word)
		}
	}
	return finalMessage
}

func selectWord(words []string) string {

	l := []WordType{}

	for _, w := range words {

		if w == "" {
			continue
		}
		val, index := findWord(w, l)

		if val {
			wt := l[index]
			wt.Count = wt.Count + 1
			l[index] = wt
		} else {
			l = append(l, WordType{Word: w, Count: 1})
		}
	}

	return getCommonWord(l)
}

// WordType struct used to keep the frecuency of the word
type WordType struct {
	Word  string
	Count int16
}

func findWord(w string, words []WordType) (bool, int) {

	for i, ww := range words {
		if ww.Word == w {
			return true, i
		}
	}

	return false, -1
}

// GetCommonWord function to retrieve the most common word
func getCommonWord(words []WordType) string {

	if len(words) < 1 {
		return ""
	}

	var (
		index  = -1
		bigOne = int16(-1)
	)

	for i, w := range words {
		if w.Count > bigOne {
			index = i
			bigOne = w.Count
		}
	}

	return words[index].Word
}

func trilateracion(satelites []structs.Satellite) (x, y float32, err error) {
	/*
		if !validarSatelites(satelites) {
			err := fmt.Errorf("Satellites not in range")

			return 0, 0, err
		}
	*/
	//Coordenadas de los satélites y distancias a la nave
	//kenobi A
	x1, y1, d1 := satelites[0].Location.X, satelites[0].Location.Y, satelites[0].Distance
	//skywalker B
	x2, y2, d2 := satelites[1].Location.X, satelites[1].Location.Y, satelites[1].Distance
	//sato C
	x3, y3, d3 := satelites[2].Location.X, satelites[2].Location.Y, satelites[2].Distance

	//Distancia de kenobi a skywalker
	ABDistance := float32(math.Pow(math.Pow(float64(x2-x1), 2)+math.Pow(float64(y2-y1), 2), 0.5))

	ex := structs.Point{
		X: (x2 - x1) / ABDistance,
		Y: (y2 - y1) / ABDistance,
	}
	aux := structs.Point{
		X: x3 - x1,
		Y: y3 - y1,
	}

	//signed magnitude of the x component
	i := ex.X*aux.X + ex.Y*aux.Y

	//the unit vector in the y direction.
	aux2 := structs.Point{
		X: x3 - x1 - i*ex.X,
		Y: y3 - y1 - i*ex.Y,
	}
	ey := structs.Point{
		X: aux2.X / normalize(aux2),
		Y: aux2.Y / normalize(aux2),
	}
	//the signed magnitude of the y component
	j := ey.X*aux.X + ey.Y*aux.Y

	//coordinates
	x = float32(math.Pow(float64(d1), 2)-math.Pow(float64(d2), 2)+math.Pow(float64(ABDistance), 2)) / (2 * ABDistance)
	y = float32(math.Pow(float64(d1), 2)-math.Pow(float64(d3), 2)+math.Pow(float64(i), 2)+math.Pow(float64(j), 2))/(2*j) - (i*x)/j

	return x, y, nil
}

func ObtenerUnSatelite(nombreSatelite string) (structs.Satellite, error) {
	return BuscarSatelitePorNombre(nombreSatelite, onlineSatelites)
}

func MofificarSatelite(sateliteModificado structs.Satellite) []structs.Satellite {
	nuevoSateltite := []structs.Satellite{}

	for _, s := range onlineSatelites {
		if s.Name == sateliteModificado.Name {
			s = sateliteModificado
		}

		nuevoSateltite = append(nuevoSateltite, s)
	}

	onlineSatelites = nuevoSateltite

	return onlineSatelites
}

func BuscarSatelitePorNombre(nombreSatelite string, satelites []structs.Satellite) (structs.Satellite, error) {

	for _, s := range satelites {
		if s.Name == nombreSatelite {
			return s, nil
		}
	}

	err := fmt.Errorf("Satélite %s no existe", nombreSatelite)

	return structs.Satellite{}, err
}

func GetSatelitesOnlineSplit() []structs.Satellite {

	return onlineSatelites
}

func GetLocationSplit(distances ...float32) (x, y float32) {

	x, y, err := trilateracion(onlineSatelites)

	if err != nil {
		return 0, 0
	}

	return x, y
}
