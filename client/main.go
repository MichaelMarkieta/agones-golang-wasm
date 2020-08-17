package main

func main() {
	// THIS WORKS, BUT NOT IN WASM, BECAUSE UDP IS NOT SUPPORTED BY BROWSERS
	//log.Printf("Connecting to game server...")
	//addr, err := net.ResolveUDPAddr("udp4", "35.203.44.130:7777")
	//conn, err := net.DialUDP("udp", nil, addr)
	//if err != nil {
	//	log.Fatalf("Error connecting to game server: %s", err)
	//}
	//log.Printf("... connected")
	//_, err2 := conn.Write([]byte("hello\n"))
	//if err2 != nil {
	//	log.Fatalf("Error connecting to game server: %s", err)
	//}
	//buffer := make([]byte, 1024)
	//n, _, err := conn.ReadFromUDP(buffer)
	//log.Printf("%s", string(buffer[0:n]))
	//defer conn.Close()
	//log.Printf("Connection closed.")
}
