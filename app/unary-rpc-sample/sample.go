package unaryrpcsample

import (
	"app/unary-rpc-sample/pb"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	directoryPath string = "./unary-rpc-sample/storage/"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {

	paths, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Fatalln("error: ", err)
	}

	var files []string

	for _, path := range paths {
		if !path.IsDir() {
			files = append(files, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Filenames: files,
	}

	return res, nil

}

func (*server) DonwloadFile(req *pb.DownloadRequest, stream pb.FileService_DonwloadFileServer) error {
	fmt.Println("download starts...")

	filename := req.GetFilname()
	path := directoryPath + filename

	file, err := os.Open(path)

	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		res := &pb.DownloadResponse{Data: buf[:n]}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}

		time.Sleep(1 * time.Second)
	}
	return nil

}

func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("upload starts")

	var buf bytes.Buffer
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &pb.UploadResponse{Size: int32(buf.Len())}
			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Println("received data: ", string(data))
		buf.Write(data)
	}
}

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("start upload notify")

	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		data := req.GetData()
		log.Println("received data", data)
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("received %vbytes", size),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}

// 一旦未使用
func RunServerGinSample() {
	// ginの実装
	router := gin.Default()
	router.Run(fmt.Sprintf(":%d", 8084))

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})
}

func SwichClientOrServer() {
	var callType string = ""
	args := os.Args

	if len(args) == 2 {
		callType = "server"
	}

	switch {
	case callType == "server":
		RunServerSample()
	default:
		RunClientServer()

	}

}

func RunClientServer() {
	fmt.Println("start run client grpc")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("error client sever: ", err)
	}

	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	// callListFiles(client)
	// callDownloadFiles(client)
	// callUpload(client)
	callUploadAndNotifyProgress(client)

}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})

	if err != nil {
		log.Fatalln("error call list", err)
	}
	fmt.Println(res.GetFilenames())
}

func callDownloadFiles(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filname: "sample.txt"}
	stream, err := client.DonwloadFile(context.Background(), req)
	if err != nil {
		log.Fatalln("error stream: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("otger stream error: ", err)
		}

		log.Println("response: ", string(res.GetData()))
	}

}

func RunServerSample() {
	fmt.Println("run server grpc sample")
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalln("failed backend server", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("error: ", err)
	}
}

func callUpload(client pb.FileServiceClient) {
	filename := "sample.txt"
	path := directoryPath + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("error: ", err)
	}
	defer file.Close()

	stream, err := client.Upload(context.Background())

	if err != nil {
		log.Fatalln("error: ", err)
	}

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln("error: ", err)
		}

		req := &pb.UploadRequest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln("error: ", err)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("error: ", err)
	}

	log.Println("get size: ", res.GetSize())
}

func callUploadAndNotifyProgress(client pb.FileServiceClient) {
	filename := "sample.txt"
	path := directoryPath + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("error: ", err)
	}
	defer file.Close()

	stream, _ := client.UploadAndNotifyProgress(context.Background())

	buf := make([]byte, 5)
	go func() {
		for {
			n, err := file.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("gorouting error: ", err)
			}

			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]}
			sendErr := stream.Send(req)
			if sendErr != nil {
				log.Fatalln(sendErr)
			}

			time.Sleep(1 * time.Second)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("received message: %v", res.GetMsg())
		}
		close(ch)
	}()
	<-ch
}
