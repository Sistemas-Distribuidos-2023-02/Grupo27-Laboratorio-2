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
	"sync"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo27-Laboratorio-2/protos"
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
		host="dist107.inf.santiago.usm.cl"
		puerto ="50052"
		nombre ="DataNode1"
	}else if host == "DataNode2"{
		host="dist108.inf.santiago.usm.cl"
		puerto ="50052"
		nombre ="DataNode2"
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
		response, err := c.OmsToDataNode(context.Background(), &pb.Message{Body: mensaje})
		if err != nil {
			log.Println("Server "+nombre+" not responding: ")
			log.Println("Trying again in 10 seconds. . .")
			time.Sleep(10 * time.Second)
			continue
		}
		log.Printf("Response from server "+nombre+": "+"%s\n\n", response.Body)
		return response.Body
	}
}


func (s *Server)RegionalToOms(ctx context.Context, in *pb.Message)(*pb.Message, error){
	inMessage:=string(in.Body)
	split:=strings.Split(inMessage,":") //nombre-apellido-status
	server:=split[0]
	inMessage=split[1]

	log.Printf("Received message from %s: %s", server,inMessage)

	
	//OBTENER DIRECTORIO ACTUAL
	directorioActual, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
	}

	if strings.Contains(inMessage, "-") && (strings.Contains(inMessage, "infectado") || strings.Contains(inMessage, "muerto")){
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

		err = file.Close()
		if err != nil {
			fmt.Println("Error al cerrar el archivo:", err)
		}

		return &pb.Message{Body: "OK"}, nil
	}else{
		return &pb.Message{Body: "Mensaje no valido"}, nil
	}
}

func (s *Server)OnuToOms(ctx context.Context, in *pb.Message)(*pb.Message, error){
	log.Printf("Received message from ONU: %s", in.Body)

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
			if len(lineas[i]) <= 0{
				continue
			}
			//fmt.Println("Linea de DATA.txt:",lineas[i])
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
		wg:=sync.WaitGroup{}
		for i := 0; i < len(infectados_id); i++ {
			split:=strings.Split(infectados_id[i],"-") //id-datanode
			id:=split[0]
			datanode:=split[1]

			if datanode == "1"{
				//Pedir a DataNode1
				wg.Add(1)
				go func(){
				response:=ConexionGRPC(id,"DataNode1")
				infectados = append(infectados, response)
				defer wg.Done()
				}()
				
			}else if datanode == "2"{
				//Pedir a DataNode2
				wg.Add(1)
				go func(){
				response:=ConexionGRPC(id,"DataNode2")
				infectados = append(infectados, response)
				defer wg.Done()
				}()
			}
		}
		wg.Wait()
		infectados_response:=strings.Join(infectados, "\n")
		log.Println("Sending message to ONU: INFECTADOS")
		return &pb.Message{Body:"\n"+infectados_response+"\n"}, nil

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
			if len(lineas[i]) <= 0{
				continue
			}

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
		wg:=sync.WaitGroup{}
		for i := 0; i < len(infectados_id); i++ {
			split:=strings.Split(infectados_id[i],"-") //id-datanode
			id:=split[0]
			datanode:=split[1]
			
			if datanode == "1"{
				//Pedir a DataNode1
				wg.Add(1)
				go func(){
				response:=ConexionGRPC(id,"DataNode1")
				infectados = append(infectados, response)
				defer wg.Done()
				}()
				
			}else if datanode == "2"{
				//Pedir a DataNode2
				wg.Add(1)
				go func(){
				response:=ConexionGRPC(id,"DataNode2")
				infectados = append(infectados, response)
				defer wg.Done()
				}()
			}
		}
		wg.Wait()
		infectados_response:=strings.Join(infectados, "\n")
		log.Println("Sending message to ONU: MUERTOS")
		return &pb.Message{Body:"\n"+infectados_response+"\n"}, nil

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
func RecoverIds() (string){
	fmt.Println("Recovering id's. . .\n")
	directorioActual, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
		return ""
	}
	content, err := os.ReadFile(filepath.Join(directorioActual,"OMS","DATA.txt"))
	if err != nil {
		log.Fatal(err)
	}
	lineas := strings.Split(string(content), "\n")

	for i := 0; i < len(lineas); i++ {
		if len(lineas[i]) <= 0{
			continue
		}
		split:=strings.Split(lineas[i],"-") //id-datanode-status
				
		id,_:=strconv.Atoi(split[0])
		ids=append(ids,id)
	}
	return ""
}


var server_name string
func main() {
	
	server_name="OMS"
	fmt.Println("Starting "+server_name+" . . .\n")

	RecoverIds()

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