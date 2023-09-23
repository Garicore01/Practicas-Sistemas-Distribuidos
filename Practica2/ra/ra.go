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
    OurSeqNum   int
    HigSeqNum   int
    OutRepCnt   int
    ReqCS       bool
    RepDefd     []int
    ms          *ms.MessageSystem
    done        chan bool
    chrep       chan bool
    Mutex       sync.Mutex // mutex para proteger concurrencia sobre las variables
    // TODO: completar
    exclude map[Exclusion]bool
    req chan Request
}

const{
    N = 6
}


func New(me int, usersFile string) (*RASharedDB) {
    messageTypes := []Message{Request, Reply}
    msgs = ms.New(me, usersFile , messageTypes)
    ra := RASharedDB{0, 0, 0, false, []int{}, &msgs,  make(chan bool),  make(chan bool), &sync.Mutex{}}
    // TODO completar
    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(){
    ra.Mutex.Lock()
    ra.ReqCS = true
    OurSeqNum = HigSeqNum + 1 //CUIDADO CON ESTO, NO SE SI NECESARIO
    ra.Mutex.Unlock()
    ra.OutRepCnt =  N-1
    messagePayload := []byte("sample-payload")
    for i := 1; i <= N; i++ {
        if i != me {
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
