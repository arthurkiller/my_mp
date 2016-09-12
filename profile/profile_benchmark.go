package profilem

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/arthurkiller/my_mp/grpc/profile"
)

func makec() (profile.ProfileClient, error) {
	conn, err := grpc.Dial("127.0.0.1:23589", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cli := profile.NewProfileClient(conn)
	return cli, nil
}

func Benchmark_GetUserInfo(t *testing.B) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := profile.GetUserInfoRequest{"a1"}
	for i := 0; i < t.N; i++ {
		_, err = cli.GetUserInfo(context.Background(), &req)
		if err != nil {
			t.Error(err)
		}
	}
}

func Benchmark_GetFans(t *testing.B) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := profile.GetFansRequest{
		Uid: "a1",
	}
	for i := 0; i < t.N; i++ {
		_, err = cli.GetFans(context.Background(), &req)
		if err != nil {
			t.Error(err)
		}
	}
}

func Benchmark_GetFollow(t *testing.B) {
	cli, err := makec()
	if err != nil {
		t.Error(err)
	}
	req := profile.GetFollowRequest{
		Uid: "a1",
	}
	for i := 0; i < t.N; i++ {
		_, err = cli.GetFollow(context.Background(), &req)
		if err != nil {
			t.Error(err)
		}
	}
}

//func Benchmark_AddFollow(t *testing.B) {
//	cli, err := makec()
//	if err != nil {
//		t.Error(err)
//	}
//	req := profile.AddFollowRequest{
//		Uid:     "a1",
//		DestUid: "d4",
//	}
//	for i := 0; i < t.N; i++ {
//		_, err = cli.AddFollow(context.Background(), &req)
//		if err != nil {
//			t.Error(err)
//		}
//	}
//}
//
//func Benchmark_DeleteFollow(t *testing.B) {
//	cli, err := makec()
//	if err != nil {
//		t.Error(err)
//	}
//	req := profile.DeleteFollowRequest{
//		Uid:     "a1",
//		DestUid: "d4",
//	}
//	for i := 0; i < t.N; i++ {
//		_, err = cli.DeleteFollow(context.Background(), &req)
//		if err != nil {
//			t.Error(err)
//		}
//	}
//}
