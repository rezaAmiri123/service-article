package handler

import (
	"context"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/rezaAmiri123/service-article/gen/pb"
	"github.com/rezaAmiri123/service-article/internal/model"
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
		Slug:        slug.Make(req.GetTitle()),
		Description: req.GetDescription(),
		Body:        req.GetBody(),
		UserID:      user.Id,
		Tags:        tags,
	}
	if err = article.Validate(); err != nil {
		return nil, err
	}
	if err = h.repo.Create(ctx, &article); err != nil {
		return nil, err
	}
	return article.ProtoArticle(true), nil
}
func (h *articleHandler) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.GetArticle")
	defer span.Finish()

	//user, err := h.getUser(ctx)
	user, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}

	article, err := h.repo.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		return nil, err
	}

	favorited, err := h.repo.IsFavorited(ctx, article, user.Id)
	if err != nil {
		msg := fmt.Sprintf("failded to get user favorited")
		return nil, status.Error(codes.Aborted, msg)
	}
	return article.ProtoArticle(favorited), nil
}

func (h *articleHandler) getUser(ctx context.Context) (*userPb.UserResponse, error) {
	empty := userPb.Empty{}
	md, _ := metadata.FromIncomingContext(ctx)
	newMD := metadata.New(map[string]string{"authorization": md.Get("authorization")[0]})
	ctx = metadata.NewOutgoingContext(ctx, newMD)

	return h.userClient.GetUser(ctx, &empty)
}

func (h *articleHandler) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.UpdateArticle")
	defer span.Finish()

	user, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}

	article, err := h.repo.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		return nil, err
	}

	if article.UserID != user.Id {
		msg := fmt.Sprintf("wrong user")
		return nil, status.Error(codes.PermissionDenied, msg)
	}

	article.Overwrite(
		req.GetTitle(),
		req.GetDescription(),
		req.GetBody(),
	)

	if err = article.Validate(); err != nil {
		err = fmt.Errorf("validation error: %w", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = h.repo.Update(ctx, article); err != nil {
		msg := fmt.Sprintf("database error: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return article.ProtoArticle(true), nil
}

func (h *articleHandler) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.DeleteArticle")
	defer span.Finish()

	user, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}

	article, err := h.repo.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		return nil, err
	}

	if article.UserID != user.Id {
		msg := fmt.Sprintf("wrong user")
		return nil, status.Error(codes.PermissionDenied, msg)
	}

	if err = h.repo.Delete(ctx, article); err != nil {
		msg := fmt.Sprintf("database error: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	return &pb.Empty{}, nil
}

func (h *articleHandler) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.CreateComment")
	defer span.Finish()

	user, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}

	article, err := h.repo.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		return nil, err
	}
	comment := model.Comment{
		Body:      req.GetBody(),
		ArticleID: article.ID,
		UserID:    user.Id,
	}
	if err := comment.Validate(); err != nil {
		return nil, err
	}
	if err := h.repo.CreateComment(ctx, &comment); err != nil {
		return nil, err
	}
	return comment.ProtoComment(), nil
}

func (h *articleHandler) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.DeleteComment")
	defer span.Finish()

	user, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}

	article, err := h.repo.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		return nil, err
	}

	comment, err := h.repo.GetCommentByID(ctx, req.GetId())
	if err != nil {
		msg := fmt.Sprintf("database error: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	if article.UserID != user.Id || comment.ArticleID != article.ID {
		msg := fmt.Sprintf("wrong user")
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	err = h.repo.DeleteComment(ctx, comment)
	if err != nil {
		msg := fmt.Sprintf("database error: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return &pb.Empty{}, nil
}

func (h *articleHandler) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.CommentsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.GetComments")
	defer span.Finish()

	//user, err := h.getUser(ctx)
	//if err != nil {
	//	return nil, err
	//}

	article, err := h.repo.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		msg := fmt.Sprintf("database error: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	comments, err := h.repo.GetComments(ctx, article)
	if err != nil {
		msg := fmt.Sprintf("database error: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	pcs := make([]*pb.Comment, 0, len(comments))
	for _, c := range comments {
		pcs = append(pcs, c.ProtoComment())
	}
	return &pb.CommentsResponse{Comments: pcs}, nil
}

func (h *articleHandler) FavoriteArticle(ctx context.Context, req *pb.FavoriteArticleRequest) (*pb.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.FavoriteArticle")
	defer span.Finish()

	user, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}

	article, err := h.repo.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		return nil, err
	}
	err = h.repo.AddFavorite(ctx, article, user.Id)
	if err != nil {
		msg := fmt.Sprintf("failed to add favorite: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return article.ProtoArticle(true), nil
}

func (h *articleHandler) UnfavoriteArticle(ctx context.Context, req *pb.FavoriteArticleRequest) (*pb.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "articleHandler.UnfavoriteArticle")
	defer span.Finish()
	user, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}

	article, err := h.repo.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		return nil, err
	}
	err = h.repo.DeleteFavorite(ctx, article, user.Id)
	if err != nil {
		msg := fmt.Sprintf("failed to delete favorite: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return article.ProtoArticle(false), nil
}
