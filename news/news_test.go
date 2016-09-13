package newsm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/arthurkiller/my_mp/grpc/news"
)

func makec() (news.NewsClient, error) {
	conn, err := grpc.Dial("127.0.0.1:23588", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cli := news.NewNewsClient(conn)
	return cli, nil
}

func Test_GetNews(t *testing.T) {
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
	t.Log("first post a news", result)

	req1 := news.GetNewsRequest{
		Uid:   "a1",
		Index: uint64(0),
	}
	result1, err := cli.GetNews(context.Background(), &req1)
	if err != nil {
		t.Error(err)
	}
	t.Log(result1)
}

func Test_GetMyNews(t *testing.T) {
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
	t.Log("first post a news", result)

	req1 := news.GetNewsRequest{
		Uid:   "a1",
		Index: uint64(0),
	}
	result1, err := cli.GetMyNews(context.Background(), &req1)
	if err != nil {
		t.Error(err)
	}
	t.Log(result1)
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
	t.Log("first post a news", result)

	req1 := news.RecallNewsRequest{result.Newsid}
	result1, err := cli.RecallNews(context.Background(), &req1)
	if err != nil {
		t.Error(err)
	}
	t.Log(result1)
}

func Test_LikeNews(t *testing.T) {
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
	t.Log("first post a news", result)

	req1 := news.LikeNewsRequest{result.Newsid}
	result1, err := cli.LikeNews(context.Background(), &req1)
	if err != nil {
		t.Error(err)
	}
	t.Log(result1)
}

func Benchmark_GetNews(b *testing.B) {
	cli, err := makec()
	if err != nil {
		b.Error(err)
	}
	req := news.GetNewsRequest{
		Uid:   "a1",
		Index: uint64(0),
	}
	for i := 0; i < b.N; i++ {
		_, err := cli.GetNews(context.Background(), &req)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_GetMyNews(b *testing.B) {
	cli, err := makec()
	if err != nil {
		b.Error(err)
	}
	req := news.GetNewsRequest{
		Uid:   "a1",
		Index: uint64(0),
	}
	for i := 0; i < b.N; i++ {
		_, err := cli.GetMyNews(context.Background(), &req)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_PostNews(b *testing.B) {
	cli, err := makec()
	if err != nil {
		b.Error(err)
	}
	req := news.PostNewsRequest{
		Uid:       "a1",
		Devid:     "iphone",
		TimeStamp: fmt.Sprint(time.Now().UnixNano()),
		MeipaiID:  "m1",
		Values:    []byte("hello,world"),
	}
	for i := 0; i < b.N; i++ {
		_, err := cli.PostNews(context.Background(), &req)
		if err != nil {
			b.Error(err)
		}
	}
}
