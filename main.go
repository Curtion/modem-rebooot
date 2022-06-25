package main

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

func login() (string, error) {
	resp, err := http.Post("http://192.168.1.1/login.cgi", "application/x-www-form-urlencoded",
		strings.NewReader("username=useradmin&password=ayak6&save=%B5%C7%C2%BC"),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	compileRegex := regexp.MustCompile("sessionKey=(.*?)\"")
	matchArr := compileRegex.FindStringSubmatch(string(body))
	sessionKey := matchArr[len(matchArr)-1]
	return sessionKey, nil
}

func reboot(sessionKey string) (string, error) {
	resp, err := http.Get("http://192.168.1.1/rebootinfo.cgi?sessionKey=" + sessionKey)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func start_reboot() {
	sessionKey, err := login()
	if err != nil {
		log.Fatal(err)
	}
	_, err = reboot(sessionKey)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	for range time.Tick(time.Duration(300) * time.Second) {
		pinger, err := ping.NewPinger("baidu.com")
		pinger.SetPrivileged(true)
		if err != nil {
			log.Fatal(err)
		}
		pinger.Count = 5
		err = pinger.Run()
		if err != nil {
			log.Fatal(err)
		}
		stats := pinger.Statistics()
		count := stats.PacketsRecv
		if count == 0 {
			start_reboot()
		}
	}
}
