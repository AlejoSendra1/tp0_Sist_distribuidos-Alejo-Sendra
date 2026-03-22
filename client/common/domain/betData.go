package domain

import (
	"os"
	"strconv"
	"time"

	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

var log = logging.MustGetLogger("log")


// ClientPersonalData data from the client used to keep register of transactions
type BetData struct {
	Agency uint8
	Document  uint32
	BetNumber uint64
	Birthdate time.Time
	FirstName string
	LastName string
}

// Gets and validates all fields to create the client BetData
func GetBetDataFromEnv() (BetData, error) {
	var err error
	data := BetData{}
	
	//CLI_ID / agency
	agencyIDStr := os.Getenv("CLI_ID")
	agencyVal, err := strconv.Atoi(agencyIDStr)
	if err != nil {
		return data, errors.Wrapf(err, "Could not parse CLI_ID env var as number")
	}
	data.Agency = uint8(agencyVal)
	

	//DOCUMENTO
	documentStr := os.Getenv("DOCUMENTO")
	docVal, err := strconv.Atoi(documentStr)
	if err != nil {
		return data, errors.Wrapf(err ,"Invalid document")
	}
	data.Document = uint32(docVal)
	
	//NUMERO
	betStrNumber := os.Getenv("NUMERO")
	betNumber, err := strconv.Atoi(betStrNumber)
	if err != nil {
		return data, errors.Wrapf(err, "Invalid bet number")
	}
	data.BetNumber = uint64(betNumber)

	//BIRTHDATE -> (1999-03-17)
	birthStr := os.Getenv("NACIMIENTO")
	data.Birthdate, err = time.Parse("2006-01-02", birthStr)
	if err != nil {
		return data, errors.Wrapf(err, "Invalid birthdate")
	}
	
	//NOMBRE Y APELLIDO
	data.FirstName = os.Getenv("NOMBRE")
	data.LastName = os.Getenv("APELLIDO")	
	if len(data.FirstName) == 0 || len(data.LastName) == 0 {
		return data, errors.New("Name fields can't be empty")
	} 

	return data, nil
}
