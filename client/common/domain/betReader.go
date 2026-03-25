package domain

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

// ClientPersonalData data from the client used to keep register of transactions
type Bet struct {
	Document  uint32
	BetNumber uint64
	Birthdate time.Time
	FirstName string
	LastName string
}

type Reader {
	file *os.file
	scanner *bufio.Scanner
	batchAmount int
}

const (
	FIRSTNAME_CSV_POSITION = 0
	LASTNAME_CSV_POSITION = 1
	DOCUMENT_CSV_POSITION = 2
	BIRTHDATE_CSV_POSITION = 3
	BETNUMBER_CSV_POSITION = 4
	BET_STATIC_SIZE = 12 // 4 B for the document + 4 B for the birth date + 8 B for the bet number
)

var log = logging.MustGetLogger("log")

func Close() {
	file.close()
}

func NewReader(agencyNum string, batchAmount int) (*Reader,error) {
	file, err := os.Open(fmt.Sprintf("/.data/agency-%v.csv",agencyNum))
    if err != nil {
		return nil, err
    }
	log.Infof("action: opening_csv | result: success")

	reader := &Reader{
		file: file
		scanner: bufio.NewScanner(file)
		batchAmount: batchAmount
	}
	return reader, nil
}

// Gets as many bets from the file as indicated by the batchAmount
func (r *Reader)GetBatch() ([]Bet, error) {
	var bets []Bet
	
    for len(bets) < batchAmount && scanner.Scan() {
		line := scanner.Text()
        if err != nil {
			return bets,errors.Wrapf(err ,"error reading from file")
        }
		new_bet, _ := createBetFromStr(line)
		bets = append(bets, new_bet)
    }
	
	log.Infof("action: processing_csv | result: success")
	log.Infof("CANTIDAD DE BETS LEIDAS DE ARCHIVO %v", len(bets))
	
    return bets, nil
}

func createBetFromStr(betFields string) (Bet,error) {
	fields := strings.Split(betFields, ",")
	bet := Bet{}

	if len(fields) < 5 {
		return bet, errors.New("Unsufficient fields to create a Bet, error in csv")
	}

	//DOCUMENTO
	docNum, err := strconv.Atoi(fields[DOCUMENT_CSV_POSITION])
	if err != nil {
		return bet, errors.Wrapf(err ,"Invalid document")
	}
	bet.Document = uint32(docNum)

	//NUMERO
	betNumber, err := strconv.Atoi(fields[BETNUMBER_CSV_POSITION])
	if err != nil {
		return bet, errors.Wrapf(err, "Invalid bet number")
	}
	bet.BetNumber = uint64(betNumber)

	//BIRTHDATE 
	bet.Birthdate, err = time.Parse("2006-01-02", fields[BIRTHDATE_CSV_POSITION])
	if err != nil {
		return bet, errors.Wrapf(err, "Invalid birthdate")
	}

	bet.FirstName = fields[FIRSTNAME_CSV_POSITION]
	bet.LastName = fields[LASTNAME_CSV_POSITION]
	//NOMBRE Y APELLIDO
	if len(bet.FirstName) == 0 || len(bet.LastName) == 0 {
		return bet, errors.New("Name fields can't be empty")
	} 

	return bet, nil
}

func (bet *Bet) BytesSize() int{
	return len(bet.FirstName) + len(bet.LastName) + BET_STATIC_SIZE
}