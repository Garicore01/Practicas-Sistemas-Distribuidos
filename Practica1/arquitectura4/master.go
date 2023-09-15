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
	"fmt"
	"net"
	"os"
	"os/exec"
	"practica1/com"
	//"strconv"
	"time"
	"bufio"
	"strings"
)
type Cola struct {
	items []string
}

// Enqueue agrega un elemento al final de la cola
func (c *Cola) Enqueue(item string) {
	c.items = append(c.items, item)
}

// Dequeue remueve y devuelve el elemento en la parte frontal de la cola
func (c *Cola) Dequeue() string {
	if len(c.items) == 0 {
		return ""
	}
	item := c.items[0]
	c.items = c.items[1:]
	return item
}

// IsEmpty devuelve true si la cola está vacía, de lo contrario, false
func (c *Cola) IsEmpty() bool {
	return len(c.items) == 0
}

// Size devuelve el número de elementos en la cola
func (c *Cola) Size() int {
	return len(c.items)
}

/*
* PRE: true
* POST: Finaliza el programa si solo si <err> es un error.
*/
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

/* 
* PRE: <endpoint> debe ser una @ip valida, <request> debe ser la request solicitada por el cliente.
* POST: Devuelve el Reply (Primes) que devuelve el Worker
*/
func makeConnWorker(endpoint string, request com.Request) ( com.Reply ){
	tcpAddr, err := net.ResolveTCPAddr("tcp", endpoint)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	// Envio la Request al Worker
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	err = encoder.Encode(request) //Aqui se envia el request
	checkError(err)

	// Recibir el Reply del Worker con los numeros primos solicitados
	var reply com.Reply
	err = decoder.Decode(&reply) 
	checkError(err)
	conn.Close()
	return reply
}

/*
* PRE: conn debe ser una conexión valida.
* POST: handleRequestsSec, busca y envia al cliente los números primos encontrados en el intervalo solicitado.
 */
func handleRequestsSec(jobs chan *net.TCPConn, endpoint string) {
	
	/* Bucle infinito para no perder ninguna Gorutine */ 
	for {
		// Nos bloqueamos hasta recibir una petición
		conn := <- jobs
		// Levantamos el proceso del Worker
		sshConn(endpoint)

		encoder := gob.NewEncoder(conn)
		decoder := gob.NewDecoder(conn)
		var request com.Request
		//Recibo el mensaje
		err := decoder.Decode(&request) 
		checkError(err)
		var reply com.Reply
		// Esperamos 3 segundos para que le tiempo al Worker a arrancar el proceso.
		time.Sleep(time.Duration(3000) * time.Millisecond)

		// Envio al Worker que me calcule los primos
		reply = makeConnWorker(endpoint,request)

		// Mandar al cliente los datos calculados
		err = encoder.Encode(reply) 
		checkError(err)
		defer conn.Close()
	}
}

/* 
* PRE: <endpoint> debe ser una @ip valida
* POST: El proceso Worker de la maquina <endpoint>, inicia su ejecución en el puerto que especifica <endpoint>
*/
func sshConn(endpoint string){
	comando := "/usr/bin/ssh"
	// Separo la @IP del puerto
	ip_worker := strings.Split(endpoint, ":")
	argument1 := "root@" + ip_worker[0]
	puerto := ip_worker[1]
	
	//goCommand := "cd /home/a848905/Practicas/Distribuidos/practica1/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a848905/Practicas/Distribuidos/practica1/worker.go " + strconv.Itoa(puerto)
	//goCommand := "cd /home/a849183/Desktop/practica1/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a849183/Desktop/practica1/worker.go " + strconv.Itoa(puerto)
	goCommand := "export PATH=$PATH:/usr/local/go/bin; cd /root/practica1/; go mod tidy; nohup go run /root/practica1/worker.go " + puerto
    cmd := exec.Command(comando,argument1,goCommand)
	err := cmd.Start()

	if err != nil {
        fmt.Printf("Error al ejecutar el comando 1: %v\n", err)
        return
    }
}

func leerIpsFichero(nombre string) (Cola) {
	cola:= Cola{}
	file, err := os.Open(nombre)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		
	} else { 
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		
		for scanner.Scan() {

			cola.Enqueue(scanner.Text())
		}
	}
	return cola
	
}

func main() {
	// Declaramos los parametros de la conexión.
	CONN_TYPE 	:= "tcp"
	CONN_HOST 	:= "192.168.1.139"
	CONN_PORT 	:= "31010"
	MAX_WORKER 	:= 10

	//Leer endpoint
	
	queue := leerIpsFichero("listaWorkers.txt")

	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	checkError(err)
	defer listener.Close()

	/* Creo un canal de capacidad 10 */
	jobs := make(chan *net.TCPConn,MAX_WORKER)

	/* Lanzo el pool de Gorutines */ 
	for j:= 0; j < MAX_WORKER && !queue.IsEmpty() ; j++  { 
		go handleRequestsSec(jobs,queue.Dequeue())
	}
	
	/* Voy recibiendo  peticiones */
	for {
		conn, err := listener.Accept()
		checkError(err)
		print("Nueva conexión ", conn.RemoteAddr, "\n")
		/* En esta arquitectura concurrente, tenemos que aceptar varias peticiones simultaneamente, pero ya tenemos un pool de Gorutines esperando a recibir 
		   trabajo, por lo que tenemos que enviar la información que necesita mediante el canal sincrono "jobs" */ 
		jobs <- conn.(*net.TCPConn)
		//conn.Close()
	}

}
