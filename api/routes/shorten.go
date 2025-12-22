package routes

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type request struct {
	URL          string        `json:"url"`
	CustomeShort string        `json:"short"`
	Expiry       time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomeShort    string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fibre.StatusBadRequest).JSON(fibre.Map{"error": "cannot parse JSON"})
	}

	//implement rate limiting
	//check if the input is an actual url

	if !govalidator.IsURL(body.URL) {
		return c.Status(fibre.StatusBadRequest).JSON(fibre.Map{"error": "Invalid URL"})
	}

	//check for the domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fibre.StatusServiceUnavailable).JSON(fibre.Map{"error": "You can't hack the system!!"})
	}

	//enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)
}
