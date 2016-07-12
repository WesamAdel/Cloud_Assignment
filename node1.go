// TODO: In this file you need to build a graph of nodes
// and each node contains a set of files and print all files
// you get, graph and shortest path for specified file.

package main

import (
	"./core/student"
	"fmt"
	"time"
	"encoding/json"
	"strconv"
	"os"
)

// Message struct.
type Message struct {
	From     int
	To       int
	UserName string
	Content  string
}

// Global Declarations.
var masterAddr string = "10.0.0.4:46321"
var connectedNodes = []int{2}
var fileList = []string{
	"1939620_437577509_n.jpg",
	"5978610_937577509_n.jpg",
	"6436120_737577509_n.jpg"}

var student1 *student.Student

// TODO: Change this to your current password.
var studentPassword string = "CS7hYr"

// Implementing ReceiveHandler for student package.
type RcvHandler struct{}

// Handle a message received.
func (rcvHand *RcvHandler) ReceiveHandler(from int, to int, username string,
	content string) {	
	// DONOT CHANGE PARAMENTERS OR FUNCTION HEADER.
	// TODO: Implement handling a message received.
	//if (username == "user16") && (content != "") {
		go updateGraph(content)
		go sendAgain(from,content)
		//student := new(student.Student)

	//}
}

func sendAgain(from int, content string){
	
	for _ , i:= range connectedNodes{
		if 	i != from {
			student1.SendMsg(i,content)
		}
	}
}

var adjacencyList map[int][]int
var nodeImgs map[int][]string

func updateGraph(content string) {
	var strArr []string
	_ = json.Unmarshal([]byte(content), &strArr)
	//fmt.Printf("%v\n",strArr)
	if len(strArr) > 0 {
		a,_ := strconv.Atoi(strArr[0])
		if _, ok := nodeImgs[a]; !ok {
			nodeImgs[a] = append(nodeImgs[a],strArr[1],strArr[2],strArr[3])
			s,_ := strconv.Atoi(strArr[4])
			for i:= 0 ; i < s ; i++{
				b , _ := strconv.Atoi(strArr[5+i])
				adjacencyList[a] = append(adjacencyList[a],b)
				//test
				//fmt.Printf("%v\n",adjacencyList[a])
			}
		}
	}
}

// BFS ------------------
var parent []int

// queue functions 
func push (q []int, x int) []int{
	q = append(q, x)
	return q
}

func top (q []int) int{
	return q[0]
}

func pop (q []int) []int{
	q = q[1:]
	return q
}

func empty (q []int) bool{
	if len(q) == 0 {
		return true
	}
	return false
}

// main function to return the shortest path from 1 to the target node which contains the file 
func shortest_path (file string) []int{
	tar := get_tar(file)
	bfs(tar)
	path := get_path(1,tar)
	return path
}

// returns the target node that contains the file 
func get_tar (file string) int{
	for k,_ := range nodeImgs {
		for _,f := range nodeImgs[k]{
			if f == file{
				return k
			}
		}
	}
	
	return -1
}

//perfoms BFS on the graph
func bfs (tar int){
		var Q []int
		var vis []bool
		
		// initialize visited & parent arrays
		for i := 0; i <= len(adjacencyList); i++ {
			vis = append(vis,false)
			parent = append(parent,i)
		}
		
		Q = make([]int,1)
		Q [0] = 1
		vis[1] = true
		
		for !empty(Q) {
			u := top(Q)
			Q = pop(Q)
			
			// loop on node's adjList
			for i := 0; i< len(adjacencyList[u]); i++{
				p := adjacencyList[u][i]
				//fmt.Printf("p:%d\n",p)
				if vis[p] {
					continue 
				}
				vis[p] = true
				parent[p] = u
				Q = push(Q,p)
			}	
		}
		
}

// get the shortest path from node 1 to the target node
func get_path (s int, d int) []int{
	var path []int 
	
	for d != s {
		path = append(path,d)
		d = parent [d]
	}
	
	path = append(path,s)
	
	return path
}


func main() {
	// Setup connection to master of current node.
	student1 = new(student.Student)
	error := student1.Connect(masterAddr, studentPassword)
	if error != nil {
		fmt.Println("Failed to connect to master node:", error)
		return
	}

	// Link implementation of ReceiveHandler to student.
	rcv := new(RcvHandler)
	go student1.Receive(rcv)
	// End of Setup.
	
	adjacencyList = make(map[int][]int)
	nodeImgs = make(map[int][]string)
	
	N := 10
	time.Sleep(time.Second * time.Duration(N))
	
	// TODO: Broadcast your files to neighbours.
	sentArr := []string{"1"}
	sentArr = append(sentArr,fileList...)
	sentArr = append(sentArr,strconv.Itoa(len(connectedNodes)))
	for _ , value := range connectedNodes{
		sentArr = append(sentArr,strconv.Itoa(value))
	}
	
    msgStr, _ := json.Marshal(sentArr)
    
	//fmt.Printf("%s\n",string(msgStr))
	
	for _ , i:= range connectedNodes{
		student1.SendMsg(i,string(msgStr))
	}
	
	updateGraph(string(msgStr))
	
	// TODO: It's expected to converge after N second
	// To be able to print a stable graph and shortest
	// path for file.
	N = 40
	time.Sleep(time.Second * time.Duration(N))
	fmt.Printf("%v\n",adjacencyList)
	// TODO: Print results in output file.
	file, _ := os.OpenFile("output", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0777)
	defer file.Close()
	//fmt.Printf("%v\n",shortest_path("5834591_818124870_n.jpg"))
	_ , _ = file.WriteString("Graph\n")
	for node , val := range adjacencyList{
		str := "Node " + strconv.Itoa(node)+" :"
		for _ , adjNode := range val {
			str = str + " " + strconv.Itoa(adjNode)
		}
		_ , _ = file.WriteString(str)
		_ , _ = file.WriteString("\n")
	}
		
	_ , _ = file.WriteString("\n")
	//test
	//fmt.Printf("TEST:: %v",)	
	
	v := shortest_path("5834591_818124870_n.jpg")
	
	//fmt.Printf("PATH ::%v\n",v)
	
	_ , _ = file.WriteString("shortest path from source to destination is :\n")
	str := ""
	i := len(v) -1
	
	for indx := i ; indx > -1; indx = indx-1{
		str = str + strconv.Itoa(v[indx]) + " "
	}
	_ , _ = file.WriteString(str)
}
