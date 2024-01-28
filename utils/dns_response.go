package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type DNSAnswer struct {
	Name    string
	Type    DNSRecordType
	Class   uint16
	TTL     uint32
	Addr    net.IP
	Cname   string
	MXPref  uint16   // For MX records, preference value
	MXHost  string   // For MX records, host name
	TXTData []string // For TXT records, can be multiple strings
}

type DNSResponse struct {
	Header    DNSHeader
	Questions []DNSQuestion
	Answers   []DNSAnswer
}

func ParseDNSResponse(data []byte) DNSResponse {
	packet := DNSResponse{}
	header := parseDNSHeader(data)
	packet.Header = header
	numberOfQuestions := int(header.Qdcount)
	numberOfAnswers := int(header.Ancount)

	offset := HEADER_SIZE
	for i := 0; i < numberOfQuestions; i++ {
		question, nextOffset := parseDNSQuestion(data[offset:])
		packet.Questions = append(packet.Questions, question)
		offset += nextOffset

	}
	for i := 0; i < numberOfAnswers; i++ {
		answer, nextOffset := parseDNSAnswer(data[offset:])
		packet.Answers = append(packet.Answers, answer)
		offset += nextOffset

	}

	return packet
}

func parseDNSAnswer(data []byte) (DNSAnswer, int) {
	var answer DNSAnswer
	var offset int

	name, n := parseDNSName(data)
	answer.Name = name
	offset = n

	answer.Type = DNSRecordType(binary.BigEndian.Uint16(data[offset : offset+2]))
	answer.Class = binary.BigEndian.Uint16(data[offset+2 : offset+4])
	answer.TTL = binary.BigEndian.Uint32(data[offset+4 : offset+8])
	dataLength := binary.BigEndian.Uint16(data[offset+8 : offset+10])
	offset += 10

	switch answer.Type {
	case TypeA: // A record
		answer.Addr = net.IP(data[offset : offset+int(dataLength)]).To4()

	case TypeAAAA: // AAAA record
		answer.Addr = net.IP(data[offset : offset+int(dataLength)]).To16()

	case TypeCNAME: // CNAME record
		cname, _ := parseDNSName(data[offset:])
		answer.Cname = cname

	case TypeMX: // MX record
		answer.MXPref = binary.BigEndian.Uint16(data[offset : offset+2])
		exchange, n := parseDNSName(data[offset+2:])
		answer.MXHost = exchange
		offset += n

	case TypeTXT: // TXT record
		end := offset + int(dataLength)
		var txtParts []string
		for offset < end {
			txtLength := int(data[offset])
			offset++
			txtParts = append(txtParts, string(data[offset:offset+txtLength]))
			offset += txtLength
		}
		answer.TXTData = txtParts
	}

	offset += int(dataLength)
	return answer, offset
}

func (response DNSResponse) Serialize() []byte {
	buffer := new(bytes.Buffer)

	// Write header
	binary.Write(buffer, binary.BigEndian, response.Header.ID)
	binary.Write(buffer, binary.BigEndian, response.Header.Flags)
	binary.Write(buffer, binary.BigEndian, response.Header.Qdcount)
	binary.Write(buffer, binary.BigEndian, response.Header.Ancount)
	binary.Write(buffer, binary.BigEndian, response.Header.Nscount)
	binary.Write(buffer, binary.BigEndian, response.Header.Arcount)

	// Write questions
	for _, question := range response.Questions {
		writeDNSName(buffer, question.Name)
		binary.Write(buffer, binary.BigEndian, question.Type)
		binary.Write(buffer, binary.BigEndian, question.Class)
	}

	// Write answers
	for _, answer := range response.Answers {
		writeDNSName(buffer, answer.Name)
		binary.Write(buffer, binary.BigEndian, answer.Type)
		binary.Write(buffer, binary.BigEndian, answer.Class)
		binary.Write(buffer, binary.BigEndian, answer.TTL)
		switch answer.Type {
		case TypeA:
			binary.Write(buffer, binary.BigEndian, uint16(4)) // Length of an IPv4 address is 4 bytes
			buffer.Write(answer.Addr.To4())
			break
		case TypeAAAA:
			binary.Write(buffer, binary.BigEndian, uint16(16)) // Length of an IPv6 address
			buffer.Write(answer.Addr.To16())
			break
		case TypeCNAME: // CNAME record
			cname := serializeDNSName(answer.Cname) // Assume answer.CNAME is a string with the canonical name
			binary.Write(buffer, binary.BigEndian, uint16(len(cname)))
			buffer.Write(cname)

		case TypeMX: // MX record
			binary.Write(buffer, binary.BigEndian, uint16(2+len(serializeDNSName(answer.MXHost)))) // 2 bytes for priority + length of host
			binary.Write(buffer, binary.BigEndian, answer.MXPref)                                  // MX priority
			buffer.Write(serializeDNSName(answer.MXHost))

		case TypeTXT: // TXT record
			/*txt := []byte(answer.TXTData)
			binary.Write(buffer, binary.BigEndian, uint16(len(txt)))
			buffer.Write(txt)*/

		}

	}

	return buffer.Bytes()
}
func serializeDNSName(name string) []byte {
	var buffer bytes.Buffer
	writeDNSName(&buffer, name)
	return buffer.Bytes()
}

// ToString creates a string representation of the DNSResponse
func (response DNSResponse) ToString() string {
	var sb strings.Builder

	// Append header information
	sb.WriteString("Header:\n")
	sb.WriteString(fmt.Sprintf("ID: %d, Flags: %d, Qdcount: %d, Ancount: %d, Nscount: %d, Arcount: %d\n",
		response.Header.ID, response.Header.Flags, response.Header.Qdcount,
		response.Header.Ancount, response.Header.Nscount, response.Header.Arcount))

	// Append questions
	sb.WriteString("Questions:\n")
	for _, question := range response.Questions {
		sb.WriteString(fmt.Sprintf("Name: %s, Type: %d, Class: %d\n",
			question.Name, question.Type, question.Class))
	}

	// Append answers
	sb.WriteString("Answers:\n")
	for _, answer := range response.Answers {
		sb.WriteString(fmt.Sprintf("Name: %s, Type: %d, Class: %d, TTL: %d, Data: ",
			answer.Name, answer.Type, answer.Class, answer.TTL))

		switch answer.Type {
		case TypeA, TypeAAAA:
			sb.WriteString(fmt.Sprintf("IP Address: %s\n", answer.Addr))
		case TypeCNAME:
			sb.WriteString(fmt.Sprintf("CNAME: %s\n", answer.Cname))
		case TypeMX:
			sb.WriteString(fmt.Sprintf("MX Preference: %d, MX Host: %s\n", answer.MXPref, answer.MXHost))
		case TypeTXT:
			sb.WriteString(fmt.Sprintf("TXT: %s\n", strings.Join(answer.TXTData, ", ")))
			// Add more cases as needed for other types
		}
	}

	return sb.String()
}
