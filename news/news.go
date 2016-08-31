package news

import (
	"github.com/arthurkiller/my_mp/db"
	"github.com/arthurkiller/my_mp/grpc/news"
)

type NewsServer interface {
	GetNews(news.GetNewsRequest) news.GetNewsReply
	GetMyNews(news.GetNewsRequest) news.GetNewsReply
	PostNews(news.PostNewsRequest) news.PsotNewsReply
	RecallNews(news.RecallNewsRequest) news.RecallNewsReply
	LikeNews(news.LikeNewsRequest) news.LikeNewsReply
	ForwardNews(news.ForwardNewsRequest) news.ForwardNewsReply
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

func (n *newsServer) GetNews(req news.GetNewsRequest) news.GetNewsReply {
}

func (n *newsServer) GetMyNews(req news.GetNewsRequest) news.GetNewsReply
func (n *newsServer) PostNews(req news.PostNewsRequest) news.PsotNewsReply
func (n *newsServer) RecallNews(req news.RecallNewsRequest) news.RecallNewsReply
func (n *newsServer) LikeNews(req news.LikeNewsRequest) news.LikeNewsReply
func (n *newsServer) ForwardNews(req news.ForwardNewsRequest) news.ForwardNewsReply
