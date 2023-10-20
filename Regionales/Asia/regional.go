package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo27-Laboratorio-1/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func ConexionGRPC(mensaje string ){
	
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
		break
	}
}

func ObtenerNombre() string{
	directorioActual, err := os.Getwd()
    if err != nil {
        fmt.Println("Error al obtener el directorio actual:", err)
        return ""
    }
    content, err := os.ReadFile(directorioActual+"\\Regionales\\names.txt")
	if err != nil {
		log.Fatal(err)
	}

	lineas := strings.Split(string(content), "\n")
	rand_num:=rand.Intn(len(lineas))

	linea:=lineas[rand_num] 
	nombre,apellido:=strings.Split(linea," ")[0],strings.Split(linea," ")[1]
	return nombre+"-"+apellido
}

func ObtenerStatus() string{
	rand_num:=rand.Intn(100)
	if rand_num > 55{
		return "muerto"
	}else {
		return "infectado"
	}
}

var server_name string
func main() {
	
	nombresUsados := make(map[string]bool)
	//LEER EL ARCHIVO
	
	//OBTENER NOMBRE
	server_name = "Asia"
	fmt.Println("Iniciando regional "+server_name+" . . .")
	//MANDAR 5 DATOS
	var nombre_apellido string
	var status string
	for i := 0; i < 5; i++ {
		nombre_apellido=ObtenerNombre()
		
		for {
			if nombresUsados[nombre_apellido]{
				nombre_apellido=ObtenerNombre()
			}else{
				nombresUsados[nombre_apellido]=true
				break
			}
		}

		status=ObtenerStatus()
		ConexionGRPC(nombre_apellido+"-"+status)
		//fmt.Println(nombre_apellido+"-"+status)
	}
	fmt.Println("Termino 5")
	for{
		time.Sleep(3*time.Second)
		nombre_apellido=ObtenerNombre()

		for {
			if nombresUsados[nombre_apellido]{
				nombre_apellido=ObtenerNombre()
			}else{
				nombresUsados[nombre_apellido]=true
				break
			}
		}

		status=ObtenerStatus()
		ConexionGRPC(nombre_apellido+"-"+status)
		//fmt.Println(nombre_apellido+"-"+status)
	}
	
}
