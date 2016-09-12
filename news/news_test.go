package newsm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/arthurkiller/my_mp/grpc/news"
)

func Test_GetNews(t *testing.T) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := news.GetNewsRequest{
		Uid:   "a1",
		Index: uint64(0),
	}
	result, err := cli.GetNews(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Test_GetMyNews(t *testing.T) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := news.GetNewsRequest{
		Uid:   "a1",
		Index: uint64(0),
	}
	result, err := cli.GetMyNews(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Test_PostNews(t *testing.T) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := news.PostNewsRequest{
		Uid:       "a1",
		Devid:     "iphone",
		TimeStamp: fmt.Sprint(time.Now().UnixNano()),
		MeipaiID:  "m1",
		Values:    []byte("hello,world"),
	}
	result, err := cli.PostNews(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Test_RecallNews(t *testing.T) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := news.RecallNewsRequest{}
	result, err := cli.RecallNews(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Test_LikeNews(t *testing.T) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := news.LikeNewsRequest{}
	result, err := cli.LikeNews(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
