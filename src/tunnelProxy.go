package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var lock2 sync.Mutex
var httpI []ProxyIp
var httpS []ProxyIp

var httpIp string
var httpsIp string

func httpSRunTunnelProxyServer() {
	log.Println("HTTP 隧道代理启动 - 监听IP端口 -> ", conf.Config.Ip+":"+conf.Config.HttpTunnelPort)
	server := &http.Server{
		Addr:      conf.Config.Ip + ":" + conf.Config.HttpTunnelPort,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpsIp = getHttpsIp()
			httpIp = gethttpIp()

			if authorization(w, r) {
				return
			}
			if r.Method == http.MethodConnect {
				//log.Printf("隧道代理 | HTTPS 请求：%s 使用代理: %s", r.URL.String(), httpsIp)
				destConn, err := net.DialTimeout("tcp", httpsIp, 20*time.Second)
				if err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}
				destConn.SetReadDeadline(time.Now().Add(20 * time.Second))
				var req []byte
				req = MergeArray([]byte(fmt.Sprintf("%s %s %s%s", r.Method, r.Host, r.Proto, []byte{13, 10})), []byte(fmt.Sprintf("Host: %s%s", r.Host, []byte{13, 10})))
				for k, v := range r.Header {
					req = MergeArray(req, []byte(fmt.Sprintf(
						"%s: %s%s", k, v[0], []byte{13, 10})))
				}
				req = MergeArray(req, []byte{13, 10})
				io.ReadAll(r.Body)
				all, err := io.ReadAll(r.Body)
				if err == nil {
					req = MergeArray(req, all)
				}
				destConn.Write(req)
				w.WriteHeader(http.StatusOK)
				hijacker, ok := w.(http.Hijacker)
				if !ok {
					http.Error(w, "not supported", http.StatusInternalServerError)
					return
				}
				clientConn, _, err := hijacker.Hijack()
				if err != nil {
					return
				}
				clientConn.SetReadDeadline(time.Now().Add(20 * time.Second))
				destConn.Read(make([]byte, 1024)) //先读取一次
				go io.Copy(destConn, clientConn)
				go io.Copy(clientConn, destConn)

			} else {
				//log.Printf("隧道代理 | HTTP 请求：%s 使用代理: %s", r.URL.String(), httpIp)
				tr := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
				//配置代理
				proxyUrl, parseErr := url.Parse("http://" + httpIp)
				if parseErr != nil {
					return
				}
				tr.Proxy = http.ProxyURL(proxyUrl)
				client := &http.Client{Timeout: 20 * time.Second, Transport: tr}
				request, err := http.NewRequest(r.Method, "", r.Body)
				//增加header选项
				request.URL = r.URL
				request.Header = r.Header
				//处理返回结果
				res, err := client.Do(request)
				if err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}
				defer res.Body.Close()

				for k, vv := range res.Header {
					for _, v := range vv {
						w.Header().Add(k, v)
					}
				}
				var bodyBytes []byte
				bodyBytes, _ = io.ReadAll(res.Body)
				res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				w.WriteHeader(res.StatusCode)
				io.Copy(w, res.Body)
				res.Body.Close()

			}
		}),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

//代理身份验证
func authorization(w http.ResponseWriter, r *http.Request) bool {
	var auth string
	if r.Header.Get("Authorization") != "" {
		auth = r.Header.Get("Authorization")
	} else if r.Header.Get("Proxy-Authorization") != "" {
		auth = r.Header.Get("Proxy-Authorization")
	} else {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Private\"")
		w.WriteHeader(401)
		w.Write([]byte("401 Proxy Authentication Required"))
		return true
	}
	s := strings.SplitN(auth, " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return true
	}
	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return true
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return true
	}
	//conf.Auth.User
	//log.Printf("解码Basic: %s", pair)
	if pair[0] != conf.Auth.User && pair[1] != conf.Auth.Pass {
		return true
	}
	return false
}

// MergeArray 合并数组
func MergeArray(dest []byte, src []byte) (result []byte) {
	result = make([]byte, len(dest)+len(src))
	//将第一个数组传入result
	copy(result, dest)
	//将第二个数组接在尾部，也就是 len(dest):
	copy(result[len(dest):], src)
	return
}

func gethttpIp() string {
	lock2.Lock()
	defer lock2.Unlock()
	if len(ProxyPool) == 0 {
		return ""
	}
	for _, v := range ProxyPool {
		if v.Type == "HTTP" {
			is := true
			for _, vv := range httpI {
				if v.Ip == vv.Ip && v.Port == vv.Port {
					is = false
				}
			}
			if is {
				httpI = append(httpI, v)
				return v.Ip + ":" + v.Port
			}
		}
	}
	var addr string
	if len(httpI) != 0 {
		addr = httpI[0].Ip + ":" + httpI[0].Port
	}
	httpI = make([]ProxyIp, 0)
	if addr == "" {
		addr = httpsIp
	}
	return addr
}

func getHttpsIp() string {
	lock2.Lock()
	defer lock2.Unlock()
	if len(ProxyPool) == 0 {
		return ""
	}
	for _, v := range ProxyPool {
		if v.Type == "HTTPS" || v.Type == "SOCKET5" {
			is := true
			for _, vv := range httpS {
				if v.Ip == vv.Ip && v.Port == vv.Port {
					is = false
				}
			}
			if is {
				httpS = append(httpS, v)
				return v.Ip + ":" + v.Port
			}
		}
	}
	var addr string
	if len(httpS) != 0 {
		addr = httpS[0].Ip + ":" + httpS[0].Port
	}
	httpS = make([]ProxyIp, 0)
	return addr
}
