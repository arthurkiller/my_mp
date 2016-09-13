package profilem

import (
	"log"

	"github.com/arthurkiller/my_mp/db"
	"github.com/arthurkiller/my_mp/grpc/profile"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/context"
)

//TODO: SCAN returns may not shorter than the array given to the redis.Scan()

type ProfileServer interface {
	GetUserInfo(context.Context, *profile.GetUserInfoRequest) (*profile.GetUserInfoReply, error)
	GetFans(context.Context, *profile.GetFansRequest) (*profile.GetFansReply, error)
	GetFollow(context.Context, *profile.GetFollowRequest) (*profile.GetFollowReply, error)
	AddFollow(context.Context, *profile.AddFollowRequest) (*profile.AddFollowReply, error)
	DeleteFollow(context.Context, *profile.DeleteFollowRequest) (*profile.DeleteFollowReply, error)
}

type profileServer struct {
	redisPoll redism.Redism
}

func NewProfileServer(redisPoll redism.Redism) ProfileServer {
	server := new(profileServer)
	server.redisPoll = redisPoll

	return server
}

func (p *profileServer) GetUserInfo(ctx context.Context, req *profile.GetUserInfoRequest) (*profile.GetUserInfoReply, error) {
	conn := p.redisPoll.Get("uid-info", "R", req.Uid)
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETALL", req.Uid))
	if err != nil {
		log.Println("error in do redis get userinfo hgetall:", err)
		return &profile.GetUserInfoReply{
			Status: 1,
			Info:   &profile.UserInfo{},
		}, err
	}

	var info profile.UserInfo
	err = redis.ScanStruct(v, &info)
	if err != nil {
		log.Println("Error in scan struct:", err)
		return &profile.GetUserInfoReply{
			Status: 1,
			Info:   &profile.UserInfo{},
		}, err
	}

	result := profile.GetUserInfoReply{
		Status: 0,
		Info:   &info,
	}
	return &result, nil
}

func (p *profileServer) GetFans(ctx context.Context, req *profile.GetFansRequest) (*profile.GetFansReply, error) {
	conn := p.redisPoll.Get("uid-fans", "R", req.Uid)
	defer conn.Close()

	index := req.Index
	uidList := make([]string, 0, 20)
	v, err := redis.Values(conn.Do("SSCAN", req.Uid, index, "COUNT", 10))
	if err != nil {
		log.Println("error in do redis get user fans SSCAN:", err)
		return &profile.GetFansReply{
			Status: 1,
			Index:  0,
			Fanss:  make([]*profile.UserInfo, 1),
		}, err
	}
	_, err = redis.Scan(v, &index, &uidList)
	if err != nil {
		log.Println("error in do redis scan", err)
		return &profile.GetFansReply{
			Status: 1,
			Index:  0,
			Fanss:  make([]*profile.UserInfo, 1),
		}, err
	}

	result := profile.GetFansReply{}
	result.Status = 0
	result.Index = index
	result.Fanss = make([]*profile.UserInfo, 0, 20)
	for _, u := range uidList {
		connUID := p.redisPoll.Get("uid-info", "R", u)
		vv, err := redis.Values(connUID.Do("HGETALL", u))
		defer connUID.Close()
		if err != nil {
			log.Println("Error in Get userinfo hgetall:", err)
			return &profile.GetFansReply{
				Status: 1,
				Index:  0,
				Fanss:  make([]*profile.UserInfo, 1),
			}, err
		}

		var info profile.UserInfo
		err = redis.ScanStruct(vv, &info)
		if err != nil {
			log.Println("Error in scan struct:", err)
			return &profile.GetFansReply{
				Status: 1,
				Index:  0,
				Fanss:  make([]*profile.UserInfo, 1),
			}, err
		}
		result.Fanss = append(result.Fanss, &info)
	}
	return &result, nil
}

func (p *profileServer) GetFollow(ctx context.Context, req *profile.GetFollowRequest) (*profile.GetFollowReply, error) {
	conn := p.redisPoll.Get("uid-follow", "R", req.Uid)
	defer conn.Close()

	index := int64(req.Index)
	uidList := make([]string, 0, 20)
	v, err := redis.Values(conn.Do("SSCAN", req.Uid, index, "COUNT", 10))
	if err != nil {
		log.Println("error in do redis get user follow hscan:", err)
		return &profile.GetFollowReply{
			Status:  1,
			Index:   0,
			Follows: make([]*profile.UserInfo, 1),
		}, err
	}
	_, err = redis.Scan(v, &index, &uidList)
	if err != nil {
		log.Println("error in do redis scan", err)
		return &profile.GetFollowReply{
			Status:  1,
			Index:   0,
			Follows: make([]*profile.UserInfo, 1),
		}, err
	}

	result := profile.GetFollowReply{}
	result.Status = 0
	result.Index = index
	result.Follows = make([]*profile.UserInfo, 0, 20)
	for _, u := range uidList {
		connUID := p.redisPoll.Get("uid-info", "R", u)
		defer connUID.Close()
		vv, err := redis.Values(connUID.Do("HGETALL", u))
		if err != nil {
			log.Println("Error in Getuserinfo :", err)
			return &profile.GetFollowReply{
				Status:  1,
				Index:   0,
				Follows: make([]*profile.UserInfo, 1),
			}, err
		}
		var info profile.UserInfo
		err = redis.ScanStruct(vv, &info)
		if err != nil {
			log.Println("Error in scan struct:", err)
			return &profile.GetFollowReply{
				Status:  1,
				Index:   0,
				Follows: make([]*profile.UserInfo, 1),
			}, err
		}
		result.Follows = append(result.Follows, &info)
	}
	return &result, nil
}

func (p *profileServer) AddFollow(ctx context.Context, req *profile.AddFollowRequest) (*profile.AddFollowReply, error) {
	conn := p.redisPoll.Get("uid-follow", "W", req.Uid)
	defer conn.Close()

	n, err := conn.Do("SADD", req.Uid, req.DestUid)
	if err != nil {
		log.Println("error in do redis add follows:", err)
		return &profile.AddFollowReply{Status: 1}, err
	}

	if n == 0 {
		return &profile.AddFollowReply{Status: 1}, err
	}
	return &profile.AddFollowReply{Status: 0}, nil
}

func (p *profileServer) DeleteFollow(ctx context.Context, req *profile.DeleteFollowRequest) (*profile.DeleteFollowReply, error) {
	conn := p.redisPoll.Get("uid-follow", "W", req.Uid)
	defer conn.Close()

	n, err := conn.Do("SREM", req.Uid, req.DestUid)
	if err != nil {
		log.Println("error in do redis remove follow:", err)
		return &profile.DeleteFollowReply{Status: 1}, err
	}

	if n == 0 {
		return &profile.DeleteFollowReply{Status: 1}, err
	}
	return &profile.DeleteFollowReply{Status: 0}, nil
}
