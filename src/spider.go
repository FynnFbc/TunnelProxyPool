package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var wg2 sync.WaitGroup
var wgp sync.WaitGroup
var ch2 = make(chan int, 10000)

// 是否抓取中
var run = false

func spiderRun() {
	run = true
	defer func() {
		run = false
	}()

	count = 0
	//log.Println("开始抓取代理...")

	Functions := []func(){
		get_qydaili,     //齐云代理
		get_89daili,     //89代理
		get_kxdaili,     //开心代理
		get_kdaili,      //快代理
		get_gkydaili,    //高可用代理 可用0
		get_xsdaili,     //小舒代理
		get_lmydaili,    //命运零代理
		get_dbdaili,     //db代理
		get_hidemydaili, //hidemy代理
		get_scrapedaili, //scrape代理
		get_mydaili,     //my代理
		get_prodaili,    //Proxy代理
		get_padaili,     //爬代理
		get_freshdaili,  //fresh代理
		get_p11daili,    //Proxy11代理
		get_66daili,     //66ip代理
		get_github,      //代理列表
		get_opendaili,   //open代理
	}
	for i := range Functions {
		wg2.Add(1)
		go Functions[i]()
	}

	wg2.Wait()
	//log.Printf("\r%s 代理抓取结束  可用IP: %s\n", time.Now().Format("2006-01-02 15:04:05"), len(ProxyPool))

	//导出代理到文件
	export()
}

func get_qydaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "齐云代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := "\\s*?<td data-title=\"IP\">(.*?)</td>\\s*?<td data-title=\"PORT\">(.*?)</td>"
	for i := 1; i <= 20; i++ {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, "https://proxy.ip3366.net/free/?action=china&page="+strconv.Itoa(i), Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}
	wgp.Wait()
}
func get_89daili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "89代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `<tr>\s*?<td>\s*?(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s*?</td>\s*?<td>\s*?(\d+?)\s*?</td>`
	for i := 1; i <= 25; i++ {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, "https://www.89ip.cn/index_"+strconv.Itoa(i)+".html", Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}
	wgp.Wait()
}
func get_kxdaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "开心代理"
	Method := "GET"
	Body := ""
	Interval := 1
	Re := "<tr[\\s\\S]*?<td>(.*?)</td>\\s*?<td>(.*?)</td>"
	for i := 1; i <= 10; i++ {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, "http://www.kxdaili.com/dailiip/1/"+strconv.Itoa(i)+".html", Re)
		time.Sleep(time.Duration(Interval) * time.Second)
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, "http://www.kxdaili.com/dailiip/2/"+strconv.Itoa(i)+".html", Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}
	// 获取每日最新ip
	result := GetResp(Method, Body, "http://www.kxdaili.com/daili.html")
	urls := regexp.MustCompile("<a class=\"title\" href=\"(.*?)\">").FindAllStringSubmatch(result, -1)
	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, "http://www.kxdaili.com"+urls[0][1], "](.*?)@HTTP")
	wgp.Wait()
}
func get_kdaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "快代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := "<td>(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})</td>[\\s\\S]*?<td>(\\d+?)</td>"
	for s := 1; s <= 5; s++ {
		for i := 1; i <= 10; i++ {
			wgp.Add(1)
			go SpiderProxy(Name, Method, Body, "http://www.ip3366.net/free/?stype="+strconv.Itoa(s)+"&page="+strconv.Itoa(i), Re)
			time.Sleep(time.Duration(Interval) * time.Second)
		}
	}
	wgp.Wait()
}
func get_gkydaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "高可用代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}):(\\d+)@HTTP"
	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, "https://ip.jiangxianli.com/api/proxy_ips", "\"ip\":\"(.*?)\",\"port\":\"(.*?)\"")
	time.Sleep(time.Duration(Interval) * time.Second)
	result := GetResp(Method, Body, "https://ip.jiangxianli.com/blog.html")
	urls := regexp.MustCompile("<h3><a href=\"(.*?)\">").FindAllStringSubmatch(result, -1)
	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, urls[0][1], Re)
	time.Sleep(time.Duration(Interval) * time.Second)
	wgp.Wait()
}
func get_xsdaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "小舒代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}):(\\d+)@HTTP"
	result := GetResp(Method, Body, "https://www.xsdaili.cn/")
	urls := regexp.MustCompile("<div class=\"title\">\\s*?<a href=\"(.*?)\">").FindAllStringSubmatch(result, -1)
	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, "https://www.xsdaili.cn"+urls[0][1], Re)
	time.Sleep(time.Duration(Interval) * time.Second)
	wgp.Wait()
}
func get_lmydaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "命运零代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := "\"host\": \"(.*?)\", \"port\": (.*?),"

	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, "http://proxylist.fatezero.org/proxy.list", Re)
	time.Sleep(time.Duration(Interval) * time.Second)
	wgp.Wait()
}
func get_dbdaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "db代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := "\">(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:\\d+?)</a>"
	Urls := []string{
		"http://proxydb.net/?protocol=http&anonlvl=%s&country=CN",
		"http://proxydb.net/?protocol=https&anonlvl=%s&country=CN",
		"http://proxydb.net/?protocol=socks4&anonlvl=%s&country=CN",
		"http://proxydb.net/?protocol=socks5&anonlvl=%s&country=CN",
		"http://proxydb.net/?protocol=http&anonlvl=%s&country=KH",
		"http://proxydb.net/?protocol=https&anonlvl=%s&country=KH",
		"http://proxydb.net/?protocol=socks4&anonlvl=%s&country=KH",
		"http://proxydb.net/?protocol=socks5&anonlvl=%s&country=KH",
		"http://proxydb.net/?protocol=http&anonlvl=%s&country=HK",
		"http://proxydb.net/?protocol=https&anonlvl=%s&country=HK",
		"http://proxydb.net/?protocol=socks4&anonlvl=%s&country=HK",
		"http://proxydb.net/?protocol=socks5&anonlvl=%s&country=HK",
		"http://proxydb.net/?protocol=http&anonlvl=%s&country=TW",
		"http://proxydb.net/?protocol=https&anonlvl=%s&country=TW",
		"http://proxydb.net/?protocol=socks4&anonlvl=%s&country=TW",
		"http://proxydb.net/?protocol=socks5&anonlvl=%s&country=TW",
		"http://proxydb.net/?protocol=http&anonlvl=%s&country=SG",
		"http://proxydb.net/?protocol=https&anonlvl=%s&country=SG",
		"http://proxydb.net/?protocol=socks4&anonlvl=%s&country=SG",
		"http://proxydb.net/?protocol=socks5&anonlvl=%s&country=SG",
		"http://proxydb.net/?protocol=http&anonlvl=%s&country=JP",
		"http://proxydb.net/?protocol=https&anonlvl=%s&country=JP",
		"http://proxydb.net/?protocol=socks4&anonlvl=%s&country=JP",
		"http://proxydb.net/?protocol=socks5&anonlvl=%s&country=JP",
	}
	for u := range Urls {
		for i := 1; i <= 4; i++ {
			wgp.Add(1)
			go SpiderProxy(Name, Method, Body, fmt.Sprintf(Urls[u], strconv.Itoa(i)), Re)
			time.Sleep(time.Duration(Interval) * time.Second)
		}
	}
	wgp.Wait()
}
func get_hidemydaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "hidemy代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `<tr><td>(.*?)</td><td>(\d+?)</td>`
	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, "https://hidemy.name/cn/proxy-list/?maxtime=1000&type=45#list", Re)
	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, "https://hidemy.name/cn/proxy-list/?maxtime=1000&type=45&start=64#list", Re)

	Types := []string{"h", "s"}
	for t := range Types {
		for i := 0; i <= 210; i++ {
			wgp.Add(1)
			go SpiderProxy(Name, Method, Body, "https://hidemy.name/cn/proxy-list/?maxtime=1000&type="+Types[t]+"&start="+strconv.Itoa(i*64)+"#list", Re)
			time.Sleep(time.Duration(Interval) * time.Second)
		}
	}
	wgp.Wait()
}
func get_scrapedaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "scrape代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+?)\n`

	Urls := []string{
		"https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all&simplified=true",
		"https://api.proxyscrape.com/v2/?request=getproxies&protocol=socks4&timeout=10000&country=all&simplified=true",
		"https://api.proxyscrape.com/v2/?request=getproxies&protocol=socks5&timeout=10000&country=all&simplified=true",
	}
	for i := range Urls {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, Urls[i], Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}

	wgp.Wait()
}
func get_mydaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "my代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `>(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+?)#`

	Urls := []string{
		"https://www.my-proxy.com/free-socks-5-proxy.html",
		"https://www.my-proxy.com/free-socks-4-proxy.html",
	}
	for i := range Urls {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, Urls[i], Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}

	for i := 1; i <= 10; i++ {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, `https://www.my-proxy.com/free-proxy-list-`+strconv.Itoa(i)+`.html`, Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}

	wgp.Wait()
}
func get_prodaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "Proxy代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+?)\n`

	Urls := []string{
		"https://www.socks-proxy.net/",
		"https://free-proxy-list.net/",
		"https://www.us-proxy.org/",
		"https://free-proxy-list.net/uk-proxy.html",
		"https://www.sslproxies.org/",
		"https://free-proxy-list.net/anonymous-proxy.html",
	}
	for i := range Urls {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, Urls[i], Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}

	wgp.Wait()
}
func get_padaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "爬代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+?)<br/>`

	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, "http://www.padaili.com/proxyapi?api=vld845sXw5OmQBa00y4tLb5maonvSwct&num=100&type=3&xiangying=2&order=jiance", Re)
	time.Sleep(time.Duration(Interval) * time.Second)

	wgp.Wait()
}
func get_freshdaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "fresh代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `<td>(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})</td>[\s\S]*?<td>(\d+)</td>`

	Urls := []string{
		"https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1",
		"https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-2",
		"https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-3",
		"https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-4",
		"https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-5",
		"https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-6",
		"https://list.proxylistplus.com/Socks-List-1",
		"https://list.proxylistplus.com/Socks-List-2",
		"https://list.proxylistplus.com/SSL-List-1",
		"https://list.proxylistplus.com/SSL-List-2",
	}
	for i := range Urls {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, Urls[i], Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}

	wgp.Wait()
}
func get_p11daili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "Proxy11代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{2,5})`

	Urls := []string{
		"https://proxy11.com/api/proxy.txt?key=NTI5NA.Y2U2qw.EWL_l4huIsL15P8dQrfPszzZglY&country=Hong%20Kong&speed=3",
		"https://proxy11.com/api/proxy.txt?key=NTI5NA.Y2U2qw.EWL_l4huIsL15P8dQrfPszzZglY&country=Singapore&speed=3",
		"https://proxy11.com/api/proxy.txt?key=NTI5NA.Y2U2qw.EWL_l4huIsL15P8dQrfPszzZglY&country=Japan&speed=3",
		"https://proxy11.com/api/proxy.txt?key=NTI5NA.Y2U2qw.EWL_l4huIsL15P8dQrfPszzZglY&country=United%20States%20of%20America&speed=3",
	}
	for i := range Urls {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, Urls[i], Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}

	wgp.Wait()
}
func get_66daili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "66ip代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{2,5})<br`

	wgp.Add(1)
	go SpiderProxy(Name, Method, Body, "http://www.66ip.cn/mo.php?sxb=&tqsl=4200&port=&export=&ktip=&sxa=", Re)
	time.Sleep(time.Duration(Interval) * time.Second)

	wgp.Wait()
}
func get_github() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "代理列表"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{2,5})`

	Urls := []string{
		"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/socks4.txt",
		"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/socks5.txt",
		"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/http.txt",
		"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/https.txt",
		"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/http.txt",
		"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/socks4.txt",
		"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/socks5.txt",
		"https://github.com/monosans/proxy-list/blob/main/proxies/socks5.txt",
		"https://github.com/monosans/proxy-list/blob/main/proxies/socks4.txt",
		"https://github.com/monosans/proxy-list/blob/main/proxies/http.txt",
		"https://proxylist.live/nodes/free_2.php",
		"https://www.proxy-list.download/api/v1/get?type=http",
		"https://www.proxyscan.io/download?type=http",
		"https://api.openproxylist.xyz/http.txt",
		"http://alexa.lr2b.com/proxylist.txt",
		"http://rootjazz.com/proxies/proxies.txt",
		"https://www.freeproxychecker.com/result/http_proxies.txt",
		"https://multiproxy.org/txt_all/proxy.txt",
		"https://proxy-spider.com/api/proxies.example.txt",
		"http://spys.me/proxy.txt",
	}
	for i := range Urls {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, Urls[i], Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}

	wgp.Wait()
}
func get_opendaili() {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...", Name)
	}()
	Name := "open代理"
	Method := "GET"
	Body := ""
	Interval := 0
	Re := `items:["(.*?)"]`

	Urls := []string{
		"https://openproxy.space/list/socks4",
		"https://openproxy.space/list/socks5",
		"https://openproxy.space/list/http",
	}
	for i := range Urls {
		wgp.Add(1)
		go SpiderProxy(Name, Method, Body, Urls[i], Re)
		time.Sleep(time.Duration(Interval) * time.Second)
	}

	wgp.Wait()
}

func SpiderProxy(Name string, Method string, Body string, Url string, Re string) {
	defer func() {
		wgp.Done()
		//log.Printf("%s 结束...", Name)
	}()
	//log.Printf("%s 开始... %s", Name, Url)
	var pis []ProxyIp
	result := GetResp(Method, Body, Url)
	proxy := regexp.MustCompile(Re).FindAllStringSubmatch(result, -1)
	if len(proxy) == 0 {
		return
	}
	var _ip string
	var _port string
	for i := range proxy {
		if strings.Contains(proxy[i][1], ":") {
			tmp := strings.Split(proxy[i][1], ":")
			_ip, _ = url.QueryUnescape(tmp[0])
			_port, _ = url.QueryUnescape(tmp[1])
		} else {
			_ip, _ = url.QueryUnescape(proxy[i][1])
			_port, _ = url.QueryUnescape(proxy[i][2])
		}
		_is := true
		for pi := range ProxyPool {
			if ProxyPool[pi].Ip == _ip && ProxyPool[pi].Port == _port {
				_is = false
				break
			}
		}
		if _is {
			pis = append(pis, ProxyIp{Ip: _ip, Port: _port, Source: Name})
		}
	}
	pis = uniquePI(pis)
	countAdd(len(pis))
	for i := range pis {
		wg.Add(1)
		ch2 <- 1
		go Verify(&pis[i], &wg, ch2, true)
	}
	wg.Wait()

}
func GetResp(Method string, Body string, Url string) string {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	//if ProxyIs {
	//	proxyUrl, parseErr := url.Parse("http://" + conf.Proxy.Host + ":" + conf.Proxy.Port)
	//	if parseErr != nil {
	//		log.Println("代理地址错误: \n" + parseErr.Error())
	//		continue
	//	}
	//	tr.Proxy = http.ProxyURL(proxyUrl)
	//}
	client := http.Client{Timeout: 20 * time.Second, Transport: tr}
	request, _ := http.NewRequest(Method, Url, strings.NewReader(Body))
	//设置请求头
	SetHeadersConfig(map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36", "Connection": "close"}, &request.Header)
	//处理返回结果
	res, err := client.Do(request)
	if err != nil {
		return ""
	}
	dataBytes, _ := io.ReadAll(res.Body)
	return string(dataBytes)
}
