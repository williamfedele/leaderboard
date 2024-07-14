package main

import (
	"context"
	"log"
	"net"

	pb "github.com/williamfedele/leaderboard/protos"

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
	entries, err := s.redisClient.ZRevRangeWithScores(ctx, "leaderboard", 0, int64(req.Limit-1)).Result()
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

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	s := grpc.NewServer()
	pb.RegisterLeaderboardServiceServer(s, &server{redisClient: redisClient})

	log.Println("Server started on port :8080")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
