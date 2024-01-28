package utils

import (
	"bytes"
	"encoding/binary"
	"strings"
)

const HEADER_SIZE = 12

// DNSRecordType represents the DNS record type
type DNSRecordType uint16

// Constants for different DNS record types
const (
	TypeA     DNSRecordType = 1  // A record (IPv4 address)
	TypeAAAA  DNSRecordType = 28 // AAAA record (IPv6 address)
	TypeCNAME DNSRecordType = 5  // CNAME record
	TypeMX    DNSRecordType = 15 // MX record
	TypeTXT   DNSRecordType = 16 // TXT record
)

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	Qdcount uint16
	Ancount uint16
	Nscount uint16
	Arcount uint16
}

type DNSQuestion struct {
	Name  string
	Type  DNSRecordType
	Class uint16
}

// DNSPacket represents a full DNS packet
type DNSPacket struct {
	Header    DNSHeader
	Questions []DNSQuestion
}

// Serialize converts the DNSPacket into a byte slice
func (packet *DNSPacket) Serialize() []byte {
	buffer := new(bytes.Buffer)

	// Write header
	binary.Write(buffer, binary.BigEndian, packet.Header.ID)
	binary.Write(buffer, binary.BigEndian, packet.Header.Flags)
	binary.Write(buffer, binary.BigEndian, packet.Header.Qdcount)
	binary.Write(buffer, binary.BigEndian, packet.Header.Ancount)
	binary.Write(buffer, binary.BigEndian, packet.Header.Nscount)
	binary.Write(buffer, binary.BigEndian, packet.Header.Arcount)

	// Write questions
	for _, question := range packet.Questions {
		writeDNSName(buffer, question.Name)
		binary.Write(buffer, binary.BigEndian, question.Type)
		binary.Write(buffer, binary.BigEndian, question.Class)
	}

	return buffer.Bytes()
}

// writeDNSName writes a domain name in DNS packet format
func writeDNSName(buffer *bytes.Buffer, name string) {
	for _, part := range strings.Split(name, ".") {
		buffer.WriteByte(byte(len(part)))
		buffer.WriteString(part)
	}
	buffer.WriteByte(0) // Null byte to end the name
}
func ParseDNSPacket(data []byte) DNSPacket {
	packet := DNSPacket{}
	header := parseDNSHeader(data)
	packet.Header = header
	numberOfQuestions := int(header.Qdcount)
	offset := HEADER_SIZE
	for i := 0; i < numberOfQuestions; i++ {
		question, nextOffset := parseDNSQuestion(data[offset:])
		packet.Questions = append(packet.Questions, question)
		offset += nextOffset

	}

	return packet
}

func parseDNSHeader(data []byte) DNSHeader {
	return DNSHeader{
		ID:      binary.BigEndian.Uint16(data[:2]),
		Flags:   binary.BigEndian.Uint16(data[2:4]),
		Qdcount: binary.BigEndian.Uint16(data[4:6]),
		Ancount: binary.BigEndian.Uint16(data[6:8]),
		Nscount: binary.BigEndian.Uint16(data[8:10]),
		Arcount: binary.BigEndian.Uint16(data[10:12]),
	}
}

func parseDNSQuestion(data []byte) (DNSQuestion, int) {
	name, offset := parseDNSName(data)
	qtype := binary.BigEndian.Uint16(data[offset : offset+2])
	qclass := binary.BigEndian.Uint16(data[offset+2 : offset+4])
	return DNSQuestion{Name: name, Type: DNSRecordType(qtype), Class: qclass}, offset + 4
}

func parseDNSName(data []byte) (string, int) {
	var name string
	offset := 0

	for {
		length := int(data[offset])
		if length == 0 {
			offset++ // Move past the null byte
			break
		}

		if len(name) > 0 {
			name += "." // Add dot only between labels, not at the beginning
		}

		name += string(data[offset+1 : offset+1+length])
		offset += 1 + length
	}

	return name, offset
}
