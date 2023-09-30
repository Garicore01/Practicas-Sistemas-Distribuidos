package practica2
import (
	"fmt"
	"os"
	"time"
	"os/exec"
	"bufio"
	"strings"
	"strconv"
)



func leerUsers(path string) (arr []string){
	f , _ := os.Open(path)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		arr.append(arr,scanner.Scan())
	}
}





/* 
* PRE: <endpoint> debe ser una @ip valida
* POST: El proceso Worker de la maquina <endpoint>, inicia su ejecución en el puerto que especifica <endpoint>
*/
func encenderProceso(pid int,endpoint string){
	comando := "/usr/bin/ssh"
	// Separo la @IP del puerto
	ip := strings.Split(endpoint, ":")
	credentials := "root@" + ip[0]
	
	//goCommand := "cd /home/a848905/Practicas/Distribuidos/practica1/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a848905/Practicas/Distribuidos/practica1/worker.go " + strconv.Itoa(puerto)
	//goCommand := "cd /home/a849183/Desktop/practica1/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a849183/Desktop/practica1/worker.go " + strconv.Itoa(puerto)
	
	if pid > ra.N/2 {
		goCommand := "cd /home/a849183/Practicas/Distribuidos/practica2/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a849183/Practicas/Distribuidos/practica2/escritor.go " + strconv.Itoa(pid)

	} else {
		goCommand := "cd /home/a849183/Practicas/Distribuidos/practica2/; /usr/local/go/bin/go mod tidy; nohup /usr/local/go/bin/go run /home/a849183/Practicas/Distribuidos/practica2/lector.go " + strconv.Itoa(pid)
	}

    cmd := exec.Command(comando,credentials,goCommand)
	err := cmd.Start()

	if err != nil {
        fmt.Printf("Error al ejecutar el comando: %v\n", err)
        return
    }
}


func main(){
	ruta := "./ms/users.txt"
	dir := leerUsers(ruta)
	for i := ra.N; i > 0; i-- {
		go encenderProceso(i,dir[i])
		if i== ra.N {
			time.Sleep(2*time.Second)
		}
	}
}