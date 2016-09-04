package news

import (
	"fmt"
	"hash/crc32"
	"log"

	"github.com/arthurkiller/my_mp/db"
	"github.com/arthurkiller/my_mp/grpc/news"
	"github.com/garyburd/redigo/redis"
)

type NewsServer interface {
	GetNews(news.GetNewsRequest) news.GetNewsReply
	GetMyNews(news.GetNewsRequest) news.GetNewsReply
	PostNews(news.PostNewsRequest) news.PostNewsReply
	RecallNews(news.RecallNewsRequest) news.RecallNewsReply
	LikeNews(news.LikeNewsRequest) news.LikeNewsReply
}

type newsServer struct {
	redisPoll redism.Redism
}

func NewNewsServer(c redism.RedismConf, redisPoll redism.Redism) NewsServer {
	server := new(newsServer)
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

func (s *newsServer) GetNews(req news.GetNewsRequest) news.GetNewsReply {
	conn := s.redisPoll.Get("uid-box", "R", req.Uid)
	defer conn.Close()

	index := req.Index
	val, err := redis.Values(conn.Do("HSCAN", req.Uid, "COUNT", 10))
	if err != nil {
		log.Println("error in do redis get userbox :", err)
		return news.GetNewsReply{Status: 1}
	}
	//get the news-id list
	newsidList := make([]string, 20)
	_, err = redis.Scan(val, &index, newsidList)
	if err != nil {
		log.Println("error in scan newslist:", err)
		return news.GetNewsReply{Status: 1}
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
			return news.GetNewsReply{Status: 1}
		}
		newsMap[newsinf.Uid] = newsinf
		userlist[i] = newsinf.Uid
	}

	//check the fans-uid list
	conn = s.redisPoll.Get("uid-fans", "R", req.Uid)
	scr := s.redisPoll.GetScript("uid-check-fans")
	keyargs := make([]interface{}, 0)
	keyargs = append(keyargs, req.Uid, fmt.Sprintf("%s", len(userlist)))
	for _, vv := range userlist {
		keyargs = append(keyargs, vv)
	}
	fanslist, err := redis.Strings(scr.Do(conn, keyargs...))
	if err != nil {
		log.Println("error in redis do check fans : ", err)
		return news.GetNewsReply{Status: 1}
	}

	result := news.GetNewsReply{}
	result.Status = 0
	result.Index = index
	result.News = make([]*news.NewsInfo, 20)
	for i, v := range fanslist {
		result.News[i] = newsMap[v]
	}
	return result
}

func (s *newsServer) GetMyNews(req news.GetNewsRequest) news.GetNewsReply {
	conn := s.redisPoll.Get("self-box", "R", req.Uid)
	defer conn.Close()

	index := req.Index
	val, err := redis.Values(conn.Do("HSCAN", req.Uid, "COUNT", 10))
	if err != nil {
		log.Println("error in do redis get userbox :", err)
		return news.GetNewsReply{Status: 0}
	}
	newsList := make([]string, 10)
	_, err = redis.Scan(val, &index, newsList)
	if err != nil {
		log.Println("error in scan newslist:", err)
		return news.GetNewsReply{Status: 0}
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
	return result
}

func (s *newsServer) PostNews(req news.PostNewsRequest) news.PostNewsReply {
	//the rule of gengeric a newsid use uid + devid + timestamp to generic a sha265 for the newsid
	newsID := crc32.ChecksumIEEE([]byte(req.Uid + req.Devid + req.TimeStamp))

	conn := s.redisPoll.Get("newsid-info", "W", fmt.Sprint(newsID))
	defer conn.Close()

	_, err := conn.Do("HMSET", fmt.Sprint(newsID))
	if err != nil {
		log.Println("error in hmset the message with messageid")
		return news.PostNewsReply{Status: 0}
	}

	//TODO this should have a cache
	conn = s.redisPoll.Get("uid-fans", "R", req.Uid)
	index := int64(0)
	uidList := make([]string, 100)
	for {
		v, err := redis.Values(conn.Do("SSCAN", req.Uid, &index, "COUNT", 100))
		if err != nil {
			log.Println("error in do redis get user fans hscan:", err)
			return news.PostNewsReply{Status: 0}
		}
		redis.Scan(v, &index, &uidList)

		for _, v := range uidList {
			connuid := s.redisPoll.Get("uid-box", "W", v)
			defer connuid.Close()
			_, err := connuid.Do("LPUSH", v, newsID)
			if err != nil {
				log.Println("error in lpush message into box : ", err)
				return news.PostNewsReply{Status: 0}
			}
		}
		if index == 0 {
			break
		}
	}
	return news.PostNewsReply{Status: 1, Newsid: fmt.Sprint(newsID)}
}

func (s *newsServer) RecallNews(req news.RecallNewsRequest) news.RecallNewsReply {
	conn := s.redisPoll.Get("newsid-info", "W", req.Newsid)
	defer conn.Close()

	_, err := conn.Do("DEL", req.Newsid)
	if err != nil {
		log.Println("error in delet the newsid key: ", err)
		return news.RecallNewsReply{Status: 0}
	}
	return news.RecallNewsReply{Status: 1}
}

func (s *newsServer) LikeNews(req news.LikeNewsRequest) news.LikeNewsReply {
	conn := s.redisPoll.Get("newsid-info", "W", req.Newsid)
	defer conn.Close()

	_, err := conn.Do("HINCRBY", req.Newsid, "Likes", 1)
	if err != nil {
		log.Println("error in increase the newsid : ", err)
		return news.LikeNewsReply{Status: 0}
	}
	return news.LikeNewsReply{Status: 1}
}
