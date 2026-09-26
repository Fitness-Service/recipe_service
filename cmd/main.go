package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Fitness-Service/recipe-service/internal/models"
	"github.com/Fitness-Service/recipe-service/internal/repository"
	pb "github.com/Fitness-Service/shared/proto/recipe"
	"google.golang.org/grpc"
)

type RecipeServer struct {
	repo *repository.RecipeRepo
	pb.UnimplementedRecipeServiceServer
}

func (s *RecipeServer) GetAllRecipes(ctx context.Context, req *pb.Empty) (*pb.RecipeList, error) {
	recipes, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return &pb.RecipeList{Recipes: toProtoRecipes(recipes)}, nil
}

func toProtoRecipes(recipes []models.Recipe) []*pb.Recipe {
	result := make([]*pb.Recipe, 0, len(recipes))
	for _, r := range recipes {
		result = append(result, &pb.Recipe{
			Id:          r.ID,
			Name:        r.Name,
			Ingredients: r.Ingredients,
			Calories:    int32(r.Calories),
			Protein:     r.Protein,
			Fats:        r.Fats,
			Carbs:       r.Carbs,
			Goal:        r.Goal,
			Tags:        r.Tags,
			PrepTime:    int32(r.PrepTime),
			CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		})
	}
	return result
}

func main() {
	// DSN для подключения к PostgreSQL
	connString := "postgres://postgres:223244@localhost:5432/fitness?sslmode=disable"

	repo := repository.NewRecipeRepo(connString)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRecipeServiceServer(grpcServer, &RecipeServer{repo: repo})

	log.Println("recipe-service gRPC started on :8082")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
