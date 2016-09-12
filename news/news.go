package newsm

import (
	"fmt"
	"hash/crc32"
	"log"

	"github.com/arthurkiller/my_mp/db"
	"github.com/arthurkiller/my_mp/grpc/news"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/context"
)

type NewsServer interface {
	GetNews(context.Context, *news.GetNewsRequest) (*news.GetNewsReply, error)
	GetMyNews(context.Context, *news.GetNewsRequest) (*news.GetNewsReply, error)
	PostNews(context.Context, *news.PostNewsRequest) (*news.PostNewsReply, error)
	RecallNews(context.Context, *news.RecallNewsRequest) (*news.RecallNewsReply, error)
	LikeNews(context.Context, *news.LikeNewsRequest) (*news.LikeNewsReply, error)
}

type newsServer struct {
	redisPoll redism.Redism
}

func NewNewsServer(redisPoll redism.Redism) NewsServer {
	server := new(newsServer)
	server.redisPoll = redisPoll

	return server
}

func (s *newsServer) GetNews(ctx context.Context, req *news.GetNewsRequest) (*news.GetNewsReply, error) {
	conn := s.redisPoll.Get("uid-box", "R", req.Uid)
	defer conn.Close()

	index := req.Index
	val, err := redis.Values(conn.Do("HSCAN", req.Uid, "COUNT", 10))
	if err != nil {
		log.Println("error in do redis get userbox :", err)
		return &news.GetNewsReply{Status: 1}, err
	}
	//get the news-id list
	newsidList := make([]string, 20)
	_, err = redis.Scan(val, &index, newsidList)
	if err != nil {
		log.Println("error in scan newslist:", err)
		return &news.GetNewsReply{Status: 1}, err
	}

	//get the newsid -> news-info
	userlist := make([]string, 20)
	newsMap := make(map[string]*news.NewsInfo, 20)
	for i, v := range newsidList {
		conn = s.redisPoll.Get("newsid-info", "R", v)
		newsinf := new(news.NewsInfo)
		vs, _ := redis.Values(conn.Do("HGETALL", v))
		err = redis.ScanStruct(vs, &newsinf)
		if err != nil {
			log.Println("error in scan struct", err)
			return &news.GetNewsReply{Status: 1}, err
		}
		newsMap[newsinf.Uid] = newsinf
		userlist[i] = newsinf.Uid
	}

	//check the fans-uid list
	conn = s.redisPoll.Get("uid-fans", "R", req.Uid)
	scr := s.redisPoll.GetScript("uid-check-fans")
	var keyargs []interface{}
	keyargs = append(keyargs, req.Uid, fmt.Sprint(len(userlist)))
	for _, vv := range userlist {
		keyargs = append(keyargs, vv)
	}
	fanslist, err := redis.Strings(scr.Do(conn, keyargs...))
	if err != nil {
		log.Println("error in redis do check fans : ", err)
		return &news.GetNewsReply{Status: 1}, err
	}

	result := news.GetNewsReply{}
	result.Status = 0
	result.Index = index
	result.News = make([]*news.NewsInfo, 20)
	for i, v := range fanslist {
		result.News[i] = newsMap[v]
	}
	return &result, nil
}

func (s *newsServer) GetMyNews(ctx context.Context, req *news.GetNewsRequest) (*news.GetNewsReply, error) {
	conn := s.redisPoll.Get("uid-selfbox", "R", req.Uid)
	defer conn.Close()

	index := req.Index
	val, err := redis.Values(conn.Do("HSCAN", req.Uid, "COUNT", 10))
	if err != nil {
		log.Println("error in do redis get userbox :", err)
		return &news.GetNewsReply{Status: 0}, err
	}
	newsList := make([]string, 10)
	_, err = redis.Scan(val, &index, newsList)
	if err != nil {
		log.Println("error in scan newslist:", err)
		return &news.GetNewsReply{Status: 0}, err
	}

	newsMap := make(map[string]*news.NewsInfo, 10) // newsid -> news{}
	for _, v := range newsList {
		conn = s.redisPoll.Get("newsid-info", "R", v)
		news := news.NewsInfo{}
		vs, _ := redis.Values(conn.Do("HGETALL", v))
		redis.ScanStruct(vs, &news)
		newsMap[v] = &news
	}

	result := news.GetNewsReply{}
	result.Status = 0
	result.Index = index
	result.News = make([]*news.NewsInfo, 10)
	for i, v := range newsList {
		result.News[i] = newsMap[v]
	}
	return &result, nil
}

func (s *newsServer) PostNews(ctx context.Context, req *news.PostNewsRequest) (*news.PostNewsReply, error) {
	//the rule of gengeric a newsid use uid + devid + timestamp to generic a sha265 for the newsid
	newsID := crc32.ChecksumIEEE([]byte(req.Uid + req.Devid + req.TimeStamp))

	conn := s.redisPoll.Get("newsid-info", "W", fmt.Sprint(newsID))
	defer conn.Close()

	_, err := conn.Do("HMSET", fmt.Sprint(newsID))
	if err != nil {
		log.Println("error in hmset the message with messageid")
		return &news.PostNewsReply{Status: 1}, err
	}

	//TODO this should have a cache
	conn = s.redisPoll.Get("uid-fans", "R", req.Uid)
	index := int64(0)
	uidList := make([]string, 100)
	for {
		v, err := redis.Values(conn.Do("SSCAN", req.Uid, &index, "COUNT", 100))
		if err != nil {
			log.Println("error in do redis get user fans hscan:", err)
			return &news.PostNewsReply{Status: 1}, err
		}
		redis.Scan(v, &index, &uidList)

		for _, v := range uidList {
			connuid := s.redisPoll.Get("uid-box", "W", v)
			defer connuid.Close()
			_, err := connuid.Do("LPUSH", v, newsID)
			if err != nil {
				log.Println("error in lpush message into box : ", err)
				return &news.PostNewsReply{Status: 1}, err
			}
		}
		if index == 0 {
			break
		}
	}
	return &news.PostNewsReply{Status: 0, Newsid: fmt.Sprint(newsID)}, nil
}

func (s *newsServer) RecallNews(ctx context.Context, req *news.RecallNewsRequest) (*news.RecallNewsReply, error) {
	conn := s.redisPoll.Get("newsid-info", "W", req.Newsid)
	defer conn.Close()

	_, err := conn.Do("DEL", req.Newsid)
	if err != nil {
		log.Println("error in delet the newsid key: ", err)
		return &news.RecallNewsReply{Status: 1}, err
	}
	return &news.RecallNewsReply{Status: 0}, nil
}

func (s *newsServer) LikeNews(ctx context.Context, req *news.LikeNewsRequest) (*news.LikeNewsReply, error) {
	conn := s.redisPoll.Get("newsid-info", "W", req.Newsid)
	defer conn.Close()

	_, err := conn.Do("HINCRBY", req.Newsid, "Likes", 1)
	if err != nil {
		log.Println("error in increase the newsid : ", err)
		return &news.LikeNewsReply{Status: 1}, err
	}
	return &news.LikeNewsReply{Status: 0}, nil
}
