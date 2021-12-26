/**
 *  Author: SongLee24
 *  Email: lisong.shine@qq.com
 *  Date: 2018-08-15
 *
 *
 *  prometheus.Desc是指标的描述符，用于实现对指标的管理
 *
 */

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"prometheus-exporter/mysql"
	"sync"
	"prometheus-exporter/utils"
)

// 指标结构体
type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

/**
 * 函数：newGlobalMetric
 * 功能：创建指标描述符
 */
func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}


/**
 * 工厂方法：NewMetrics
 * 功能：初始化指标信息，即Metrics结构体
 */
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"slow_sql_metric": newGlobalMetric(namespace, "slow_sql_metric", "show slow sql ", []string{"node","user","sql","host","db"}),

		},
	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()  // 加锁
	defer c.mutex.Unlock()

	slowSqlMetricData := c.GenerateData()
	for _,data := range slowSqlMetricData {
		ch <-prometheus.MustNewConstMetric(c.metrics["slow_sql_metric"], prometheus.CounterValue, float64(data.TIME), utils.GetEnvWithDefault("MYSQL_HOST","127.0.0.1"),data.USER,data.INFO,data.HOST,data.DB)
	}

}
/***
 * 函数：GenerateData
 * 功能：获取数据
 */
func (c *Metrics) GenerateData() [] mysql.SolwLog{
	my := mysql.Init()
	my.Connect()
	result, _ := my.RunSql(mysql.SHOW_SLOW_SQL)

	if(len(result)!= 0){
		log.Println(result)
		return result
	}
	var ss [] mysql.SolwLog
	return ss

}


///**
// * 函数：GenerateMockData
// * 功能：生成模拟数据
// */
// func (c *Metrics) GenerateMockData() (mockCounterMetricData map[string]int, mockGaugeMetricData map[string]int) {
// 	mockCounterMetricData = map[string]int{
//		"yahoo.com": int(rand.Int31n(1000)),
//		"google.com": int(rand.Int31n(1000)),
//	}
//	mockGaugeMetricData = map[string]int{
//		"yahoo.com": int(rand.Int31n(10)),
//		"google.com": int(rand.Int31n(10)),
//	}
//	return
// }

