package main
import (
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"strings"
	"strconv"
  "practica2/ra"

)



func leerUsers(path string) (arr []string){
	f , _ := os.Open(path)
	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan(){
		arr = append(arr,scanner.Text())
	}
  return arr
}





/* 
* PRE: <endpoint> debe ser una @ip valida
* POST: El proceso Worker de la maquina <endpoint>, inicia su ejecución en el puerto que especifica <endpoint>
*/
func encenderProceso(pid int,endpoint string,espera chan bool){
	comando := "/usr/bin/ssh"
	// Separo la @IP del puerto
	ip := strings.Split(endpoint, ":")
	credentials := "a849183@" + ip[0]
	
	//goCommand := "cd /home/a848905/Practicas/Distribuidos/practica1/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a848905/Practicas/Distribuidos/practica1/worker.go " + strconv.Itoa(puerto)
	//goCommand := "cd /home/a849183/Desktop/practica1/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a849183/Desktop/practica1/worker.go " + strconv.Itoa(puerto)
	var goCommand string
	if pid > ra.N/2 {
		goCommand = "cd /home/a849183/Documents/Distribuidos/practica2/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a849183/Documents/Distribuidos/practica2/escritor.go " + strconv.Itoa(pid)

	} else {
		goCommand = "cd /home/a849183/Documents/Distribuidos/practica2/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a849183/Documents/Distribuidos/practica2/lector.go " + strconv.Itoa(pid)
	}
  
  cmd := exec.Command(comando,credentials,goCommand)
	err := cmd.Start()
 	 fmt.Printf("Lanzado/n")
	espera<-true
	if err != nil {
        fmt.Printf("Error al ejecutar el comando: %v\n", err)
        return
    }

}


func main(){
	ruta := "./ms/users.txt"
	dir := leerUsers(ruta)
	espera := make(chan bool)
	for i := ra.N; i > 0; i-- {
		go encenderProceso(i,dir[i-1],espera)
	}
	// Espero a que todas las Goroutines "encenderProceso" acaben.
	for j:=0 ; j<ra.N;j++ {
		<-espera
	}
}