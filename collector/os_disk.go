package collector

import (
	log "github.com/sirupsen/logrus"
)

type osDiskCollector struct {
	Query map[string]string
}

type osDisk struct {
	Instance  string
	Device    string
	Fstype 	  string
	TotalSize int64
	FreeSize  int64
	UsedRate  float64
	Status    string
	MetricStatus string
	ErrMsg    string
}

func init() {
	registerCollector("osDisk", NewOsDiskCollector)
}

func NewOsDiskCollector() (Collector, error) {
	queries := map[string]string{
		// metric node_filesystem_size
		"nodeFilesystemSize": "node_filesystem_size{mountpoint=\"/etc/hostname\"}",
		// metric node_filesystem_free
		"nodeFilesystemFree": "node_filesystem_free{mountpoint=\"/etc/hostname\"}",
	}
	return &osDiskCollector{
		Query: queries,
	}, nil
}

func (os *osDiskCollector) Merge(ch chan *CollectResult) error {
	data := osDiskData(os.Query)
	//log.Infof("+v", data)
	html, _ := mergeTpl("tpl/os_disk.html", data)

	//log.Infof("+v", *results["nodeFilesystemSize"])

	ch <- &CollectResult{Data: html}
	return nil
}

func (os *osDiskCollector) FileData(ch chan *CollectResult) error {
	data := osDiskData(os.Query)
	//log.Infof("+v", data)
	txt, _ := mergeTpl("tpl/os_disk.txt", data)

	//log.Infof("+v", *results["nodeFilesystemSize"])

	ch <- &CollectResult{Data: txt}
	return nil
}

// trans prometheus api result to osDisk
//use osDisk to merge with tpl
func transOsDisk(data map[string]*HttpGetRes) (res []osDisk) {
	nodeFilesystemSize := *data["nodeFilesystemSize"]

	if nodeFilesystemSize.Status == requestSuccess {
		result := nodeFilesystemSize.Data.Result
		for _, val := range result {
			totalSize, err := handleSizeUnit(val.Value.value, "GB")
			// build runtime error msg
			if e, disk := checkError(err); e{
				res = append(res, disk)
				continue
			}
			//get file system free info by instance
			instanceRes, err := helpTransDisk(data["nodeFilesystemFree"], val.Metric["instance"])
			if e, disk := checkError(err); e{
				res = append(res, disk)
				continue
			}
			// all check pass
			freeSize, err := handleSizeUnit(instanceRes.Value.value, "GB")
			// build runtime error msg
			if e, disk := checkError(err); e{
				res = append(res, disk)
				continue
			}
			useRate := useRate(totalSize, freeSize, 3)
			osDisk := osDisk{
				Status:    nodeFilesystemSize.Status,
				Instance:  val.Metric["instance"],
				TotalSize: totalSize,
				FreeSize: freeSize,
				UsedRate: useRate,
				MetricStatus: osDiskMetricStatus(useRate),
				Device: val.Metric["device"],
				Fstype: val.Metric["fstype"],
			}
			res = append(res, osDisk)
		}
	} else {
		osDisk := osDisk{
			Status: nodeFilesystemSize.Status,
			MetricStatus: metricStatusBad,
			ErrMsg: nodeFilesystemSize.Message,
		}
		res = append(res, osDisk)
	}

	return res
}

func checkError(err error) (b bool, disk osDisk){
	b = false
	if err != nil {
		disk = osDisk{
			Status: runtimeError,
			MetricStatus: metricStatusBad,
			ErrMsg: err.Error(),
		}
		b = true
	}
	return b, disk
}

func helpTransDisk(fileFree *HttpGetRes, instance string) (res *Result, err error){
	if fileFree.Status == requestSuccess{
		result := fileFree.Data.Result
		for _, val := range result {

			if val.Metric["instance"] == instance{
				res = &val
				break
			}
		}
	}else {
		err = &CustomError{
			errMsg: fileFree.Message,
		}
	}
	return res, err
}

func osDiskData(query map[string]string) (res []osDisk){
	q := &PromeQuery{
		Query: query,
	}
	results, err := multiPrometheusRequest(q)
	if err != nil {
		log.Errorf("osDisk metric error,", err.Error())
	}
	return transOsDisk(results)
}

func osDiskMetricStatus(f float64) (res string) {
	if f < 70{
		return metricStatusOK
	}else if f > 90 {
		return metricStatusBad
	}else {
		return metricStatusWarn
	}
}