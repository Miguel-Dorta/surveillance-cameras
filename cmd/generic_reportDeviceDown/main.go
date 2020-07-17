package main

import (
	"flag"
	"github.com/digineo/go-ping"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	pinger *ping.Pinger
	ip *net.IPAddr
	debug bool
	botToken, chatID, camID string
	httpClient = &http.Client{Timeout: time.Second * 5}
)

func init()  {
	var ipv4 string
	flag.StringVar(&ipv4, "ip", "127.0.0.1", "IP to ping")
	flag.StringVar(&botToken, "bot-token", "", "Telegram Bot Token")
	flag.StringVar(&chatID, "chat-id", "", "Telegram Chat ID")
	flag.StringVar(&camID, "cam-id", "", "Camera ID")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.Parse()

	var err error
	pinger, err = ping.New("0.0.0.0", "")
	if err != nil {
		log.Fatalf("error creating pinger: %s\n", err)
	}

	ip, err = net.ResolveIPAddr("ip4", ipv4)
	if err != nil {
		log.Fatalf("error resolving IPv4: %s\n", err)
	}
}

func main() {
	for {
		_, err := pinger.Ping(ip, time.Second)
		if err != nil {
			debugPrintf("ping failed")
			errRetry(err)
		}
		debugPrintf("ping success")
		time.Sleep(time.Second)
	}
}

func errRetry(err error) {
	errs := make([]error, 0, 5)
	errs = append(errs, err)

	for i:=0; i<5; i++ {
		debugPrintf("retry %d", i+1)
		timeout := time.Second
		if i==4 {
			timeout = time.Second * 5
		}
		_, err = pinger.Ping(ip, timeout)
		if err == nil {
			debugPrintf("retry success")
			return
		}
		errs = append(errs, err)
	}

	debugPrintf("more than 5 retries failed, reporting...")
	sendMessage(errs)
	return
}

func sendMessage(errs []error) {
	form := url.Values{}
	form.Set("text", createMessage("La comunicación con la cámara " + camID + " ha fallado. Motivos:", errs))
	form.Set("chat_id", chatID)

	resp, err := httpClient.PostForm("https://api.telegram.org/bot" + botToken + "/sendMessage", form)
	if err != nil {
		log.Println("Error sending post request to Telegram: " + err.Error())
	} else {
		resp.Body.Close()
	}

	disconnectedRetry()
	return
}

func disconnectedRetry() {
	secondsBetweenRequests := time.Second
	for {
		secondsBetweenRequests *= 2
		if secondsBetweenRequests > time.Minute * 5 {
			secondsBetweenRequests = time.Minute * 5
		}
		debugPrintf("waiting %s for new retry", secondsBetweenRequests.String())
		time.Sleep(secondsBetweenRequests)

		_, err := pinger.Ping(ip, time.Second * 5)
		if err == nil {
			debugPrintf("retry success")
			return
		}
		debugPrintf("retry failed")
	}
}

func createMessage(s string, errs []error) string {
	sb := new(strings.Builder)
	sb.WriteString(s)
	for _, err := range errs {
		sb.WriteString("\n- ")
		sb.WriteString(err.Error())
	}
	return sb.String()
}

func debugPrintf(format string, v ...interface{}) {
	if debug {
		log.Printf(format + "\n", v...)
	}
}
