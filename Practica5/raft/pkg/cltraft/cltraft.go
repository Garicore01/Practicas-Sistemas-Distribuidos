package cltraft
import(
	"fmt"
	"raft/internal/raft"
	"raft/internal/comun/rpctimeout"
	"strconv"
	"time"
	"raft/internal/comun/check"
)

func main(){
	dns := "raftGA-service.default.svc.cluster.local"
	
	name := "raft"
	puerto := "2900"
	var direcciones []string
	// Inicializo los nombres los tres nodos que van a participar en el cluster.
	for i := 0; i < 3; i++ {
		nodo := name + "-" + strconv.Itoa(i) + "." + dns + ":" + puerto
		direcciones = append(direcciones, nodo)
	}
	var nodos []rpctimeout.HostPort
	// Creo los tres HostPort para poder realizar las peticiones.
	for _, endpoint := range direcciones {
		nodos = append(nodos, rpctimeout.HostPort(endpoint))
	}
	// Espero a que los nodos se levanten.
	time.Sleep(10 * time.Second)
	var reply raft.ResultadoRemoto
	// Creo las operaciones que voy a realizar, posteriormente las someto al lider.
	op1:= raft.Operacion{Operacion: "escribir", Clave: "y", Valor: "9"}
	op2:= raft.Operacion{Operacion: "leer", Clave: "y", Valor: ""}
	fmt.Println("Buscando lider")
	
	// Busco quien es el lider actual.
	// Con un Proxy nos evitariamos tener que hacer esto.
	idLider := -1
	for i := 0; i < 3 && idLider != -1; i++ {
		err := nodos[i].CallTimeout("NodoRaft.ObtenerEstadoNodo",
		raft.Vacio{}, &reply, 10*time.Millisecond)
		if err == nil {
			if reply.Lider {
				idLider = i
			}
		}
	}
	fmt.Println("Lider encontrado, id: ", idLider)
	// Someto las operaciones al lider.
	fmt.Println("Sometiendo operacion1 a lider")
	err := nodos[idLider].CallTimeout("Raft.SometerOperacionRaft", op1, &reply, 5000 * time.Second)
	check.CheckError(err, "Error en llamada RPC al lider ")

	fmt.Println("Sometiendo operacion2 a lider")
	err := nodos[idLider].CallTimeout("Raft.SometerOperacionRaft", op2, &reply, 5000 * time.Second)
	check.CheckError(err, "Error en llamada RPC Para nodo ")




}