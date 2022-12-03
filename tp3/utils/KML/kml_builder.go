package kml

import "fmt"
import "os"
import "tp3/utils"
import "tp3/cola"

const KML = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<kml xmlns=\"http://earth.google.com/kml/2.1\">\n\t<Document>\n\t\t<name>%s</name>\n%s\n\t</Document>\n</kml>\n"

const KML_START = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<kml xmlns=\"http://earth.google.com/kml/2.1\">\n\t<Document>\n\t\t<name>%s</name>\n"
const KML_END = "\n\t</Document>\n</kml>\n"

const KML_POINT = "\t\t<Placemark>\n\t\t\t<name>%s</name>\n\t\t\t<Point>\n\t\t\t\t<coordinates>%s, %s</coordinates>\n\t\t\t</Point>\n\t\t</Placemark>\n"

const KML_LINE = "\t\t<Placemark>\n\t\t\t<name>%s</name>\n\t\t\t<LineString>\n\t\t\t\t<coordinates>%s, %s %s, %s</coordinates>\n\t\t\t</LineString>\n\t\t</Placemark>\n"

func FORMAT_KML(title string, elements string) string {
	return fmt.Sprintf(KML, title, elements)
}

func POINT_KML(name, latitud, longitud string) string {
	return fmt.Sprintf(KML_POINT, name, latitud, longitud)
}

func LINE_KML(name, fromLatitud, fromLongitud, toLatitud, toLongitud string) string {
	return fmt.Sprintf(KML_LINE, name, fromLatitud, fromLongitud, toLatitud, toLongitud)
}

type KMLBuilder struct {
	file         *os.File
	started      bool
	colaDeLineas cola.Cola[func()]
}

func CrearKML(outFile string) (*KMLBuilder, error) {
	archivo, err := utils.AbrirOCrearArchivo(outFile)

	if err != nil {
		return nil, err
	}

	builder := new(KMLBuilder)

	builder.file = archivo
	return builder, nil
}
func (builder *KMLBuilder) StartKML(title string) {
	if builder.started {
		return
	}
	builder.started = true
	builder.file.WriteString(fmt.Sprintf(KML_START, title))

	builder.colaDeLineas = cola.CrearColaEnlazada[func()]() // Para asegurar que las lineas esten al final, se encolan funciones, para no guardar en structs
}

func (builder *KMLBuilder) CloseKML() {
	if !builder.started {
		return
	}

	for !builder.colaDeLineas.EstaVacia() {
		builder.colaDeLineas.Desencolar()()
	}

	builder.file.WriteString(KML_END)
	builder.Close()
}

func (builder *KMLBuilder) Close() {
	builder.file.Close()
}

func (builder *KMLBuilder) AddPoint(name, latitud, longitud string) {
	if !builder.started {
		return
	}

	builder.file.WriteString(POINT_KML(name, latitud, longitud))
}

func (builder *KMLBuilder) AddLine(name, fromLatitud, fromLongitud, toLatitud, toLongitud string) {
	builder.colaDeLineas.Encolar(func() { builder.file.WriteString(LINE_KML(name, fromLatitud, fromLongitud, toLatitud, toLongitud)) })
}
