package main

import (
	"net/rpc"
	"net"
	"log"
	"net/http"
	"time"
	"sync"
	"fmt"
)

type Args struct {
	A string
	UID string
	Stime time.Time
	Room string
}


type ReadingReply struct {
	ID int
	Content string
}

type RoomList struct {
	roomList map[string]string
	Lock sync.Mutex
}

type History struct {
	history map[string]string
	Lock sync.Mutex
}

type NameID struct {
	nameid map[string]string
	Lock sync.Mutex
}

type RoomCount struct {
	roomcount map[string]int
	Lock sync.Mutex
}

type HistoryTime struct {
	historytime map[string]time.Time
	Lock sync.Mutex
}

type Arith string
/*var History map[string]string*/
/*var RoomList map[string]string*/
/*var NameID map[string]string*/
/*var RoomCount map[string]int*/
/*var HistoryTime map[string]time.Time*/
var roomList RoomList
var history History
var nameID NameID
var roomCount RoomCount
var historyTime HistoryTime

func (rl RoomList)Get(k string) string{
	rl.Lock.Lock()
	defer rl.Lock.Unlock()
	return rl.roomList[k]
}
func (rl RoomList)Set(k,v string) {
	rl.Lock.Lock()
	defer rl.Lock.Unlock()
	rl.roomList[k]=v
}
/*func (rl RoomList)Delete(k string)  {
	rl.Lock.Lock()
	defer rl.Lock.Unlock()
	delete(rl.roomList,k)
}*/

func (hs History)Get(k string) string{
	hs.Lock.Lock()
	defer hs.Lock.Unlock()
	return hs.history[k]
}
func (hs History)Set(k,v string) {
	hs.Lock.Lock()
	defer hs.Lock.Unlock()
	hs.history[k]=v
}

func (nid NameID)Get(k string) string{
	nid.Lock.Lock()
	defer nid.Lock.Unlock()
	return nid.nameid[k]
}
func (nid NameID)Set(k,v string) {
	nid.Lock.Lock()
	defer nid.Lock.Unlock()
	nid.nameid[k]=v
}

func (rc RoomCount)Get(k string) int{
	rc.Lock.Lock()
	defer rc.Lock.Unlock()
	return rc.roomcount[k]
}
func (rc RoomCount)Set(k string,v int) {
	rc.Lock.Lock()
	defer rc.Lock.Unlock()
	rc.roomcount[k]=v
}

func (ht HistoryTime)Get(k string) time.Time{
	ht.Lock.Lock()
	defer ht.Lock.Unlock()
	return ht.historytime[k]
}
func (ht HistoryTime)Set(k string,v time.Time) {
	ht.Lock.Lock()
	defer ht.Lock.Unlock()
	ht.historytime[k]=v
}

func CheckRoom()  {
	var subTime time.Duration
	for true {
		for roomID := range historyTime.historytime{
			subTime =time.Now().Sub(historyTime.Get(roomID))
			if subTime>30*time.Second {
				fmt.Println(subTime.Seconds(),roomID)
				/*RoomCount[roomID]=-1*/
				roomCount.Set(roomID,-1)
				roomList.Set(roomID,"This room is going to be closed#*#")
				time.Sleep(time.Second)
				delete(roomList.roomList,roomID)
				delete(history.history,roomID)
				delete(nameID.nameid,roomID)
				delete(roomCount.roomcount,roomID)
				delete(historyTime.historytime,roomID)
			}
		}
		time.Sleep(1*time.Second)
	}
}

func (t *Arith) Creatroom(args *Args, reply *string) error {
	/*RoomList[args.Room]=args.A*/
	roomList.Set(args.Room,args.A)
	return nil
}
func (t *Arith) Showroom(args *Args, reply *string) error {
	var roomlist string
	for room :=range roomList.roomList{
		roomlist +="["+room+"]"
	}
	*reply=roomlist
	return nil
}
func (t *Arith) ReceiveMessage(args *Args, reply *string) error {
	/*RoomCount[args.Room]++*/
	roomCount.roomcount[args.Room]++
	/*RoomList[args.Room]=args.A*/
	roomList.Set(args.Room,args.A)
	/*History[args.Room]="\n"+History[args.Room]+args.A*/
	var hs string
	hs="\n"+history.Get(args.Room)+args.A+"\n"
	history.Set(args.Room,hs)
	/*NameID[args.Room]= args.UID*/
	nameID.Set(args.Room,args.UID)
	/*HistoryTime[args.Room]=args.Stime*/
	historyTime.Set(args.Room,args.Stime)
	return nil
}
func (t *Arith) Show(args *Args, reply *string) error {

	*reply =history.Get(args.Room)

	return nil
}
func (t *Arith) Reading(args *Args, readingreply *ReadingReply) error {
	if roomCount.Get(args.Room)!=-1 {
		if args.UID!=nameID.Get(args.Room) {
			var message1 string
			message1 = roomList.Get(args.Room)
			readingreply.ID = roomCount.Get(args.Room)
			readingreply.Content = message1
		}
	}else {
		var message1 string
		message1 = roomList.Get(args.Room)
		readingreply.ID = roomCount.Get(args.Room)
		readingreply.Content = message1
	}
	return nil
}

func main()  {
	roomList.roomList=make(map[string]string)
	history.history=make(map[string]string)
	nameID.nameid=make(map[string]string)
	roomCount.roomcount=make(map[string]int)
	historyTime.historytime=make(map[string]time.Time)
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go CheckRoom()
	go http.Serve(l,nil)
	time.Sleep(100*time.Minute)
}

