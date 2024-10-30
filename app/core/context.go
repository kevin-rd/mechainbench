package core

import "time"

type Context struct {
	// 每批次发送数量
	Quantity uint

	BatchNum uint

	Timeout time.Duration
}
