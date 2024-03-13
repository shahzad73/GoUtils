package main

import (
	"context"
	"log"
	"time"

	pb "proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
)

const (
	ADDRESS = "localhost:50051"
)

type TodoTask struct {
	Name        string
	Description string
	Done        bool
}

func orderUnaryClientInterceptor(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Preprocessor phase
	log.Println(".......Method : " + method)
	// Invoking the remote method
	err := invoker(ctx, method, req, reply, cc, opts...)
	// Postprocessor phase
	log.Println(reply)
	return err
}

func main() {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(orderUnaryClientInterceptor))

	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	defer conn.Close()

	c := pb.NewTodoServiceClient(conn)

	// Add a deadline
	clientDeadline := time.Now().Add(time.Duration(2 * time.Second))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	todos := []TodoTask{
		{Name: "Code review", Description: "Review new feature code", Done: false},
		{Name: "Make YouTube Video", Description: "Start Go for beginners series", Done: false},
		{Name: "Go to the gym", Description: "Leg day", Done: false},
		{Name: "", Description: "This is empty string test", Done: false},
		{Name: "Buy groceries", Description: "Buy tomatoes, onions, mangos", Done: false},
		{Name: "Meet with mentor", Description: "Discuss blockers in my project", Done: false},
	}

	ctxA := metadata.AppendToOutgoingContext(ctx, "k1", "v1", "k2", "v2", "k3", "v3")

	for _, todo := range todos {
		res, err := c.CreateTodo(ctxA, &pb.NewTodo{Name: todo.Name, Description: todo.Description, Done: todo.Done})

		if err != nil {
			got := status.Code(err)
			// log.Fatalf("could not create user: %v", err)
			log.Printf("Error Occured -> addOrder : , %v:", got)

			errorCode := status.Code(err)
			if errorCode == codes.InvalidArgument {
				log.Printf("Invalid Argument Error : %s", err)
				errorStatus := status.Convert(err)
				for _, d := range errorStatus.Details() {
					switch info := d.(type) {
					case *epb.BadRequest_FieldViolation:
						log.Printf("Request Field Invalid: %s", info)
					default:
						log.Printf("Unexpected error type: %s", info)
					}
				}
			} else {
				log.Printf("Unhandled error : %s ", err)
			}

		}

		log.Printf(`
           ID : %s
           Name : %s
           Description : %s
           Done : %v,
       `, res.GetId(), res.GetName(), res.GetDescription(), res.GetDone())
	}

}
