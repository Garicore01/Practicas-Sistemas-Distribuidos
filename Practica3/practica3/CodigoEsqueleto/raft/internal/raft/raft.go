// Escribir vuestro código de funcionalidad Raft en este fichero
//

package raft

//
// API
// ===
// Este es el API que vuestra implementación debe exportar
//
// nodoRaft = NuevoNodo(...)
//   Crear un nuevo servidor del grupo de elección.
//
// nodoRaft.Para()
//   Solicitar la parado de un servidor
//
// nodo.ObtenerEstado() (yo, mandato, esLider)
//   Solicitar a un nodo de elección por "yo", su mandato en curso,
//   y si piensa que es el msmo el lider
//
// nodoRaft.SometerOperacion(operacion interface()) (indice, mandato, esLider)

// type AplicaOperacion


import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "crypto/rand"
	"sync"
	"time"
	//"net/rpc"
	"math/rand"
	"raft/internal/comun/rpctimeout"
)

// Constantes para establecer los tipos de de nodos que podemos tener.
const (
	FOLLOWER = "follower"
	LEADER = "leader"
	CANDIDATE = "candidate"
)
const (
	// Constante para fijar valor entero no inicializado
	IntNOINICIALIZADO = -1

	//  false deshabilita por completo los logs de depuracion
	// Aseguraros de poner kEnableDebugLogs a false antes de la entrega
	kEnableDebugLogs = true

	// Poner a true para logear a stdout en lugar de a fichero
	kLogToStdout = false

	// Cambiar esto para salida de logs en un directorio diferente
	kLogOutputDir = "./logs_raft/"
)

type TipoOperacion struct {
	Operacion string  // La operaciones posibles son "leer" y "escribir"
	Clave string
	Valor string    // en el caso de la lectura Valor = ""
}


// A medida que el nodo Raft conoce las operaciones de las  entradas de registro
// comprometidas, envía un AplicaOperacion, con cada una de ellas, al canal
// "canalAplicar" (funcion NuevoNodo) de la maquina de estados 
type AplicaOperacion struct {
	Indice int  // en la entrada de registro
	Operacion TipoOperacion
}

// Tipo de dato Go que representa un solo nodo (réplica) de raftRequestVote
//
type NodoRaft struct {
	Mux   sync.Mutex       // Mutex para proteger acceso a estado compartido

	// Host:Port de todos los nodos (réplicas) Raft, en mismo orden
	Nodos []rpctimeout.HostPort
	Yo    int           // indice de este nodos en campo array "nodos"
	IdLider int
	// Utilización opcional de este logger para depuración
	// Cada nodo Raft tiene su propio registro de trazas (logs)
	Logger *log.Logger

	// Vuestros datos aqui.
	
	RequestVote chan bool
 
	AppendEntriesChan chan bool

	FollowerChan chan bool

	CandidateChan chan bool

	LeaderChan chan  bool

	

	// mirar figura 2 para descripción del estado que debe mantenre un nodo Raft

	VotedFor int

	Voted bool

	VotosRecibidos int
	CurrentTerm int

	Rol string

	Votos int
	
}



// Creacion de un nuevo nodo de eleccion
//
// Tabla de <Direccion IP:puerto> de cada nodo incluido a si mismo.
//
// <Direccion IP:puerto> de este nodo esta en nodos[yo]
//
// Todos los arrays nodos[] de los nodos tienen el mismo orden

// canalAplicar es un canal donde, en la practica 5, se recogerán las
// operaciones a aplicar a la máquina de estados. Se puede asumir que
// este canal se consumira de forma continúa.
//
// NuevoNodo() debe devolver resultado rápido, por lo que se deberían
// poner en marcha Gorutinas para trabajos de larga duracion
func NuevoNodo(nodos []rpctimeout.HostPort, yo int, 
						canalAplicarOperacion chan AplicaOperacion) *NodoRaft {
	nr := &NodoRaft{}
	nr.Nodos = nodos
	nr.Yo = yo
	nr.IdLider = -1

	if kEnableDebugLogs {
		nombreNodo := nodos[yo].Host() + "_" + nodos[yo].Port() 
		logPrefix := fmt.Sprintf("%s", nombreNodo)
		
		fmt.Println("LogPrefix: ", logPrefix)

		if kLogToStdout {
			nr.Logger = log.New(os.Stdout, nombreNodo + " -->> ",
								log.Lmicroseconds|log.Lshortfile)
		} else {
			err := os.MkdirAll(kLogOutputDir, os.ModePerm)
			if err != nil {
				panic(err.Error())
			}
			logOutputFile, err := os.OpenFile(fmt.Sprintf("%s/%s.txt",
			  kLogOutputDir, logPrefix), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				panic(err.Error())
			}
			nr.Logger = log.New(logOutputFile, 
						   logPrefix + " -> ", log.Lmicroseconds|log.Lshortfile)
		}
		nr.Logger.Println("logger initialized")
	} else {
		nr.Logger = log.New(ioutil.Discard, "", 0)
	}

	// Añadir codigo de inicialización
	// configuro los datos del nuevo Nodo.
	nr.VotedFor = -1 		// -1 indica que aun no he votado.

	nr.CurrentTerm = -1

	nr.Rol = FOLLOWER

	go nr.raftProtocol()
	
	return nr
}


// Metodo Para() utilizado cuando no se necesita mas al nodo
//
// Quizas interesante desactivar la salida de depuracion
// de este nodo
//
func (nr *NodoRaft) para() {
	go func() {time.Sleep(5 * time.Millisecond); os.Exit(0) } ()
}

// Devuelve "yo", mandato en curso y si este nodo cree ser lider
//
// Primer valor devuelto es el indice de este  nodo Raft el el conjunto de nodos 
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) obtenerEstado() (int, int, bool, int) {
	var yo int = nr.Yo
	var mandato int
	var esLider bool
	var idLider int = nr.IdLider
	

	// Vuestro codigo aqui
	
	mandato = nr.CurrentTerm

	esLider = nr.IdLider == nr.Yo

	idLider = nr.IdLider

	return yo, mandato, esLider, idLider
}

// El servicio que utilice Raft (base de datos clave/valor, por ejemplo)
// Quiere buscar un acuerdo de posicion en registro para siguiente operacion
// solicitada por cliente.

// Si el nodo no es el lider, devolver falso
// Sino, comenzar la operacion de consenso sobre la operacion y devolver en
// cuanto se consiga
// 
// No hay garantia que esta operacion consiga comprometerse en una entrada de
// de registro, dado que el lider puede fallar y la entrada ser reemplazada
// en el futuro.
// Primer valor devuelto es el indice del registro donde se va a colocar 
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) someterOperacion(operacion TipoOperacion) (int, int,
															bool, int, string) {
	indice := -1
	mandato := -1
	EsLider := false
	idLider := -1
	valorADevolver := ""
	

	// Vuestro codigo aqui
	

	return indice, mandato, EsLider, idLider, valorADevolver
}


// -----------------------------------------------------------------------
// LLAMADAS RPC al API
//
// Si no tenemos argumentos o respuesta estructura vacia (tamaño cero)
type Vacio struct{}

func (nr * NodoRaft) ParaNodo(args Vacio, reply *Vacio) error {
	defer nr.para()
	return nil
}

type EstadoParcial struct {
	Mandato	int
	EsLider bool
	IdLider	int
}

type EstadoRemoto struct {
	IdNodo	int
	EstadoParcial
}

func (nr *NodoRaft) ObtenerEstadoNodo(args Vacio, reply *EstadoRemoto) error {
	reply.IdNodo,reply.Mandato,reply.EsLider,reply.IdLider = nr.obtenerEstado()
	return nil
}

type ResultadoRemoto struct {
	ValorADevolver string
	IndiceRegistro int
	EstadoParcial
}

func (nr *NodoRaft) SometerOperacionRaft(operacion TipoOperacion,
												reply *ResultadoRemoto) error {
	reply.IndiceRegistro,reply.Mandato, reply.EsLider,
			reply.IdLider,reply.ValorADevolver = nr.someterOperacion(operacion)
	return nil
}

// -----------------------------------------------------------------------
// LLAMADAS RPC protocolo RAFT
//
// Structura de ejemplo de argumentos de RPC PedirVoto.
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
//
type ArgsPeticionVoto struct {
	Term int
	CandidateId int
	//Falta lo de los Logs
}

// Structura de ejemplo de respuesta de RPC PedirVoto,
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
//
//
type RespuestaPeticionVoto struct {
	Term int
	VoteGranted bool
}


//Método para mandar los RPC
func requestVotes(nr *NodoRaft) {
    var reply RespuestaPeticionVoto
    for i := 0; i < len(nr.Nodos); i++ {
        if i != nr.Yo {
            go nr.enviarPeticionVoto(i, &ArgsPeticionVoto{nr.CurrentTerm, nr.Yo}, &reply)
        }
    }
}


// Metodo para RPC PedirVoto
//
func (nr *NodoRaft) PedirVoto(peticion *ArgsPeticionVoto,
										reply *RespuestaPeticionVoto) error {
	
	if nr.Voted == false { // Si aun no he votado, compruebo los mandatos.
		if peticion.Term < nr.CurrentTerm { 
			// Devuelvo falso
			reply.VoteGranted = false
			reply.Term = nr.CurrentTerm
		} else if peticion.Term > nr.CurrentTerm {
			reply.VoteGranted = true
			nr.Voted = true
			nr.CurrentTerm = peticion.Term
			nr.VotedFor = peticion.CandidateId
		} else { // En un futuro tendremos que modificarlo.
			reply.VoteGranted = true
			nr.Voted = true
			nr.VotedFor = peticion.CandidateId
		}
		if peticion.Term < nr.CurrentTerm { 
			// Devuelvo falso
			reply.VoteGranted = false
			reply.Term = nr.CurrentTerm
		}
	}
	return nil	// Todo funciona correctamente.
}

type ArgAppendEntries struct {
	Term int
	LeaderId int
	// PrevLogIndex int
	// PrevLogTerm int
	//Entries []
	//LeaderCommit int
}

type Results struct {
	Term int
	Success bool
}


// Metodo de tratamiento de llamadas RPC AppendEntries
func (nr *NodoRaft) AppendEntries(args *ArgAppendEntries,
													  results *Results) error {
	if args.Term < nr.CurrentTerm {
		results.Term = nr.CurrentTerm
		results.Success = false
	} else if args.Term > nr.CurrentTerm { // Si el mandato que me mandan es mayor, tengo que actualizar el de todos.
		nr.CurrentTerm = args.Term // Actualizo el CurrentTerm
		results.Success = true
		results.Term = args.Term
		if nr.Rol == LEADER || nr.IdLider == nr.Yo {
			nr.FollowerChan <- true
		} else { 
			nr.AppendEntriesChan <- true
		}
	} else { // Caso en el que los mandatos son iguales. Caso ideal.
		results.Success = true
		results.Term = args.Term
		nr.IdLider = args.LeaderId
		nr.AppendEntriesChan <- true
	}
	return nil
}


// ----- Metodos/Funciones a utilizar como clientes
//
//

// Ejemplo de código enviarPeticionVoto
//
// nodo int -- indice del servidor destino en nr.nodos[]
//
// args *RequestVoteArgs -- argumentos par la llamada RPC
//
// reply *RequestVoteReply -- respuesta RPC
//
// Los tipos de argumentos y respuesta pasados a CallTimeout deben ser
// los mismos que los argumentos declarados en el metodo de tratamiento
// de la llamada (incluido si son punteros
//
// Si en la llamada RPC, la respuesta llega en un intervalo de tiempo,
// la funcion devuelve true, sino devuelve false
//
// la llamada RPC deberia tener un timout adecuado.
//
// Un resultado falso podria ser causado por una replica caida,
// un servidor vivo que no es alcanzable (por problemas de red ?),
// una petición perdida, o una respuesta perdida
//
// Para problemas con funcionamiento de RPC, comprobar que la primera letra
// del nombre  todo los campos de la estructura (y sus subestructuras)
// pasadas como parametros en las llamadas RPC es una mayuscula,
// Y que la estructura de recuperacion de resultado sea un puntero a estructura
// y no la estructura misma.
//
func (nr *NodoRaft) enviarPeticionVoto(nodo int, args *ArgsPeticionVoto,
											reply *RespuestaPeticionVoto) bool {
	err := nr.Nodos[nodo].CallTimeout("NodoRaft.PedirVoto",args,reply,20*time.Millisecond) // Pido el voto a los demas servidores, mandando mi request y un TimeOut de espera.
	if err != nil {
		return false
	} else {
		if reply.VoteGranted {
			nr.VotosRecibidos++;
			if nr.VotosRecibidos > len(nr.Nodos)/2 {
				nr.LeaderChan <- true
			}
		} else if reply.Term > nr.CurrentTerm { // El mandato que me mandan es mayor que el mio.
			nr.CurrentTerm = reply.Term;
			nr.Rol = FOLLOWER;
		}
	}												
	
												
	if nr.CurrentTerm > args.Term { // Si el mandato que me mandan es menor, devuelvo false.
		reply.VoteGranted = false
	} else {
		reply.VoteGranted = true
	}
	
	return true
}

func (nr *NodoRaft) mandarHeartbeat(){
	//Mandar heartbeat a todos los nodos
	var reply Results
	var mandar ArgAppendEntries
	mandar.Term =  nr.CurrentTerm
	mandar.LeaderId = nr.Yo
	
	for i := 0; i < len(nr.Nodos); i++ {
		if i != nr.Yo {
			_ = nr.Nodos[i].CallTimeout("NodoRaft.AppendEntries",&mandar,&reply,20*time.Millisecond)
			if reply.Term > nr.CurrentTerm {
				nr.FollowerChan <- true
				nr.CurrentTerm = reply.Term
				nr.IdLider = -1 // Aun no se sabe quien es el lider
			}
		}
	}
}


func (nr *NodoRaft) raftProtocol(){
	for { 
		switch nr.Rol {
			case FOLLOWER:

				select {			
					case <-nr.AppendEntriesChan: // Me bloqueo hasta recibir un mensaje por el canal.
						nr.Rol = FOLLOWER // Me convierto en FOLLOWER.
					case <-time.After(getRandomTimeout()): // Pasa mi TimeOut.
						nr.IdLider = -1 
						nr.Rol = CANDIDATE // Me convierto en candidato.
				}

			case CANDIDATE:
				nr.CurrentTerm++
				nr.VotedFor = nr.Yo
				nr.Votos = 1
				requestVotes(nr)

				select {
					case <- nr.AppendEntriesChan: //Ha llegado un heartbeat
						nr.Rol = FOLLOWER
					case <- nr.FollowerChan:
						nr.Rol = FOLLOWER

					case <- time.After(getRandomTimeout()): //Timeout, nueva eleccion
						nr.Rol = CANDIDATE
					case <- nr.LeaderChan:
						nr.Rol = LEADER
				}

			case LEADER:
				nr.IdLider = nr.Yo
				nr.mandarHeartbeat() //Mandar heartbeat a todos los nodos
				select {
					case <- nr.FollowerChan:
						nr.Rol = FOLLOWER
					
					case <- time.After(50 * time.Millisecond):  //pasado el timeout mando heartbeat
						nr.Rol = LEADER
				}
			
			
		}
	}
	
}
func getRandomTimeout() time.Duration { // Devuelve un timeout aleatorio entre 150 y 300 ms.
	return time.Duration(150 + rand.Intn(300)) * time.Millisecond
}