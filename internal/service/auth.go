package service

import (
	"context"

	"github.com/XVNDEX/blackrocksouth_test/internal/entity"
	"github.com/XVNDEX/blackrocksouth_test/pb"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer is the server that provides auth services
type AuthServer struct {
	Credentials *entity.Credentials
	pb.UnimplementedAuthServiceServer
}

// NewAuthServer returns a new AuthServer
func NewAuthServer(Credentials *entity.Credentials) *AuthServer {
	return &AuthServer{
		Credentials: Credentials,
	}
}

// SignInUser checks if user provided valid login and password
func (c *AuthServer) SignInUser(ctx context.Context,
	req *pb.SignInUserInput) (*pb.SignInUserResponse, error) {
	login := req.GetLogin()
	pwd := req.GetPwd()

	if len(login) == 0 || len(pwd) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "login or password empty")
	}

	if !c.Credentials.CheckCredentials(login, pwd) {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect login or password")
	}

	token := FetchToken("some-secret-token")
	return &pb.SignInUserResponse{Token: token.AccessToken}, nil
}

// FetchToken simulates a token lookup and omits the details of proper token acquisition
func FetchToken(s string) *oauth2.Token {
	return &oauth2.Token{
		AccessToken: s,
	}
}
