package mr

import (
	"practica2/ms"
	"practica2/gf"
	"practica2/ra"
)

type Update struct {
	Text string 
}

type Barrier struct {}

func ReceiveMessage(msgs *ms.MessageSystem, File string, reqChan chan ra.Request,
	repChan chan ra.Reply){
	for {
		message:= msgs.Receive()
		// Establecemos los posibles casos del mensaje, para enviar por el canal correspondiente el mensaje.
		switch messageType := message.(type) {
		case ra.Request:
			reqChan <- messageType
		case ra.Reply:
			repChan <- messageType
		case Update:
			gf.EscribirFichero(File, messageType.Text)
		}
	}
}		