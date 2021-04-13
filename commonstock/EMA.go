package commonstock

import (
	"time"

	"github.com/idoall/TokenExchangeCommon/commonmodels"
)

// EMA struct
type EMA struct {
	Alpha    float64
	Period   int //默认计算几天的EMA
	points   []*EMAPoint
	kline    []*commonmodels.Kline
	LastTime int64
}

type EMAPoint struct {
	point
}

// NewEMA new Func
func NewEMA(list []*commonmodels.Kline, period int) *EMA {
	m := &EMA{kline: list, Period: period, Alpha: 2.0 / (float64(period) + 1.0)}
	return m
}

// Calculation Func
func (e *EMA) Calculation() *EMA {
	for _, v := range e.kline {
		e.Add(v.KlineTime, v.Close)
	}
	return e
}

// GetPoints return Point
func (e *EMA) GetPoints() []*EMAPoint {
	return e.points
}

// Add adds a new Value to Ema
// 使用方法，先添加最早日期的数据,最后一条应该是当前日期的数据，结果与 AICoin 对比完全一致
func (e *EMA) Add(timestamp time.Time, value float64) {
	var p *EMAPoint
	emaTminusOne := value
	lastTime := timestamp.Unix()
	if lastTime == e.LastTime { //相同时间只是计算并更新最后一个point
		p = e.points[len(e.points)-1]
		emaTminusOne = e.points[len(e.points)-2].Value
	} else { //不同时间则是往后计算point
		if len(e.points) > 0 {
			emaTminusOne = e.points[len(e.points)-1].Value
		}
		p = new(EMAPoint)
		p.Time = timestamp
		e.points = append(e.points, p)
	}

	//平滑指数，一般取作2/(N+1)
	//alpha := 2.0 / (float64(e.Period) + 1.0)

	// fmt.Println(alpha)

	// 计算 EMA指数
	emaT := e.Alpha*value + (1-e.Alpha)*emaTminusOne
	p.Value = emaT
	e.LastTime = e.points[len(e.points)-1].Time.Unix()
}
