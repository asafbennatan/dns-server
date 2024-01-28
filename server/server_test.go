package server

import (
	"dnsServer/client"
	"dnsServer/utils"
	"net"
	"testing"
	"time"
)

func Test_DNSServer(t *testing.T) {
	// Start the DNS server and get the stop channel
	server, err := NewDNSServer(":53")
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	server.Start()
	defer server.Stop()

	dnsClient, err := client.NewDNSClient("127.0.0.1:53")
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	defer dnsClient.Close()
	// Wait a bit to ensure the server is ready
	time.Sleep(time.Second)

	// Now run your tests
	// ...

	type args struct {
		question utils.DNSQuestion
	}
	tests := []struct {
		name    string
		args    args
		want    utils.DNSAnswer
		wantErr bool
	}{
		{
			name: "Test A Record Query",
			args: args{
				question: utils.DNSQuestion{Name: "google.com", Type: utils.TypeA},
			},
			want: utils.DNSAnswer{
				Name:  "google.com",
				Type:  utils.TypeA,
				Class: 1,
				TTL:   300,
				Addr:  net.IPv4(1, 2, 3, 4), // Expected IP address
			},
			wantErr: false,
		},

		{
			name: "Test AAAA Record Query",
			args: args{
				question: utils.DNSQuestion{Name: "google.com", Type: utils.TypeAAAA},
			},
			want: utils.DNSAnswer{
				Name:  "google.com",
				Type:  utils.TypeAAAA,
				Class: 1,
				TTL:   300,
				Addr:  net.ParseIP("::1"), // Expected IP address
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dnsClient.SendQuery(tt.args.question.Name, tt.args.question.Type)
			if (err != nil) != tt.wantErr {
				t.Errorf("DNS Query error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Answers) == 0 {
				t.Errorf("want %v , did not get any answer", tt.want)
				return
			}
			ans := got.Answers[0]
			if !answersEqual(ans, tt.want) {
				t.Errorf("DNS Query = %v, want %v", ans, tt.want)
			}
		})
	}

}

func answersEqual(a, b utils.DNSAnswer) bool {
	return a.Name == b.Name &&
		a.Type == b.Type &&
		a.Class == b.Class &&
		a.TTL == b.TTL &&
		a.Addr.Equal(b.Addr) // Use .Equal for net.IP comparison
	// Add more comparisons for other fields if necessary
}
