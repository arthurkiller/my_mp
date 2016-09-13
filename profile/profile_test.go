package profilem

import (
	"context"
	"testing"

	"github.com/arthurkiller/my_mp/grpc/profile"
)

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
