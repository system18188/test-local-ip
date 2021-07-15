package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	url string
	ip string
	level string
)

func init()  {
	flag.StringVar(&url, "url", "http://myip.ipip.net", "URL地址")
	flag.StringVar(&ip,"localIp", "", "检查指定坊问IP")
	flag.StringVar(&level,"level","debug","日记等级 可选值:panic,fatal,error,warn,warning,info,debug")
}

func main()  {
	flag.Parse()

	logrus.SetFormatter(&logrus.TextFormatter{})
	l, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.SetLevel(l)

	var client *http.Client
	if ip != "" {
		// 批定IP
		localAddr, err := net.ResolveIPAddr("ip", ip)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		localTCPAddr := net.TCPAddr{
			IP: localAddr.IP,
		}

		d := net.Dialer{
			LocalAddr: &localTCPAddr,
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}

		tr := &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DialContext: d.DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		}
		client = &http.Client{Transport: tr}
	} else {
		client = new(http.Client)
	}
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err)
		return
	}
	client.Timeout, _ = time.ParseDuration("30s")
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		logrus.Error(err)
		return
	}
	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Debug("GET:",url,"\n\r", "BODY:", string(body))
}
