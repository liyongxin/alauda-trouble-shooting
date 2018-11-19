package collector

import (
	"encoding/json"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"sync"
	"bytes"
	"bufio"
	"html/template"
	log "github.com/sirupsen/logrus"
	"strconv"
	"math"
)

const (
	requestSuccess = "success"
	requestError = "req_error"
	runtimeError = "runtime_error"
)

type HttpGetRes struct {
	Status string `json:"status"`
	Data ResultValue `json:"data"`
	Message string `json:"message"`
}

type ResultValue struct {
	ResultType string `json:"resultType"`
	Result []Result `json:"result"`
}


type Result struct {
	Metric map[string]string `json:"metric"`
	Value  Value             `json:"value"`
}

type Value struct {
	time float64
	value string
}

//rewrite for parsing prometheus metric value 
func (v *Value)UnmarshalJSON(data []byte) error  {
	var raw []interface{}
	err := json.Unmarshal(data,&raw)
	if err != nil {
		log.Errorln("Unmarshal error:", err.Error())
	}
	floatValue, ok:=raw[0].(float64)
	if !ok {
		log.Errorln("Unmarshal error:", err.Error())
	}
	stringValue,ok:=raw[1].(string)
	if !ok {
		log.Errorln("Unmarshal error:", err.Error())
	}
	v.time = floatValue
	v.value = stringValue
	return nil
}

type CustomError struct {
	errMsg string
}

func (ce *CustomError) Error() string {
	return ce.errMsg
}

/*
	Desc: call prometheus api
	Returns: HttpGetRes
*/
func PrometheusHttpGet(query string) (res *HttpGetRes, err error){
	url := fmt.Sprintf("%s%s",PrometheusConfig.Address, query)
	res = &HttpGetRes{}
	defer func() {
		if err := recover();err != nil {
			errMsg := fmt.Sprintf("request error,%v, url=%s,", err, url)
			log.Errorln(errMsg)
			res.Status = "recover_error"
			res.Message = errMsg
			return
		}
	}()
	httpCLi := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpCLi.Get(url)
	if resp == nil && err != nil{
		errMsg := fmt.Sprintf("request error,%v, url=%s,", err, url)
		log.Errorln(errMsg)
		res.Status = "error"
		res.Message = errMsg
		return res, nil
	}
	defer resp.Body.Close()
	if err != nil {
		log.Errorln(err)
	}
	bts, _:= ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bts, res)
	if err != nil {
		log.Errorln(err)
	}
	return res, nil
}

type PromeQuery struct {
	Query map[string]string
}

//multi get prometheus request
func multiPrometheusRequest(query *PromeQuery) (map[string]*HttpGetRes, error){
	res := make(map[string]*HttpGetRes)
	defer func() {
		if err := recover();err != nil {
			errMsg := fmt.Sprintf("multi request error,%v", err)
			log.Errorln(errMsg)
			return
		}
	}()
	//sync
	wg := sync.WaitGroup{}
	length := len(query.Query)
	wg.Add(length)
	resCh := make(chan map[string]*HttpGetRes, length)
	for metric, path := range query.Query {
		go func(ch chan map[string]*HttpGetRes, metric, path string) {
			metricVals, err := PrometheusHttpGet(path)
			if err != nil {
				log.Errorf("get metric %s data error, ",metric, err.Error())
			}
			res := map[string]*HttpGetRes{
				metric: metricVals,
			}
			ch <- res
			wg.Done()
		}(resCh, metric, path)
	}
	wg.Wait()

	close(resCh)

	for ch := range resCh{
		for metric, value := range ch{
			res[metric] = value
		}
	}
	return res,nil
}

//merge tpl and data
//returns html strings
func mergeTpl(path string, data interface{}) (string, error){
	t, err := template.ParseFiles(path)
	if err != nil {
		return err.Error(), err
	}
	var buffer  bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	err = t.Execute(writer, data)
	writer.Flush()
	return buffer.String(), nil
}

func handleSizeUnit(val string, unit string) (ival int64, err error) {
	ival, err = strconv.ParseInt(val, 10,64)
	if err != nil {
		log.Errorln("trans unit error,", err)
		return -1 ,err
	}
	switch unit {
	case "GB":
		ival = ival >> 30
	case "MB":
		ival = ival >> 20
	case "KB":
		ival = ival >> 10
	}

	return ival, nil
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func useRate(total,free int64, n int) float64{
	sRate := float64(free)/float64(total)
	return (1 - Round(sRate, n)) * 100
}