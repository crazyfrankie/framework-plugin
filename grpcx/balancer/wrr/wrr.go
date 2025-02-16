package wrr

/*
平滑加权负载均衡算法
*/
import (
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const Name = "crazyfrankie_wrr"

func newBuilder() balancer.Builder {
	// NewBalancerBuilder 将一个 Pick Builder 转化为了一个 BalancerBuilder
	return base.NewBalancerBuilder(Name, &PickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type PickerBuilder struct{}

func (p PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*conn, 0, len(info.ReadySCs))
	for sc, sci := range info.ReadySCs {
		conn := &conn{
			cc: sc,
		}
		md, ok := sci.Address.Metadata.(map[string]any)
		if ok {
			weightVal := md["weight"].(float64)
			conn.weight = int(weightVal)
			conn.currWeight = conn.weight
		}
		if conn.weight == 0 {
			conn.weight = 10
			conn.currWeight = 10
		}
		conns = append(conns, conn)
	}
	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	conns []*conn
	mux   sync.Mutex
}

// Pick 真正实现加权负载均衡算法的地方
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	total := 0
	// 计算总权重 当前权重 挑选节点
	target := p.conns[0]
	for _, c := range p.conns {
		total += c.weight
		c.currWeight = c.currWeight + c.weight
		if target == nil || target.currWeight < c.currWeight {
			target = c
		}
	}

	target.currWeight = target.currWeight - total

	return balancer.PickResult{
		SubConn: target.cc,
		//Done: func(info balancer.DoneInfo) {
		// 回调函数
		// 很多动态算法，就是在这里根据调用结果来调整权重
		//},
	}, nil
}

// conn 是节点的抽象
type conn struct {
	weight     int              // 权重
	currWeight int              // 当前权重
	cc         balancer.SubConn // 真正的 grpc 中对节点的抽象
}
