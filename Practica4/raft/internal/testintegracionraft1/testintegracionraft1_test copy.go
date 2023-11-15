package testintegracionraft1

import (
	"fmt"
	"raft/internal/comun/check"

	//"log"
	//"crypto/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"raft/internal/comun/rpctimeout"
	"raft/internal/despliegue"
	"raft/internal/raft"
)

const (
	//hosts
	MAQUINA1 = "127.0.0.1"
	MAQUINA2 = "127.0.0.1"
	MAQUINA3 = "127.0.0.1"

	//puertos
	PUERTOREPLICA1 = "29001"
	PUERTOREPLICA2 = "29002"
	PUERTOREPLICA3 = "29003"

	//nodos replicas
	REPLICA1 = MAQUINA1 + ":" + PUERTOREPLICA1
	REPLICA2 = MAQUINA2 + ":" + PUERTOREPLICA2
	REPLICA3 = MAQUINA3 + ":" + PUERTOREPLICA3

	// paquete main de ejecutables relativos a PATH previo
	EXECREPLICA = "cmd/srvraft/main.go"

	// comandos completo a ejecutar en máquinas remota con ssh. Ejemplo :
	// 				cd $HOME/raft; go run cmd/srvraft/main.go 127.0.0.1:29001

	// Ubicar, en esta constante, nombre de fichero de vuestra clave privada local
	// emparejada con la clave pública en authorized_keys de máquinas remotas

	PRIVKEYFILE = "id_rsa_test"
)

// PATH de los ejecutables de modulo golang de servicio Raft
var PATH string = filepath.Join(os.Getenv("HOME"), "Documentos", "Universidad", 
	"SistemasDistribuidos", "Practicas-Sistemas-Distribuidos", "Practica3", 
	"practica3", "CodigoEsqueleto", "raft")

// go run cmd/srvraft/main.go 0 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003
var EXECREPLICACMD string = "cd " + PATH + "; go run " + EXECREPLICA
// TEST primer rango
func TestPrimerasPruebas(t *testing.T) { // (m *testing.M) {
	// <setup code>
	// Crear canal de resultados de ejecuciones ssh en maquinas remotas
	cfg := makeCfgDespliegue(t,
		3,
		[]string{REPLICA1, REPLICA2, REPLICA3},
		[]bool{true, true, true})

	// tear down code
	// eliminar procesos en máquinas remotas
	defer cfg.stop()

	// Run test sequence

	// Test1 : No debería haber ningun primario, si SV no ha recibido aún latidos
	t.Run("T1:soloArranqueYparada",
		func(t *testing.T) { cfg.soloArranqueYparadaTest1(t) })

	// Test2 : No debería haber ningun primario, si SV no ha recibido aún latidos
	t.Run("T2:ElegirPrimerLider",
		func(t *testing.T) { cfg.elegirPrimerLiderTest2(t) })

	// Test3: tenemos el primer primario correcto
	t.Run("T3:FalloAnteriorElegirNuevoLider",
		func(t *testing.T) { cfg.falloAnteriorElegirNuevoLiderTest3(t) })

	// Test4: Tres operaciones comprometidas en configuración estable
	t.Run("T4:tresOperacionesComprometidasEstable",
		func(t *testing.T) { cfg.tresOperacionesComprometidasEstable(t) })
}

// TEST primer rango
func TestAcuerdosConFallos(t *testing.T) { // (m *testing.M) {
	// <setup code>
	// Crear canal de resultados de ejecuciones ssh en maquinas remotas
	cfg := makeCfgDespliegue(t,
		3,
		[]string{REPLICA1, REPLICA2, REPLICA3},
		[]bool{true, true, true})

	// tear down code
	// eliminar procesos en máquinas remotas
	defer cfg.stop()

	// Test5: Se consigue acuerdo a pesar de desconexiones de seguidor
	t.Run("T5:AcuerdoAPesarDeDesconexionesDeSeguidor ",
		func(t *testing.T) { cfg.AcuerdoApesarDeSeguidor(t) })

	t.Run("T6:SinAcuerdoPorFallos ",
		func(t *testing.T) { cfg.SinAcuerdoPorFallos(t) })

	t.Run("T7:SometerConcurrentementeOperaciones ",
		func(t *testing.T) { cfg.SometerConcurrentementeOperaciones(t) })

}

// ---------------------------------------------------------------------
//
// Canal de resultados de ejecución de comandos ssh remotos
type canalResultados chan string

func (cr canalResultados) stop() {
	close(cr)

	// Leer las salidas obtenidos de los comandos ssh ejecutados
	for s := range cr {
		fmt.Println(s)
	}
}

// ---------------------------------------------------------------------
// Operativa en configuracion de despliegue y pruebas asociadas
type configDespliegue struct {
	t           *testing.T
	conectados  []bool
	numReplicas int
	nodosRaft   []rpctimeout.HostPort
	cr          canalResultados
	idLider     int
}

// Crear una configuracion de despliegue
func makeCfgDespliegue(t *testing.T, n int, nodosraft []string,
	conectados []bool) *configDespliegue {
	cfg := &configDespliegue{}
	cfg.t = t
	cfg.conectados = conectados
	cfg.numReplicas = n
	cfg.nodosRaft = rpctimeout.StringArrayToHostPortArray(nodosraft)
	cfg.cr = make(canalResultados, 2000)

	return cfg
}

func (cfg *configDespliegue) stop() {
	time.Sleep(500 * time.Millisecond)
	cfg.cr.stop()
}

// --------------------------------------------------------------------------
// FUNCIONES DE SUBTESTS

// Se pone en marcha una replica ?? - 3 NODOS RAFT
func (cfg *configDespliegue) soloArranqueYparadaTest1(t *testing.T) {
	//t.Skip("SKIPPED soloArranqueYparadaTest1")

	fmt.Println(t.Name(), ".....................")
	cfg.t = t // Actualizar la estructura de datos de tests para errores

	// Poner en marcha replicas en remoto con un tiempo de espera incluido
	cfg.startDistributedProcesses()

	// Comprobar estado replica 0
	cfg.comprobarEstadoRemoto(0, 0, false, -1)

	// Comprobar estado replica 1
	cfg.comprobarEstadoRemoto(1, 0, false, -1)

	// Comprobar estado replica 2
	cfg.comprobarEstadoRemoto(2, 0, false, -1)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses()
	fmt.Println(".............", t.Name(), "Superado")
}

// Primer lider en marcha - 3 NODOS RAFT
func (cfg *configDespliegue) elegirPrimerLiderTest2(t *testing.T) {
	// t.Skip("SKIPPED ElegirPrimerLiderTest2")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	// Se ha elegido lider ?
	fmt.Printf("Probando lider en curso\n")
	IDLider := cfg.pruebaUnLider(3)
	fmt.Printf("El nodo %d es el lider\n", IDLider)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses() // Parametros
	fmt.Println(".............", t.Name(), "Superado")
}

// Fallo de un primer lider y reeleccion de uno nuevo - 3 NODOS RAFT
func (cfg *configDespliegue) falloAnteriorElegirNuevoLiderTest3(t *testing.T) {
	// t.Skip("SKIPPED FalloAnteriorElegirNuevoLiderTest3")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	IDLider := cfg.pruebaUnLider(3)
	fmt.Printf("El nodo %d es el lider\n", IDLider)

	// Desconectar lider
	var reply raft.Vacio
	err := cfg.nodosRaft[IDLider].CallTimeout("NodoRaft.IdleNodo",
		raft.Vacio{}, &reply, 150*time.Millisecond)
	for err != nil {
		err = cfg.nodosRaft[IDLider].CallTimeout("NodoRaft.IdleNodo",
			raft.Vacio{}, &reply, 10*time.Millisecond)
	}

	fmt.Printf("Parado el nodo %d\n", IDLider)

	fmt.Printf("Comprobar nuevo lider\n")
	cfg.pruebaUnLider(3)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}

// 3 operaciones comprometidas con situacion estable y sin fallos - 3 NODOS RAFT
func (cfg *configDespliegue) tresOperacionesComprometidasEstable(t *testing.T) {
	// t.Skip("SKIPPED tresOperacionesComprometidasEstable")
	fmt.Println(t.Name(), ".....................")
	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	IDLider := cfg.pruebaUnLider(3)
	fmt.Printf("El nodo %d es el lider\n", IDLider)

	cfg.comprobarOperacion(0, "escribir", "key_1", "value_1", "ok")
	cfg.comprobarOperacion(1, "escribir", "key_2", "value_2", "ok")
	cfg.comprobarOperacion(2, "leer", "key_1", "", "value_1")
	cfg.stopDistributedProcesses()
	fmt.Println(".............", t.Name(), "Superado")
}

// Se consigue acuerdo a pesar de desconexiones de seguidor -- 3 NODOS RAFT
func (cfg *configDespliegue) AcuerdoApesarDeSeguidor(t *testing.T) {
	// t.Skip("SKIPPED AcuerdoApesarDeSeguidor")
	fmt.Println(t.Name(), ".....................")
	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	IDLider := cfg.pruebaUnLider(3)
	fmt.Printf("El nodo %d es el lider\n", IDLider)

	cfg.comprobarOperacion(0, "escribir", "key_1", "value_1", "ok")
	cfg.stopRandomNode()
	time.Sleep(2 * time.Second)

	cfg.pruebaUnLider(3)

	cfg.comprobarOperacion(1, "escribir", "key_2", "value_2", "ok")
	cfg.comprobarOperacion(2, "escribir", "key_3", "value_3", "ok")
	cfg.comprobarOperacion(3, "leer", "key_2", "", "value_2")

	cfg.wakeIdleNodes()
	cfg.comprobarOperacion(4, "escribir", "key_4", "value_4", "ok")
	cfg.comprobarOperacion(5, "leer", "key_4", "", "value_4")
	cfg.comprobarOperacion(6, "leer", "key_3", "", "value_3")
	cfg.stopDistributedProcesses()
	fmt.Println(".............", t.Name(), "Superado")
}

// NO se consigue acuerdo al desconectarse mayoría de seguidores -- 3 NODOS RAFT
func (cfg *configDespliegue) SinAcuerdoPorFallos(t *testing.T) {
	// t.Skip("SKIPPED SinAcuerdoPorFallos")
	fmt.Println(t.Name(), ".....................")
	cfg.startDistributedProcesses()
	fmt.Printf("Probando líder en curso\n")
	cfg.pruebaUnLider(3)
	// Comprometer una entrada
	cfg.comprobarOperacion(0, "escribir", "key_1", "value_1", "ok")
	// Obtener un lider y, a continuación desconectar 2 de los nodos Raft
	cfg.stopFollowers()
	// Comprobar varios acuerdos con 2 réplicas desconectada
	cfg.someterOperacionFallo("escribir", "key_2", "value_2")
	cfg.someterOperacionFallo("escribir", "key_3", "value_3")
	cfg.someterOperacionFallo("leer", "key_2", "")
	// reconectar los nodos Raft  desconectados y probar varios acuerdos
	cfg.wakeIdleNodes()
	cfg.comprobarOperacion(4, "escribir", "key_4", "value_4", "ok")
	cfg.comprobarOperacion(5, "leer", "key_2", "", "value_2")
	cfg.comprobarOperacion(6, "escribir", "key_5", "value_5", "ok")
	cfg.stopDistributedProcesses()
	fmt.Println(".............", t.Name(), "Superado")

}

// Se somete 5 operaciones de forma concurrente -- 3 NODOS RAFT
func (cfg *configDespliegue) SometerConcurrentementeOperaciones(t *testing.T) {
	t.Skip("SKIPPED SometerConcurrentementeOperaciones")

	fmt.Println(t.Name(), ".....................")
	cfg.startDistributedProcesses()
	// Obtener un lider y, a continuación someter una operacion
	fmt.Printf("Probando líder en curso\n")
	cfg.pruebaUnLider(3)
	cfg.someterOperacion("escribir", "key_1", "value_1")
	// Someter 5  operaciones concurrentes
	go cfg.someterOperacion("escribir", "key_2", "value_2")
	go cfg.someterOperacion("escribir", "key_3", "value_3")
	go cfg.someterOperacion("leer", "key_2", "")
	go cfg.someterOperacion("escribir", "key_4", "value_4")
	go cfg.someterOperacion("leer", "key_2", "")
	time.Sleep(5 * time.Second)
	// Comprobar estados de nodos Raft, sobre todo
	// el avance del mandato en curso e indice de registro de cada uno
	// que debe ser identico entre ellos
	cfg.comprobarEstadoLog(5)
	cfg.stopDistributedProcesses()
	fmt.Println(".............", t.Name(), "Superado")
}

// --------------------------------------------------------------------------
// FUNCIONES DE APOYO

// Comprobar que hay un solo lider
// probar varias veces si se necesitan reelecciones
func (cfg *configDespliegue) pruebaUnLider(numreplicas int) int {
	for iters := 0; iters < 10; iters++ {
		time.Sleep(500 * time.Millisecond)
		mapaLideres := make(map[int][]int)
		for i := 0; i < numreplicas; i++ {
			if cfg.conectados[i] {
				if _, mandato, eslider, _ := cfg.obtenerEstadoRemoto(i); eslider {
					mapaLideres[mandato] = append(mapaLideres[mandato], i)
				}
			}
		}

		ultimoMandatoConLider := -1
		for mandato, lideres := range mapaLideres {
			if len(lideres) > 1 {
				cfg.t.Fatalf("mandato %d tiene %d (>1) lideres",
					mandato, len(lideres))
			}
			if mandato > ultimoMandatoConLider {
				ultimoMandatoConLider = mandato
			}
		}

		if len(mapaLideres) != 0 {
			return mapaLideres[ultimoMandatoConLider][0] // Termina
		}
	}
	cfg.t.Fatalf("un lider esperado, ninguno obtenido")

	return -1 // Termina
}

func (cfg *configDespliegue) obtenerEstadoRemoto(
	indiceNodo int) (int, int, bool, int) {
	var reply raft.EstadoRemoto
	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEstadoNodo",
		raft.Vacio{}, &reply, 10*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC ObtenerEstadoRemoto")
	fmt.Println("nodo: ", reply.IdNodo, " mandato: ", reply.Mandato,
		" eslider: ", reply.EsLider, " idLider: ", reply.IdLider)
	return reply.IdNodo, reply.Mandato, reply.EsLider, reply.IdLider
}

func (cfg *configDespliegue) obtenerEstadoRegistro(indexNode int) (int, int) {
	var reply raft.EstadoRegistro
	err := cfg.nodosRaft[indexNode].CallTimeout("NodoRaft.ObtenerEstadoRegistro",
		raft.Vacio{}, &reply, 20*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC ObtenerEstadoRegistro")
	return reply.IndiceRegistro, reply.Mandato
}

func (cfg *configDespliegue) someterOperacion(op string, key string,
	value string) (int, int, bool, int, string) {
	var reply raft.ResultadoRemoto
	fmt.Printf("Someter operación a %d\n", cfg.idLider)
	err := cfg.nodosRaft[cfg.idLider].CallTimeout("NodoRaft.SometerOperacionRaft",
		raft.TipoOperacion{op, key, value}, &reply, 5000*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC SometerOperacionRaft")
	cfg.idLider = reply.IdLider
	return reply.IndiceRegistro, reply.Mandato, reply.EsLider, reply.IdLider,
		reply.ValorADevolver
}

func (cfg *configDespliegue) someterOperacionFallo(op string, key string,
	value string) {
	var reply raft.ResultadoRemoto
	err := cfg.nodosRaft[cfg.idLider].CallTimeout("NodoRaft.SometerOperacionRaft",
		raft.TipoOperacion{op, key, value}, &reply, 5000*time.Millisecond)
	if err != nil {
		cfg.t.Fatalf(
			"Se ha conseguido acuerdo sin mayoría simple en subtest %s",
			cfg.t.Name())
	} else {
		fmt.Printf("Acuerdo no conseguido\n")
	}
}

// start  gestor de vistas; mapa de replicas y maquinas donde ubicarlos;
// y lista clientes (host:puerto)
func (cfg *configDespliegue) startDistributedProcesses() {
	//cfg.t.Log("Before starting following distributed processes: ", cfg.nodosRaft)

	for i, endPoint := range cfg.nodosRaft {
		despliegue.ExecMutipleHosts(EXECREPLICACMD+
			" "+strconv.Itoa(i)+" "+
			rpctimeout.HostPortArrayToString(cfg.nodosRaft),
			[]string{endPoint.Host()}, cfg.cr, PRIVKEYFILE)

		// dar tiempo para se establezcan las replicas
		//time.Sleep(500 * time.Millisecond)
	}

	// aproximadamente 500 ms para cada arranque por ssh en portatil
	time.Sleep(2500 * time.Millisecond)
}

func (cfg *configDespliegue) stopDistributedProcesses() {
	var reply raft.Vacio

	for _, endPoint := range cfg.nodosRaft {
		err := endPoint.CallTimeout("NodoRaft.ParaNodo",
			raft.Vacio{}, &reply, 10*time.Millisecond)
		check.CheckError(err, "Error en llamada RPC Para nodo")
	}
}

// Se para al líder
func (cfg *configDespliegue) stopLeader(IDLeader int) {
	var reply raft.Vacio
	for i, endPoint := range cfg.nodosRaft {
		if i == IDLeader {
			err := endPoint.CallTimeout("NodoRaft.ParaNodo", raft.Vacio{},
				&reply, 10*time.Millisecond)
			check.CheckError(err, "ParaNodo RPC")
			cfg.conectados[i] = false
		}
	}
}

func (cfg *configDespliegue) stopFollowers() {
	var reply raft.Vacio
	for i, endPoint := range cfg.nodosRaft {
		if i != cfg.idLider {
			err := endPoint.CallTimeout("NodoRaft.ParaNodo", raft.Vacio{},
				&reply, 20*time.Millisecond)
			check.CheckError(err, "Error en llamada RPC ParaNodo")
			cfg.conectados[i] = false
			fmt.Printf("Nodo %d parado\n", i)
		}
	}
}

// Comprobar estado remoto de un nodo con respecto a un estado prefijado
func (cfg *configDespliegue) comprobarEstadoRemoto(idNodoDeseado int,
	mandatoDeseado int, esLiderDeseado bool, IdLiderDeseado int) {
	idNodo, _, _, _ := cfg.obtenerEstadoRemoto(idNodoDeseado)

	//cfg.t.Log("Estado replica 0: ", idNodo, mandato, esLider, idLider, "\n")

	if idNodo != idNodoDeseado {
		cfg.t.Fatalf("Estado incorrecto en replica %d en subtest %s",
			idNodoDeseado, cfg.t.Name())
	}
}

func (cfg *configDespliegue) comprobarOperacion(indiceRegistro int, op string,
	key string, value string, returnedValue string) {

	index, _, _, _, returnValue := cfg.someterOperacion(op, key, value)
	if index != indiceRegistro || returnValue != returnedValue {
		cfg.t.Fatalf("Operacion incorrecta en replica %d en subtest %s",
			indiceRegistro, cfg.t.Name())
	}
}

func (cfg *configDespliegue) comprobarEstadoLog(indiceRegistro int) {
	indices := make([]int, cfg.numReplicas)
	terms := make([]int, cfg.numReplicas)
	for i := range cfg.nodosRaft {
		indices[i], terms[i] = cfg.obtenerEstadoRegistro(i)
	}
}

func (cfg *configDespliegue) stopRandomNode() {
	var reply raft.Vacio
	node := rand.Intn(3)
	err := cfg.nodosRaft[node].CallTimeout("NodoRaft.IdleNode", raft.Vacio{},
		&reply, 20*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC IdleNode")
	cfg.conectados[node] = false
	fmt.Printf("Parado el nodo %d\n", node)
}

func (cfg *configDespliegue) wakeIdleNodes() {
	var reply raft.Vacio
	for i, endPoint := range cfg.nodosRaft {
		if !cfg.conectados[i] {
			err := endPoint.CallTimeout("NodoRaft.WakeNodo", raft.Vacio{},
				&reply, 20*time.Millisecond)
			check.CheckError(err, "Error en llamada RPC WakeNodo")
			cfg.conectados[i] = true
			fmt.Printf("Arrancado el nodo %d\n", i)
		}
	}
}