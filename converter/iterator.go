package converter

import (
	"encoding/csv"
	"encoding/xml"
	"log"
	"os"

	"github.com/SkobelevIgor/stackexchange-xml-to-csv/encoders"
)

func iterate(typeName string, xmlFile *os.File, csvFile *os.File, skipHTMLDecoding bool) (totalCounter int64, convertedCounter int64, err error) {
	xmlDecoder := xml.NewDecoder(xmlFile)
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	encoderType, err := encoders.NewEncoder(typeName)
	if err != nil {
		return
	}

	err = csvWriter.Write(encoderType.GetCSVHeaderRow())
	if err != nil {
		return
	}

	for {
		t, _ := xmlDecoder.Token()
		if t == nil {
			break
		}

		switch ty := t.(type) {
		case xml.StartElement:
			if ty.Name.Local == "row" {
				totalCounter++
				encoder, _ := encoders.NewEncoder(typeName)
				err = xmlDecoder.DecodeElement(&encoder, &ty)
				checkError(err, xmlFile)

				err = csvWriter.Write(encoder.GETCSVRow(skipHTMLDecoding))
				checkError(err, xmlFile)
				convertedCounter++
			}
		}

	}

	return
}

func checkError(err error, file *os.File) {
	if err != nil {
		log.Printf("[%s] Error: %s",
			file.Name(), err)
	}
}
