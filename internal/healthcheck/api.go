package healthcheck

import (
	// routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/gofiber/fiber/v2"
)

// RegisterHandlers registers the handlers that perform healthchecks.
func RegisterHandlers(r *fiber.App, version string) {
	// r.To("GET,HEAD", "/healthcheck", healthcheck(version))
	r.Get("/healthcheck", healthcheck)
}

// healthcheck responds to a healthcheck request.
func healthcheck(version string) {
	// return func(c *fiber.Context) error {
	// 	return c.Write("OK " + version)
	// }
	func(c *fiber.Ctx) error {
	  return c.SendString("OK " + version)
	  // return c.Write("OK " + version)

	  return nil
	}
}
