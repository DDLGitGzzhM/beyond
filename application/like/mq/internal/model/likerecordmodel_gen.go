// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	likeRecordFieldNames          = builder.RawFieldNames(&LikeRecord{})
	likeRecordRows                = strings.Join(likeRecordFieldNames, ",")
	likeRecordRowsExpectAutoSet   = strings.Join(stringx.Remove(likeRecordFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	likeRecordRowsWithPlaceHolder = strings.Join(stringx.Remove(likeRecordFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	likeRecordModel interface {
		Insert(ctx context.Context, data *LikeRecord) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*LikeRecord, error)
		FindOneByBizIdObjIdUserId(ctx context.Context, bizId string, objId int64, userId int64) (*LikeRecord, error)
		Update(ctx context.Context, data *LikeRecord) error
		Delete(ctx context.Context, id int64) error
	}

	defaultLikeRecordModel struct {
		conn  sqlx.SqlConn
		table string
	}

	LikeRecord struct {
		Id         int64     `db:"id"`          // 主键ID
		BizId      string    `db:"biz_id"`      // 业务ID
		ObjId      int64     `db:"obj_id"`      // 点赞对象id
		UserId     int64     `db:"user_id"`     // 用户ID
		LikeType   int64     `db:"like_type"`   // 类型 0:点赞 1:点踩
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 最后修改时间
	}
)

func newLikeRecordModel(conn sqlx.SqlConn) *defaultLikeRecordModel {
	return &defaultLikeRecordModel{
		conn:  conn,
		table: "`like_record`",
	}
}

func (m *defaultLikeRecordModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultLikeRecordModel) FindOne(ctx context.Context, id int64) (*LikeRecord, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", likeRecordRows, m.table)
	var resp LikeRecord
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLikeRecordModel) FindOneByBizIdObjIdUserId(ctx context.Context, bizId string, objId int64, userId int64) (*LikeRecord, error) {
	var resp LikeRecord
	query := fmt.Sprintf("select %s from %s where `biz_id` = ? and `obj_id` = ? and `user_id` = ? limit 1", likeRecordRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, bizId, objId, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLikeRecordModel) Insert(ctx context.Context, data *LikeRecord) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, likeRecordRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.BizId, data.ObjId, data.UserId, data.LikeType)
	return ret, err
}

func (m *defaultLikeRecordModel) Update(ctx context.Context, newData *LikeRecord) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, likeRecordRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.BizId, newData.ObjId, newData.UserId, newData.LikeType, newData.Id)
	return err
}

func (m *defaultLikeRecordModel) tableName() string {
	return m.table
}
