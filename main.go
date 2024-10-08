package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
	url2 "net/url"
)

type Image struct {
	URL        string
	Secure_URL string
	Type       string
	Width      string
	Height     string
	Alt        string
}

type Video struct {
	URL        string
	Secure_URL string
	Type       string
	Width      string
	Height     string
}

type Audio struct {
	URL        string
	Secure_URL string
	Type       string
}

type Result struct {
	Title string
	Type  string
	Image Image
	Video Video
	Audio Audio
	URL   string
}

// TODO: Add caching
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/Mapetr/embed-api")
	})

	app.Post("/embed", func(c *fiber.Ctx) error {
		link := string(c.Body()[:])

		url, urlErr := url2.Parse(link)
		if urlErr != nil {
			return c.SendStatus(400)
		}

		res, err := http.Get(url.String())
		if err != nil {
			log.Errorf(err.Error())
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			c.Status(500)
			return c.SendString(fmt.Sprintf("status code error %d", res.StatusCode))
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			println(fmt.Sprintf("Error: %s", err.Error()))
			return c.SendStatus(500)
		}

		result := new(Result)

		//TODO: Add support for arrays

		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
			property := s.AttrOr("property", "")
			content := s.AttrOr("content", "")

			if property == "" {
				property = s.AttrOr("name", "")
			}

			switch property {
			case "og:title":
				result.Title = content
				break
			case "og:type":
				result.Type = content
				break
			case "og:image":
				result.Image.URL = content
				break
			case "og:image:url":
				result.Image.URL = content
				break
			case "og:image:secure_url":
				result.Image.Secure_URL = content
				break
			case "og:image:type":
				result.Image.Type = content
				break
			case "og:image:width":
				result.Image.Width = content
				break
			case "og:image:height":
				result.Image.Height = content
				break
			case "og:image:alt":
				result.Image.Alt = content
				break
			case "og:video":
				result.Video.URL = content
				break
			case "og:video:url":
				result.Video.URL = content
				break
			case "og:video:secure_url":
				result.Video.Secure_URL = content
				break
			case "og:video:type":
				result.Video.Type = content
				break
			case "og:video:width":
				result.Video.Width = content
				break
			case "og:video:height":
				result.Video.Height = content
				break
			case "og:audio":
				result.Audio.URL = content
				break
			case "og:audio:url":
				result.Audio.URL = content
				break
			case "og:audio:secure_url":
				result.Audio.Secure_URL = content
				break
			case "og:audio:type":
				result.Audio.Type = content
				break
			case "og:url":
				result.URL = content
				break
			}
		})

		return c.JSON(result)
	})

	app.Listen(":3000")
}
