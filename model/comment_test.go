package model

import (
	"testing"
	"time"

	"github.com/nathanbizkit/article-management/message"
	"github.com/stretchr/testify/assert"
)

func TestUnit_CommentModel(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping unit tests.")
	}

	t.Run("Validate", func(t *testing.T) {
		tests := []struct {
			title    string
			c        *Comment
			hasError bool
		}{
			{
				"validate comment: success",
				&Comment{
					Body:      "A text body.",
					ArticleID: 1,
					UserID:    1,
				},
				false,
			},
			{
				"validate comment: no body",
				&Comment{
					Body:      "",
					ArticleID: 1,
					UserID:    1,
				},
				true,
			},
			{
				"validate comment: no article id",
				&Comment{
					Body:      "A text body.",
					ArticleID: 0,
					UserID:    1,
				},
				true,
			},
			{
				"validate comment: no user id",
				&Comment{
					Body:      "A text body.",
					ArticleID: 1,
					UserID:    0,
				},
				true,
			},
		}

		for _, tt := range tests {
			err := tt.c.Validate()

			if tt.hasError {
				assert.Error(t, err, tt.title)
			} else {
				assert.NoError(t, err, tt.title)
			}
		}
	})

	t.Run("ResponseComment", func(t *testing.T) {
		createdAt := time.Now()
		updatedAt := time.Now().Add(10 * time.Hour)
		updatedAtString := updatedAt.Format(time.RFC3339Nano)

		following := false
		expected := message.CommentResponse{
			ID:   1,
			Body: "This is my comment.",
			Author: message.ProfileResponse{
				Username:  "foo_user",
				Name:      "FooUser",
				Bio:       "This is my bio.",
				Image:     "https://imgur.com/image.jpeg",
				Following: following,
			},
			CreatedAt: createdAt.Format(time.RFC3339Nano),
			UpdatedAt: &updatedAtString,
		}

		c := Comment{
			ID:     1,
			Body:   "This is my comment.",
			UserID: 1,
			Author: User{
				ID:        1,
				Username:  "foo_user",
				Email:     "foo@example.com",
				Password:  "encrypted_password",
				Name:      "FooUser",
				Bio:       "This is my bio.",
				Image:     "https://imgur.com/image.jpeg",
				CreatedAt: time.Now(),
				UpdatedAt: nil,
			},
			ArticleID: 1,
			CreatedAt: createdAt,
			UpdatedAt: &updatedAt,
		}

		actual := c.ResponseComment(following)
		assert.Equal(t, expected, actual)
	})
}
