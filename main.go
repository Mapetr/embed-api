package main

import (
	json2 "encoding/json"
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

type Result struct {
	Title string
	Type  string
	Image Image
	URL   string
}

func main() {
	app := fiber.New()

	app.Get("/embed", func(c *fiber.Ctx) error {
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

		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
			property := s.AttrOr("property", "")
			content := s.AttrOr("content", "")

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
			case "og:url":
				result.URL = content
				break
			}
		})

		json, jsonErr := json2.Marshal(result)
		if jsonErr != nil {
			println(jsonErr.Error())
			return c.SendStatus(500)
		}

		return c.Send(json)
	})

	app.Listen(":3000")
}
