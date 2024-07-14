package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/williamfedele/leaderboard/protogen"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	c := pb.NewLeaderboardServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	random := rand.New(rand.NewSource(2))

	// Add 100 random scores
	log.Println("Adding 100 random scores...")
	for i := 0; i < 100; i++ {
		_, err = c.AddScore(ctx, &pb.AddScoreRequest{PlayerScore: &pb.PlayerScore{
			PlayerId: fmt.Sprintf("Player%d", i),
			Score:    int64(random.Intn(100)),
		}})
		if err != nil {
			log.Fatalf("failed to add score: %v", err)
		}
	}

	// Get top 5 scores
	response, err := c.GetTopScores(ctx, &pb.GetTopScoresRequest{Limit: 5})
	if err != nil {
		log.Fatalf("failed to get top scores: %v", err)
	}
	log.Println("Top Scores:")
	for i, score := range response.Scores {
		log.Printf("%d. Player: %s, Score: %d\n", i+1, score.PlayerId, score.Score)
	}

	// Get scores around Player7
	r, err := c.GetScoresAroundPlayer(ctx, &pb.GetScoresAroundPlayerRequest{PlayerId: "Player7", Radius: 2})
	if err != nil {
		log.Fatalf("failed to get scores around player: %v", err)
	}
	log.Println("Scores around Player3:")
	for i, score := range r.Scores {
		log.Printf("%d. Player: %s, Score: %d\n", i+1, score.PlayerId, score.Score)
	}

}
