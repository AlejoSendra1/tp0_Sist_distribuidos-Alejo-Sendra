package common

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	FIRST_NAME_MAX_SIZE = 32
	LAST_NAME_MAX_SIZE = 32
)

// ClientPersonalData data from the client used to keep register of transactions
type BetData struct {
	Agency uint8
	FirstName string
	LastName string
	Document  uint32
	Birthdate time.Time
	BetNumber uint64
}

// Gets and validates all fields to create the client BetData
func GetBetDataFromEnv() (*BetData, error) {
	var err error
	data := &BetData{}

	//NOMBRE Y APELLIDO
	data.FirstName = os.Getenv("NOMBRE")
	data.LastName = os.Getenv("APELLIDO")	
	if len(data.FirstName) == 0 || len(data.LastName) == 0 {
		return nil, errors.New("Name fields can't be empty")
	} 
	if len(data.FirstName) > FIRST_NAME_MAX_SIZE {
		data.FirstName = data.FirstName[:FIRST_NAME_MAX_SIZE]
	}
	if len(data.LastName) > LAST_NAME_MAX_SIZE {
		data.LastName = data.LastName[:LAST_NAME_MAX_SIZE]
	}

	//CLI_ID
	agencyIDStr := os.Getenv("CLI_ID")
	agencyVal, err := strconv.Atoi(agencyIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_ID env var as number")
	}
	data.Agency = uint8(agencyVal)

	//DOCUMENTO
	documentStr := os.Getenv("DOCUMENTO")
	docVal, err := strconv.Atoi(documentStr)
	if err != nil {
		return nil, errors.Wrapf(err ,"Invalid document")
	}
	data.Document = uint32(docVal)

	//NUMBER
	betStrNumber := os.Getenv("NUMBER")
	betNumber, err := strconv.Atoi(betStrNumber)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid bet number")
	}
	data.BetNumber = uint64(betNumber)

	//BIRTHDATE -> (1999-03-17)
	birthStr := os.Getenv("NACIMIENTO")
	data.Birthdate, err = time.Parse("1999-03-17", birthStr)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid birthdate")
	}

	return data, nil
}