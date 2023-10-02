package mr

import (
	"practica2/ms"
	"practica2/gf"
	"practica2/ra"
)

type Update struct {
	Text string 
}



func ReceiveMessage(msgs *ms.MessageSystem, File string, reqChan chan ra.Request,repChan chan ra.Reply){
	for {
		message:= msgs.Receive()
		// Establecemos los posibles casos del mensaje, para enviar por el canal correspondiente el mensaje.
		switch messageType := message.(type) {
		case ra.Request:
			reqChan <- messageType // Envio una request al proceso "requestReceived"
		case ra.Reply:
			repChan <- messageType // Envio un reply al proceso "replyRecieved"
		case Update:
			gf.EscribirFichero(File, messageType.Text)
		}
	}
}		