package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/frank5487/pcbook-test/client"
	"github.com/frank5487/pcbook-test/pb"
	"github.com/frank5487/pcbook-test/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

//func createLaptop(laptopClient pb.LaptopServiceClient, laptop *pb.Laptop) {
//	req := &pb.CreateLaptopRequest{Laptop: laptop}
//
//	// set timeout
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	res, err := laptopClient.CreateLaptop(ctx, req)
//	if err != nil {
//		st, ok := status.FromError(err)
//		if ok && st.Code() == codes.AlreadyExists {
//			// not a big deal
//			log.Print("laptop already exists")
//		} else {
//			log.Fatal("cannot create laptop: ", err)
//		}
//		return
//	}
//
//	log.Printf("created laptop with id: %s", res.Id)
//}
//
//func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
//	log.Printf("search filter: ", filter)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	req := &pb.SearchLaptopRequest{Filter: filter}
//	stream, err := laptopClient.SearchLaptop(ctx, req)
//	if err != nil {
//		log.Fatal("cannot search laptop: ", err)
//	}
//
//	for {
//		res, err := stream.Recv()
//		if err == io.EOF {
//			return
//		}
//		if err != nil {
//			log.Fatal("cannot receive response: ", err)
//		}
//
//		laptop := res.GetLaptop()
//		log.Print("- found: ", laptop.GetId())
//	}
//}
//
//func uploadImage(laptopClient pb.LaptopServiceClient, laptopID string, imagePath string) {
//	file, err := os.Open(imagePath)
//	if err != nil {
//		log.Fatal("cannot open image file: ", err)
//	}
//	defer file.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	stream, err := laptopClient.UploadImage(ctx)
//	if err != nil {
//		log.Fatal("cannot upload image: ", err)
//	}
//
//	req := &pb.UploadImageRequest{
//		Data: &pb.UploadImageRequest_Info{
//			Info: &pb.ImageInfo{
//				LaptopId:  laptopID,
//				ImageType: filepath.Ext(imagePath),
//			},
//		},
//	}
//
//	err = stream.Send(req)
//	if err != nil {
//		log.Fatal("cannot send image info: ", err, stream.RecvMsg(nil))
//	}
//
//	reader := bufio.NewReader(file)
//	buffer := make([]byte, 1024)
//
//	for {
//		n, err := reader.Read(buffer)
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatal("cannot read chunk to buffer: ", err)
//		}
//
//		req := &pb.UploadImageRequest{
//			Data: &pb.UploadImageRequest_ChunkData{
//				ChunkData: buffer[:n],
//			},
//		}
//
//		err = stream.Send(req)
//		if err != nil {
//			err2 := stream.RecvMsg(nil)
//			log.Fatal("cannot send chunk to server: ", err, err2)
//		}
//	}
//
//	res, err := stream.CloseAndRecv()
//	if err != nil {
//		log.Fatal("cannot receive response: ", err)
//	}
//
//	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
//}
//
//func rateLaptop(laptopClient pb.LaptopServiceClient, laptopIDs []string, scores []float64) error {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	stream, err := laptopClient.RateLaptop(ctx)
//	if err != nil {
//		return fmt.Errorf("cannot rate laptop: %v", err)
//	}
//
//	waitResponse := make(chan error)
//	// go routine to receive responses
//	go func() {
//		for {
//			res, err := stream.Recv()
//			if err == io.EOF {
//				log.Print("no more responses")
//				waitResponse <- nil
//				return
//			}
//			if err != nil {
//				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
//				return
//			}
//
//			log.Print("received response: ", res)
//		}
//	}()
//
//	// send requests
//	for i, laptopID := range laptopIDs {
//		req := &pb.RateLaptopRequest{
//			LaptopId: laptopID,
//			Score:    scores[i],
//		}
//
//		err := stream.Send(req)
//		if err != nil {
//			return fmt.Errorf("cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
//		}
//
//		log.Print("sent request: ", req)
//	}
//
//	err = stream.CloseSend()
//	if err != nil {
//		return fmt.Errorf("cannot close send: %v", err)
//	}
//
//	err = <-waitResponse
//	return err
//}

func testCreateLaptop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := -1; i < 10; i++ {
		laptop := sample.NewLaptop()
		laptopClient.CreateLaptop(laptop)
	}

	filter := &pb.Filter{
		MaxPriceUsd: 2999,
		MinCpuCores: 3,
		MinCpuGhz:   1.5,
		MinRam: &pb.Memory{
			Value: 7,
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *client.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(), "tmp/laptop.jpg")
}

func testRateLaptop(laptopClient *client.LaptopClient) {
	n := 3
	laptopIDs := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop (y/n)? ")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}

		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}

		err := laptopClient.RateLaptop(laptopIDs, scores)
		if err != nil {
			log.Fatal(err)
		}

	}
}

const (
	username = "admin1"
	//username        = "user1"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const laptopServicePath = "/frank5487.pcbook.LaptopService/"

	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadImage":  true,
		laptopServicePath + "RateLaptop":   true,
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// load certificate of the CA who signed the server's certificate
	pemServerCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	// create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	cc1, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.Dial(
		*serverAddress,
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := client.NewLaptopClient(cc2)
	//testCreateLaptop(laptopClient)
	//testSearchLaptop(laptopClient)
	//testUploadImage(laptopClient)
	testRateLaptop(laptopClient)
}
