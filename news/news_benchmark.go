package newsm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/arthurkiller/my_mp/grpc/news"
	"google.golang.org/grpc"
)

func makec() (news.NewsClient, error) {
	conn, err := grpc.Dial("127.0.0.1:23588", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cli := news.NewNewsClient(conn)
	return cli, nil
}

func Benchmark_GetNews(t *testing.B) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := news.GetNewsRequest{
		Uid:   "a1",
		Index: uint64(0),
	}
	for i := 0; i < t.N; i++ {
		_, err := cli.GetNews(context.Background(), &req)
		if err != nil {
			t.Error(err)
		}
	}
}

func Benchmark_GetMyNews(t *testing.B) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := news.GetNewsRequest{
		Uid:   "a1",
		Index: uint64(0),
	}
	for i := 0; i < t.N; i++ {
		_, err := cli.GetMyNews(context.Background(), &req)
		if err != nil {
			t.Error(err)
		}
	}
}

func Benchmark_PostNews(t *testing.B) {
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
	for i := 0; i < t.N; i++ {
		_, err := cli.PostNews(context.Background(), &req)
		if err != nil {
			t.Error(err)
		}
	}
}

func Benchmark_RecallNews(t *testing.B) {
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

func Benchmark_LikeNews(t testing.B) {
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
