package httpserver

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/ipoluianov/aneth_eth/an"
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/images"
	"github.com/ipoluianov/aneth_eth/pages"
	"github.com/ipoluianov/aneth_eth/static"
	"github.com/ipoluianov/aneth_eth/utils"
	"github.com/ipoluianov/aneth_eth/views"
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
		w.Header().Add("Content-Type", "text/css")
		w.Write([]byte(static.FileStyleCss))
		return
	}

	if strings.Contains(r.RequestURI, "single_chart.js") {
		// w.Header().Add("Content-Type", "text/css")
		w.Write([]byte(static.FileSingleChart))
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

	if reqType == "v" {
		if len(parts) == 2 {
			c.processView(w, r, parts[1])
			return
		}
	}

	if reqType == "p" {
		if len(parts) == 2 {
			c.processPage(w, r, parts[1])
			return
		}
	}

	if reqType == "images" {
		if len(parts) == 2 {
			item, err := images.Instance.Get(parts[1])
			if err != nil {
				w.WriteHeader(404)
				return
			}
			w.Header().Add("Content-Type", "image/png")
			w.Write(item.Data)
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
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(result))
		return
	}

	// STATIC HTML
	html := string(static.FileIndex)
	html = strings.ReplaceAll(html, "%CONTENT%", "UNKNOWN QUERY")
	result = []byte(html)
	w.Write(result)
}

func (c *HttpServer) processView(w http.ResponseWriter, _ *http.Request, viewCode string) {
	result := string(static.FileIndex)

	title := common.GlobalSiteName
	description := common.GlobalSiteDescription
	content := ""

	if viewCode == "index" {
		content = c.getHomePage()
	}

	if viewCode == "legal_user_agreement" {
		content = static.FileUserAgreement
	}

	if viewCode == "legal_policy" {
		content = static.FilePrivatePolicy
	}

	if viewCode == "map" {
		content = pages.BuildMap()
		title = "Site Map - " + common.GlobalSiteName
		description = "Site Map. " + common.GlobalSiteDescription
	}

	if viewCode == "state" {
		content = static.FileState
		title = "State - " + common.GlobalSiteName
		description = "State of the site. " + common.GlobalSiteDescription
	}

	if len(content) == 0 {
		content, title, description = views.GetView(viewCode, title, description, "default", 400, true, true, true, true)
	}

	result = strings.ReplaceAll(result, "%TITLE%", title)
	result = strings.ReplaceAll(result, "%DESCRIPTION%", description)
	result = strings.ReplaceAll(result, "%CONTENT%", content)

	w.Write([]byte(result))
}

func (c *HttpServer) processPage(w http.ResponseWriter, _ *http.Request, pageCode string) {
	result := string(static.FileIndex)

	pageRes := pages.Instance.GetPage(pageCode)

	if pageRes == nil {
		w.WriteHeader(404)
		return
	}

	result = strings.ReplaceAll(result, "%TITLE%", pageRes.Name)
	result = strings.ReplaceAll(result, "%DESCRIPTION%", pageRes.Description)
	result = strings.ReplaceAll(result, "%CONTENT%", pageRes.Content)

	w.Write([]byte(result))
}

func (c *HttpServer) getHomePage() string {
	result := ""

	pageRes := pages.Instance.GetPage("index")
	result += pageRes.Content

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

	fAddHeader2 := func(name string) {
		tmp := `    <h2>%NAME%</h2>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		result += tmp
	}

	fAddHeader3 := func(name string) {
		tmp := `    <h3>%NAME%</h3>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		result += tmp
	}

	fAddHeader2("ETH.U00.IO")
	fAddText("ETH analytics")

	tasks := an.Instance.GetTasks()
	groups := an.Instance.GetTaskGroups()

	taskWithGroup := make(map[string]struct{})

	for _, gr := range groups {
		fAddHeader3(gr.Name)
		for _, task := range tasks {
			found := false
			for _, taskInGroup := range gr.Tasks {
				if task.Code == taskInGroup {
					found = true
					break
				}
			}
			if found {
				fAddItem(task.Name, "/v/"+task.Code)
				taskWithGroup[task.Code] = struct{}{}
			}
		}
	}

	fAddHeader3("Other reports")
	for _, task := range tasks {
		if _, ok := taskWithGroup[task.Code]; ok {
			continue
		}
		fAddItem(task.Name, "/v/"+task.Code)
	}

	return result
}
