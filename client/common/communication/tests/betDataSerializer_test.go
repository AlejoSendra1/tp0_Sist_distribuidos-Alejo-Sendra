package communication_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/communication"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/domain"
)

func TestSerialize_CorrectLayout(t *testing.T) {
	bet := &domain.Bet{
		FirstName: "John",
		LastName:  "Doe",
		Document:  12345678,
		Birthdate: time.Date(1990, time.March, 15, 0, 0, 0, 0, time.UTC),
		BetNumber: 9876543210,
	}

	result := communication.SerializeBet(bet)

	offset := 0

	// Document (4 bytes)
	doc := binary.BigEndian.Uint32(result[offset : offset+4])
	if doc != 12345678 {
		t.Errorf("Document: expected 12345678, got %d", doc)
	}
	offset += 4

	// BetNumber (8 bytes)
	betNum := binary.BigEndian.Uint64(result[offset : offset+8])
	if betNum != 9876543210 {
		t.Errorf("BetNumber: expected 9876543210, got %d", betNum)
	}
	offset += 8

	// Birthdate year (2 bytes)
	year := binary.BigEndian.Uint16(result[offset : offset+2])
	if year != 1990 {
		t.Errorf("Birthdate year: expected 1990, got %d", year)
	}
	offset += 2

	// Birthdate month (1 byte)
	if result[offset] != byte(time.March) {
		t.Errorf("Birthdate month: expected %d, got %d", time.March, result[offset])
	}
	offset += 1

	// Birthdate day (1 byte)
	if result[offset] != 15 {
		t.Errorf("Birthdate day: expected 15, got %d", result[offset])
	}
	offset += 1

	// FirstName length + content
	firstNameLen := int(result[offset])
	offset += 1
	if firstNameLen != len("John") {
		t.Errorf("FirstName length: expected %d, got %d", len("John"), firstNameLen)
	}
	if string(result[offset:offset+firstNameLen]) != "John" {
		t.Errorf("FirstName: expected John, got %s", string(result[offset:offset+firstNameLen]))
	}
	offset += firstNameLen

	// LastName length + content
	lastNameLen := int(result[offset])
	offset += 1
	if lastNameLen != len("Doe") {
		t.Errorf("LastName length: expected %d, got %d", len("Doe"), lastNameLen)
	}
	if string(result[offset:offset+lastNameLen]) != "Doe" {
		t.Errorf("LastName: expected Doe, got %s", string(result[offset:offset+lastNameLen]))
	}
}

func TestSerialize_CorrectTotalLength(t *testing.T) {
	bet := &domain.Bet{
		FirstName: "John",
		LastName:  "Doe",
		Document:  1,
		Birthdate: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		BetNumber: 1,
	}

	result := communication.SerializeBet(bet)

	// 4 + 8 + 2 + 1 + 1 + 1 + len("John") + 1 + len("Doe")
	expectedLen := 4 + 8 + 2 + 1 + 1 + 1 + len("John") + 1 + len("Doe")
	if len(result) != expectedLen {
		t.Errorf("Total length: expected %d, got %d", expectedLen, len(result))
	}
}

func TestSerialize_NamesAreTrimmedAtMaxSize(t *testing.T) {
	longName := string(make([]byte, communication.FIRST_NAME_MAX_SIZE+10))
	bet := &domain.Bet{
		FirstName: longName,
		LastName:  longName,
		Document:  1,
		Birthdate: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		BetNumber: 1,
	}

	result := communication.SerializeBet(bet)

	// jump fixed fields: 4 + 8 + 2 + 1 + 1 = 16
	offset := 16
	firstNameLen := int(result[offset])
	if firstNameLen != communication.FIRST_NAME_MAX_SIZE {
		t.Errorf("FirstName should be trimmed to %d, got %d", communication.FIRST_NAME_MAX_SIZE, firstNameLen)
	}
	offset += 1 + firstNameLen

	lastNameLen := int(result[offset])
	if lastNameLen != communication.LAST_NAME_MAX_SIZE {
		t.Errorf("LastName should be trimmed to %d, got %d", communication.LAST_NAME_MAX_SIZE, lastNameLen)
	}
}

func TestSerialize_EmptyNames(t *testing.T) {
	bet := &domain.Bet{
		FirstName: "",
		LastName:  "",
		Document:  1,
		Birthdate: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		BetNumber: 1,
	}

	result := communication.SerializeBet(bet)

	// jump fixed fields
	offset := 16
	if result[offset] != 0 {
		t.Errorf("FirstName length should be 0, got %d", result[offset])
	}
	offset += 1
	if result[offset] != 0 {
		t.Errorf("LastName length should be 0, got %d", result[offset])
	}
}