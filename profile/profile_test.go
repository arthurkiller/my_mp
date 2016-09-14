package profilem

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/arthurkiller/my_mp/grpc/profile"
)

func Makec() (profile.ProfileClient, error) {
	conn, err := grpc.Dial("127.0.0.1:23589", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cli := profile.NewProfileClient(conn)
	return cli, nil
}

func Test_GetUserInfo(t *testing.T) {
	cli, err := Makec()
	if err != nil {
		t.Error(err)
	}
	req := profile.GetUserInfoRequest{
		Uid: "a1",
	}
	result, err := cli.GetUserInfo(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Test_GetFans(t *testing.T) {
	cli, err := Makec()
	if err != nil {
		t.Error(err)
	}
	req := profile.GetFansRequest{
		Uid:   "a1",
		Index: 0,
	}
	result, err := cli.GetFans(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Test_GetFollow(t *testing.T) {
	cli, err := Makec()
	if err != nil {
		t.Error(err)
	}
	req := profile.GetFollowRequest{
		Uid:   "a1",
		Index: 0,
	}
	result, err := cli.GetFollow(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Test_AddFollow(t *testing.T) {
	cli, err := Makec()
	if err != nil {
		t.Error(err)
	}
	req := profile.AddFollowRequest{
		Uid:     "a1",
		DestUid: "c4",
	}
	result, err := cli.AddFollow(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Test_DeleteFollow(t *testing.T) {
	cli, err := Makec()
	if err != nil {
		t.Error(err)
	}
	req := profile.DeleteFollowRequest{
		Uid:     "a1",
		DestUid: "c4",
	}
	result, err := cli.DeleteFollow(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Benchmark_GetUserInfo(b *testing.B) {
	cli, err := Makec()
	if err != nil {
		b.Error(err)
	}
	req := profile.GetUserInfoRequest{"a1"}
	for i := 0; i < b.N; i++ {
		_, err = cli.GetUserInfo(context.Background(), &req)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_GetFans(b *testing.B) {
	cli, err := Makec()
	if err != nil {
		b.Error(err)
	}
	req := profile.GetFansRequest{
		Uid: "a1",
	}
	for i := 0; i < b.N; i++ {
		_, err = cli.GetFans(context.Background(), &req)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_GetFollow(b *testing.B) {
	cli, err := Makec()
	if err != nil {
		b.Error(err)
	}
	req := profile.GetFollowRequest{
		Uid: "a1",
	}
	for i := 0; i < b.N; i++ {
		_, err = cli.GetFollow(context.Background(), &req)
		if err != nil {
			b.Error(err)
		}
	}
}
func Benchmark_GetUserInfo_Parallel(b *testing.B) {
	cli, err := Makec()
	if err != nil {
		b.Error(err)
	}
	req := profile.GetUserInfoRequest{"a1"}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err = cli.GetUserInfo(context.Background(), &req)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func Benchmark_GetFans_Parallel(b *testing.B) {
	cli, err := Makec()
	if err != nil {
		b.Error(err)
	}
	req := profile.GetFansRequest{
		Uid: "a1",
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err = cli.GetFans(context.Background(), &req)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func Benchmark_GetFollow_Parallel(b *testing.B) {
	cli, err := Makec()
	if err != nil {
		b.Error(err)
	}
	req := profile.GetFollowRequest{
		Uid: "a1",
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err = cli.GetFollow(context.Background(), &req)
			if err != nil {
				b.Error(err)
			}
		}
	})
}
