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

type Result struct {
	Title string
	Type  string
	Image string
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

			switch property {
			case "og:title":
				result.Title = s.AttrOr("content", "")
				break
			case "og:type":
				result.Type = s.AttrOr("content", "")
				break
			case "og:image":
				result.Image = s.AttrOr("content", "")
				break
			case "og:url":
				result.URL = s.AttrOr("content", "")
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
