/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
*/
package ra //cambiar a ra

import (
    "os"
    "fmt"
    "practica2/ms"
    "practica2/gf" //quitar luego
    "sync"
    "github.com/DistributedClocks/GoVector/govec"
    "github.com/DistributedClocks/GoVector/govec/vclock"
    "strconv"
)

type Request struct{
    Buffer  []byte
    Pid     int
    Op      string
}

type Reply struct{}

type Exclusion struct{
    Operation1 string
    Operation2 string
}

type RASharedDB struct {
    OutRepCnt   int
    ReqCS       bool
    RepDefd     []bool
    ms          *ms.MessageSystem
    done        chan bool
    chrep       chan bool
    Mutex       sync.Mutex
    // TODO: completar
    exclude     map[Exclusion]bool
    req         chan Request
    rep         chan Reply
    logger      *govec.GoLog
    op          string
    VC          vclock.VClock
}

const (
    N = 6
)
var fichero_pruebas string = "logPruebas"  //quitar luego
func  replyRecieved(ra *RASharedDB){
    for{
        <-ra.rep
        gf.EscribirFichero(fichero_pruebas,"recibo mensaje "+strconv.Itoa(ra.OutRepCnt)+"\n")
        ra.OutRepCnt = ra.OutRepCnt-1
        if ra.OutRepCnt == 0{
            ra.chrep <- true
        }
    }

}

func New(msgs *ms.MessageSystem,me int, req chan Request, rep chan Reply, op string) (*RASharedDB) {

    logger :=   govec.InitGoVector(strconv.Itoa(me), strconv.Itoa(me), govec.GetDefaultConfig())

   fichero_pruebas = fichero_pruebas+strconv.Itoa(me)+".txt" //quitar luego
   gf.CrearFichero(fichero_pruebas) //quitar luego
    ra := RASharedDB{0, false, make([]bool, N), msgs, make(chan bool), make(chan bool), sync.Mutex{}, 
                        make(map[Exclusion] bool),  req,  rep, logger,op, vclock.New()}
    
    // TODO completar
    // Posibles casos.
    ra.exclude[Exclusion{"write","write"}]  = true
    ra.exclude[Exclusion{"write","read"}]   = true
    ra.exclude[Exclusion{"read" ,"write"}]  = true
    ra.exclude[Exclusion{"read" ,"read"}]   = false

    go requestReceived(&ra)
    go replyRecieved(&ra)

    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(){
    // Parte que tiene que ser atomica, para ello utilizo el Mutex
    ra.Mutex.Lock()

    ra.ReqCS = true
    ra.OutRepCnt =  N-1
    
    // Guardo copia de mi reloj actual
    currentVC := ra.logger.GetCurrentVC().Copy()
    currentVC.Tick(strconv.Itoa(ra.ms.Me))
    mensaje := currentVC.Bytes()
	encodedVCPayload := ra.logger.PrepareSend("Send request to "+ra.op, mensaje, govec.GetDefaultLogOptions())
    // Me guardo el reloj actual
    ra.VC = ra.logger.GetCurrentVC().Copy()

    ra.Mutex.Unlock()
    for i := 1; i <= N; i++ {
        if i != ra.ms.Me {
            ra.ms.Send(i, Request{encodedVCPayload, ra.ms.Me, ra.op})
        }
    }
    fmt.Print("Salgo del Preprotocol, con el reloj: ")
    ra.logger.GetCurrentVC().PrintVC()
    <-ra.chrep
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(){
    ra.Mutex.Lock()
    ra.ReqCS = false
    ra.Mutex.Unlock()
    gf.EscribirFichero(fichero_pruebas,"Salgo de la SC\n") //quitar luego
    for j := 1; j <= N; j++ {
        if ra.RepDefd[j-1] {
            ra.RepDefd[j-1] = false
            ra.ms.Send(j, Reply{})
        }
    }

}

func requestReceived(ra *RASharedDB) {
    for { 
        request := <-ra.req
        var defer_it bool
        var message []byte

        ra.logger.UnpackReceive("Recibir request", request.Buffer, &message,govec.GetDefaultLogOptions()) //Introducimos en el logger.
        
        reqClock,err := vclock.FromBytes(message) //Obtenemos el reloj del logger.
	checkError(err)
        ra.Mutex.Lock()
        defer_it = ra.ReqCS && HappensBefore(ra.VC, reqClock, ra.ms.Me, request.Pid) && ra.exclude[Exclusion{ra.op,request.Op}]
        ra.Mutex.Unlock()
        
        if defer_it {
            ra.RepDefd[request.Pid-1] = true 
        } else {
            ra.ms.Send(request.Pid,Reply{})
        }
    }
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}

/*
* PRE : <i> y <j> son los indices de los dos procesos de los cuales vamos a comparas sus relojes vectoriales. <ra> es 
* POST: Devuelve true si el reloj de <v1> > que el reloj de <v2> o si reloj <v1> == reloj <v2> e <i> > <j>, sino devuelve false.
*/
func HappensBefore(v1 vclock.VClock, v2 vclock.VClock ,i int, j int ) (bool){
    if v1.Compare(v2, vclock.Descendant) {
        return true
    } else if v1.Compare(v2, vclock.Concurrent) {
        return i < j
    } else {
        return false
    }
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}


