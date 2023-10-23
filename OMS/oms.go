package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"os"
	"path/filepath"
	"strings"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo27-Laboratorio-1/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct{
	pb.UnimplementedChatServiceServer
}
func ConexionGRPC(mensaje string, host string ) (string){
	
	//Uno de estos debe cambiar quizas por "regional:50052" ya que estara en la misma VM que el central
	//host :="localhost"
	var puerto, nombre string
	if host == "DataNode1"{
		host="dist108.inf.santiago.usm.cl"
		puerto ="50055"
		nombre ="OMS"
	}else if host == "DataNode2"{
		host="dist108.inf.santiago.usm.cl"
		puerto ="50055"
		nombre ="OMS"
	}
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
	directorioActual, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
	}
	
	//ESTO SE DEBBE CAMBIAR
	if inMessage == "I"{
		//Pedir infectados a DataNodes y devolverlos a ONU

		//Leer Archivo DATA.txt
		content, err := os.ReadFile(filepath.Join(directorioActual,"OMS","DATA.txt"))
		if err != nil {
			log.Fatal(err)
		}
		lineas := strings.Split(string(content), "\n")

		var infectados_id []string
		for i := 0; i < len(lineas); i++ {
			split:=strings.Split(lineas[i],"-") //id-datanode-status
			id:=split[0]
			datanode:=split[1]
			status:=split[2]

			if status == "infectado"{
				infectados_id = append(infectados_id, id+"-"+datanode)
				//fmt.Println("Infectados: ",infectados_id)
			}
		}


		var infectados []string
		for i := 0; i < len(infectados_id); i++ {
			split:=strings.Split(infectados_id[i],"-") //id-datanode
			id:=split[0]
			datanode:=split[1]

			if datanode == "1"{
				//Pedir a DataNode1
				response:=ConexionGRPC(id,"DataNode1")
				infectados = append(infectados, response)
				
			}else if datanode == "2"{
				//Pedir a DataNode2
				response:=ConexionGRPC(id,"DataNode2")
				infectados = append(infectados, response)
			}
		}
		
		infectados_response:=strings.Join(infectados, "\n")
		return &pb.Message{Body: infectados_response}, nil

	}else if inMessage == "M"{
		//Pedir muertos a DataNodes y devolverlos a ONU

		//Leer Archivo DATA.txt
		content, err := os.ReadFile(filepath.Join(directorioActual,"OMS","DATA.txt"))
		if err != nil {
			log.Fatal(err)
		}
		lineas := strings.Split(string(content), "\n")

		var infectados_id []string
		for i := 0; i < len(lineas); i++ {
			split:=strings.Split(lineas[i],"-") //id-datanode-status
			id:=split[0]
			datanode:=split[1]
			status:=split[2]

			if status == "muerto"{
				infectados_id = append(infectados_id, id+"-"+datanode)
				//fmt.Println("Infectados: ",infectados_id)
			}
		}


		var infectados []string
		for i := 0; i < len(infectados_id); i++ {
			split:=strings.Split(infectados_id[i],"-") //id-datanode
			id:=split[0]
			datanode:=split[1]

			if datanode == "1"{
				//Pedir a DataNode1
				response:=ConexionGRPC(id,"DataNode1")
				infectados = append(infectados, response)
				
			}else if datanode == "2"{
				//Pedir a DataNode2
				response:=ConexionGRPC(id,"DataNode2")
				infectados = append(infectados, response)
			}
		}
		
		infectados_response:=strings.Join(infectados, "\n")
		return &pb.Message{Body: infectados_response}, nil

	}else if strings.Contains(inMessage, "-") && (strings.Contains(inMessage, "infectado") || strings.Contains(inMessage, "muerto")){
		//Crear id no existente
		//Agregar nombre y apellido en DataNodes
		nuevo_id:=CrearId()
		split:=strings.Split(inMessage,"-") //nombre-apellido-status
		nombre:=split[0]
		apellido:=split[1]
		status:=split[2]
		var datanode string

		primera_letra:=string(apellido[0])
		if primera_letra < "M"{
			datanode="1"
		}else{
			datanode="2"
		}

		//Agregar a DATA.txt
		file, err := os.OpenFile(filepath.Join(directorioActual,"OMS","DATA.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
		if err != nil{
			fmt.Println("Ha ocurrido un error en la creacion del archivo: ",err)
		}
		fmt.Fprintln(file, nuevo_id+"-"+datanode+"-"+status)

		//Agregar a DataNode
		if datanode == "1"{
			ConexionGRPC(nuevo_id+"-"+nombre+"-"+apellido,"DataNode1")
		}else if datanode == "2"{
			ConexionGRPC(nuevo_id+"-"+nombre+"-"+apellido,"DataNode2")
		}

		return &pb.Message{Body: "OK"}, nil
	}else{
		return &pb.Message{Body: "Mensaje no valido"}, nil
	}
}

var ids[] int
func CrearId() (string){
	total:=len(ids)
	id:=total+1
	ids=append(ids,id)

	id_string:=strconv.Itoa(id)

	return id_string
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

	go func ()  {
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
	pb.RegisterChatServiceServer(grpcServer_onu, server_onu)

	go func ()  {
		if err_onu := grpcServer_onu.Serve(lis_onu); err_onu != nil {
			panic(err_onu)
		}
	}()
	
	select {}
}