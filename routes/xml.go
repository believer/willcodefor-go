package routes

import (
	"log"
	"time"

	"github.com/believer/willcodefor-go/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
	"github.com/jmoiron/sqlx"
)

type PostWithParsedDate struct {
	Post
	UpdatedAtParsed string
}

func FeedHandler(c *fiber.Ctx, db *sqlx.DB) error {
	posts := []PostWithParsedDate{}
	engineXML := mustache.New("./xmls", ".xml")

	if err := engineXML.Load(); err != nil {
		log.Fatal(err)
	}

	q := `
    SELECT title, slug, body, updated_at
    FROM post
    WHERE published = true
    ORDER BY created_at DESC
  `

	err := db.Select(&posts, q)

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for _, post := range posts {
		body := utils.MarkdownToHTML([]byte(post.Body))
		post.Body = body.String()
		post.UpdatedAtParsed = post.UpdatedAt.Format(time.RFC3339)
	}

	c.Type("xml")

	return engineXML.Render(c, "feed", fiber.Map{
		"Metadata": fiber.Map{
			"Title":       "willcodefor.beer",
			"URL":         "https://willcodefor.beer/",
			"Description": "Things I learn while browsing the web",
			"Author": fiber.Map{
				"Name":  "Rickard Natt och Dag",
				"Email": "rickard@willcodefor.dev",
			},
		},
		"Posts":            posts,
		"LatestPostUpdate": posts[0].UpdatedAt.Format(time.RFC3339),
	})
}

func SitemapHandler(c *fiber.Ctx, db *sqlx.DB) error {
	posts := []PostWithParsedDate{}
	engineXML := mustache.New("./xmls", ".xml")

	if err := engineXML.Load(); err != nil {
		log.Fatal(err)
	}

	err := db.Select(&posts, "SELECT slug, updated_at FROM post WHERE published = true ORDER BY created_at DESC")

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for _, post := range posts {
		post.UpdatedAtParsed = post.UpdatedAt.Format(time.RFC3339)
	}

	c.Type("xml")

	return engineXML.Render(c, "sitemap", fiber.Map{
		"URL":   "https://willcodefor.beer/",
		"Posts": posts,
	})

}
