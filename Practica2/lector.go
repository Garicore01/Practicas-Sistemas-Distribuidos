package practica2



import (
	"fmt"
	"os"
	"practica2/ra"

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
}
