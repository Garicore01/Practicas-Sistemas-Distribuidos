package mr

import (
	"ms"
	"gf"
	"ra"
)

type Update struct {
	Text string 
}

type Barrier struct {}

func ReceiveMessage(msgs *ms.MessageSystem, File string, reqChan chan ra.Request,
	repChan chan ra.Reply, barChan chan bool){
	for {
		message:= msgs.Receive()
		// Establecemos los posibles casos del mensaje, para enviar por el canal correspondiente el mensaje.
		switch messageType := message.(type) {
		case re.Request:
			reqChan <- messageType
		case re.Reply:
			repChan <- messageType
		case Update:
			gf.EscribirFichero(File, messageType.text)
		case Barrier:
			barChan <- true
		}
	}
}		