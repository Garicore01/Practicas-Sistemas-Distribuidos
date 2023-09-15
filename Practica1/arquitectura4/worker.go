package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"practica1/com"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
// PRE: verdad
// POST: IsPrime devuelve verdad si n es primo y falso en caso contrario
func IsPrime(n int) (foundDivisor bool) {
	foundDivisor = false
	for i := 2; (i < n) && !foundDivisor; i++ {
		foundDivisor = (n%i == 0)
	}
	return !foundDivisor
}

// PRE: interval.A < interval.B
// POST: FindPrimes devuelve todos los números primos comprendidos en el
//
//	intervalo [interval.A, interval.B]
func FindPrimes(interval com.TPInterval) (primes []int) {
	for i := interval.A; i <= interval.B; i++ {
		if IsPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}


func handleRequestsSec(conn *net.TCPConn) {
	/* Bucle infinito para no perder ninguna Gorutine */ 
	for {
		encoder := gob.NewEncoder(conn)
		decoder := gob.NewDecoder(conn)
		var request com.Request
		
		err := decoder.Decode(&request) //Recibo el mensaje
		checkError(err)
		primes := FindPrimes(request.Interval)              //Busco los primos del intervalo recibido.
		err = encoder.Encode(com.Reply{request.Id, primes}) //Envio los numeros primos encontrados.
		checkError(err)
		defer conn.Close()
	}
}

func main() {
	// Declaramos los parametros de la conexión.
	CONN_TYPE := "tcp"
	CONN_HOST := "155.210.154.206"
	CONN_PORT := os.Args[1]
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	checkError(err)
	

	/* Espero a recibir la petición */
	print("Espero conexion\n")
	conn, err := listener.Accept()
	print("Recibo conexion\n")
	checkError(err)
	handleRequestsSec(conn.(*net.TCPConn))
	defer listener.Close()
	print("Conexión ", conn.RemoteAddr, "\n")
	print("Cierro conexion ", conn.RemoteAddr, "\n")
	//conn.Close()


}