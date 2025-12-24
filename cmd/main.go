package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JuD4Mo/go_api_web_course/internal/course"
	"github.com/JuD4Mo/go_api_web_course/pkg/bootstrap"
	"github.com/JuD4Mo/go_api_web_course/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {

	//Cargamos las variables de entorno que están en el archivo .env por medio del package godotenv
	_ = godotenv.Load()

	//Instanciamos un logger propio
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	ctx := context.Background()

	courseRepo := course.NewRepo(db, l)
	courseService := course.NewService(l, courseRepo)
	h := handler.NewCourseHttpServer(ctx, course.MakeEndpoints(courseService, course.Config{LimitPage: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	//Se crea una instancia de un servidor
	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST, PATCH, OPTIONS, DELETE, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
