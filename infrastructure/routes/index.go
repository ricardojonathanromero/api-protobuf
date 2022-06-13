package routes

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw "github.com/ricardojonathanromero/api-protobuf/proto/gateway"
	"github.com/ricardojonathanromero/api-protobuf/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
)

type IRouter interface {
	ProxyServer(gRPCAddr, httpAddr string) error
}

type router struct{}

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

func isHeaderAllowed(s string) (string, bool) {
	// check if allowedHeaders contain the header
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		// send uppercase header
		return strings.ToUpper(s), true
	}
	// if not in allowed header, don't send the header
	return s, false
}

func (r *router) ProxyServer(gRPCAddr, httpAddr string) error {
	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux(
		// convert header in response(going from gateway) from metadata received.
		runtime.WithOutgoingHeaderMatcher(isHeaderAllowed),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaller runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			//creating a new HTTTPStatusError with a custom status, and passing error
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			// using default handler to do the rest of heavy lifting of marshaling error and adding headers
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaller, writer, request, &newError)
		}),
	)

	// setting up a dail up for gRPC service by specifying endpoint/target url
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	err := gw.RegisterPostsHandlerFromEndpoint(context.Background(), mux, gRPCAddr, []grpc.DialOption{opts})
	if err != nil {
		return err
	}

	// configure debug
	gin.SetMode(gin.DebugMode)
	if env := utils.GetEnv("ENV", "local"); env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Creating a normal HTTP server
	server := gin.New()
	server.Use(gin.Logger())

	// configure cors
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// configure group
	server.Group("v1/*{grpc_gateway}").Any("", gin.WrapH(mux))

	// configure health
	server.GET("/health_tz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	server.GET("/health_rz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// expose swagger doc
	server.StaticFile("/docs/swagger.json", "./docs/post.swagger.json")

	// start server
	return server.Run(httpAddr)
}

var _ IRouter = (*router)(nil)

func New() IRouter {
	return &router{}
}
