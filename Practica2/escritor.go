package main

import (
	"fmt"
	"os"
	"strconv"
	"practica2/gf"
	"practica2/mr"
	"practica2/ra"
	"practica2/ms"
   	"time"  
)

func escritor (msgs *ms.MessageSystem, radb *ra.RASharedDB, File string, text string, me int) {
	
	for {
		
		radb.PreProtocol() // Solicito entrar a la zona critica
		// Zona critica.
		gf.EscribirFichero(File, text)
		// Actualizo todos los ficheros de los demas.
		for i := 1; i <= ra.N; i++ {
			if i != me {
				msgs.Send(i, mr.Update{text})
			}
		}
		radb.PostProtocol() // Realizo el postProtocol
	}
}

func main() {
	meString := os.Args[1]
	fmt.Println("Escritor con PID " + meString)
	me, _ := strconv.Atoi(meString)
	text := "Puta Gari"+meString+"/n"
	File := "fichero_" + meString + ".txt"
	usersFile := "./ms/users.txt"
	gf.CrearFichero(File)
	// Canales necesarios
	reqChan := make(chan ra.Request) // Canal para las solicitudes
	repChan := make(chan ra.Reply) // Canal para las respuestas
	
	messageType := []ms.Message{ra.Request{}, ra.Reply{}, mr.Update{}}
	msgs := ms.New(me, usersFile, messageType)
	
	go mr.ReceiveMessage(&msgs, File, reqChan, repChan)
	time.Sleep(10*time.Second)
	radb := ra.New(&msgs, me, reqChan, repChan,"write")
	go escritor (&msgs, radb, File, text, me)
	fin := make(chan bool)
	<-fin
}