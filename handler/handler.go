package handler

import (
	"context"
	"fmt"
	"github.com/dy-gopkg/kit/dao/redis"
	"github.com/dy-platform/user-srv-session/idl"
	pb "github.com/dy-platform/user-srv-session/idl/platform/user/srv-session"
	rgo "github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
}

//Record 登记用户session
func (h *Handler) Record(ctx context.Context, req *pb.RecordReq, resp *pb.RecordResp) error {
	uuid := uuid.NewV4()
	resp.BaseResp = new(base.Resp)
	_, err := redis.Do("set", uuid, req.UID)
	if err != nil {
		resp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		resp.BaseResp.Msg = err.Error()
		logrus.Errorf("redis error %v", err)
		return nil
	}
	if req.ExpireTime != -1 {
		_, err = redis.Do("expire", uuid, req.ExpireTime)
		resp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		resp.BaseResp.Msg = err.Error()

		logrus.Errorf("redis error %v", err)
		return nil
	}

	resp.BaseResp.Code = int32(base.CODE_OK)
	resp.BaseResp.Msg = "success"
	return nil
}

//Refresh 刷新token
func (h *Handler) Refresh(ctx context.Context, req *pb.RefreshReq, resp *pb.RefreshResp) error {
	_, err := redis.Do("expire", req.Token, req.ExpireTime)
	resp.BaseResp = new(base.Resp)
	if err != nil {
		resp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		resp.BaseResp.Msg = err.Error()

		logrus.Errorf("redis error %v", err)
		return nil
	}
	return nil
}

//Query 通过token 获得用户ID
func (h *Handler) Query(ctx context.Context, req *pb.QueryReq, resp *pb.QueryResp) error {
	value, err := redis.Do("get", req.Token)
	uid, err := rgo.Int64(value, err)
	resp.BaseResp = new(base.Resp)
	if err != nil {
		resp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		resp.BaseResp.Msg = fmt.Sprintf("query redis error %v", err.Error())
		logrus.Errorf("redis error %v", err)
		return nil
	}
	resp.UID = uid
	resp.BaseResp.Msg = "success"
	resp.BaseResp.Code = int32(base.CODE_OK)
	return nil
}
