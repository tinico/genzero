// Code generated by genzero. DO NOT EDIT.

package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var JwtBlacklistTableName = "jwt_blacklist"

type jwtBlacklist_model interface {
	FindCount(ctx context.Context) (int64, error)
	FindAll(ctx context.Context) ([]*JwtBlacklist, error)
	FindList(ctx context.Context, pageSize, page int64, keyword string, jwtBlacklist *JwtBlacklist) (resp []*JwtBlacklist, total int64, err error)
	TableName() string
	FindById(ctx context.Context, id int64) (*JwtBlacklist, error)
	FindByAdminerId(ctx context.Context, adminerId int64) (*JwtBlacklist, error)
	FindByUuid(ctx context.Context, uuid string) (*JwtBlacklist, error)
	FindByToken(ctx context.Context, token string) (*JwtBlacklist, error)
	FindByPlatform(ctx context.Context, platform string) (*JwtBlacklist, error)
	FindByIp(ctx context.Context, ip string) (*JwtBlacklist, error)
	FindByExpireAt(ctx context.Context, expireAt time.Time) (*JwtBlacklist, error)
	FindsByIds(ctx context.Context, ids []int64) ([]*JwtBlacklist, error)
	FindsByAdminerIds(ctx context.Context, adminerIds []int64) ([]*JwtBlacklist, error)
	FindsByUuids(ctx context.Context, uuids []string) ([]*JwtBlacklist, error)
	FindsByTokens(ctx context.Context, tokens []string) ([]*JwtBlacklist, error)
	FindsByPlatforms(ctx context.Context, platforms []string) ([]*JwtBlacklist, error)
	FindsByIps(ctx context.Context, ips []string) ([]*JwtBlacklist, error)
	FindsByExpireAts(ctx context.Context, expireAts []time.Time) ([]*JwtBlacklist, error)

	FindsById(ctx context.Context, id int64) ([]*JwtBlacklist, error)
	FindsByAdminerId(ctx context.Context, adminerId int64) ([]*JwtBlacklist, error)
	FindsByUuid(ctx context.Context, uuid string) ([]*JwtBlacklist, error)
	FindsByToken(ctx context.Context, token string) ([]*JwtBlacklist, error)
	FindsByPlatform(ctx context.Context, platform string) ([]*JwtBlacklist, error)
	FindsByIp(ctx context.Context, ip string) ([]*JwtBlacklist, error)
	FindsByExpireAt(ctx context.Context, expireAt time.Time) ([]*JwtBlacklist, error)

	formatUuidKey(uuid string) string
}

func (m *defaultJwtBlacklistModel) FindCount(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) as count from %s", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query)
	return count, err
}

func (m *defaultJwtBlacklistModel) FindAll(ctx context.Context) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	query := fmt.Sprintf("select %s from %s limit 99999", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJwtBlacklistModel) FindList(ctx context.Context, pageSize, page int64, keyword string, jwtBlacklist *JwtBlacklist) (resp []*JwtBlacklist, total int64, err error) {
	sq := squirrel.Select(jwtBlacklistRows).From(m.table)
	if jwtBlacklist != nil {
		if jwtBlacklist.Id >= 0 {
			sq = sq.Where("`id` = ?", jwtBlacklist.Id)
		}
		if jwtBlacklist.AdminerId > 0 {
			sq = sq.Where("`adminer_id` = ?", jwtBlacklist.AdminerId)
		}
		if jwtBlacklist.Uuid != "" {
			sq = sq.Where("`uuid` = ?", jwtBlacklist.Uuid)
		}
		if jwtBlacklist.Token != "" {
			sq = sq.Where("`token` = ?", jwtBlacklist.Token)
		}
		if jwtBlacklist.Platform != "" {
			sq = sq.Where("`platform` = ?", jwtBlacklist.Platform)
		}
		if jwtBlacklist.Ip != "" {
			sq = sq.Where("`ip` = ?", jwtBlacklist.Ip)
		}
		if !jwtBlacklist.ExpireAt.IsZero() && !jwtBlacklist.ExpireAt.Equal(time.Unix(0, 0)) {
			sq = sq.Where("`expire_at` = ?", jwtBlacklist.ExpireAt.Format("2006-01-02 15:04:05"))
		}
	}
	if pageSize > 0 && page > 0 {
		sqCount := sq.RemoveLimit().RemoveOffset()
		sq = sq.Limit(uint64(pageSize)).Offset(uint64((page - 1) * pageSize))
		queryCount, agrsCount, e := sqCount.ToSql()
		if e != nil {
			err = e
			return
		}
		queryCount = strings.ReplaceAll(queryCount, jwtBlacklistRows, "COUNT(*)")
		if err = m.QueryRowNoCacheCtx(ctx, &total, queryCount, agrsCount...); err != nil {
			return
		}
	}
	query, agrs, err := sq.ToSql()
	if err != nil {
		return
	}
	resp = make([]*JwtBlacklist, 0)
	if err = m.QueryRowsNoCacheCtx(ctx, &resp, query, agrs...); err != nil {
		return
	}
	return
}

func (m *defaultJwtBlacklistModel) TableName() string {
	return m.table
}

func (m *defaultJwtBlacklistModel) FindById(ctx context.Context, id int64) (*JwtBlacklist, error) {
	var resp JwtBlacklist
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", jwtBlacklistRows, m.table)
	jwtBlacklistIdKey := fmt.Sprintf("%s%v", cacheJwtBlacklistIdPrefix, id)
	err := m.QueryRowCtx(ctx, &resp, jwtBlacklistIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJwtBlacklistModel) FindByAdminerId(ctx context.Context, adminerId int64) (*JwtBlacklist, error) {
	var resp JwtBlacklist
	query := fmt.Sprintf("select %s from %s where `adminer_id` = ? limit 1", jwtBlacklistRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, adminerId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJwtBlacklistModel) FindByUuid(ctx context.Context, uuid string) (*JwtBlacklist, error) {
	var resp JwtBlacklist
	query := fmt.Sprintf("select %s from %s where `uuid` = ? limit 1", jwtBlacklistRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, uuid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJwtBlacklistModel) FindByToken(ctx context.Context, token string) (*JwtBlacklist, error) {
	var resp JwtBlacklist
	query := fmt.Sprintf("select %s from %s where `token` = ? limit 1", jwtBlacklistRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, token)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJwtBlacklistModel) FindByPlatform(ctx context.Context, platform string) (*JwtBlacklist, error) {
	var resp JwtBlacklist
	query := fmt.Sprintf("select %s from %s where `platform` = ? limit 1", jwtBlacklistRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, platform)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJwtBlacklistModel) FindByIp(ctx context.Context, ip string) (*JwtBlacklist, error) {
	var resp JwtBlacklist
	query := fmt.Sprintf("select %s from %s where `ip` = ? limit 1", jwtBlacklistRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, ip)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJwtBlacklistModel) FindByExpireAt(ctx context.Context, expireAt time.Time) (*JwtBlacklist, error) {
	var resp JwtBlacklist
	query := fmt.Sprintf("select %s from %s where `expire_at` = ? limit 1", jwtBlacklistRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, expireAt.Format("2006-01-02 15:04:05"))
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJwtBlacklistModel) FindsByIds(ctx context.Context, ids []int64) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	if len(ids) == 0 {
		return resp, nil
	}
	query := fmt.Sprintf("select %s from %s where `id` in (?)", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, ids)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByAdminerIds(ctx context.Context, adminerIds []int64) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	if len(adminerIds) == 0 {
		return resp, nil
	}
	query := fmt.Sprintf("select %s from %s where `adminer_id` in (?)", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, adminerIds)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByUuids(ctx context.Context, uuids []string) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	if len(uuids) == 0 {
		return resp, nil
	}
	query := fmt.Sprintf("select %s from %s where `uuid` in (?)", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uuids)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByTokens(ctx context.Context, tokens []string) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	if len(tokens) == 0 {
		return resp, nil
	}
	query := fmt.Sprintf("select %s from %s where `token` in (?)", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, tokens)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByPlatforms(ctx context.Context, platforms []string) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	if len(platforms) == 0 {
		return resp, nil
	}
	query := fmt.Sprintf("select %s from %s where `platform` in (?)", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, platforms)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByIps(ctx context.Context, ips []string) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	if len(ips) == 0 {
		return resp, nil
	}
	query := fmt.Sprintf("select %s from %s where `ip` in (?)", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, ips)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByExpireAts(ctx context.Context, expireAts []time.Time) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	if len(expireAts) == 0 {
		return resp, nil
	}
	query := fmt.Sprintf("select %s from %s where `expire_at` in (?)", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, expireAts)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsById(ctx context.Context, id int64) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	query := fmt.Sprintf("select %s from %s where `id` = ? ", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, id)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByAdminerId(ctx context.Context, adminerId int64) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	query := fmt.Sprintf("select %s from %s where `adminer_id` = ? ", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, adminerId)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByUuid(ctx context.Context, uuid string) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	query := fmt.Sprintf("select %s from %s where `uuid` = ? ", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uuid)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByToken(ctx context.Context, token string) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	query := fmt.Sprintf("select %s from %s where `token` = ? ", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, token)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByPlatform(ctx context.Context, platform string) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	query := fmt.Sprintf("select %s from %s where `platform` = ? ", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, platform)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByIp(ctx context.Context, ip string) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	query := fmt.Sprintf("select %s from %s where `ip` = ? ", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, ip)
	return resp, err
}

func (m *defaultJwtBlacklistModel) FindsByExpireAt(ctx context.Context, expireAt time.Time) ([]*JwtBlacklist, error) {
	var resp = make([]*JwtBlacklist, 0)
	query := fmt.Sprintf("select %s from %s where `expire_at` = ? ", jwtBlacklistRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, expireAt)
	return resp, err
}

func (m *defaultJwtBlacklistModel) formatUuidKey(uuid string) string {
	return fmt.Sprintf("cache:jwtBlacklist:uuid:%v", uuid)
}
