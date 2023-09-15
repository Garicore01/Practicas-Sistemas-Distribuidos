/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: server.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*				correspondientes a la práctica 1
 */
package main

import (
	"encoding/gob"
	//"golang.org/x/crypto/ssh"
	//"golang.org/x/crypto/ssh/agent"
	"fmt"
	"net"
	"os"
	"os/exec"
	"practica1/com"
	"strconv"
)


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}


func makeConnWorker(endpoint string, request com.Request) ( com.Reply ){
	tcpAddr, err := net.ResolveTCPAddr("tcp", endpoint)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	fmt.Printf("Mando a worker")
	err = encoder.Encode(request) //Aqui se envia el request
	checkError(err)

	var reply com.Reply
	err = decoder.Decode(&reply) //Espero respueta del servidor
	checkError(err)
	conn.Close()
	return reply
}

/*
* PRE: conn debe ser una conexión valida.
* POST: handleRequestsSec, busca y envia al cliente los números primos encontrados en el intervalo solicitado.
 */
func handleRequestsSec(jobs chan *net.TCPConn,id int,endpoint string) {
	/* Enciendo el worker, mediante una conexión SSH */
	sshConn(id,endpoint)
	
	/* Bucle infinito para no perder ninguna Gorutine */ 
	for {
		conn := <- jobs
		encoder := gob.NewEncoder(conn)
		decoder := gob.NewDecoder(conn)
		var request com.Request
		err := decoder.Decode(&request) //Recibo el mensaje
		checkError(err)
		/* Envio al Worker que me calcule los primos */
		var reply com.Reply
		reply = makeConnWorker(endpoint,request)

		/*Mandar al cliente los datos calculados*/
		err = encoder.Encode(reply) //Envio los numeros primos encontrados.
		checkError(err)
		defer conn.Close()
	}
}



func sshConn(puerto int,endpoint string){
    // Comando que deseas ejecutar en tu máquina.

	comando := "/usr/bin/ssh"

	argument1:= "a848905@"+endpoint
	
	goCommand := "cd /home/a848905/Practicas/Distribuidos/practica1/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a848905/Practicas/Distribuidos/practica1/worker.go " + strconv.Itoa(puerto)

    cmd := exec.Command(comando,argument1,goCommand)
	err := cmd.Start()
	


	if err != nil {
        fmt.Printf("Error al ejecutar el comando 1: %v\n", err)
        return
    }

	fmt.Printf("hola\n")
}

/*
func sshConn(puerto int, endpoint string){

	server := endpoint+":22"//+strconv.Itoa(puerto)
	// Configuración de la conexión SSH 
	sshConfig := &ssh.ClientConfig{
		User: "root",//"a849183",
		Auth: []ssh.AuthMethod{
			ssh.Password("Welcome1."),
		},
	}
	//Realizo la conexión SSH 
	client, err := ssh.Dial("tcp", server, sshConfig)

	if err != nil {
		fmt.Printf("Error al conectar: %v", err)
	}

	defer client.Close()
	session, err := client.NewSession()

	if err != nil {
		fmt.Printf("Error al crear la sesión: %v", err)
	}

	defer session.Close()
	cmd := "go run /home/gari/Documentos/worker.go "+strconv.Itoa(puerto)
	output, err := session.CombinedOutput(cmd)
	
	if err != nil {
		fmt.Printf("Error al ejecutar la instrucción: %v", err)
	}

	// Imprime la salida de la instrucción
	fmt.Println(string(output))

}*/


func main() {
	// Declaramos los parametros de la conexión.
	CONN_TYPE := "tcp"
	CONN_HOST := "127.0.0.1"
	CONN_PORT := "31010"
	MAX_WORKER := 10
	endpoint:= "155.210.154.206"
	
	

	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	checkError(err)
	defer listener.Close()
	/* Creo un canal de capacidad 10 */
	jobs := make(chan *net.TCPConn,MAX_WORKER)
	/* Lanzo el pool de Gorutines */ 
	for j:= 0; j < MAX_WORKER; j++ { 
		go handleRequestsSec(jobs,j+31010,endpoint)
	}
	/* Voy recibiendo  peticiones */
	for {
		conn, err := listener.Accept()
		checkError(err)
		print("Conexión ", conn.RemoteAddr, "\n")
		/* En esta arquitectura concurrente, tenemos que aceptar varias peticiones simultaneamente, pero ya tenemos un pool de Gorutines esperando a recibir 
		   trabajo, por lo que tenemos que enviar la información que necesita mediante el canal sincrono "jobs" */ 
		jobs <- conn.(*net.TCPConn)
		print("Cierro conexion ", conn.RemoteAddr, "\n")
		//conn.Close()
	}

}
