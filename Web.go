/*
MIT License

Copyright (c) 2020 Joshua Glass (SuperNinja_4965)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package GoHelpers

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//Starts http and https Web server
func StartWebServer(HTTPPort string, HTTPSPort string, OpenBrowserOnLoad bool) {
	ExecPath, err2 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err2 != nil {
		log.Fatal(err2)
	}
	go func() {
		// check for https certificates
		if _, err := os.Stat(ExecPath + "/HTTPS-key/server.crt"); os.IsNotExist(err) {
			fmt.Printf("server.crt does not exist. HTTPS NOT STARTED\n")
		} else if _, err := os.Stat(ExecPath + "/HTTPS-key/server.key"); os.IsNotExist(err) {
			fmt.Printf("server.key does not exist. HTTPS NOT STARTED\n")
		} else {
			// begin https server
			err_https := http.ListenAndServeTLS(":"+HTTPSPort, ExecPath+"/HTTPS-key/server.crt", ExecPath+"/HTTPS-key/server.key", nil)
			Helper.OpenBrowser("https://localhost:" + HTTPSPort)
			if err_https != nil {
				log.Fatal("Web server (HTTPS): \n", err_https)
			}
		}
	}()

	// begin http server
	err_http := http.ListenAndServe(":"+HTTPPort, nil)
	Helper.OpenBrowser("http://localhost:" + HTTPPort)
	if err_http != nil {
		log.Fatal("Web server (HTTP): ", err_http)
	}
}

// OpenBrowser When given a URL it will open a Web browser to it
func OpenBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// GetServerIp Gets the IP address of the device.
// If the device has more than one IP it will select the first IP.
func GetServerIP(ipNum int) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	var ips [255]string
	var i int
	i = 0
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips[i] = ipnet.IP.String()
				i = i + 1
			}
		}
	}
	return ips[ipNum]
}
