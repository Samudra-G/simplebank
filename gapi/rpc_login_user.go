package gapi

import (
	"context"
	"database/sql"

	db "github.com/Samudra-G/simplebank/db/sqlc"
	"github.com/Samudra-G/simplebank/pb"
	"github.com/Samudra-G/simplebank/util"
	"github.com/Samudra-G/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password: %v", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}
	
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %v", err)
	}

	mtdt := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:            refreshPayload.ID,
		Username:      user.Username,
		RefreshToken:  refreshToken,
		UserAgent:     mtdt.UserAgent,
		ClientIp:      mtdt.ClientIP, 
		IsBlocked:     false,
		ExpiredAt:     refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	rsp := &pb.LoginUserResponse{
		User: convertUser(user),
		SessionId: session.ID.String(),
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}