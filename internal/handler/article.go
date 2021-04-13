package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/rezaAmiri123/service-article/internal/model"
	"github.com/rezaAmiri123/service-article/pkg/utils"
	"google.golang.org/grpc/metadata"

	pb "github.com/rezaAmiri123/service-article/gen/pb"
	"github.com/rezaAmiri123/service-article/internal/repository"
	"github.com/rezaAmiri123/service-article/pkg/logger"
	userPb "github.com/rezaAmiri123/service-user/gen/pb"
)

type articleHandler struct {
	repo       repository.ArticleRepository
	logger     logger.Logger
	userClient userPb.UsersClient
}

func NewArticleHandler(repo repository.ArticleRepository, logger logger.Logger, userClient userPb.UsersClient) *articleHandler {
	return &articleHandler{repo: repo, logger: logger, userClient: userClient}
}

func (h *articleHandler) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.CreateArticle")
	defer span.Finish()
	user, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}
	tags := make([]model.Tag, 0, len(req.GetTagList()))
	for _, t := range req.GetTagList() {
		tags = append(tags, model.Tag{Name: t})
	}
	article := model.Article{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Body:        req.GetBody(),
		UserID:      utils.UintToString(user.Id),
		Tags:        tags,
	}
	if err = article.Validate();err != nil{
		return nil, err
	}
	if err = h.repo.Create(ctx, &article);err != nil{
		return nil, err
	}
	return article.ProtoArticle(), nil
}

func (h *articleHandler) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.Comment, error) {
	return nil, nil
}

func (h *articleHandler) getUser(ctx context.Context) (*userPb.UserResponse, error) {
	empty := userPb.Empty{}
	md, _ := metadata.FromIncomingContext(ctx)
	newMD := metadata.New(map[string]string{"authorization": md.Get("authorization")[0]})
	ctx = metadata.NewOutgoingContext(ctx, newMD)

	return h.userClient.GetUser(ctx, &empty)
}