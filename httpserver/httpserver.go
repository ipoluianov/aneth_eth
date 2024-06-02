package httpserver

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/ipoluianov/aneth_eth/an"
	"github.com/ipoluianov/aneth_eth/static"
	"github.com/ipoluianov/aneth_eth/utils"
	"github.com/ipoluianov/gomisc/logger"
)

var Instance *HttpServer

type HttpServer struct {
	srv    *http.Server
	srvTLS *http.Server
}

func NewHttpServer() *HttpServer {
	var c HttpServer
	return &c
}

func init() {
	Instance = NewHttpServer()
}

func (c *HttpServer) Start() {
	go c.thListen()
	go c.thListenTLS()
}

func (c *HttpServer) portHttp() string {
	if utils.IsRoot() {
		return ":80"
	}
	return ":8080"
}

func (c *HttpServer) portHttps() string {
	if utils.IsRoot() {
		return ":443"
	}
	return ":8443"
}

func (c *HttpServer) thListen() {
	c.srv = &http.Server{
		Addr: c.portHttp(),
	}

	c.srv.Handler = c

	logger.Println("HttpServer thListen begin")
	err := c.srv.ListenAndServe()
	if err != nil {
		logger.Println("HttpServer thListen error: ", err)
	}
	logger.Println("HttpServer thListen end")
}

func (c *HttpServer) thListenTLS() {
	logger.Println("HttpServer::thListenTLS begin")
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = make([]tls.Certificate, 0)
	pathToBundle := logger.CurrentExePath() + "/bundle.crt"
	pathToPrivate := logger.CurrentExePath() + "/private.key"
	logger.Println("HttpServer::thListenTLS bundle.crt path:", pathToBundle)
	logger.Println("HttpServer::thListenTLS private.key path:", pathToPrivate)
	logger.Println("HttpServer::thListenTLS loading certificates ...")
	cert, err := tls.LoadX509KeyPair(pathToBundle, pathToPrivate)
	if err == nil {
		logger.Println("HttpServer::thListenTLS certificates is loaded SUCCESS")
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	} else {
		logger.Println("HttpServer::thListenTLS loading certificates ERROR", err)
		return
	}

	serverAddress := c.portHttps()
	c.srvTLS = &http.Server{
		Addr:      serverAddress,
		TLSConfig: tlsConfig,
	}
	c.srvTLS.Handler = c

	logger.Println("HttpServer::thListenTLS starting server at", serverAddress)
	listener, err := tls.Listen("tcp", serverAddress, tlsConfig)
	if err != nil {
		logger.Println("HttpServer::thListenTLS starting server ERROR", err)
		return
	}

	logger.Println("HttpServer::thListenTLS starting server SUCCESS")
	err = c.srvTLS.Serve(listener)
	if err != nil {
		logger.Println("HttpServerTLS thListen error: ", err)
		return
	}
	logger.Println("HttpServer::thListenTLS end")
}

func (c *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		logger.Println("ProcessHTTP host: ", r.Host)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Request-Method", "GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		redirectUrl := ""
		if utils.IsRoot() {
			host := strings.ReplaceAll(r.Host, c.portHttp(), "")
			redirectUrl = "https://" + host + r.RequestURI

		} else {
			host := strings.ReplaceAll(r.Host, c.portHttp(), "")
			redirectUrl = "https://" + host + c.portHttps() + r.RequestURI
		}
		logger.Println("Redirect to HTTPS:", redirectUrl)
		http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Request-Method", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	if strings.Contains(r.RequestURI, "style.css") {
		w.Write([]byte(static.FileStyleCss))
		return
	}

	parts := strings.FieldsFunc(r.RequestURI, func(r rune) bool {
		return r == '/'
	})

	var result []byte

	if len(parts) == 0 {
		c.processPage(w, r, "index")
		return
	}

	reqType := parts[0]

	if reqType == "p" {
		if len(parts) == 2 {
			c.processPage(w, r, parts[1])
			return
		}
	}

	if reqType == "d" {
		if len(parts) < 2 {
			w.WriteHeader(500)
			w.Write([]byte("wrong request: api - missing argument"))
			return
		}
		pageCode := parts[1]
		result := GetData(pageCode)
		w.Write([]byte(result))
		return
	}

	// STATIC HTML
	html := string(static.FileIndex)
	html = strings.ReplaceAll(html, "%CONTENT%", "UNKNOWN QUERY")
	result = []byte(html)
	w.Write(result)
}

func (c *HttpServer) processPage(w http.ResponseWriter, r *http.Request, pageCode string) {
	result := []byte(static.FileIndex)
	str := string(result)

	content := ""

	if pageCode == "index" {
		content = c.getHomePage()
	}

	if pageCode == "map" {
		content = c.getMap()
	}

	if pageCode == "state" {
		content = static.FileState
	}

	if len(content) == 0 {
		content = c.getPage(pageCode)
	}

	result = []byte(strings.ReplaceAll(str, "%CONTENT%", content))
	w.Write(result)
}

func (c *HttpServer) getMap() string {
	result := ""

	fAddItem := func(name string, url string) {
		tmp := `    <li><a href="%URL%">%NAME%</a></li>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%URL%", url)
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		result += tmp
	}

	fAddHeader := func(name string) {
		tmp := `    <h2>%NAME%</h2>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		result += tmp
	}

	fAddHeader("Main")
	fAddItem("INDEX", "/")
	fAddItem("SITE MAP", "/p/map")

	anResults := an.Instance.GetResultsCodes()
	fAddHeader("Analytics")
	for _, r := range anResults {
		fAddItem(r, "/p/"+r)
	}

	fAddHeader("JSON-RPC")
	fAddItem("STATE", "/d/state")

	for _, r := range anResults {
		fAddItem(r, "/d/"+r)
	}

	return result
}

func (c *HttpServer) getPage(code string) string {
	result := ""
	task := an.Instance.GetTask(code)
	if task == nil {
		return ""
	}

	if task.Type == "timechart" {
		result = static.FileViewChart
	}
	if task.Type == "table" {
		result = static.FileViewTable
	}

	result = strings.ReplaceAll(result, "%VIEW_CODE%", task.Code)
	result = strings.ReplaceAll(result, "%VIEW_NAME%", task.Name)
	result = strings.ReplaceAll(result, "%VIEW_DESC%", task.Description)
	result = strings.ReplaceAll(result, "VIEW_INSTANCE", "default")

	return result
}

func (c *HttpServer) getHomePage() string {
	result := ""

	fAddItem := func(name string, url string) {
		tmp := `    <li><a href="%URL%">%NAME%</a></li>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%URL%", url)
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		result += tmp
	}

	fAddText := func(text string) {
		tmp := `<div>%TEXT%</div>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%TEXT%", text)
		result += tmp
	}

	fAddHeader := func(name string) {
		tmp := `    <h2>%NAME%</h2>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		result += tmp
	}

	fAddHeader("ETH.U00.IO")
	fAddText("ETH analytics")

	anResults := an.Instance.GetResultsCodes()
	for _, r := range anResults {
		fAddItem(r, "/p/"+r)
	}

	return result
}
