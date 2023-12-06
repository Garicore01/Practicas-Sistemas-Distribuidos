package main

import (
	"fmt"
	"raft/internal/raft"
	"raft/internal/comun/rpctimeout"
	"strconv"
	"time"
	"log"
	"raft/internal/comun/check"
	"os"

)

func main(){
	dns := "raft.default.svc.cluster.local:6000"
	logFile, err1 := os.OpenFile("/cliente/cltraft.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err1 != nil {
		fmt.Printf("Error al abrir el archivo de log: %v\n", err1)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	
	name := "raft"
	var direcciones []string
	// Inicializo los nombres los tres nodos que van a participar en el cluster.
	for i := 0; i < 3; i++ {
		nodo := name + "-" + strconv.Itoa(i) + "." + dns
		direcciones = append(direcciones, nodo)
	}
	var nodos []rpctimeout.HostPort
	// Creo los tres HostPort para poder realizar las peticiones.
	for _, endpoint := range direcciones {
		nodos = append(nodos, rpctimeout.HostPort(endpoint))
	}
	// Espero a que los nodos se levanten.
	var reply raft.ResultadoRemoto
	// Creo las operaciones que voy a realizar, posteriormente las someto al lider.
	op1:= raft.TipoOperacion{Operacion: "escribir", Clave: "y", Valor: "9"}
	op2:= raft.TipoOperacion{Operacion: "leer", Clave: "y", Valor: ""}
	log.Println("Buscando lider")
	

	var err error
	time.Sleep(20 * time.Second)
	err = nodos[0].CallTimeout("NodoRaft.SometerOperacionRaft", op1, &reply, 5000 * time.Millisecond)
	check.CheckError(err, "Error en SometerOperacion 1")
	for reply.IdLider == -1 {
		err = nodos[0].CallTimeout("NodoRaft.SometerOperacionRaft", op1, &reply, 5000 * time.Millisecond)
		check.CheckError(err, "Error en SometerOperacion 1 bucle")
	}
	if !reply.EsLider{
		log.Println("Operacion sometada a ", reply.IdLider)
		err = nodos[reply.IdLider].CallTimeout("NodoRaft.SometerOperacionRaft", op1, &reply, 5000 * time.Millisecond)
		check.CheckError(err, "Error en SometerOperacion 1 fueraBucle")
	}

	log.Println("Resultado de la operacion1: ", reply.ValorADevolver)
	log.Println("Operacion1 acabada")
	log.Println("Operacion2 sometida a ", reply.IdLider)
	err = nodos[reply.IdLider].CallTimeout("NodoRaft.SometerOperacionRaft", op2, &reply, 5000 * time.Millisecond)
	check.CheckError(err, "Error en SometerOperacion 2")
	log.Println("Resultado de la operacion2: ", reply.ValorADevolver)
	log.Println("Operacion2 acabada")
	log.Println("Vamos a empezar a mandar operaciones consecutivas sin parar")
	numero := 0
	var opx raft.TipoOperacion
	for 1 == 1{
		numero++
		opx = raft.TipoOperacion{Operacion: "escribir", Clave: "y", Valor: strconv.Itoa(numero)}
		err = nodos[reply.IdLider].CallTimeout("NodoRaft.SometerOperacionRaft", opx, &reply, 5000 * time.Millisecond)
		check.CheckError(err, "Error en SometerOperacion x")

		log.Println("Resultado de la operacionx: ", reply.ValorADevolver)
	}


}