package collector

import (
	"sync"
	"strings"
	"html/template"
	log "github.com/sirupsen/logrus"
)

// troubleCollector implements the Collector interface.
type troubleCollector struct {
	Collectors map[string]Collector
}

type CollectResult struct {
	Data string
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Merge data & template.
	Merge(chan *CollectResult) error
	FileData(chan *CollectResult) error
}

var (
	factories = make(map[string]func() (Collector, error))
)

//for specific collector init called
func registerCollector(collector string, factory func() (Collector, error)){
	factories[collector] = factory
}

// for main to get all registered collectors
func newTroubleCollector() (*troubleCollector, error) {
	collectors := make(map[string]Collector)
	for col := range factories {
		collector, err := factories[col]()
		if err != nil {
			return nil, err
		}
		collectors[col] = collector
	}
	return &troubleCollector{Collectors: collectors}, nil
}

var PrometheusConfig struct {
	Address string
	Cmd  string
}

//core collect
func Collect(cmd string) (res string){
	webServer := "webServer"
	collectors, err := newTroubleCollector()
	if err != nil {
		log.Errorln("get trouble collector err,", err)
	}
	wg := sync.WaitGroup{}
	length := len(collectors.Collectors)
	wg.Add(length)
	resCh := make(chan *CollectResult, length)
	for _, col := range collectors.Collectors {
		go func(ch chan *CollectResult, c Collector) {
			if cmd == webServer{
				c.Merge(resCh)
			}else{
				c.FileData(resCh)
			}
			wg.Done()
		}(resCh, col)
	}
	wg.Wait()

	close(resCh)
	var resArr  []string
	for res := range resCh{
		resArr = append(resArr, res.Data)
		//log.Error(res.Data)
	}
	if cmd == webServer{
		res, _ = mergeTpl("tpl/common.html", template.HTML(strings.Join(resArr,"")))
	}else{
		res, _ = mergeTpl("tpl/common.txt", template.HTML(strings.Join(resArr,"\n\n")))
		filename := createFile(res)
		log.Info("diagnose data file ", filename)
	}
	return res
}
