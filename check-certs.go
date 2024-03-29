// Copyright 2013 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tag Modification History :
// 0.1 -cat- adding timeout
// 0.2 -cat- exitcode
// 0.3 -cat- retry on connection failure
// 0.4 -cat- short options (t,c,v)
// 0.41 - cat - : space separator
//

package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultConcurrency = 512

const (
	errExpiringShortly = "%s : ** '%s' (S/N %X) expires in %d hours! **"
	errExpiringSoon    = "%s : '%s' (S/N %X) expires in roughly %d days."
	errSunsetAlg       = "%s : '%s' (S/N %X) expires after the sunset date for its signature algorithm '%s'."
)

type sigAlgSunset struct {
	name      string    // Human readable name of signature algorithm
	sunsetsAt time.Time // Time the algorithm will be sunset
}

// sunsetSigAlgs is an algorithm to string mapping for signature algorithms
// which have been or are being deprecated.  See the following links to learn
// more about SHA1's inclusion on this list.
//
// - https://technet.microsoft.com/en-us/library/security/2880823.aspx
// - http://googleonlinesecurity.blogspot.com/2014/09/gradually-sunsetting-sha-1.html
var sunsetSigAlgs = map[x509.SignatureAlgorithm]sigAlgSunset{
	x509.MD2WithRSA: {
		name:      "MD2 with RSA",
		sunsetsAt: time.Now(),
	},
	x509.MD5WithRSA: {
		name:      "MD5 with RSA",
		sunsetsAt: time.Now(),
	},
	x509.SHA1WithRSA: {
		name:      "SHA1 with RSA",
		sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	x509.DSAWithSHA1: {
		name:      "DSA with SHA1",
		sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	x509.ECDSAWithSHA1: {
		name:      "ECDSA with SHA1",
		sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
	},
}

var (
	hostsFile   = flag.String("hosts", "Services.txt", "The path to the file containing a list of hosts to check.")
	warnYears   = flag.Int("years", 0, "Warn if the certificate will expire within this many years.")
	warnMonths  = flag.Int("months", 0, "Warn if the certificate will expire within this many months.")
	warnDays    = flag.Int("days", 0, "Warn if the certificate will expire within this many days.")
	checkSigAlg = flag.Bool("check-sig-alg", true, "Verify that non-root certificates are using a good signature algorithm.")
	concurrency = flag.Int("c", defaultConcurrency, "concurrency : Maximum number of hosts to check at once.")
	timeout     = flag.Int("t", 5, "Timeout for TLS connection")
	noexitcode  = flag.Bool("noexitcode", false, "Don't return exit code.")
	verbose     = flag.Bool("v", false, "verbose mode, show insecure TLS < 1.1")
)

var versions = map[uint16]string{
	tls.VersionSSL30: "SSL",
	tls.VersionTLS10: "TLS 1.0",
	tls.VersionTLS11: "TLS 1.1",
	tls.VersionTLS12: "TLS 1.2", // secure
	tls.VersionTLS13: "TLS 1.3",
}

type certErrors struct {
	commonName string
	errs       []error
}

type hostResult struct {
	host  string
	err   error
	certs []certErrors
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate)

	flag.Parse()

	if len(*hostsFile) == 0 {
		flag.Usage()
		os.Exit(-1)
	}
	if *warnYears < 0 {
		*warnYears = 0
	}
	if *warnMonths < 0 {
		*warnMonths = 0
	}
	if *warnDays < 0 {
		*warnDays = 0
	}
	if *warnYears == 0 && *warnMonths == 0 && *warnDays == 0 {
		*warnDays = 30
	}
	if *concurrency < 0 {
		*concurrency = defaultConcurrency
	}

	err := processHosts()
	if !(*noexitcode) && err != 0 {
		os.Exit(err)
	}
}

func processHosts() int {
	done := make(chan struct{})
	defer close(done)

	hosts := queueHosts(done)
	results := make(chan hostResult)

	var wg sync.WaitGroup
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go func() {
			processQueue(done, hosts, results)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	errorexists := 0

	for r := range results {
		if r.err != nil {
			log.Printf("%s : %v\n", r.host, r.err)
			errorexists = 1
			continue
		}
		if len(r.certs) > 0 {
			for _, cert := range r.certs {
				for _, err := range cert.errs {
					log.Println(err)
					errorexists = 1
				}
			}
		}
	}
	return errorexists
}

func queueHosts(done <-chan struct{}) <-chan string {
	hosts := make(chan string)
	go func() {
		defer close(hosts)

		fileContents, err := ioutil.ReadFile(*hostsFile)
		if err != nil {
			return
		}
		lines := strings.Split(string(fileContents), "\n")
		for _, line := range lines {
			host := strings.TrimSpace(line)
			if len(host) == 0 || host[0] == '#' {
				continue
			}
			select {
			case hosts <- host:
			case <-done:
				return
			}
		}
	}()
	return hosts
}

func processQueue(done <-chan struct{}, hosts <-chan string, results chan<- hostResult) {
	for host := range hosts {
		select {
		case results <- checkHost(host):
		case <-done:
			return
		}
	}
}

func checkHost(host string) (result hostResult) {
	var dialer net.Dialer

	result = hostResult{
		host:  host,
		certs: []certErrors{},
	}

	dialer.Timeout = (time.Duration)(*timeout) * time.Second
	conn, err := tls.DialWithDialer(&dialer, "tcp", host, nil)
	if err != nil {
		for attempts := 3; attempts > 0 && err != nil; attempts-- {
			time.Sleep(time.Second)
			conn, err = tls.DialWithDialer(&dialer, "tcp", host, nil)
		}
		if err != nil {
			result.err = err
			return
		}
	}

	defer conn.Close()

	if *verbose && conn.ConnectionState().Version < tls.VersionTLS12 {
		log.Printf("WARNING: insecure %s version with %s\n", versions[conn.ConnectionState().Version], host)
	}

	timeNow := time.Now()
	checkedCerts := make(map[string]struct{})
	for _, chain := range conn.ConnectionState().VerifiedChains {
		for certNum, cert := range chain {
			if _, checked := checkedCerts[string(cert.Signature)]; checked {
				continue
			}
			checkedCerts[string(cert.Signature)] = struct{}{}
			cErrs := []error{}

			// Check the expiration.
			if timeNow.AddDate(*warnYears, *warnMonths, *warnDays).After(cert.NotAfter) {
				expiresIn := int64(cert.NotAfter.Sub(timeNow).Hours())
				if expiresIn <= 48 {
					cErrs = append(cErrs, fmt.Errorf(errExpiringShortly, host, cert.Subject.CommonName, cert.SerialNumber, expiresIn))
				} else {
					cErrs = append(cErrs, fmt.Errorf(errExpiringSoon, host, cert.Subject.CommonName, cert.SerialNumber, expiresIn/24))
				}
			}

			// Check the signature algorithm, ignoring the root certificate.
			if alg, exists := sunsetSigAlgs[cert.SignatureAlgorithm]; *checkSigAlg && exists && certNum != len(chain)-1 {
				if cert.NotAfter.Equal(alg.sunsetsAt) || cert.NotAfter.After(alg.sunsetsAt) {
					cErrs = append(cErrs, fmt.Errorf(errSunsetAlg, host, cert.Subject.CommonName, cert.SerialNumber, alg.name))
				}
			}

			result.certs = append(result.certs, certErrors{
				commonName: cert.Subject.CommonName,
				errs:       cErrs,
			})
		}
	}

	return
}
