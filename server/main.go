package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	pb "proto"

	"github.com/google/uuid"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// Port for gRPC server to listen to
	PORT    = ":50051"
	crtFile = "private_key.pem"
	keyFile = "public_key.pem"
)

type TodoServer struct {
	pb.UnimplementedTodoServiceServer
}

func (s *TodoServer) CreateTodo(ctx context.Context, in *pb.NewTodo) (*pb.Todo, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	//header := metadata.Pairs("header-key", "val")

	if in.Name == "" || !ok {
		log.Printf("Name cannot be null field ....")

		errorStatus := status.New(codes.InvalidArgument, "Invalid Name info received .... 2")
		ds, err := errorStatus.WithDetails(
			&epb.BadRequest_FieldViolation{
				Field: "Name",
				Description: fmt.Sprintf(
					"Name cannot be empty string %s : %s",
					in.Name, in.Description),
			},
		)
		if err != nil {
			return nil, errorStatus.Err()
		}
		return nil, ds.Err()
	}

	for key, value := range md {
		fmt.Printf("  Key: %s", key)
		for counter, strVal := range value {
			fmt.Printf("  count: %d   value: %s", counter, strVal)
		}
	}

	//log.Printf("Going to delay for 5 seconds")
	//time.Sleep(5 * time.Second)
	//log.Printf("Delayed function completed after 5 seconds")

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("Deadline reached on the client side and requet is no more valid so abort sending results ......")
		return nil, nil
	} else {
		log.Printf("Received: %v", in.GetName())
		todo := &pb.Todo{
			Name:        in.GetName(),
			Description: in.GetDescription(),
			Done:        false,
			Id:          uuid.New().String(),
		}

		return todo, nil
	}

}

// Server :: Unary Interceptor 1
func todoUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Pre-processing logic
	// Gets info about the current RPC call by examining the args passed in
	log.Println("======= [Server Interceptor] ", info.FullMethod)
	log.Printf(" Pre Proc Message : %s", req)

	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)

	// Post processing logic
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}

func main() {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	log.Fatalf(string(cert.OCSPStaple))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	/*opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}*/

	lis, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}

	//s := grpc.NewServer()
	s := grpc.NewServer(
		grpc.UnaryInterceptor(todoUnaryServerInterceptor),
	)

	pb.RegisterTodoServiceServer(s, &TodoServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}

/*

You can have only one unary interceptor that will receiev or send request BUT if you need
a chain of interceptors then you can do this trick


func chainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // Execute interceptors in order
        for _, interceptor := range interceptors {
            req, err := interceptor(ctx, req, info, handler)
            if err != nil {
                return nil, err
            }
        }

        // Call the actual RPC handler
        return handler(ctx, req)
    }
}

// ...

server := grpc.NewServer(
    grpc.UnaryInterceptor(chainUnaryInterceptors(
        myFirstUnaryInterceptor,
        mySecondUnaryInterceptor,
    )),
)

here you register only interfactor but in the constructore you pass references to multiple interceptor
functions. So when a call is intercepted then in the interceptor function you loop through each interctor
and execute them one by one


*/
