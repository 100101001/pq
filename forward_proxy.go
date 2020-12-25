package pq

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

var proxyServerAddr string

const (
	HL_ServerAddr = "10.23.123.63:8080"
	LF_ServerAddr = "10.17.127.72:8080"
)

// 正向代理拨号
func forwardProxyDial(ctx context.Context, network, addr string) (conn net.Conn, e error) {
	proxyServerAddr = LF_ServerAddr
	defer func() {
		if e != nil {
			fmt.Println(e)
		}
		fmt.Println(fmt.Sprintf("forwardProxyDial proxyServerAddr=%s end", proxyServerAddr))
	}()

	proxyUser := "edu.prek.wx_data_aliyunrds_liyixuan.01"
	proxyPwd := "LL3pqgDb6KchzX1m"

	connectReq := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: addr},
		Host:   addr,
		Header: make(http.Header),
	}
	basic_str := base64.StdEncoding.EncodeToString([]byte(proxyUser + ":" + proxyPwd))
	connectReq.Header.Add("Proxy-Authorization", "Basic "+basic_str)

	conn, err := net.Dial(network, proxyServerAddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.Dial")
	}
	connectReq.Write(conn)

	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, connectReq)
	if err != nil {
		conn.Close()
		return nil, errors.Wrap(err, "http.ReadResponse")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("resp.StatusCode != 200 is :%d", resp.StatusCode))
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "ioutil.ReadAll")
		}
		fmt.Println(fmt.Sprintf("resp_body is :%s", string(resp_body)))
		conn.Close()
		return nil, errors.New("http proxy response is " + strconv.Itoa(resp.StatusCode))
	}
	return conn, nil
}
