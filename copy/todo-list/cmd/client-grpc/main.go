package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	v12 "github.com/sko00o/go-lab/copy/todo-list/pkg/api/v1"
)

const (
	apiVersion = "v1"
)

func main() {
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v12.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// create
	res1 := create(ctx, c)

	// read
	res2 := read(ctx, c, res1.Id)

	// update
	update(ctx, c, res2.ToDo)

	// read all
	readAll(ctx, c)

	// delete
	delete(ctx, c, res1.Id)
}

func create(ctx context.Context, c v12.ToDoServiceClient) *v12.CreateResponse {
	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	req1 := v12.CreateRequest{
		Api: apiVersion,
		ToDo: &v12.ToDo{
			Title:       "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
			Reminder:    reminder,
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("create failed: %v", err)
	}
	log.Printf("create result: <%+v>\n\n", res1)

	return res1
}

func read(ctx context.Context, c v12.ToDoServiceClient, id int64) *v12.ReadResponse {
	req2 := v12.ReadRequest{
		Api: apiVersion,
		Id:  id,
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("read failed: %v", err)
	}
	log.Printf("read result: <%+v>\n\n", res2)

	return res2
}

func update(ctx context.Context, c v12.ToDoServiceClient, todo *v12.ToDo) *v12.UpdateResponse {
	req3 := v12.UpdateRequest{
		Api: apiVersion,
		ToDo: &v12.ToDo{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description + " + updated",
			Reminder:    todo.Reminder,
		},
	}
	res3, err := c.Update(ctx, &req3)
	if err != nil {
		log.Fatalf("update failed: %v", err)
	}
	log.Printf("update result: <%+v>\n\n", res3)

	return res3
}

func readAll(ctx context.Context, c v12.ToDoServiceClient) *v12.ReadAllResponse {
	req4 := v12.ReadAllRequest{
		Api: apiVersion,
	}
	res4, err := c.ReadAll(ctx, &req4)
	if err != nil {
		log.Fatalf("read all failed: %v", err)
	}
	log.Printf("read all result: <%+v>\n\n", res4)

	return res4
}

func delete(ctx context.Context, c v12.ToDoServiceClient, id int64) *v12.DeleteResponse {
	req5 := v12.DeleteRequest{
		Api: apiVersion,
		Id:  id,
	}
	res5, err := c.Delete(ctx, &req5)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Delete result: <%+v>\n\n", res5)

	return res5
}
