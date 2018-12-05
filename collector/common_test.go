package collector

import (
	"testing"
	"encoding/json"
	"fmt"
)

var rawJsonStr = "[1542200396.592, \"0.39266490459946607\"]"
func TestCollect(t *testing.T) {
	t.Run("tesst unmarshal", func(t *testing.T) {
		abc := Value{}
		json.Unmarshal([]byte(rawJsonStr),&abc)
		println(fmt.Sprintf("%+v",abc))
	})
}

var testRwa = `{"status":"success"
,"data":{"resultType":"vector","result":[{"metric":{"container_name":"prometheus-exporter-demo","cpu":"total","endpoint":"https-metrics","id":"/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod8aaca0da_d908_11e8_83a5_000c2948e532.slice/docker-4c1fb7831500975e05698e6ce554ab441fd78fd16e37512da88cf4390c5ab177.scope","image":"192.168.8.33:60080/alaudak8s/prometheus-exporter-demo@sha256:deb71047721eacc5fced15de1674aca084fcd3ad49a83e3a6be4387d21249e68","instance":"192.168.8.4:10250","job":"kubelet","name":"k8s_prometheus-exporter-demo_prometheus-exporter-demo-7f85b4df64-9xs72_default_8aaca0da-d908-11e8-83a5-000c2948e532_0","namespace":"default","pod_name":"prometheus-exporter-demo-7f85b4df64-9xs72","service":"kubelet"},"value":[1542200396.592,"0.39266490459946607"]}]}}
`
func TestResult(t *testing.T) {
	t.Run("test result", func(t *testing.T) {
		result := HttpGetRes{}
		err:=json.Unmarshal([]byte(testRwa),&result)
		if err != nil {
			println(err.Error())
		}
		t.Logf("%+v",result)
	})

}

func TestByte(t *testing.T)  {
	t.Run("test bytes to mb", func(t *testing.T) {
		ival, err := handleSizeUnit("536608213121", "GBQ")
		if err != nil{
			t.Logf(err.Error())
		}
		t.Log(ival)
	})
}

func TestDemcial(t *testing.T)  {
	a := float64(44152348672)
	b := float64(12961140736)
	c := b/a
	t.Log(b/a)
	e := Round(c,3)
	t.Log(e)
}
