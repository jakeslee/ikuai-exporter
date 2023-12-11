package pkg

import (
	"fmt"
	"github.com/jakeslee/ikuai"
	"github.com/jakeslee/ikuai/action"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strconv"
	"time"
)

type IKuaiExporter struct {
	ikuai       *ikuai.IKuai
	versionDesc *prometheus.Desc // ikuai 版本

	// CPU
	cpuUsageRatioDesc *prometheus.Desc // CPU 使用
	cpuTempDesc       *prometheus.Desc // CPU 温度

	// 内存
	memSizeDesc    *prometheus.Desc // 内存指标
	memUsageDesc   *prometheus.Desc // 内存指标
	memCachedDesc  *prometheus.Desc // 内存指标
	memBuffersDesc *prometheus.Desc // 内存指标

	// 终端
	lanDeviceDesc      *prometheus.Desc // 内网终端信息
	lanDeviceCountDesc *prometheus.Desc // 内网终端数量
	ifaceInfoDesc      *prometheus.Desc // 接口信息
	UpDesc             *prometheus.Desc // 在线状态，host/link
	UpTimeDesc         *prometheus.Desc // 在线时间，host/link

	// 网络，device/host/iface
	streamUpBytesDesc   *prometheus.Desc // 流量上行数据包
	streamDownBytesDesc *prometheus.Desc // 流量上行数据包
	streamUpSpeedDesc   *prometheus.Desc // 流量上行速度
	streamDownSpeedDesc *prometheus.Desc // 流量上行速度
	connCountDesc       *prometheus.Desc // 连接数指标
}

func NewIKuaiExporter(kuai *ikuai.IKuai) *IKuaiExporter {
	return &IKuaiExporter{
		ikuai: kuai,
		versionDesc: prometheus.NewDesc("ikuai_version", "IKuai version info",
			[]string{"version", "arch", "verstring"}, nil),
		cpuUsageRatioDesc: prometheus.NewDesc("ikuai_cpu_usage_ratio", "IKuai CPU usage ratio",
			[]string{"id"}, nil),
		cpuTempDesc: prometheus.NewDesc("ikuai_cpu_temperature", "",
			nil, nil),
		memSizeDesc: prometheus.NewDesc("ikuai_memory_size_bytes", "",
			[]string{}, nil),
		memUsageDesc: prometheus.NewDesc("ikuai_memory_usage_bytes", "",
			[]string{}, nil),
		memCachedDesc: prometheus.NewDesc("ikuai_memory_cached_bytes", "",
			[]string{}, nil),
		memBuffersDesc: prometheus.NewDesc("ikuai_memory_buffers_bytes", "",
			[]string{}, nil),
		lanDeviceDesc: prometheus.NewDesc("ikuai_device_info", "ikuai_device_info",
			[]string{"id", "mac", "hostname", "ip_addr", "comment"}, nil),
		lanDeviceCountDesc: prometheus.NewDesc("ikuai_device_count", "",
			[]string{}, nil),
		ifaceInfoDesc: prometheus.NewDesc("ikuai_iface_info", "",
			[]string{"id", "interface", "comment", "internet", "parent_interface", "ip_addr"}, nil),
		UpDesc: prometheus.NewDesc("ikuai_up", "",
			[]string{"id"}, nil),
		UpTimeDesc: prometheus.NewDesc("ikuai_uptime", "",
			[]string{"id"}, nil),
		streamUpBytesDesc: prometheus.NewDesc("ikuai_network_send_bytes", "",
			[]string{"id"}, nil),
		streamDownBytesDesc: prometheus.NewDesc("ikuai_network_recv_bytes", "",
			[]string{"id"}, nil),
		streamUpSpeedDesc: prometheus.NewDesc("ikuai_network_send_kbytes_per_second", "",
			[]string{"id"}, nil),
		streamDownSpeedDesc: prometheus.NewDesc("ikuai_network_recv_kbytes_per_second", "",
			[]string{"id"}, nil),
		connCountDesc: prometheus.NewDesc("ikuai_network_conn_count", "",
			[]string{"id"}, nil),
	}
}

func (i *IKuaiExporter) Describe(descs chan<- *prometheus.Desc) {
	descs <- i.versionDesc
	descs <- i.cpuUsageRatioDesc
	descs <- i.cpuTempDesc
	descs <- i.memSizeDesc
	descs <- i.memUsageDesc
	descs <- i.memCachedDesc
	descs <- i.memBuffersDesc
	descs <- i.lanDeviceDesc
	descs <- i.lanDeviceCountDesc
	descs <- i.ifaceInfoDesc
	descs <- i.UpDesc
	descs <- i.UpTimeDesc
	descs <- i.streamUpBytesDesc
	descs <- i.streamDownBytesDesc
	descs <- i.streamUpSpeedDesc
	descs <- i.streamDownSpeedDesc
	descs <- i.connCountDesc
}

func (i *IKuaiExporter) Collect(metrics chan<- prometheus.Metric) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("collect ikuai panic, %v", err)

			metrics <- prometheus.MustNewConstMetric(i.UpDesc, prometheus.GaugeValue, 0,
				"host")
		}
	}()

	stat, err := i.ikuai.ShowSysStat()

	if isFail(&stat.Result, err) {
		log.Printf("ikuai ShowSysStat: %v, %+v", err, stat.Result)
		panic(stat.Result)
	}

	sysStat := stat.Data.SysStat

	metrics <- prometheus.MustNewConstMetric(i.versionDesc, prometheus.GaugeValue, 1,
		sysStat.Verinfo.Version,
		sysStat.Verinfo.Arch,
		sysStat.Verinfo.Verstring)

	if len(sysStat.Cputemp) > 0 {
		metrics <- prometheus.MustNewConstMetric(i.cpuTempDesc, prometheus.GaugeValue, float64(sysStat.Cputemp[0]))
	} else {
		log.Printf("sysStat.Cputemp is empty")
	}

	for idx, item := range sysStat.Cpu {
		s := item[:len(item)-1]
		per, _ := strconv.ParseFloat(s, 64)

		metrics <- prometheus.MustNewConstMetric(i.cpuUsageRatioDesc, prometheus.GaugeValue, per/100,
			fmt.Sprintf("core/%v", idx))
	}

	metrics <- prometheus.MustNewConstMetric(i.memSizeDesc, prometheus.GaugeValue, float64(sysStat.Memory.Total))
	metrics <- prometheus.MustNewConstMetric(i.memUsageDesc, prometheus.GaugeValue,
		float64(sysStat.Memory.Total-sysStat.Memory.Available))
	metrics <- prometheus.MustNewConstMetric(i.memCachedDesc, prometheus.GaugeValue, float64(sysStat.Memory.Cached))
	metrics <- prometheus.MustNewConstMetric(i.memBuffersDesc, prometheus.GaugeValue, float64(sysStat.Memory.Buffers))

	lanDevice, err := i.ikuai.ShowMonitorLan()

	if isFail(&lanDevice.Result, err) {
		log.Printf("ikuai ShowMonitorLan: %v, %+v", err, lanDevice.Result)
	} else {
		devices := map[string]action.LanDeviceInfo{}

		for _, device := range lanDevice.Data.Data {
			deviceId := fmt.Sprintf("device/%v", device.IPAddr)

			if _, ok := devices[deviceId]; !ok {
				devices[deviceId] = device
			}
		}

		for deviceId, device := range devices {
			metrics <- prometheus.MustNewConstMetric(i.lanDeviceDesc, prometheus.GaugeValue, 1,
				deviceId, device.Mac, device.Hostname, device.IPAddr, device.Comment)

			metrics <- prometheus.MustNewConstMetric(i.streamUpBytesDesc, prometheus.GaugeValue, float64(device.TotalUp),
				deviceId)

			metrics <- prometheus.MustNewConstMetric(i.streamDownBytesDesc, prometheus.GaugeValue, float64(device.TotalDown),
				deviceId)

			metrics <- prometheus.MustNewConstMetric(i.streamUpSpeedDesc, prometheus.GaugeValue, float64(device.Upload),
				deviceId)

			metrics <- prometheus.MustNewConstMetric(i.streamDownSpeedDesc, prometheus.GaugeValue, float64(device.Download),
				deviceId)

			metrics <- prometheus.MustNewConstMetric(i.connCountDesc, prometheus.GaugeValue, float64(device.ConnectNum),
				deviceId)
		}
	}

	metrics <- prometheus.MustNewConstMetric(i.lanDeviceCountDesc, prometheus.GaugeValue, float64(sysStat.OnlineUser.Count))

	monitorInterface, err := i.ikuai.ShowMonitorInterface()

	if isFail(&monitorInterface.Result, err) {
		log.Printf("ikuai ShowMonitorInterface: %v, %+v", err, monitorInterface.Result)
	} else {
		i.interfaceMetrics(metrics, monitorInterface)
	}

	// Host metric
	metrics <- prometheus.MustNewConstMetric(i.UpTimeDesc, prometheus.GaugeValue, float64(sysStat.Uptime),
		"host")

	metrics <- prometheus.MustNewConstMetric(i.streamUpBytesDesc, prometheus.GaugeValue, float64(sysStat.Stream.TotalUp),
		"host")

	metrics <- prometheus.MustNewConstMetric(i.streamDownBytesDesc, prometheus.GaugeValue, float64(sysStat.Stream.TotalDown),
		"host")

	metrics <- prometheus.MustNewConstMetric(i.streamUpSpeedDesc, prometheus.GaugeValue, float64(sysStat.Stream.Upload),
		"host")

	metrics <- prometheus.MustNewConstMetric(i.streamDownSpeedDesc, prometheus.GaugeValue, float64(sysStat.Stream.Download),
		"host")

	metrics <- prometheus.MustNewConstMetric(i.connCountDesc, prometheus.GaugeValue, float64(sysStat.Stream.ConnectNum),
		"host")

	// 无报错，up
	metrics <- prometheus.MustNewConstMetric(i.UpDesc, prometheus.GaugeValue, 1,
		"host")
}

func (i *IKuaiExporter) interfaceMetrics(metrics chan<- prometheus.Metric, monitorInterface *action.ShowMonitorInterfaceResult) {
	for _, iface := range monitorInterface.Data.IfaceStream {
		internet := ""
		parentIface := ""
		ifaceUp := 1
		ifaceId := fmt.Sprintf("iface/%v", iface.Interface)
		ifaceUptime := int64(0)

		for _, ifaceCheck := range monitorInterface.Data.IfaceCheck {
			if ifaceCheck.Interface == iface.Interface {
				internet = ifaceCheck.Internet
				parentIface = ifaceCheck.ParentInterface

				if ifaceCheck.Result != "success" {
					ifaceUp = 0
				} else {
					updateTime, err := strconv.ParseInt(ifaceCheck.Updatetime, 10, 64)
					if err == nil {
						ifaceUptime = time.Now().Unix() - updateTime
					}
				}
			}
		}

		metrics <- prometheus.MustNewConstMetric(i.ifaceInfoDesc, prometheus.GaugeValue, 1,
			ifaceId, iface.Interface, iface.Comment, internet, parentIface, iface.IPAddr)

		metrics <- prometheus.MustNewConstMetric(i.UpDesc, prometheus.GaugeValue, float64(ifaceUp),
			ifaceId)

		metrics <- prometheus.MustNewConstMetric(i.UpTimeDesc, prometheus.GaugeValue, float64(ifaceUptime),
			ifaceId)

		metrics <- prometheus.MustNewConstMetric(i.streamUpBytesDesc, prometheus.GaugeValue, float64(iface.TotalUp),
			ifaceId)

		metrics <- prometheus.MustNewConstMetric(i.streamDownBytesDesc, prometheus.GaugeValue, float64(iface.TotalDown),
			ifaceId)

		metrics <- prometheus.MustNewConstMetric(i.streamUpSpeedDesc, prometheus.GaugeValue, float64(iface.Upload),
			ifaceId)

		metrics <- prometheus.MustNewConstMetric(i.streamDownSpeedDesc, prometheus.GaugeValue, float64(iface.Download),
			ifaceId)

		ifaceConn, nErr := strconv.ParseInt(iface.ConnectNum, 10, 8)
		if nErr != nil {
			ifaceConn = 0
		}

		metrics <- prometheus.MustNewConstMetric(i.connCountDesc, prometheus.GaugeValue, float64(ifaceConn),
			ifaceId)
	}
}

func isFail(result *action.Result, err error) bool {
	if err != nil || result.ErrMsg != "Success" {
		return true
	}
	return false
}
