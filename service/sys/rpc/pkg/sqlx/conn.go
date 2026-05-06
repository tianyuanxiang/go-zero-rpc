package sqlx

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// TimeoutConn 是对 sqlx.SqlConn 的包装器。
//
// 所有带 Ctx 后缀的方法均会忽略调用方传入的 context（通常是 HTTP 请求上下文，
// 默认超时 3s），改用以 context.Background() 为父的独立超时上下文。
//
// 这样做的目的：go-zero REST 框架默认给每个 HTTP 请求附加较短的超时
// （RestConf.Timeout，默认 3000ms），而数据库查询（特别是远程 DB 或连接池
// 冷启动阶段）可能超出该限制。通过此包装器，DB 操作获得独立、可配置的超时，
// 不再受 HTTP 请求生命周期影响。
//
// 注意：此包装器会屏蔽调用方通过 ctx 传递的取消信号（cancel）。
// 对于绝大多数内部查询场景这是可以接受的。
type TimeoutConn struct {
	sqlx.SqlConn
	timeout time.Duration
}

// NewTimeoutConn 创建一个带统一超时的 SqlConn 包装器。
//
// @param conn    原始 go-zero sqlx 连接
// @param timeout 每次 DB 操作的最大等待时间，建议设置为 10s
func NewTimeoutConn(conn sqlx.SqlConn, timeout time.Duration) sqlx.SqlConn {
	return &TimeoutConn{
		SqlConn: conn,
		timeout: timeout,
	}
}

// newCtx 创建一个以 context.Background() 为父、带固定超时的新上下文。
func (c *TimeoutConn) newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

// ExecCtx 覆盖原始方法，使用独立超时上下文执行写操作。
func (c *TimeoutConn) ExecCtx(ctx context.Context, q string, args ...any) (sql.Result, error) {
	newCtx, cancel := c.newCtx()
	defer cancel()
	return c.SqlConn.ExecCtx(newCtx, q, args...)
}

// PrepareCtx 覆盖原始方法，使用独立超时上下文准备 SQL 语句。
func (c *TimeoutConn) PrepareCtx(ctx context.Context, q string) (sqlx.StmtSession, error) {
	newCtx, cancel := c.newCtx()
	defer cancel()
	return c.SqlConn.PrepareCtx(newCtx, q)
}

// QueryRowCtx 覆盖原始方法，使用独立超时上下文查询单行。
func (c *TimeoutConn) QueryRowCtx(ctx context.Context, v any, q string, args ...any) error {
	newCtx, cancel := c.newCtx()
	defer cancel()
	return c.SqlConn.QueryRowCtx(newCtx, v, q, args...)
}

// QueryRowPartialCtx 覆盖原始方法，使用独立超时上下文查询单行（部分字段）。
func (c *TimeoutConn) QueryRowPartialCtx(ctx context.Context, v any, q string, args ...any) error {
	newCtx, cancel := c.newCtx()
	defer cancel()
	return c.SqlConn.QueryRowPartialCtx(newCtx, v, q, args...)
}

// QueryRowsCtx 覆盖原始方法，使用独立超时上下文查询多行。
func (c *TimeoutConn) QueryRowsCtx(ctx context.Context, v any, q string, args ...any) error {
	newCtx, cancel := c.newCtx()
	defer cancel()
	return c.SqlConn.QueryRowsCtx(newCtx, v, q, args...)
}

// QueryRowsPartialCtx 覆盖原始方法，使用独立超时上下文查询多行（部分字段）。
func (c *TimeoutConn) QueryRowsPartialCtx(ctx context.Context, v any, q string, args ...any) error {
	newCtx, cancel := c.newCtx()
	defer cancel()
	return c.SqlConn.QueryRowsPartialCtx(newCtx, v, q, args...)
}

// TransactCtx 覆盖原始方法，使用独立超时上下文执行事务。
//
// 注意：事务内部的所有操作共享此超时，请确保 timeout 足够覆盖整个事务。
func (c *TimeoutConn) TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	newCtx, cancel := c.newCtx()
	defer cancel()
	return c.SqlConn.TransactCtx(newCtx, fn)
}
