package main

import (
	"context"
	"log"
	"net"

	pb "github.com/williamfedele/leaderboard/protogen"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedLeaderboardServiceServer
	redisClient *redis.Client
}

func (s *server) AddScore(ctx context.Context, req *pb.AddScoreRequest) (*pb.AddScoreResponse, error) {
	err := s.redisClient.ZAdd(ctx, "leaderboard", redis.Z{
		Score:  float64(req.PlayerScore.Score),
		Member: req.PlayerScore.PlayerId,
	}).Err()
	if err != nil {
		return nil, err
	}
	return &pb.AddScoreResponse{Success: true}, nil
}

func (s *server) GetTopScores(ctx context.Context, req *pb.GetTopScoresRequest) (*pb.GetTopScoresResponse, error) {
	entries, err := s.redisClient.ZRevRangeWithScores(ctx, "leaderboard", 0, req.Limit-1).Result()
	if err != nil {
		return nil, err
	}
	var topScores []*pb.PlayerScore
	for _, entry := range entries {
		topScores = append(topScores, &pb.PlayerScore{
			PlayerId: entry.Member.(string),
			Score:    int64(entry.Score),
		})
	}
	return &pb.GetTopScoresResponse{Scores: topScores}, nil
}

func (s *server) GetScoresAroundPlayer(ctx context.Context, req *pb.GetScoresAroundPlayerRequest) (*pb.GetScoresAroundPlayerResponse, error) {

	rank, err := s.redisClient.ZRevRank(ctx, "leaderboard", req.PlayerId).Result()
	if err != nil {
		return nil, err
	}

	entries, err := s.redisClient.ZRevRangeWithScores(ctx, "leaderboard", rank-req.Radius, rank+req.Radius).Result()
	if err != nil {
		return nil, err
	}
	var scores []*pb.PlayerScore
	for _, entry := range entries {
		scores = append(scores, &pb.PlayerScore{
			PlayerId: entry.Member.(string),
			Score:    int64(entry.Score),
		})
	}
	return &pb.GetScoresAroundPlayerResponse{Scores: scores}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	s := grpc.NewServer()
	pb.RegisterLeaderboardServiceServer(s, &server{redisClient: redisClient})

	log.Println("Server started on port :8080")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
