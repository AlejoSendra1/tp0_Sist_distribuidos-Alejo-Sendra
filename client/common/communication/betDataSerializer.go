package communication

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/domain"
)

const (
	FIRST_NAME_MAX_SIZE = 64
	LAST_NAME_MAX_SIZE = 64
)

// Serilize the given BetData Struct to be send to the server
// based on the defined protocol
func SerializeBetData(betData *domain.BetData) []byte {
	
    // in case the names exceed the wanted size they will be cut
	if len(betData.FirstName) > FIRST_NAME_MAX_SIZE {
		betData.FirstName = betData.FirstName[:FIRST_NAME_MAX_SIZE]
	}
	if len(betData.LastName) > LAST_NAME_MAX_SIZE {
		betData.LastName = betData.LastName[:LAST_NAME_MAX_SIZE]
	}
	
    var bytesToSend []byte
    bytesToSend = append(bytesToSend, byte(betData.Agency))

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