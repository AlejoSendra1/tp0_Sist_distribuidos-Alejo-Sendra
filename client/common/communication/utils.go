package communication

import (
	"encoding/binary"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/domain"
)

const (
	FIRST_NAME_MAX_SIZE = 64
	LAST_NAME_MAX_SIZE = 64
    HEADER_SIZE = 3 // 2 B for the packet size + 1 B for the agency number
)

// Serilize the given Bet Struct to be send to the server
// based on the defined protocol
func SerializeBet(betData *domain.Bet) []byte {
	
    // in case the names exceed the wanted size they will be cut
	if len(betData.FirstName) > FIRST_NAME_MAX_SIZE {
		betData.FirstName = betData.FirstName[:FIRST_NAME_MAX_SIZE]
	}
	if len(betData.LastName) > LAST_NAME_MAX_SIZE {
		betData.LastName = betData.LastName[:LAST_NAME_MAX_SIZE]
	}
	
    var bytesToSend []byte

    tmp_buff := make([]byte, 4)
    binary.BigEndian.PutUint32(tmp_buff, betData.Document)
    bytesToSend = append(bytesToSend, tmp_buff...)

    tmp_buff = make([]byte, 8)
    binary.BigEndian.PutUint64(tmp_buff, betData.BetNumber)
    bytesToSend = append(bytesToSend, tmp_buff...)

    tmp_buff = make([]byte, 2)
    binary.BigEndian.PutUint16(tmp_buff, uint16(betData.Birthdate.Year()))
    bytesToSend = append(bytesToSend, tmp_buff...)

    bytesToSend = append(bytesToSend, byte(betData.Birthdate.Month())) 
    bytesToSend = append(bytesToSend, byte(betData.Birthdate.Day())) 

    // Strings need to be handled separately
    bytesToSend = append(bytesToSend, byte(uint8(len(betData.FirstName))))
    for _, char := range betData.FirstName {
        bytesToSend = append(bytesToSend, byte(char))
    }

    bytesToSend = append(bytesToSend, byte(uint8(len(betData.LastName))))
    for _, char := range betData.LastName {
        bytesToSend = append(bytesToSend, byte(char))
    }

    return bytesToSend
}

func createBatch(bets []domain.Bet, betsOffset int, agencyNum string) ([]byte, int) {
    var batch []byte
    
    for (len(batch) + bets[betsOffset].BytesSize() < 8000 - HEADER_SIZE) && len(bets) > betsOffset { // luego tomar de entorno
        betSerialization := SerializeBet(&bets[betsOffset])
        batch = append(batch, betSerialization...)
        betsOffset += 1
    }

    header := make([]byte, HEADER_SIZE)
    binary.BigEndian.PutUint16(header[0:2], uint16(len(batch)))
    agencyNumAsInt, _ := strconv.Atoi(agencyNum)
    header = append(header, byte(agencyNumAsInt))

    return append(header, batch...), betsOffset
}