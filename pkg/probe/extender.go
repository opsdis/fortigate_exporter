package probe

import (
	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"reflect"
	_ "reflect"
	"regexp"
	"strconv"
)

func validate(i interface{}) bool {
	switch i.(type) {
	case float64:
		return true
	case string:
		return true
	default:
		return false
	}
}

func convertMapToStruct(m map[string]interface{}, s interface{}) error {
	stValue := reflect.ValueOf(s).Elem()
	sType := stValue.Type()
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		tagName := field.Tag.Get("json")
		if value, ok := m[tagName]; ok {
			if validate(value) {
				stValue.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

func probeExtender(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {

	var (
		extenderExist = prometheus.NewDesc(
			"fortigate_extender_exists",
			"Information if extender exists or not 1=exists 0=do not exists",
			[]string{"vdom"}, nil,
		)

		extenderInfo = prometheus.NewDesc(
			"fortigate_extender_info",
			"Infos about a extender",
			[]string{"vdom", "extender_name", "software_version", "hardware_version"}, nil,
		)
		extenderCpu = prometheus.NewDesc(
			"fortigate_extender_cpu",
			"Cpu utilization",
			[]string{"vdom", "extender_name"}, nil,
		)
		extenderMemory = prometheus.NewDesc(
			"fortigate_extender_memory",
			"Memory utilization",
			[]string{"vdom", "extender_name"}, nil,
		)
		modemInfo = prometheus.NewDesc(
			"fortigate_extender_modem_info",
			"Info about the modem",
			[]string{"vdom", "extender_name", "modem", "data_plan", "manufacturer", "product", "model", "service",
				"esn_imei", "band", "modem_type", "wireless_operator", "lte_physical_cellid"}, nil,
		)
		signalRsrq = prometheus.NewDesc(
			"fortigate_extender_signal_rssq",
			"Reference Signal Received Quality",
			[]string{"vdom", "extender_name", "modem"}, nil,
		)

		signalRsrp = prometheus.NewDesc(
			"fortigate_extender_signal_rsrp",
			"Reference Signal Received Power",
			[]string{"vdom", "extender_name", "modem"}, nil,
		)

		lteSinr = prometheus.NewDesc(
			"fortigate_extender_lte_sinr",
			"LTE sinr",
			[]string{"vdom", "extender_name", "modem"}, nil,
		)

		lteRssi = prometheus.NewDesc(
			"fortigate_extender_lte_rssi",
			"LTE rssi",
			[]string{"vdom", "extender_name", "modem"}, nil,
		)
		connectStatus = prometheus.NewDesc(
			"fortigate_extender_connect_status",
			"Connect status 1=up 0=N/A",
			[]string{"vdom", "extender_name", "modem"}, nil,
		)
		signalStrength = prometheus.NewDesc(
			"fortigate_extender_signal_strength",
			"Signal strength",
			[]string{"vdom", "extender_name", "modem"}, nil,
		)
		simInfo = prometheus.NewDesc(
			"fortigate_extender_sim_info",
			"Infos about the extender modem sim card",
			[]string{"vdom", "extender_name", "modem", "sim", "ismi", "iccid"}, nil,
		)
		simDataUsage = prometheus.NewDesc(
			"fortigate_extender_sim_data_usage",
			"Sim card data usage",
			[]string{"vdom", "extender_name", "modem", "sim"}, nil,
		)
		simActive = prometheus.NewDesc(
			"fortigate_extender_sim_active",
			"Sim card active",
			[]string{"vdom", "extender_name", "modem", "sim"}, nil,
		)
		simStatus = prometheus.NewDesc(
			"fortigate_extender_sim_status",
			"Sim card status",
			[]string{"vdom", "extender_name", "modem", "sim"}, nil,
		)
	)

	type Sim struct {
		Carrier            string  `json:"carrier"`
		PhoneNumber        string  `json:"phone_number"`
		Status             string  `json:"status"`
		IsActive           float64 `json:"is_active"`
		Imsi               string  `json:"imsi"`
		Iccid              string  `json:"iccid"`
		MaximumAllowedData float64 `json:"maximum_allowed_data"`
		OverageAllowed     string  `json:"overage_allowed"`
		NextBillingDate    string  `json:"next_billing_date"`
		DataUsage          float64 `json:"data_usage"`
		Slot               float64 `json:"slot"`
		Modem              float64 `json:"modem"`
	}

	type CdmaProfile struct {
		Nai         string `json:"NAI"`
		Idx         string `json:"idx"`
		Status      string `json:"status"`
		HomeAddr    string `json:"home_addr"`
		PrimaryHa   string `json:"primary_ha"`
		SecondaryHa string `json:"secondary_ha"`
		AaaSpi      string `json:"aaa_spi"`
		HaSpi       string `json:"ha_spi"`
	}

	type Modem struct {
		DataPlan          string       `json:"data_plan"`
		PhysicalPort      string       `json:"physical_port"`
		Manufacturer      string       `json:"manufacturer"`
		Product           string       `json:"product"`
		Model             string       `json:"model"`
		Revision          string       `json:"revision"`
		Imsi              string       `json:"imsi"`
		PinStatus         string       `json:"pin_status"`
		Service           string       `json:"service"`
		SignalStrength    string       `json:"signal_strength"`
		Rssi              string       `json:"rssi"`
		ConnectStatus     string       `json:"connect_status"`
		GsmProfile        []any        `json:"gsm_profile"`
		CdmaProfile       *CdmaProfile `json:"cdma_profile"`
		EsnImei           string       `json:"esn_imei"`
		ActivationStatus  string       `json:"activation_status"`
		RoamingStatus     string       `json:"roaming_status"`
		UsimStatus        string       `json:"usim_status"`
		OmaDmVersion      string       `json:"oma_dm_version"`
		Plmn              string       `json:"plmn"`
		Band              string       `json:"band"`
		SignalRsrq        string       `json:"signal_rsrq"`
		SignalRsrp        string       `json:"signal_rsrp"`
		LteSinr           string       `json:"lte_sinr"`
		LteRssi           string       `json:"lte_rssi"`
		LteRsThroughput   string       `json:"lte_rs_throughput"`
		LteTsThroughput   string       `json:"lte_ts_throughput"`
		LtePhysicalCellid string       `json:"lte_physical_cellid"`
		ModemType         string       `json:"modem_type"`
		DrcCdmaEvdo       string       `json:"drc_cdma_evdo"`
		CurrentSnr        string       `json:"current_snr"`
		WirelessOperator  string       `json:"wireless_operator"`
		OperatingMode     string       `json:"operating_mode"`
		WirelessSignal    string       `json:"wireless_signal"`
		UsbWanMac         string       `json:"usb_wan_mac"`
		Sims              map[string]Sim
	}

	type System struct {
		CPU             float64 `json:"cpu"`
		Memory          float64 `json:"memory"`
		IP              string  `json:"ip"`
		SoftwareVersion string  `json:"software_version"`
		HardwareVersion string  `json:"hardware_version"`
		Mac             string  `json:"mac"`
		Netmask         string  `json:"netmask"`
		Gateway         string  `json:"gateway"`
		AddrType        string  `json:"addr_type"`
		FgtIP           string  `json:"fgt_ip"`
		GpsLat          string  `json:"gps_lat"`
		GpsLong         string  `json:"gps_long"`
	}

	type Extender struct {
		Name   string `json:"name"`
		ID     string `json:"id"`
		Vdom   string `json:"vdom"`
		System System `json:"system"`
		Modems map[string]Modem
	}

	type Results struct {
		//Name       string `json:"name"`
		Id string `json:"id"`
		//Authorized string `json:"authorized"`
	}

	type managedResponse []struct {
		Results []Results `json:"results"`
		Vdom    string    `json:"vdom"`
		Path    string    `json:"path"`
		Name    string    `json:"name"`
		Status  string    `json:"status"`
		Serial  string    `json:"serial"`
		Version string    `json:"version"`
		Build   int       `json:"build"`
	}

	type rawResponse []struct {
		Results []map[string]interface{} `json:"results"`
	}

	// Consider implementing pagination to remove this limit of 1000 entries
	var response managedResponse
	if err := c.Get("api/v2/cmdb/extender-controller/extender", "vdom=*", &response); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	var m []prometheus.Metric
	var extenders []Extender

	for _, rs := range response {

		// Check if an extender exists for the vdom
		exists := 0.0
		if len(rs.Results) > 0 {
			exists = 1.0
		}
		m = append(m, prometheus.MustNewConstMetric(extenderExist, prometheus.GaugeValue, exists, rs.Vdom))

		for _, extenderID := range rs.Results {

			var rawData rawResponse
			if err := c.Get("api/v2/monitor/extender-controller/extender", "vdom=*&id="+extenderID.Id, &rawData); err != nil {
				log.Printf("Error: %v", err)
				return nil, false
			}
			for _, e := range rawData {

				for _, extenderData := range e.Results {
					var system System
					s, _ := extenderData["system"]
					convertMapToStruct(s.(map[string]interface{}), &system)
					var extender = Extender{
						Name:   extenderData["name"].(string),
						ID:     extenderData["id"].(string),
						Vdom:   rs.Vdom,
						System: system,
						Modems: make(map[string]Modem),
					}
					extenders = append(extenders, extender)

					for extenderKey, extenderValue := range extenderData {
						// Since the response does not include an array of modems, but just keys like
						// modem1 and modem2 we must find keys in the response that match
						// modem[0-9]{1,2}
						// My understanding is that a extender can have max 2 modems, but the
						// regex support one or two numbers
						// The same logic is applied for sim cards below
						match, err := regexp.MatchString("^modem[0-9]{1,2}$", extenderKey)
						if err == nil && match {
							cdmaProfile := &CdmaProfile{}
							modem := &Modem{Sims: make(map[string]Sim), CdmaProfile: cdmaProfile}
							convertMapToStruct(extenderValue.(map[string]interface{}), modem)
							extender.Modems[extenderKey] = *modem

							for modemKey, modemValue := range extenderValue.(map[string]interface{}) {
								match, err := regexp.MatchString("^sim[0-9]{1,2}$", modemKey)
								if err == nil && match {
									// Found a sim
									sim := &Sim{}
									convertMapToStruct(modemValue.(map[string]interface{}), sim)
									modem.Sims[modemKey] = *sim
								}

								match, err = regexp.MatchString("^cdma_profile$", modemKey)
								if err == nil && match {
									// Found a cdma_profile
									//modemValue.(map[string]interface{})["status"] = "kalle"
									convertMapToStruct(modemValue.(map[string]interface{}), cdmaProfile)
								}
							}
						}
					}
				}
			}
		}
	}

	for _, extender := range extenders {
		m = append(m, prometheus.MustNewConstMetric(extenderInfo, prometheus.GaugeValue, 1, extender.Vdom,
			extender.Name, extender.System.SoftwareVersion, extender.System.HardwareVersion))
		m = append(m, prometheus.MustNewConstMetric(extenderCpu, prometheus.GaugeValue, extender.System.CPU,
			extender.Vdom, extender.Name))
		m = append(m, prometheus.MustNewConstMetric(extenderMemory, prometheus.GaugeValue, extender.System.Memory,
			extender.Vdom, extender.Name))

		for modemName, modem := range extender.Modems {
			m = append(m, prometheus.MustNewConstMetric(modemInfo, prometheus.GaugeValue, 1, extender.Vdom,
				extender.Name, modemName, modem.DataPlan, modem.Manufacturer, modem.Product, modem.Model, modem.Service,
				modem.EsnImei, modem.Band, modem.ModemType, modem.WirelessOperator, modem.LtePhysicalCellid))

			if value, err := strconv.ParseFloat(modem.SignalRsrq, 64); err == nil {
				m = append(m, prometheus.MustNewConstMetric(signalRsrq, prometheus.GaugeValue, value, extender.Vdom,
					extender.Name, modemName))
			}

			if value, err := strconv.ParseFloat(modem.SignalRsrp, 64); err == nil {
				m = append(m, prometheus.MustNewConstMetric(signalRsrp, prometheus.GaugeValue, value, extender.Vdom, extender.Name, modemName))
			}

			if value, err := strconv.ParseFloat(modem.LteSinr, 64); err == nil {
				m = append(m, prometheus.MustNewConstMetric(lteSinr, prometheus.GaugeValue, value, extender.Vdom, extender.Name, modemName))
			}

			if value, err := strconv.ParseFloat(modem.LteRssi, 64); err == nil {
				m = append(m, prometheus.MustNewConstMetric(lteRssi, prometheus.GaugeValue, value, extender.Vdom, extender.Name, modemName))
			}

			if value, err := strconv.ParseFloat(modem.SignalStrength, 64); err == nil {
				m = append(m, prometheus.MustNewConstMetric(signalStrength, prometheus.GaugeValue, value, extender.Vdom, extender.Name, modemName))
			}

			if modem.ConnectStatus == "CONN_STATE_CONNECTED" {
				m = append(m, prometheus.MustNewConstMetric(connectStatus, prometheus.GaugeValue, 1, extender.Vdom, extender.Name, modemName))
			} else {
				m = append(m, prometheus.MustNewConstMetric(connectStatus, prometheus.GaugeValue, 0, extender.Vdom, extender.Name, modemName))
			}

			for simName, sim := range modem.Sims {
				m = append(m, prometheus.MustNewConstMetric(simInfo, prometheus.GaugeValue, 1, extender.Vdom, extender.Name, modemName, simName, sim.Imsi, sim.Iccid))
				m = append(m, prometheus.MustNewConstMetric(simDataUsage, prometheus.GaugeValue, sim.DataUsage, extender.Vdom, extender.Name, modemName, simName))
				m = append(m, prometheus.MustNewConstMetric(simActive, prometheus.GaugeValue, sim.IsActive, extender.Vdom, extender.Name, modemName, simName))
				if sim.Status == "enable" {
					m = append(m, prometheus.MustNewConstMetric(simStatus, prometheus.GaugeValue, 1, extender.Vdom, extender.Name, modemName, simName))
				} else {
					m = append(m, prometheus.MustNewConstMetric(simStatus, prometheus.GaugeValue, 0, extender.Vdom, extender.Name, modemName, simName))
				}
			}

		}
	}
	return m, true
}
