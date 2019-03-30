package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dy-gopkg/kit/dao/redis"
	"github.com/dy-platform/user-srv-session/idl"
	pb "github.com/dy-platform/user-srv-session/idl/platform/user/srv-session"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
}

type User struct {
	Uid int64 `json:"uid"`
}

//Record 登记用户session
func (h *Handler) Record(ctx context.Context, req *pb.RecordReq, resp *pb.RecordResp) error {
	uuid := uuid.NewV4()
	resp.BaseResp = new(base.Resp)
	user := User{
		Uid: req.UID,
	}
	value, _ := json.Marshal(user)
	_, err := redis.Do("SET", uuid.String(), string(value))
	if err != nil {
		resp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		resp.BaseResp.Msg = err.Error()
		logrus.Errorf("redis set error %v", err)
		return nil
	}
	if req.ExpireTime != -1 {
		_, err = redis.Do("expire", uuid.String(), req.ExpireTime)
		if err != nil {
			resp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
			resp.BaseResp.Msg = err.Error()
			logrus.Errorf("redis expire error %v", err)
			return nil
		}
	}
	resp.Token = uuid.String()
	resp.BaseResp.Code = int32(base.CODE_OK)
	resp.BaseResp.Msg = "success"

	logrus.Debug("resp:", resp)
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

	resp.BaseResp.Code = int32(base.CODE_OK)
	resp.BaseResp.Msg = "success"

	return nil
}

//Query 通过token 获得用户ID
func (h *Handler) Query(ctx context.Context, req *pb.QueryReq, resp *pb.QueryResp) error {
	value, err := redis.Do("get", req.Token)
	user := &User{}
	json.Unmarshal(value.([]byte), user)
	resp.BaseResp = new(base.Resp)
	if err != nil {
		resp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		resp.BaseResp.Msg = fmt.Sprintf("query redis error %v", err.Error())
		logrus.Errorf("redis error %v", err)
		return nil
	}
	resp.UID = user.Uid
	resp.BaseResp.Msg = "success"
	resp.BaseResp.Code = int32(base.CODE_OK)
	return nil
}
