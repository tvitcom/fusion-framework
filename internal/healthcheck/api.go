package healthcheck

import (
	// routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/gofiber/fiber/v2"
)

var appver string

// RegisterHandlers registers the handlers that perform healthchecks.
func RegisterHandlers(r *fiber.App, ver string) {
	appver = ver
	r.Get("/healthcheck", check)
}

// healthcheck responds to a healthcheck request.
func check(c *fiber.Ctx) error {
	return c.SendString("OK " + appver)
}
