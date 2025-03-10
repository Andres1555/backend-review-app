package comment

import (
	"strconv"

	"github.com/NetKBs/backend-reviewapp/src/schema"
	"github.com/NetKBs/backend-reviewapp/src/social/answer"
	"github.com/NetKBs/backend-reviewapp/src/social/reaction"
)

func GetCommentsReviewCountService(id uint) (commentsCount uint, err error) {

	commentsCount, err = GetCommentsReviewCountRepository(id)
	if err != nil {
		return commentsCount, err
	}

	return commentsCount, nil
}

func GetCommentsByIdReviewService(id uint, limit int, cursor uint) ([]CommentResponseDTO, string, error) {
	reviewComments := []CommentResponseDTO{}

	comments, err := GetCommentsByIdReviewRepository(id, limit, cursor)
	if err != nil {
		return reviewComments, "", err
	}

	for _, comment := range comments {
		reactions, err := reaction.GetReactionsCountService(comment.ID, "comment")
		if err != nil {
			return reviewComments, "", err
		}

		replies, err := answer.GetCountAnswersByCommentIdService(comment.ID)
		if err != nil {
			return reviewComments, "", err
		}

		reviewComments = append(reviewComments, CommentResponseDTO{
			ID:        comment.ID,
			UserId:    comment.UserId,
			ReviewId:  comment.ReviewId,
			Text:      comment.Text,
			Likes:     reactions["likes"],
			Dislikes:  reactions["dislikes"],
			Answers:   replies,
			CreatedAt: comment.CreatedAt.String(),
			UpdatedAt: comment.UpdatedAt.String(),
		})
	}

	nextCursor := ""
	if len(comments) > 0 {
		nextCursor = strconv.FormatUint(uint64(comments[len(comments)-1].ID), 10)
	}

	return reviewComments, nextCursor, nil
}

func GetCommentByIdService(id uint) (commentDTO CommentResponseDTO, err error) {
	comment, err := GetCommentsByIdRepository(id)
	if err != nil {
		return commentDTO, err
	}

	reactions, err := reaction.GetReactionsCountService(id, "comment")
	if err != nil {
		return commentDTO, err
	}

	replies, err := answer.GetCountAnswersByCommentIdService(id)
	if err != nil {
		return commentDTO, err
	}

	commentDTO = CommentResponseDTO{
		ID:        comment.ID,
		UserId:    comment.UserId,
		ReviewId:  comment.ReviewId,
		Text:      comment.Text,
		Likes:     reactions["likes"],
		Dislikes:  reactions["dislikes"],
		Answers:   replies,
		CreatedAt: comment.CreatedAt.String(),
		UpdatedAt: comment.UpdatedAt.String(),
	}

	return commentDTO, nil
}

func CreateCommentService(comment CommentCreateDTO, userId uint) (id uint, err error) {
	commentSchema := schema.Comment{UserId: userId, ReviewId: comment.ReviewId, Text: comment.Text}

	id, err = CreateCommentRepository(commentSchema)
	if err != nil {
		return id, err
	}

	return id, nil
}

func UpdateCommentService(id uint, comment CommentUpdateDTO) (err error) {
	commentSchema := schema.Comment{Text: comment.Text}
	err = UpdateCommentRepository(id, commentSchema)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCommentService(id uint) (err error) {
	err = DeleteCommentRepository(id)
	if err != nil {
		return err
	}
	return nil
}
