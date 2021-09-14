package album

import (
	"github.com/julienschmidt/httprouter"
	// "github.com/tvitcom/fusion-framework/internal/errors"
	"github.com/tvitcom/fusion-framework/pkg/log"
	// "github.com/tvitcom/fusion-framework/pkg/pagination"
    "github.com/justinas/alice"
	"net/http"
	"context"
	"fmt"
)

// m1 is middleware 1
func m1(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //do something with m1
        fmt.Println("m1 start here")
        w.Header().Set("Server1", "Ok1")//
        next.ServeHTTP(w, r)
        fmt.Println("m1 end here")
    })
}

// m2 is middleware 2
func m2(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //do something with m2
        fmt.Println("m2 start here")
        w.Header().Set("Server2", "Ok2")//
        next.ServeHTTP(w, r)
        fmt.Println("m2 end here")
    })
}

func m3(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //do something with m2
        fmt.Println("m3 start here")
        w.Header().Set("Server3", "Ok3")//
        next.ServeHTTP(w, r)
        fmt.Println("m3 end here")
    })
}


type resource struct {
	agregator Agregator
	logger  log.Logger
}

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(
	r *httprouter.Router, 
	agregator Agregator, 
	// authHandler routing.Handler, 
	logger log.Logger,
) {
	res := resource{agregator, logger}
    // For exampple: gzip, ratelim, jwtauth, csp
    // indexMwConvey := alice.New(m1, m2, m3)
    // getMwConvey := alice.New(m1, m3)
	
	r.GET("/page/:name", middlewared(alice.New(m1, m2, m3).ThenFunc(res.index)))
	// r.Get("/albums", res.query)
	
	// the following endpoints require a valid JWT
	//- r.Use(authHandler)
	r.GET("/albums/:id", middlewared(alice.New(m1, m3).ThenFunc(res.get)))
	// r.Post("/albums", res.create)
	// r.Put("/albums/:id", res.update)
	// r.Delete("/albums/:id", res.delete)
}

// wrapper wraps http.Handler and returns httprouter.Handle
func middlewared(next http.Handler) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        //pass httprouter.Params to request context
        ctx := context.WithValue(r.Context(), "params", ps)
        //call next middleware with new context
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func (res resource) index(w http.ResponseWriter, r *http.Request) {
    // get httprouter.Params from request context
    ps := r.Context().Value("params").(httprouter.Params)
    fmt.Fprintf(w, "Pagename is:, %s", ps.ByName("name"))
    return
}

func (res resource) get(w http.ResponseWriter, r *http.Request) {
	ps := r.Context().Value("params").(httprouter.Params)
	album, _ := res.agregator.Get(r.Context(), ps.ByName("id"))
	fmt.Fprintf(w, "%v\n", album)
	return
}

// func (res resource) query(c *routing.Context) error {
// 	ctx := c.Request.Context()
// 	count, err := res.agregator.Count(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	pages := pagination.NewFromRequest(c.Request, count)
// 	albums, err := res.agregator.Query(ctx, pages.Offset(), pages.Limit())
// 	if err != nil {
// 		return err
// 	}
// 	pages.Items = albums
// 	return c.Write(pages)
// }

// func (res resource) create(c *routing.Context) error {
// 	var input CreateAlbumRequest
// 	if err := c.Read(&input); err != nil {
// 		res.logger.With(c.Request.Context()).Info(err)
// 		return errors.BadRequest("")
// 	}
// 	album, err := res.agregator.Create(c.Request.Context(), input)
// 	if err != nil {
// 		return err
// 	}

// 	return c.WriteWithStatus(album, http.StatusCreated)
// }

// func (res resource) update(c *routing.Context) error {
// 	var input UpdateAlbumRequest
// 	if err := c.Read(&input); err != nil {
// 		res.logger.With(c.Request.Context()).Info(err)
// 		return errors.BadRequest("")
// 	}

// 	album, err := res.agregator.Update(c.Request.Context(), c.Param("id"), input)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(album)
// }

// func (res resource) delete(c *routing.Context) error {
// 	album, err := res.agregator.Delete(c.Request.Context(), c.Param("id"))
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(album)
// }
/*
func writeHtmlResult(w http.ResponseWriter, res string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, res)
}

func writeJSONResult(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}*/
