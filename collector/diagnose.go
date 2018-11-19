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

func (dc *diagnoseCollector) Merge(ch chan *CollectResult) error{
	//time.Sleep(1*time.Second)
	results, err := PrometheusHttpGet("global_service_diagnose")
	if err != nil {
		log.Errorf("diagnose metric error,", err.Error())
	}
	data := transDiagnose(results)
	html, _ := mergeTpl("tpl/diagnose.html", data)
	ch <- &CollectResult{Html: html}
	return nil
}

type diagnoseRes struct {
	Details string
	GlobalServiceName string
	Status string
	Url string
	Value string
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
			}
			res = append(res, diagnoseRes)
		}
	}else {
		diagnoseRes := diagnoseRes{
			Status: data.Status,
			ErrMsg: data.Message,
		}
		res = append(res, diagnoseRes)
	}
	return res
}

func (*diagnoseCollector) Data() error{
	return nil
}