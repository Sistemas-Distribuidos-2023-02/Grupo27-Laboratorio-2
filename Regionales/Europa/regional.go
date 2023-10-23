package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo27-Laboratorio-2/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func ConexionGRPC(mensaje string ){
	
	//Uno de estos debe cambiar quizas por "regional:50052" ya que estara en la misma VM que el central
	//host :="localhost"
	var puerto, nombre, host string
	host="dist105.inf.santiago.usm.cl"
	puerto ="50052"
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
		response, err := c.RegionalToOms(context.Background(), &pb.Message{Body: mensaje})
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

var nombresUsados []int = make([]int, 407)

func ObtenerNombre() string{
	directorioActual, err := os.Getwd()
    if err != nil {
        fmt.Println("Error al obtener el directorio actual:", err)
        return ""
    }
    content, err := os.ReadFile(filepath.Join(directorioActual,"Regionales","names.txt"))
	if err != nil {
		log.Fatal(err)
	}
	remove:=func (s []int, i int) []int {
		
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	}

	try:=func()string{

		if len(nombresUsados) == 0 {
			fmt.Println("\nNo hay mas nombres disponibles")
			os.Exit(0)
		}

		rand.Seed(time.Now().UnixNano())
		lineas := strings.Split(string(content), "\n")

		var rand_num int
		if len(nombresUsados) == 1 {
			rand_num=0
		}else{
			rand_num=rand.Intn(len(nombresUsados)-1)
		}
		
		
		linea:=lineas[nombresUsados[rand_num]]

		nombresUsados=remove(nombresUsados,rand_num)
		nombre:=strings.Split(linea," ")[0]
		apellido:=strings.Split(linea," ")[1]
	
		nombre_apellido:=nombre+"-"+apellido
		nombre_apellido=strings.Replace(nombre_apellido, "\r", "", -1)

		return nombre_apellido}
	
	nombre_apellido:=try()
	return nombre_apellido
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
		
	server_name = "Europa"
	fmt.Println("Iniciando regional "+server_name+" . . .\n")

	for i := 0; i < 407; i++ {
		nombresUsados[i] = i
	}

	var nombre_apellido string
	var status string

	//MANDAR 5 DATOS INICIALES
	for i := 0; i < 5; i++ {
		nombre_apellido=ObtenerNombre()
		status=ObtenerStatus()
		ConexionGRPC(server_name+":"+nombre_apellido+"-"+status)
		//fmt.Println(nombre_apellido+"-"+status)
	}

	fmt.Println("\nSe mandaron 5 Nombres iniciales ...\nMandando datos cada 3 segundos ...\n")

	//MANDAR DATOS CADA 3 SEGUNDOS
	for{
		time.Sleep(3*time.Second)
		nombre_apellido=ObtenerNombre()
		status=ObtenerStatus()
		ConexionGRPC(server_name+":"+nombre_apellido+"-"+status)
		//fmt.Println(nombre_apellido+"-"+status)
	}
	
}
