package profile

import (
	"log"

	"github.com/arthurkiller/my_mp/db"
	"github.com/arthurkiller/my_mp/grpc/profile"
	"github.com/garyburd/redigo/redis"
)

type ProfileServer interface {
	GetUserInfo(profile.GetUserInfoRequest) profile.GetUserInfoReply
	GetFans(profile.GetFansRequest) profile.GetFansReply
	GetFollow(profile.GetFollowRequest) profile.GetFollowReply
	AddFollow(profile.AddFollowRequest) profile.AddFollowReply
	DeleteFollow(profile.DeleteFollowRequest) profile.DeleteFollowReply
}

type profileServer struct {
	redisPoll redism.Redism
}

func NewProfileServer(c redism.RedismConf, redisPoll redism.Redism) ProfileServer {
	server := new(profileServer)
	//redism.RedismConf{
	//	Maxactive:   500,
	//	Maxidle:     300,
	//	Idletimeout: 60,
	//	Masters:     []string{"127.0.0.1:6379"},
	//	Slaves:      [][]string{{"127.0.0.1:6379", "127.0.0.1:6379"}},
	//}
	server.redisPoll = redisPoll

	return server
}

func (p *profileServer) GetUserInfo(req profile.GetUserInfoRequest) profile.GetUserInfoReply {
	id := req.Uid
	conn := p.redisPoll.Get("uid-info", "R", id)
	v, err := redis.Values(conn.Do("HGETALL", id))
	if err != nil {
		log.Println("error in do redis:", err)
	}
	defer conn.Close()
	if err != nil {
		log.Println("Error in Getuserinfo :", err)
	}

	var info profile.UserInfo
	err = redis.ScanStruct(v, &info)
	if err != nil {
		log.Println("Error in scan struct:", err)
	}

	result := profile.GetUserInfoReply{
		Status: 0,
		Info:   &info,
	}
	return result
}

func (p *profileServer) GetFans(req profile.GetFansRequest) profile.GetFansReply {
	conn := p.redisPoll.Get("uid-fans", "R", req.Uid)
	defer conn.Close()
	index := int64(req.Index)
	uidList := make([]string, 10)
	v, err := redis.Values(conn.Do("HSCAN", req.Uid, &index, 10))
	if err != nil {
		log.Println("error in do redis:", err)
	}
	redis.Scan(v, &index, &uidList)

	var result profile.GetFansReply
	result.Status = 0
	result.Index = index
	result.Fanss = make([]*profile.UserInfo, 10)
	for _, u := range uidList {
		connUid := p.redisPoll.Get("uid-info", "R", u)
		vv, err := redis.Values(connUid.Do("HGETALL", u))
		defer connUid.Close()
		if err != nil {
			log.Println("Error in Getuserinfo :", err)
		}
		var info profile.UserInfo
		err = redis.ScanStruct(vv, &info)
		if err != nil {
			log.Println("Error in scan struct:", err)
		}
		result.Fanss = append(result.Fanss, &info)
	}
	return result
}

func (p *profileServer) GetFollow(req profile.GetFollowRequest) profile.GetFollowReply {
	conn := p.redisPoll.Get("uid-follow", "R", req.Uid)
	defer conn.Close()
	index := int64(req.Index)
	uidList := make([]string, 10)
	v, err := redis.Values(conn.Do("HSCAN", req.Uid, &index, 10))
	if err != nil {
		log.Println("error in do redis:", err)
	}
	redis.Scan(v, &index, &uidList)

	var result profile.GetFollowReply
	result.Status = 0
	result.Index = index
	result.Follows = make([]*profile.UserInfo, 10)
	for _, u := range uidList {
		connUid := p.redisPoll.Get("uid-info", "R", u)
		vv, err := redis.Values(connUid.Do("HGETALL", u))
		defer connUid.Close()
		if err != nil {
			log.Println("Error in Getuserinfo :", err)
		}
		var info profile.UserInfo
		err = redis.ScanStruct(vv, &info)
		if err != nil {
			log.Println("Error in scan struct:", err)
		}
		result.Follows = append(result.Follows, &info)
	}
	return result
}

func (p *profileServer) AddFollow(req profile.AddFollowRequest) profile.AddFollowReply {
	conn := p.redisPoll.Get("uid-follow", "W", req.Uid)
	defer conn.Close()

	n, err := conn.Do("SADD", req.Uid, req.DestUid)
	if err != nil {
		log.Println("error in do redis:", err)
	}

	if n == 0 {
		return profile.AddFollowReply{Status: 1}
	}
	return profile.AddFollowReply{Status: 0}

}
func (p *profileServer) DeleteFollow(req profile.DeleteFollowRequest) profile.DeleteFollowReply {
	conn := p.redisPoll.Get("uid-follow", "W", req.Uid)
	defer conn.Close()

	n, err := conn.Do("SREM", req.Uid, req.DestUid)
	if err != nil {
		log.Println("error in do redis:", err)
	}

	if n == 0 {
		return profile.DeleteFollowReply{Status: 1}
	}
	return profile.DeleteFollowReply{Status: 0}
}
