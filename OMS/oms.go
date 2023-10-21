package main

import (
	"context"
	"fmt"
	"log"
	"net"
	//"os"
	//"path/filepath"
	"strings"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo27-Laboratorio-1/protos"
	"google.golang.org/grpc"
)

type Server struct{
	pb.UnimplementedChatServiceServer
}

func (s *Server)SayHello(ctx context.Context, in *pb.Message)(*pb.Message, error){
	log.Printf("Receive message body from client: %s", in.Body)

	inMessage:=string(in.Body)

	/*directorioActual, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
	}*/
	
	//ESTO SE DEBBE CAMBIAR
	if inMessage == "I"{
		//Pedir infectados a DataNodes y devolverlos a ONU
	}else if inMessage == "M"{
		//Pedir muertos a DataNodes y devolverlos a ONU
	}else if strings.Contains(inMessage, "-") && (strings.Contains(inMessage, "infectado") || strings.Contains(inMessage, "muerto")){
		//Crear id no existente
		//Agregar nombre y apellido en DataNodes
	}



	//ESTO NO VA PERO UTIL PARA TOMAR DE REFERENCIA
	/*if len(inMessage) > 1{
		fileDataNode, err := os.OpenFile(filepath.Join(directorioActual,"DataNode","Data2","Data.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
		if err != nil{
			fmt.Println("Ha ocurrido un error en la creacion del archivo: ",err)
		}
		fmt.Fprintln(fileDataNode, inMessage)
		return &pb.Message{Body: "OK"}, nil
	}else{
		content, err := os.ReadFile(filepath.Join(directorioActual,"DataNode","Data2","Data.txt"))
		if err != nil {
			log.Fatal(err)
		}
		lineas := strings.Split(string(content), "\n")

		for i := 0; i < len(lineas); i++ {
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
	}*/
}
var server_name string
func main() {
	
	server_name="OMS"
	fmt.Println("Starting "+server_name+" . . .")

	puerto_regional:= ":50052"
	lis_regional, err_regional:= net.Listen("tcp", puerto_regional)
	fmt.Printf("Escuchando %s\n", puerto_regional)
	if err_regional != nil {
		panic(err_regional)
	}

	grpcServer_regional:= grpc.NewServer()
	server_regional:= &Server{}
	pb.RegisterChatServiceServer(grpcServer_regional, server_regional)
	if err_regional:= grpcServer_regional.Serve(lis_regional); err_regional != nil {
		panic(err_regional)
	}

	puerto_onu:= ":50053"
	lis_onu, err_onu:= net.Listen("tcp", puerto_onu)
	fmt.Printf("Escuchando %s\n", puerto_onu)
	if err_onu != nil {
		panic(err_onu)
	}

	grpcServer_onu:= grpc.NewServer()
	server_onu:= &Server{}
	pb.RegisterChatServiceServer(grpcServer_onu, server_onu)
	if err_onu := grpcServer_onu.Serve(lis_onu); err_onu != nil {
		panic(err_onu)
	}

}