package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestServer(t *testing.T) {
	gs := grpc.NewServer()
	us := &Server{}
	RegisterUserServiceServer(gs, us)

	l, err := net.Listen("tcp", ":8090")
	require.NoError(t, err)
	err = gs.Serve(l)
	t.Log(err)
}

func TestClient(t *testing.T) {
	//cc, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	cc, err := grpc.Dial("localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)
	client := NewUserServiceClient(cc)
	resp, err := client.GetByID(context.Background(), &GetByIDRequest{Id: 123})
	require.NoError(t, err)
	t.Log(resp.User)
}

func TestOneOf(t *testing.T) {
	u := &User{}
	email, ok := u.Contacts.(*User_Email)
	if ok {
		t.Log("我传入的是 email", email)
		return
	}
}
