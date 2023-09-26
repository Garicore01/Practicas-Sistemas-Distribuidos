/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
*/
package ra

import (
    "practica2/ms"
    "sync"
    "github.com/DistributedClocks/GoVector/govec"
    "github.com/DistributedClocks/GoVector/govec/vclock"
    "strconv"
)

type Request struct{
    Clock   int
    Pid     int
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
}

const (
    N = 6
)


func New(me int, usersFile string) (*RASharedDB) {

    Logger :=   govec.InitGoVector(strconv.Itoa(me), strconv.Itoa(me), govec.GetDefaultConfig())

    messageTypes := []ms.Message{Request{},Reply{}} // Defino el contenido de Message
    msgs := ms.New(me, usersFile , messageTypes)
    ra := RASharedDB{0, false, make([]bool, N), &msgs, make(chan bool), make(chan bool), sync.Mutex{}, 
                        make(map[Exclusion] bool),  make(chan Request),  make(chan Reply), Logger}
    
    // TODO completar
    // Posibles casos.
    ra.exclude[Exclusion{"write","write"}]  = true
    ra.exclude[Exclusion{"write","read"}]   = true
    ra.exclude[Exclusion{"read" ,"write"}]  = true
    ra.exclude[Exclusion{"read" ,"read"}]   = false

    



    

    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(){
    ra.Mutex.Lock()
    ra.ReqCS = true
    //OurSeqNum = HigSeqNum + 1 //CUIDADO CON ESTO, NO SE SI NECESARIO
    ra.Mutex.Unlock()
    ra.OutRepCnt =  N-1
    messagePayload := []byte("sample-payload")
    for i := 1; i <= N; i++ {
        if i != ra.ms.Me {
            ra.ms.Send(i, Request{out, ra.ms.Me, ra.op})
        }
    }
    <-ra.chrep
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(){
    ra.ReqCS = false
    for j := 1; j <= N; j++ {
        if ra.RepDefd[j-1] {
            ra.RepDefd[j-1] = false
            ra.ms.Send(j, Reply{})
        }
    }
}

func requestReceived(ra *RASharedDB) {
    for{
        request := <-ra.req
    }
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
