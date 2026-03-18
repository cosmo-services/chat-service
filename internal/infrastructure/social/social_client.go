package social_infrastructure

import (
	"context"
	"fmt"
	"main/internal/config"
	chat_domain "main/internal/domain/chat"
	social_domain "main/internal/domain/social"
	"main/pkg"
	"time"

	pb "github.com/cosmo-services/grpc-contracts/gen/api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcSocialClient struct {
	client  pb.SocialServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewGrpcSocialClient(grpcClient *pkg.GrpcClient, env config.Env) (social_domain.SocialClient, error) {
	conn, err := grpcClient.Connect(env.SocialServiceGrpcAddress)
	if err != nil {
		return nil, err
	}

	return &GrpcSocialClient{
		client:  pb.NewSocialServiceClient(conn),
		conn:    conn,
		timeout: 15 * time.Second,
	}, nil
}

func (c *GrpcSocialClient) GetProfileByUserId(userId string) (*social_domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)

	defer cancel()

	req := &pb.GetUserProfileByIdRequest{
		UserId: userId,
	}

	resp, err := c.client.GetUserProfileByUserId(ctx, req)
	if err != nil {
		return nil, c.mapGRPCError(err)
	}

	return c.mapToDomainUserProfile(resp), nil
}

func (c *GrpcSocialClient) GetProfileByUsername(username string) (*social_domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)

	defer cancel()

	req := &pb.GetUserProfileByUsernameRequest{
		Username: username,
	}

	resp, err := c.client.GetUserProfileByUsername(ctx, req)
	if err != nil {
		return nil, c.mapGRPCError(err)
	}

	return c.mapToDomainUserProfile(resp), nil
}

func (c *GrpcSocialClient) mapGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("unexpected error: %w", err)
	}

	switch st.Code() {
	case codes.NotFound:
		return chat_domain.ErrUserNotFound
	case codes.DeadlineExceeded:
		return fmt.Errorf("request timeout: %w", err)
	case codes.Unavailable:
		return fmt.Errorf("service unavailable: %w", err)
	default:
		return fmt.Errorf("gRPC error (code=%s): %w", st.Code(), err)
	}
}

func (c *GrpcSocialClient) mapToDomainUserProfile(response *pb.GetUserProfileResponse) *social_domain.UserProfile {
	return &social_domain.UserProfile{
		UserID:      response.Profile.Id,
		Username:    response.Profile.Username,
		DisplayName: response.Profile.DisplayName,
		AvatarUrl:   response.Profile.AvatarUrl,
	}
}
