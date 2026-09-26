package cmd

import (
	"context"
	"net"

	"github.com/Fitness-Service/recipe-service/internal/repository"
	pb "github.com/Fitness-Service/shared/proto/recipe"
	"google.golang.org/grpc"
)

type RecipeServer struct {
	repo *repository.RecipeRepo
	pb.UnimplementedRecipeServiceServer
}

func (s *RecipeServer) GetAllRecipes(ctx context.Context, req *pb.Empty) (*pb.RecipeList, error) {
	recipes, _ := s.repo.GetAll()
	return &pb.RecipeList{Recipes: toProtoRecipes(recipes)}, nil
}

func main() {
	lis, _ := net.Listen("tcp", ":8082")
	grpcServer := grpc.NewServer()
	pb.RegisterRecipeServiceServer(grpcServer, &RecipeServer{repo: repo})
	grpcServer.Serve(lis)
}
