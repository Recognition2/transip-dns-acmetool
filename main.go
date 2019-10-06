package main

import "fmt"
import "github.com/transip/gotransip/domain"
import "github.com/transip/gotransip"
import "os"
import "strings"

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr, `Wrong number of arguments supplied (%d/%d). Please supply:
> Hook name, e.g. 'challenge-dns-start'
> Hostname, e.g. example.com
> Filename that causes the verification to happen
> DNS TXT value that should be set
`, len(os.Args), 5)
		os.Exit(-1)
	}

	hookName := os.Args[1]
	hostName := os.Args[2]
	// Third argument is the filename that causes verification, is unused here.
	dnsTXTValue := os.Args[4]

	// We only support challenge-dns-start, challenge-dns-stop hooks
	switch hookName {
	case "challenge-dns-start", "challenge-dns-stop":
		break
	default:
		// Indicate lack of support for event type
		os.Exit(42)
	}

	// create new TransIP API SOAP client
	c, err := gotransip.NewSOAPClient(gotransip.ClientConfig{
		AccountName:    "HELP",
		PrivateKeyPath: "transip-priv.key",
	})
	if err != nil {
		fmt.Println("Error creating SOAP client")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	dom, err := domain.GetInfo(c, hostName)
	if err != nil {
		fmt.Println("Error obtaining domain information")
		fmt.Println(err.Error())
		os.Exit(2)
	}

	dnsEntries := dom.DNSEntries

	if strings.HasSuffix(hookName, "start") {
		// Append new entry
		dnsEntries = append(dnsEntries, domain.DNSEntry{
			Name:    "_acme-challenge." + hostName,
			TTL:     3600,
			Type:    domain.DNSEntryTypeTXT,
			Content: dnsTXTValue,
		})
	} else if strings.HasSuffix(hookName, "stop") {
		// Delete our newly created entry
		NewDNSEntries := make([]domain.DNSEntry)
		for _, v := range dnsEntries {
			if !strings.HasPrefix(v.Name, "_acme-challenge") {
				NewDNSEntries = append(NewDNSEntries, v)
			}
		}
		dnsEntries = NewDNSEntries
	}

	err = domain.SetDNSEntries(c, hostName, dnsEntries)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
	//A hook is invoked successfully if it exits with exit code 0. A hook which exits with exit code 42 indicates a lack of
	//support for the event type. Any other exit code indicates an error.
	os.Exit(0)
}
