package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"

	"github.com/dliluashvili/cowatchit/db"
	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/handlers"
	"github.com/dliluashvili/cowatchit/internal/interceptors"
	"github.com/dliluashvili/cowatchit/internal/middlewares"
	"github.com/dliluashvili/cowatchit/internal/repositories"
	"github.com/dliluashvili/cowatchit/internal/services"
	"github.com/dliluashvili/cowatchit/internal/shared/validators"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(5, 20)
		visitors[ip] = limiter
	}

	return limiter
}

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // Ideally, parse real IP behind proxies
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db := db.New(".env.dev")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))

	if err != nil {
		fmt.Println("Incorrect redis port")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
	})

	var validate = validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("json")
	})

	sessionService := services.NewSessionService(redisClient)
	roomRedisService := services.NewRoomRedisService(redisClient)

	userRepository := repositories.NewUserRepository(db)
	roomRepository := repositories.NewRoomRepository(db)

	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(sessionService, userService)
	roomService := services.NewRoomService(roomRepository, userService, roomRedisService)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(roomService)

	roomMessageRepository := repositories.NewRoomMessageRepository(db)
	roomMessageService := services.NewRoomMessageService(roomMessageRepository)

	webSocketManagerService := services.NewWebSocketManagerService()

	webSocketHandler := handlers.NewWebSocketHandler(validate, webSocketManagerService, sessionService, roomService, roomMessageService)

	validate.RegisterValidation("unique", validators.Unique(userRepository))
	validate.RegisterValidation("gender", validators.Gender)
	validate.RegisterValidation("date", validators.Date)
	validate.RegisterValidation("dob", validators.Dob)
	validate.RegisterValidation("username", validators.Username)
	validate.RegisterValidation("roomtitle", validators.RoomTitle)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.Get("/", handlers.HandleLanding)
	r.With(interceptors.ValidateBody[dtos.SignUpDto](validate)).Post("/auth/sign-up", authHandler.SignUp)
	r.With(interceptors.ValidateBody[dtos.SignInDto](validate)).Post("/auth/sign-in", authHandler.SignIn)
	r.With(middlewares.AuthSession(sessionService)).Get("/user/me", userHandler.Me)
	r.With(middlewares.AuthSession(sessionService)).Get("/rooms", roomHandler.HandleRoomsPage)
	r.With(middlewares.AuthSession(sessionService)).Get("/rooms/{id}", roomHandler.HandleRoomPage)
	r.With(middlewares.AuthSession(sessionService)).Get("/create-room", roomHandler.HandleCreateRoomPage)
	r.With(middlewares.AuthSession(sessionService)).With(interceptors.ValidateBody[dtos.CreateRoomDto](validate)).Post("/create-room", roomHandler.Create)
	r.Get("/ws", webSocketHandler.Handle)

	if err := http.ListenAndServe(":8080", limit(r)); err != nil {
		fmt.Printf("failed: %v\n", err)
	} else {
		fmt.Println("all good !")
	}
}
