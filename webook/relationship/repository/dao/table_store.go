package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"gorm.io/gorm"
)

const FollowRelationTableName = "follow_relations"

var (
	ErrFollowerNotFound = gorm.ErrRecordNotFound
)

type TableStoreFollowRelationDao struct {
	client *tablestore.TableStoreClient
}

func (t *TableStoreFollowRelationDao) FollowRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error) {
	request := &tablestore.SQLQueryRequest{
		// tablestore 的 API 不允许 SQL 里面带参数
		// SELECT * FROM `users` WHERE id = ?; []any{123}
		// select id,follower,followee from follow_relations where follower = 1 OR 1=1
		// follower 前端输入了 1 OR 1=1
		// 核心原则就是，千万不要用前端输入的数据来拼接 SQL
		// 最好就是用占位符
		// WHERE username = 'abc' or 1=1 AND password = '';
		// WHERE username = 'abc' AND password = ''; TRUNCATE table mysql.`users`;
		Query: fmt.Sprintf("select id,follower,followee from %s where follower = %d AND status = %d OFFSET %d LIMIT %d",
			FollowRelationTableName, follower, FollowRelationStatusActive, offset, limit),
	}
	response, err := t.client.SQLQuery(request)
	if err != nil {
		return nil, err
	}
	resultSet := response.ResultSet
	followRelations := make([]FollowRelation, 0, limit)
	for resultSet.HasNext() {
		row := resultSet.Next()
		followRelation := FollowRelation{}
		followRelation.Follower, _ = row.GetInt64ByName("follower")
		followRelation.Followee, _ = row.GetInt64ByName("followee")
		followRelations = append(followRelations, followRelation)
	}
	return followRelations, nil
}

func (t *TableStoreFollowRelationDao) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	cond := tablestore.NewCompositeColumnCondition(tablestore.LO_AND)
	cond.AddFilter(tablestore.NewSingleColumnCondition("follower", tablestore.CT_EQUAL, follower))
	cond.AddFilter(tablestore.NewSingleColumnCondition("followee", tablestore.CT_EQUAL, followee))
	req := new(tablestore.UpdateRowChange)
	req.TableName = FollowRelationTableName
	req.SetCondition(tablestore.RowExistenceExpectation_EXPECT_EXIST)
	req.SetColumnCondition(cond)
	req.PutColumn("status", int64(status))
	_, err := t.client.UpdateRow(&tablestore.UpdateRowRequest{
		UpdateRowChange: req,
	})
	return err
}

func (t *TableStoreFollowRelationDao) CntFollower(ctx context.Context, uid int64) (int64, error) {
	request := &tablestore.SQLQueryRequest{
		Query: fmt.Sprintf("SELECT COUNT(follower) as cnt from %s where followee = %d AND status = %d",
			FollowRelationTableName, uid, FollowRelationStatusActive)}
	response, err := t.client.SQLQuery(request)
	if err != nil {
		return 0, err
	}
	resultSet := response.ResultSet
	if resultSet.HasNext() {
		row := resultSet.Next()
		return row.GetInt64ByName("cnt")
	}
	return 0, ErrFollowerNotFound
}

func (t *TableStoreFollowRelationDao) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	request := &tablestore.SQLQueryRequest{
		Query: fmt.Sprintf("SELECT COUNT(followee) as cnt from %s where follower = %d AND status = %d",
			FollowRelationTableName, uid, FollowRelationStatusActive)}
	response, err := t.client.SQLQuery(request)
	if err != nil {
		return 0, err
	}
	resultSet := response.ResultSet
	if resultSet.HasNext() {
		row := resultSet.Next()
		return row.GetInt64ByName("cnt")
	}
	return 0, ErrFollowerNotFound
}

func (t *TableStoreFollowRelationDao) CreateFollowRelation(ctx context.Context, c FollowRelation) error {
	// 创建关注关系
	// 需要的是 insert or update 语义
	// 从哪里来？UpdateRowRequest + RowExistenceExpectation_IGNORE
	req := new(tablestore.UpdateRowRequest)
	change := new(tablestore.UpdateRowChange)
	change.TableName = FollowRelationTableName
	putPk := new(tablestore.PrimaryKey)
	putPk.AddPrimaryKeyColumn("follower", c.Follower)
	putPk.AddPrimaryKeyColumn("followee", c.Followee)
	change.PrimaryKey = putPk
	now := time.Now()
	change.PutColumn("status", int64(FollowRelationStatusActive))
	change.PutColumn("utime", now.Unix())
	change.PutColumn("ctime", now.Unix())
	// 如果要是冲突了就忽略掉
	change.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
	req.UpdateRowChange = change
	_, err := t.client.UpdateRow(req)
	return err
}

func (t *TableStoreFollowRelationDao) FollowRelationDetail(ctx context.Context, follower, followee int64) (FollowRelation, error) {
	request := &tablestore.SQLQueryRequest{
		Query: fmt.Sprintf("select id,follower,followee from %s where follower = %d AND followee = %d AND status = %d",
			FollowRelationTableName, follower, followee, FollowRelationStatusActive)}
	response, err := t.client.SQLQuery(request)
	if err != nil {
		return FollowRelation{}, err
	}
	resultSet := response.ResultSet
	if resultSet.HasNext() {
		row := resultSet.Next()
		return t.rowToEntity(row), nil
	}
	return FollowRelation{}, ErrFollowerNotFound
}

func (t *TableStoreFollowRelationDao) rowToEntity(row tablestore.SQLRow) FollowRelation {
	var res FollowRelation
	res.ID, _ = row.GetInt64ByName("id")
	res.Follower, _ = row.GetInt64ByName("follower")
	res.Followee, _ = row.GetInt64ByName("followee")
	status, _ := row.GetInt64ByName("status")
	res.Status = uint8(status)
	res.Ctime, _ = row.GetInt64ByName("ctime")
	res.Utime, _ = row.GetInt64ByName("utime")
	return res
}

func NewTableStoreDao(client *tablestore.TableStoreClient) *TableStoreFollowRelationDao {
	return &TableStoreFollowRelationDao{
		client: client,
	}
}
