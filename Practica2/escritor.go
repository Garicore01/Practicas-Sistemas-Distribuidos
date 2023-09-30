package practica2

import (
	"fmt"
	"os"
	"strconv"
	"practica2/gf"
	"practica2/mr"
	"practica2/ra"
	"practica2/ms"
)

func escritor (msgs *ms.MessageSystem, radb *ra.RASharedDB, File string, text string, me int) {
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
	me := os.Args[2]
	fnt.Println("Escritor con PID " + me)
	me, _ := strconv.Atoi(me)
	text := "Puta Gari"+me
	File := "fichero_" + me + ".txt"
	usersFile := "./ms/users.txt"
	gf.CrearFichero(File)
	// Canales necesarios
	reqChan := make(chan ra.Request) // Canal para las solicitudes
	repChan := make(chan ra.Reply) // Canal para las respuestas
	
	messageType := []ms.Message{ra.Request{}, ra.Reply{}, mr.Update{}}
	msgs := ms.New(me, usersFile, messageType)
	go mr.ReceiveMessage(&msgs, File, reqChan, repChan)

	radb := ra.New(&msgs, me, reqChan, repChan,"write")
	go writer(&msgs, radb, File, text, me)
}