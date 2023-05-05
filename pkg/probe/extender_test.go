package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestProbeExtender(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/cmdb/extender-controller/extender", "testdata/extender_cmdb.jsonnet")
	c.prepare("api/v2/monitor/extender-controller/extender", "testdata/extender_monitor.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeExtender, c, r) {
		t.Errorf("probeManagedSwitchStatus() returned non-success")
	}

	em := `
			# HELP fortigate_extender_connect_status Connect status 1=up 0=N/A
			# TYPE fortigate_extender_connect_status gauge
			fortigate_extender_connect_status{extender_name="FX301E5920020182",modem="modem1",vdom="root"} 1
			# HELP fortigate_extender_cpu Cpu utilization
			# TYPE fortigate_extender_cpu gauge
			fortigate_extender_cpu{extender_name="FX301E5920020182",vdom="root"} 0
			# HELP fortigate_extender_exists Information if extender exists or not 1=exists 0=do not exists
			# TYPE fortigate_extender_exists gauge
			fortigate_extender_exists{vdom="dhcp-relay"} 0
			fortigate_extender_exists{vdom="root"} 1
			# HELP fortigate_extender_info Infos about a extender
			# TYPE fortigate_extender_info gauge
			fortigate_extender_info{extender_name="FX301E5920020182",hardware_version="P23421-02",software_version="FXT201E-v4.2.2-build302",vdom="root"} 1
			# HELP fortigate_extender_lte_rssi LTE rssi
			# TYPE fortigate_extender_lte_rssi gauge
			fortigate_extender_lte_rssi{extender_name="FX301E5920020182",modem="modem1",vdom="root"} -77
			# HELP fortigate_extender_lte_sinr LTE sinr
			# TYPE fortigate_extender_lte_sinr gauge
			fortigate_extender_lte_sinr{extender_name="FX301E5920020182",modem="modem1",vdom="root"} 13
			# HELP fortigate_extender_memory Memory utilization
			# TYPE fortigate_extender_memory gauge
			fortigate_extender_memory{extender_name="FX301E5920020182",vdom="root"} 15
			# HELP fortigate_extender_modem_info Info about the modem
			# TYPE fortigate_extender_modem_info gauge
			fortigate_extender_modem_info{band="B3",data_plan="public.foo.se",esn_imei="359073068890064",extender_name="FX301E5920020182",lte_physical_cellid="01B9010F",manufacturer="Sierra Wireless, Incorporated",model="EM7455",modem="modem1",modem_type="EM7455",product="Sierra Wireless, Incorporated",service="LTE",vdom="root",wireless_operator="foo SE"} 1
			# HELP fortigate_extender_signal_rsrp Reference Signal Received Power
			# TYPE fortigate_extender_signal_rsrp gauge
			fortigate_extender_signal_rsrp{extender_name="FX301E5920020182",modem="modem1",vdom="root"} -100
			# HELP fortigate_extender_signal_rssq Reference Signal Received Quality
			# TYPE fortigate_extender_signal_rssq gauge
			fortigate_extender_signal_rssq{extender_name="FX301E5920020182",modem="modem1",vdom="root"} -8.7
			# HELP fortigate_extender_signal_strength Signal strength
			# TYPE fortigate_extender_signal_strength gauge
			fortigate_extender_signal_strength{extender_name="FX301E5920020182",modem="modem1",vdom="root"} 46
			# HELP fortigate_extender_sim_active Sim card active
			# TYPE fortigate_extender_sim_active gauge
			fortigate_extender_sim_active{extender_name="FX301E5920020182",modem="modem1",sim="sim1",vdom="root"} 1
			fortigate_extender_sim_active{extender_name="FX301E5920020182",modem="modem1",sim="sim2",vdom="root"} 0
			# HELP fortigate_extender_sim_data_usage Sim card data usage
			# TYPE fortigate_extender_sim_data_usage gauge
			fortigate_extender_sim_data_usage{extender_name="FX301E5920020182",modem="modem1",sim="sim1",vdom="root"} 436
			fortigate_extender_sim_data_usage{extender_name="FX301E5920020182",modem="modem1",sim="sim2",vdom="root"} 0
			# HELP fortigate_extender_sim_info Infos about the extender modem sim card
			# TYPE fortigate_extender_sim_info gauge
			fortigate_extender_sim_info{extender_name="FX301E5920020182",iccid="",ismi="N/A",modem="modem1",sim="sim2",vdom="root"} 1
			fortigate_extender_sim_info{extender_name="FX301E5920020182",iccid="89460860027131513443",ismi="210094716949637",modem="modem1",sim="sim1",vdom="root"} 1
			# HELP fortigate_extender_sim_status Sim card status
			# TYPE fortigate_extender_sim_status gauge
			fortigate_extender_sim_status{extender_name="FX301E5920020182",modem="modem1",sim="sim1",vdom="root"} 1
			fortigate_extender_sim_status{extender_name="FX301E5920020182",modem="modem1",sim="sim2",vdom="root"} 0
			`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}

}
