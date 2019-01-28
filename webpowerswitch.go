package main

/*
Copyright 2019 Chris Ching <chris@ching.codes>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import "flag"
import "fmt"
import "net/http"
import "os"

type WebSwitch struct {
	Address  string
	Port     int
	Username string
	Password string
	client   http.Client
}

const (
	OUTLET_FIRST = 1
	OUTLET_LAST  = 8
)

func (ws WebSwitch) Do(method string, url string) error {
	req, _ := http.NewRequest(method, fmt.Sprintf("http://%s:%d/%s", ws.Address, ws.Port, url), nil)
	req.SetBasicAuth(ws.Username, ws.Password)
	_, err := ws.client.Do(req)
	return err
}

func (ws WebSwitch) OutletOn(num int) error {
	return ws.Do("GET", fmt.Sprintf("outlet?%d=ON", num))
}

func (ws WebSwitch) OutletAllOn() error {
	// Interate over all switches since it's much faster then using the built in
	for i := OUTLET_FIRST; i <= OUTLET_LAST; i++ {
		err := ws.OutletOn(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ws WebSwitch) OutletOff(num int) error {
	return ws.Do("GET", fmt.Sprintf("outlet?%d=OFF", num))
}

func (ws WebSwitch) OutletAllOff() error {
	return ws.Do("GET", "outlet?a=OFF")
}

func (ws WebSwitch) OutletCycle(num int) error {
	return ws.Do("GET", fmt.Sprintf("outlet?%d=CCL", num))
}

// Note: this takes a while since it is using the switches 'all on' function
func (ws WebSwitch) OutletCycleAll() error {
	return ws.Do("GET", "outlet?a=CCL")
}

func main() {
	ws := WebSwitch{}
	flag.StringVar(&ws.Address, "addr", "192.168.0.100", "Address of Web Switch")
	flag.IntVar(&ws.Port, "port", 80, "HTTP port of Web Switch")
	flag.StringVar(&ws.Username, "user", "admin", "Username to login")
	flag.StringVar(&ws.Password, "password", "1234", "Password used to login")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS] CMD OUTLET|all\n\n ", os.Args[0])

		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "CMD\n"+
			"Command for switch [on|off]\n\n"+
			"OUTLET\n"+
			"Either outlet number or 'all'\n")
	}

	flag.Parse()

	args := flag.Args()

	var err error

	if len(args) == 2 {
		// Read outlet number if not set to 'all'
		allOutlets := (args[1] == "all" || args[1] == "ALL")
		var outletNum int
		if !allOutlets {
			fmt.Sscanf(args[1], "%d", &outletNum)
			if outletNum < OUTLET_FIRST || outletNum > OUTLET_LAST {
				fmt.Fprintf(os.Stderr, "Invalid outlet: %s\n", args[1])
				os.Exit(1)
			}
		}

		switch cmd := args[0]; cmd {
		case "on", "ON":
			if allOutlets {
				err = ws.OutletAllOn()
			} else {
				err = ws.OutletOn(outletNum)
			}
		case "off", "OFF":
			if allOutlets {
				err = ws.OutletAllOff()
			} else {
				err = ws.OutletOff(outletNum)
			}
		case "cycle", "CYCLE":
			if allOutlets {
				err = ws.OutletCycleAll()
			} else {
				err = ws.OutletCycle(outletNum)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
			flag.Usage()
			os.Exit(1)
		}

	} else {
		fmt.Fprintf(os.Stderr, "Missing or unexpected commands\n\n")
		flag.Usage()
		os.Exit(2)
	}

	if err != nil {
		os.Exit(1)
	}
}
