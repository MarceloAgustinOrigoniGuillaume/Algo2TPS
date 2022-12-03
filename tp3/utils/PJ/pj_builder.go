package pj

import "fmt"
import "os"
import "strconv"
import "tp3/utils"

const PJ_POINT = "%s,%s,%s\n"
const PJ_ARISTA = "%s,%s,%s\n"

type PJBuilder struct {
	file     *os.File
	started  bool
	ciudades int
	aristas  int
}

func CrearPJ(outFile string) (*PJBuilder, error) {
	archivo, err := utils.AbrirOCrearArchivo(outFile)

	if err != nil {
		return nil, err
	}

	builder := new(PJBuilder)

	builder.file = archivo

	return builder, nil
}
func (builder *PJBuilder) StartPJ(ciudades, aristas int) {
	if builder.started {
		return
	}
	builder.ciudades = ciudades
	builder.aristas = aristas
	builder.started = true
	builder.file.WriteString(fmt.Sprintf("%d\n", ciudades))
}

func (builder *PJBuilder) ClosePJ() {
	if !builder.started {
		return
	}
	builder.Close()
}

func (builder *PJBuilder) Close() {
	builder.file.Close()
}

func (builder *PJBuilder) AddCity(name, latitud, longitud string) {
	if !builder.started {
		return
	}

	if builder.ciudades <= 0 {
		fmt.Printf("\nWarning: Writing pj tried to break format, by writing more cities than said\n")
		return
	}

	builder.ciudades--
	builder.file.WriteString(fmt.Sprintf(PJ_POINT, name, latitud, longitud))

	if builder.ciudades == 0 {
		builder.file.WriteString(fmt.Sprintf("%d\n", builder.aristas))
	}
}

func (builder *PJBuilder) AddArista(desde string, hasta string, peso int) {
	if !builder.started {
		return
	}

	if builder.ciudades > 0 {
		fmt.Printf("\nWarning: Writing pj tried to break format, by writing arista before cities\n")
		return
	}

	if builder.aristas <= 0 {
		fmt.Printf("\nWarning: Writing pj tried to break format, by writing more aristas than said\n")
		return
	}

	builder.aristas--
	builder.file.WriteString(fmt.Sprintf(PJ_ARISTA, desde, hasta, strconv.Itoa(peso)))
}
