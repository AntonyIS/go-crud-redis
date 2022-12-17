package main

import (
	"fmt"
	"net/http"

	"demo-crud/cache"

	"github.com/gin-gonic/gin"
)

var (
	redisCache = cache.NewRedisCache("localhost:6379", 0, 1)
)

type Todo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func main() {
	r := gin.Default()

	r.POST("/movies", func(ctx *gin.Context) {
		var movie cache.Movie
		if err := ctx.ShouldBind(&movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := redisCache.CreateMovie(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res,
		})

	})
	r.GET("/movies", func(ctx *gin.Context) {
		movies, err := redisCache.GetMovies()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movies": movies,
		})
	})
	r.GET("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		movie, err := redisCache.GetMovie(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "movie not found",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": movie,
		})
	})
	r.PUT("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := redisCache.GetMovie(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var movie cache.Movie

		if err := ctx.ShouldBind(&movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res.Title = movie.Title
		res.Description = movie.Description
		res, err = redisCache.UpdateMovie(res)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res,
		})
	})
	r.DELETE("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := redisCache.DeleteMovie(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "movie deleted successfuly",
		})
	})
	fmt.Println(r.Run(":5000"))

}
