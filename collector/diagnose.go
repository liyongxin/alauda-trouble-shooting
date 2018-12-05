package collector

import (
	log "github.com/sirupsen/logrus"
)

type diagnoseCollector struct {
}

func init(){
	registerCollector("diagnose", NewDiagnoseCollector)
}

func NewDiagnoseCollector() (Collector, error){
	return &diagnoseCollector{}, nil
}

//send html data to channel
func (dc *diagnoseCollector) Merge(ch chan *CollectResult) error{
	//time.Sleep(1*time.Second)
	data := diagnoseData()
	html, _ := mergeTpl("tpl/diagnose.html", data)
	ch <- &CollectResult{Data: html}
	return nil
}

//send txt data to channel
func (*diagnoseCollector) FileData(ch chan *CollectResult) error{
	data := diagnoseData()
	txt, _ := mergeTpl("tpl/diagnose.txt", data)
	ch <- &CollectResult{Data: txt}
	return nil
}

type diagnoseRes struct {
	Details string
	GlobalServiceName string
	Status string
	Url string
	Value string
	MetricStatus string
	ErrMsg string
}

func transDiagnose(data *HttpGetRes) (res []diagnoseRes) {
	if data.Status == requestSuccess {
		for _, result := range data.Data.Result {
			diagnoseRes := diagnoseRes{
				Status: result.Metric["status"],
				Url: result.Metric["url"],
				GlobalServiceName: result.Metric["global_service_name"],
				Details: result.Metric["details"],
				Value: result.Value.value,
				MetricStatus: diagnoseMetricStatus(result.Metric["status"]),
			}
			res = append(res, diagnoseRes)
		}
	}else {
		diagnoseRes := diagnoseRes{
			Status: data.Status,
			ErrMsg: data.Message,
			MetricStatus: metricStatusBad,
		}
		res = append(res, diagnoseRes)
	}
	return res
}


func diagnoseData() (res []diagnoseRes){
	results, err := PrometheusHttpGet("global_service_diagnose")
	if err != nil {
		log.Errorf("diagnose metric error,", err.Error())
	}
	return transDiagnose(results)
}

func diagnoseMetricStatus(s string) string {
	if s == "OK"{
		return metricStatusOK
	}else {
		return metricStatusBad
	}
}
