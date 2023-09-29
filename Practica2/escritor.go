package main

import (
	"fmt"
	"os"
	"strconv"
	"gf"
	"mr"
	"ra"
	"ms"
)

func escritor (msgs *ms.MessageSystem, radb *ra.RASharedDB, File string text string, me int, wg *sync.WaitGroup) {
	for {
		radb.PreProtocol() // Solicito entrar a la zona critica
		// Zona critica.
		gf.EscribirFichero(File, text)
		// Actualizo todos los ficheros de los demas.
		for i := 1; i <= ra.N; i++ {
			if != me {
				msgs.Send(i, mr.Update{text})
			}
		}
		radb.PostProtocol() // Realizo el postProtocol
	}
}

func main {
	me := os.Args[1]
	fnt.Println("Escritor con PID " + me)
	me, _ := strconv.Atoi(me)
	text := "hola"
	File := "fichero_" + me + ".txt"
	usersFile := "./ms/users.txt"
	gf.CrearFichero(File)
	// Canales necesarios
	reqChan := make(chan ra.Request) // Canal para las solicitudes
	repChan := make(chan ra.Reply) // Canal para las respuestas
	
	messageType := []ms.Message{ra.Request{}, ra.Reply{}, mr.Update{}, mr.Barrier{}}
	msgs := ms.New(me, usersFile, messageType)
	go mr.ReceiveMessage(&msgs, File, reqChan, repChan, barChan)

	radb := ra.New(&msgs, me, "write", reqChan, repChan)

	var wg sync.WaitGroup
	wg.Add(1)
	go writer(&msgs, radb, File, text, me, &wg)
	wg.Wait()
}