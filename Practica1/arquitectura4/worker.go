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

func main() {
	// Declaramos los parametros de la conexión.
	CONN_TYPE := "tcp"
	CONN_HOST := "127.0.0.1"
	CONN_PORT := "31010"
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	checkError(err)
	defer listener.Close()
	/* Creo un canal de capacidad 10 */
	jobs := make(chan *net.TCPConn,MAX_JOBS)
	/* Lanzo el pool de Gorutines */ 
	for j:= 0; j < MAX_JOBS; j++ { 
		go handleRequestsSec(jobs)
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