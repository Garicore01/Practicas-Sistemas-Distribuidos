package practica2



import (
	"fmt"
	"os"
	"practica2/ra"
	"practica2/gf"
	"strconv"
)
	



func reader(radb *ra.RASharedDB, myFile string){
	for {
		radb.PreProtocol()
		_ = gf.LeerFichero(myFile)
		radb.PostProtocol()
	}
}


func main()  {
	meString := os.Args[1]
	fmt.Println("Soy el proceso ", meString)
	me, _ := strconv.Atoi(meString)
	myFile := "fichero_" + meString + ".txt"
	usersFile := "./ms/users.txt"
	gf.CrearFichero(myFile)
	reqch := make(chan ra.Request)
	repch := make(chan ra.Reply)
	repch := make(chan ra.Reply)	
	messageTypes := []ms.Message{ra.Request{}, ra.Reply{}, mr.Update{}}
	msgs := ms.New(me, usersFile, messageTypes)
	go mr.ReceiveMessage(&msgs, myFile, reqch, repch)
	radb := ra.New(&msgs, me, reqch, repch,"read")
	go reader(radb, myFile)
}
