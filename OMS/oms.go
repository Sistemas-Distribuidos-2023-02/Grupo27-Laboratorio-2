package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	//"os"
	//"path/filepath"
	"strings"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo27-Laboratorio-1/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct{
	pb.UnimplementedChatServiceServer
}
func ConexionGRPC(mensaje string ) (string){
	
	//Uno de estos debe cambiar quizas por "regional:50052" ya que estara en la misma VM que el central
	//host :="localhost"
	var puerto, nombre, host string
	host="dist108.inf.santiago.usm.cl"
	puerto ="50055"
	nombre ="OMS"
	
	log.Println("Connecting to server "+nombre+": "+host+":"+puerto+". . .")
	conn, err := grpc.Dial(host+":"+puerto,grpc.WithTransportCredentials(insecure.NewCredentials()))	
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	fmt.Printf("Esperando\n")
	defer conn.Close()

	c := pb.NewChatServiceClient(conn)
	for {
		log.Println("Sending message to server "+nombre+": "+mensaje)
		response, err := c.SayHello(context.Background(), &pb.Message{Body: mensaje})
		if err != nil {
			log.Println("Server "+nombre+" not responding: ")
			log.Println("Trying again in 10 seconds. . .")
			time.Sleep(10 * time.Second)
			continue
		}
		log.Printf("Response from server "+nombre+": "+"%s", response.Body)
		return response.Body
	}
}

func (s *Server)SayHello(ctx context.Context, in *pb.Message)(*pb.Message, error){
	log.Printf("Receive message body from client: %s", in.Body)

	inMessage:=string(in.Body)

	//OBTENER DIRECTORIO ACTUAL
	/*directorioActual, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
	}*/
	
	//ESTO SE DEBBE CAMBIAR
	if inMessage == "I"{
		//Pedir infectados a DataNodes y devolverlos a ONU
		
		response:=ConexionGRPC("1")
		return &pb.Message{Body: response}, nil

	}else if inMessage == "M"{
		response:=ConexionGRPC("1")
		return &pb.Message{Body: response}, nil
		//Pedir muertos a DataNodes y devolverlos a ONU
	}else if strings.Contains(inMessage, "-") && (strings.Contains(inMessage, "infectado") || strings.Contains(inMessage, "muerto")){
		//Crear id no existente
		//Agregar nombre y apellido en DataNodes
		response:=ConexionGRPC("1")
		return &pb.Message{Body: response}, nil
	}else{
		return &pb.Message{Body: "Mensaje no valido"}, nil
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
	go func (){pb.RegisterChatServiceServer(grpcServer_regional, server_regional)
	if err_regional:= grpcServer_regional.Serve(lis_regional); err_regional != nil {
		panic(err_regional)
	}
	}()

	puerto_onu:= ":50053"
	lis_onu, err_onu:= net.Listen("tcp", puerto_onu)
	fmt.Printf("Escuchando %s\n", puerto_onu)
	if err_onu != nil {
		panic(err_onu)
	}

	grpcServer_onu:= grpc.NewServer()
	server_onu:= &Server{}
	go func () {pb.RegisterChatServiceServer(grpcServer_onu, server_onu)
	if err_onu := grpcServer_onu.Serve(lis_onu); err_onu != nil {
		panic(err_onu)
	}
	}()
	
}