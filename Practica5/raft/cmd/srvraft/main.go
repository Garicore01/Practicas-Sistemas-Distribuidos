package main

import (

	//"fmt"
	"net"
	"net/rpc"
	"os"
	"raft/internal/raft"
	"raft/internal/comun/rpctimeout"
	"raft/internal/comun/check"
	"strconv"
	"strings"

)


func main() {
	
	dns := "raftGA-service.default.svc.cluster.local:6000"
	nombre := strings.Split(os.Args[1], "-")[0]
	me,_ := strconv.Atoi(strings.Split(os.Args[1], "-")[1])
	



	var direcciones []string
	// Inicializo los nombres los tres nodos que van a participar en el cluster.
	for i := 0; i < 3; i++ {
		nodo := nombre + "-" + strconv.Itoa(i) + "." + dns
		direcciones = append(direcciones, nodo)
	}
	var nodos []rpctimeout.HostPort
	// Creo los tres HostPort para poder realizar las peticiones.
	for _, endpoint := range direcciones {
		nodos = append(nodos, rpctimeout.HostPort(endpoint))
	}


	
	// obtener entero de indice de este nodo
	datos:= make(map[string]string)


	canalAplicarOperacion := make(chan raft.AplicaOperacion, 1000)
    canalRes := make(chan raft.AplicaOperacion, 1000)
	


	// Parte Servidor
	nr := raft.NuevoNodo(nodos, me, canalAplicarOperacion,canalRes)
	
	rpc.Register(nr)
	
	go aplicarOperacion(nr, datos, canalAplicarOperacion, canalRes)

	//fmt.Println("Replica escucha en :", me, " de ", os.Args[2:])
	
	l, err := net.Listen("tcp", direcciones[me])
	check.CheckError(err, "Main listen error:")

	rpc.Accept(l)
}

func aplicarOperacion(nr *raft.NodoRaft, almacen map[string]string, canalOper chan raft.AplicaOperacion, canalRes chan raft.AplicaOperacion) {
    for {
        operacion := <-canalOper
        nr.Logger.Println("Main ha recibido la operación: ", operacion)
        if operacion.Operacion.Operacion == "leer" {
            operacion.Operacion.Valor = almacen[operacion.Operacion.Clave]
        } else if operacion.Operacion.Operacion == "escribir" {
            almacen[operacion.Operacion.Clave] = operacion.Operacion.Valor
            operacion.Operacion.Valor = "ok"
        }
        nr.Logger.Println("Aplicando operacion: ", operacion)
        if nr.Yo == nr.IdLider {
            canalRes <- operacion
            nr.Logger.Println("Respuesta enviada: ", operacion)
        }
        nr.Logger.Println("Esperando siguiente operacion")
    }
}

