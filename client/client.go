package client

import (
	"dnsServer/utils"
	"net"
)

type DNSClient struct {
	conn *net.UDPConn
}

func NewDNSClient(serverAddress string) (*DNSClient, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &DNSClient{conn: conn}, nil
}

func (client *DNSClient) SendQuery(name string, requestType utils.DNSRecordType) (utils.DNSResponse, error) {
	header := utils.DNSHeader{
		ID:      0xABCD, // Transaction ID
		Flags:   0x0100, // Standard query
		Qdcount: 1,      // One question
	}

	packet := utils.DNSPacket{
		Header:    header,
		Questions: []utils.DNSQuestion{{name, requestType, 1}},
	}

	sentData := packet.Serialize()
	_, err := client.conn.Write(sentData)
	if err != nil {
		return utils.DNSResponse{}, err
	}

	buffer := make([]byte, 1024)
	n, _, err := client.conn.ReadFromUDP(buffer)
	if err != nil {
		return utils.DNSResponse{}, err
	}

	response := utils.ParseDNSResponse(buffer[:n])
	return response, nil
}

func (client *DNSClient) Close() {
	client.conn.Close()
}
