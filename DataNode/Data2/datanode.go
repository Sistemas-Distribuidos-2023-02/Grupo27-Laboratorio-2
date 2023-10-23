package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo27-Laboratorio-2/protos"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedChatServiceServer
}

func (s *Server)SayHello(ctx context.Context, in *pb.Message)(*pb.Message, error){
	log.Printf("Receive message body from client: %s", in.Body)

	inMessage:=string(in.Body)

	directorioActual, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
	}
	
	if len(inMessage) > 1{
		fileDataNode, err := os.OpenFile(filepath.Join(directorioActual,"DataNode","Data2","DATA.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
		if err != nil{
			fmt.Println("Ha ocurrido un error en la creacion del archivo: ",err)
		}
		fmt.Fprintln(fileDataNode, inMessage)
		return &pb.Message{Body: "OK"}, nil
	}else{
		content, err := os.ReadFile(filepath.Join(directorioActual,"DataNode","Data2","DATA.txt"))
		if err != nil {
			log.Fatal(err)
		}
		lineas := strings.Split(string(content), "\n")

		for i := 0; i < len(lineas); i++ {
			if len(lineas[i]) <= 0{
				continue
			}
			split:=strings.Split(lineas[i],"-")//id-nombre-apellido
			id:=split[0]
			nombre:=split[1]
			apellido:=split[2]

			nombre_apellido:=nombre+"-"+apellido
			nombre_apellido=strings.Replace(nombre_apellido, "\r", "", -1)

			if id == inMessage {
				return &pb.Message{Body: nombre_apellido}, nil
			}
		}
		return &pb.Message{Body: "ID no Encontrado"}, nil
	}
}

func (s *Server)OmsToDataNode(ctx context.Context, in *pb.Message)(*pb.Message, error){
	log.Printf("Receive message body from client: %s", in.Body)

	inMessage:=string(in.Body)

	directorioActual, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
	}
	
	if len(inMessage) > 1{
		fileDataNode, err := os.OpenFile(filepath.Join(directorioActual,"DataNode","Data1","DATA.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
		if err != nil{
			fmt.Println("Ha ocurrido un error en la creacion del archivo: ",err)
		}
		fmt.Fprintln(fileDataNode, inMessage)
		return &pb.Message{Body: "OK"}, nil
	}else{
		content, err := os.ReadFile(filepath.Join(directorioActual,"DataNode","Data1","DATA.txt"))
		if err != nil {
			log.Fatal(err)
		}
		lineas := strings.Split(string(content), "\n")

		for i := 0; i < len(lineas); i++ {
			if len(lineas[i]) <= 0{
				continue
			}
			split:=strings.Split(lineas[i],"-")//id-nombre-apellido
			id:=split[0]
			nombre:=split[1]
			apellido:=split[2]

			nombre_apellido:=nombre+"-"+apellido
			nombre_apellido=strings.Replace(nombre_apellido, "\r", "", -1)

			if id == inMessage {
				return &pb.Message{Body: nombre_apellido}, nil
			}
		}
		return &pb.Message{Body: "ID no Encontrado"}, nil
	}
}

var DataNode_name string
func main(){
	DataNode_name="DataNode2"
	fmt.Println("Starting "+DataNode_name+" . . .")

	puerto := ":50052"
	lis, err := net.Listen("tcp", puerto)
	fmt.Printf("Escuchando %s\n", puerto)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{}
	pb.RegisterChatServiceServer(grpcServer, server)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	
}

 