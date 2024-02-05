package server

import (
	"dnsServer/utils"
	"fmt"
	"net"
	"time"
)

type DNSServer struct {
	addr       string
	conn       *net.UDPConn
	stopSignal chan struct{}
}

func NewDNSServer(address string) (*DNSServer, error) {

	stopSignal := make(chan struct{})

	// Resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Listen for incoming connections
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return &DNSServer{addr: address, conn: conn, stopSignal: stopSignal}, nil

}

func (server *DNSServer) Start() {

	fmt.Printf("DNS Server is listening on %s\n", server.addr)
	go func() {
		for {
			// Set a timeout for reading from the connection
			server.conn.SetReadDeadline(time.Now().Add(time.Second))

			buffer := make([]byte, 1024)
			n, addr, err := server.conn.ReadFromUDP(buffer)

			// Check for timeout
			if err, ok := err.(net.Error); ok && err.Timeout() {
				// Check if a stop signal was sent
				select {
				case <-server.stopSignal:
					fmt.Println("Stopping DNS server")
					server.conn.Close()
					return
				default:
					// No stop signal, continue to the next iteration
					continue
				}
			}

			// Check for other errors
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// Handle the packet
			go handlePacket(server.conn, buffer[:n], addr)
		}
	}()

}

func (server *DNSServer) Stop() {
	server.stopSignal <- struct{}{}
}

// handlePacket processes the incoming packet and sends a response
func handlePacket(conn *net.UDPConn, data []byte, addr *net.UDPAddr) {
	//fmt.Printf("Received: %s from %s\n", string(data), addr)
	request := utils.ParseDNSPacket(data)

	//fmt.Printf("Received DNS query from %v\n", addr)
	//fmt.Printf("Header: %+v\n", request.Header)
	for i := 0; i < len(request.Questions); i++ {
		fmt.Printf("Question: %+v\n", request.Questions[i])

	}
	response := utils.DNSResponse{
		Header: utils.DNSHeader{
			ID:      request.Header.ID, // Use the same ID as the request
			Flags:   0x8180,            // Standard response, No error
			Qdcount: request.Header.Qdcount,
			Ancount: request.Header.Qdcount, // One answer
		},
		Questions: request.Questions,
	}
	for _, question := range response.Questions {
		answer := getAnswer(question)
		response.Answers = append(response.Answers, answer)
	}
	fmt.Printf(response.ToString())
	responseBytes := response.Serialize()
	_, err := conn.WriteToUDP(responseBytes, addr)
	if err != nil {
		fmt.Println(err)
	}
}

func getAnswer(question utils.DNSQuestion) utils.DNSAnswer {
	var answer utils.DNSAnswer

	// For this example, we always return an A record pointing to 1.2.3.4
	answer.Name = question.Name
	answer.Type = question.Type
	answer.Class = question.Class
	answer.TTL = 300 // TTL of 300 seconds

	// Handle different question types
	switch question.Type {
	case utils.TypeA: // Type 1: A Record (IPv4 address)
		ip := net.IPv4(1, 2, 3, 4) // Example IPv4 address
		answer.Addr = ip.To4()
		break

	case utils.TypeAAAA: // Type 28: AAAA Record (IPv6 address)
		ipv6 := net.ParseIP("::1") // Example IPv6 address
		answer.Addr = ipv6.To16()
		break

	case utils.TypeCNAME: // Type 5: CNAME Record
		answer.Cname = "example.com." // Example CNAME

	case utils.TypeMX: // Type 15: MX Record
		answer.MXPref = uint16(10)          // Example priority
		answer.MXHost = "mail.example.com." // Example mail exchange domain

	case utils.TypeTXT: // Type 16: TXT Record
		var mySlice []string

		mySlice = append(mySlice, "v=spf1 include:example.com ~all")

		answer.TXTData = mySlice

	// Add additional cases for other types as needed
	default:
		// Optionally handle unknown types or leave them unhandled
	}
	return answer
}
