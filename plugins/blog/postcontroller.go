package blog

import (
	"erp6-be-golang/core/events"
	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"
	"erp6-be-golang/response"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Initpost() {
	events.Register("BeforeGet:post", func(data interface{}) error {
		// TODO: implement before get
		return nil
	})
	events.Register("AfterGet:post", func(data interface{}) error {
		// TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:post", func(data interface{}) error {
		// TODO: implement before save
		return nil
	})
	events.Register("AfterSave:post", func(data interface{}) error {
		// TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:post", func(data interface{}) error {
		// TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:post", func(data interface{}) error {
		// TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:post", func(data interface{}) error {
		// TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:post", func(data interface{}) error {
		// TODO: implement after delete
		return nil
	})
}

func BlogPostHandler(c *fiber.Ctx, db *gorm.DB) error {
	companyId := c.FormValue("companyid")
	var ModelPost []models.Post
	if err := db.
		Preload("Useraccess").
		Preload("Postcategory.Category").
		Where("Post.recordstatus = 1 and companyid = ?", companyId).
		Find(&ModelPost).Error; err != nil {
		return err
	}

	baseUrl := os.Getenv("LOCAL_BASE_URL") + os.Getenv("LOCAL_BASE_PATH")
	BlogPost := make([]response.Post, len(ModelPost))
	for i, p := range ModelPost {
		BlogPost[i] = response.Post{
			Postid:       p.Postid,
			Title:        p.Title,
			Description:  p.Description,
			Metatag:      p.Metatag,
			Slug:         p.Slug,
			Postpic:      baseUrl + p.Postpic,
			Postupdate:   p.Postupdate,
			Created:      p.Created,
			Recordstatus: p.Recordstatus,
			Author: response.Author{
				Useraccessid: p.Useraccess.Useraccessid,
				Realname:     p.Useraccess.Realname,
			},
		}
		for _, v := range p.Postcategory {
			if BlogPost[i].Category == "" {
				BlogPost[i].Category = v.Category.Slug
			} else {
				BlogPost[i].Category += "," + v.Category.Slug
			}
		}
	}
	helpers.SuccessResponse(c, "DATA_RETRIEVED", BlogPost)
	return nil
}

func BlogSlugPostHandler(c *fiber.Ctx, db *gorm.DB) error {
	slug := c.FormValue("slug")
	var ModelPost models.Post
	if err := db.
		Preload("Useraccess").
		Preload("Postcategory.Category").
		Where("Post.recordstatus = 1 and slug = ?", slug).
		Find(&ModelPost).Error; err != nil {
		return err
	}

	baseUrl := os.Getenv("LOCAL_BASE_URL") + os.Getenv("LOCAL_BASE_PATH")
	var BlogPost = response.Post{
		Postid:       ModelPost.Postid,
		Title:        ModelPost.Title,
		Description:  ModelPost.Description,
		Metatag:      ModelPost.Metatag,
		Slug:         ModelPost.Slug,
		Postpic:      baseUrl + ModelPost.Postpic,
		Postupdate:   ModelPost.Postupdate,
		Created:      ModelPost.Created,
		Recordstatus: ModelPost.Recordstatus,
		Author: response.Author{
			Useraccessid: ModelPost.Useraccess.Useraccessid,
			Realname:     ModelPost.Useraccess.Realname,
		},
	}
	for _, v := range ModelPost.Postcategory {
		if BlogPost.Category == "" {
			BlogPost.Category = v.Category.Slug
		} else {
			BlogPost.Category += "," + v.Category.Slug
		}
	}
	helpers.SuccessResponse(c, "DATA_RETRIEVED", BlogPost)
	return nil
}
