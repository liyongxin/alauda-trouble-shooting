package main

import (
	"fmt"
	"strings"
	"os"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"gopkg.in/alecthomas/kingpin.v2"
	"alauda-trouble-shooting/collector"
	log "github.com/sirupsen/logrus"
)

const (
	prefixHttp          = "http://"
	prometheusQueryPath = "/api/v1/query?query="
	webServerCmd = "webServer"
	diagnoseCmd = "diagnose"
)

//log level
func init()  {

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

//http handler
func Handler(ctx iris.Context) {
	res := collector.Collect(webServerCmd)
	ctx.HTML(res)
}

func main() {
	var (
		prometheusAddress = kingpin.Flag("prometheus.address", "Address on which to expose metrics and web interface.").Required().String()
		webServer         = kingpin.Command(webServerCmd, "run webServer default listen on port 3322, you can run with --port to set listen port, then curl $HOSTIP:3322 to get website page.")
		webServerPort     = webServer.Flag("port", "port on webservice to expose").Default("3322").String()
		healthCheck       = kingpin.Command(diagnoseCmd, "check all functional module, include node, etcd, diagnose.")
		moduleName        = healthCheck.Arg("name", "check by module, for example \"alauda_trouble_shooting healthCheck etcd\". ").String()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	//init global prometheus url
	checkAndInitConfig(*prometheusAddress)

	switch kingpin.MustParse(kingpin.Parse(), nil) {
	// command selector
	case webServer.FullCommand():
		log.Infof("webServer starting with prometheus address %s", *prometheusAddress)

		app := iris.New()
		app.Logger().SetLevel("info")
		// Optionally, add two built'n handlers
		// that can recover from any http-relative panics
		// and log the requests to the terminal.
		app.Use(recover.New())
		app.Use(logger.New())

		// Method:   GET
		// Resource: http://localhost:8080
		app.Handle("GET", "/", Handler)
		app.Run(iris.Addr(fmt.Sprintf(":%s", *webServerPort)), iris.WithoutServerError(iris.ErrServerClosed))

	case healthCheck.FullCommand():
		if moduleName == nil || *moduleName == "" {
			log.Infoln("diagnose all")
			collector.Collect(diagnoseCmd)
		} else {
			log.Infof("diagnose %s， not supported for now!", *moduleName)
		}
	}
}

func checkAndInitConfig(prometheusUrl string) {
	if strings.HasPrefix(prometheusUrl, prefixHttp) {
		collector.PrometheusConfig.Address = fmt.Sprintf("%s%s", prometheusUrl, prometheusQueryPath)
	} else {
		log.Fatalf("prometheus.address error, must start with %s", prefixHttp)
		os.Exit(1)
	}

}
