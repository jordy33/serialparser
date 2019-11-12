// go get github.com/jacobsa/go-serial/serial
package main

import "fmt"
import "log"
import (
        "github.com/jacobsa/go-serial/serial"
        "encoding/hex"
        "os/exec"
)
import "os"

func main() {

    ch := make(chan string)
    go func(ch chan string) {
        // disable input buffering
        exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
        // do not display entered characters on the screen
        exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
        var b []byte = make([]byte, 1)
        for {
            os.Stdin.Read(b)
            ch <- string(b)
        }
    }(ch)
        // Set up options.
        options := serial.OpenOptions{
                //PortName: "/dev/ttyS0",
                PortName: "/dev/cu.usbserial-1430",
                BaudRate: 115200,
                DataBits: 8,
                StopBits: 1,
                MinimumReadSize: 4,
        }

        // Open the port.
        port, err := serial.Open(options)
        if err != nil {
                log.Fatalf("serial.Open: %v", err)
        }

        // Make sure to close it later.
        defer port.Close()

        // HandShake Response
        handshake          := []byte{0x7E,0x9c,0x00,0x00,0x00,0x01,0x67,0x34,0x7E}
        responsehandshake  := []byte{0x7E,0x9c,0x00,0x00,0x00,0x01,0x67,0x34,0x7E}

        querydoorstatus    := []byte{0x7E,0x9e,0x00,0x00,0x00,0x01,0x67,0x36,0x7E}
        responsedoorstatus := []byte{0x7E,0xa0,0x00,0x03,0x00,0x01,0x67,0x36,0x01,0x01,0x7E}

        //passengercounterdataFront := []byte{0x7E, 0xA4, 0x00, 0x00, 0x00, 0x01, 0x67, 0x38, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x03, 0x7E,0x0d,0x0a}
        passengercounterdata  := []byte{0x7E, 0xA5, 0x00, 0x00, 0x00, 0x01, 0x67, 0x38, 0x01, 0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x7E}

        cr := []byte{0x0d}
        lf := []byte{0x0a}
        buff := make([]byte, 1)
        counter :=0;
        parse := make([]byte, 0)
        for {
                select {
                    case stdin, _ := <-ch:
                        fmt.Println(":", stdin)
                        fmt.Println("------------------------")
                        _, err = port.Write(passengercounterdata)
                        if err != nil {
                                log.Fatalf("port.Write: %v", err)
                        } else {
                                log.Printf("TX:\n%v", hex.Dump(passengercounterdata))
                        }
                    default:
                        //:fmt.Println("Working..")
                }
                n, err := port.Read(buff)
                if err != nil {
                        log.Fatal(err)
                        break
                }
                if n == 0 {
                        fmt.Println("\nEOF")
                        break
                }
                if n > 0 {
                        parse = append(parse,buff...)

                        if (string(buff[:n])==string(cr)) {

                                _, err := port.Read(buff)
                                if err != nil {
                                        log.Fatal(err)
                                        break
                                }
                                parse = append(parse, buff...)
                                if (string(buff[:n]) == string(lf)) {
                                        log.Printf("RX:\n%v",hex.Dump(parse))
                                        pointer := 0;
                                        for _, x := range parse {
                                                if x == handshake[pointer] {
                                                        pointer = pointer + 1
                                                        if pointer == 9 {
                                                                responsehandshake[3] = byte(counter)
                                                                _, err := port.Write(handshake)

                                                                if err != nil {
                                                                        log.Fatalf("port.Write: %v", err)
                                                                } else {
                                                                        log.Printf("TX:\n%v", hex.Dump(responsehandshake))
                                                                }


                                                                //counter = counter + 1
                                                                break
                                                        }
                                                }
                                        }
                                        pointer = 0;
                                        for _, x := range parse {
                                                if x == querydoorstatus[pointer] {
                                                        pointer = pointer + 1
                                                        if pointer == 9 {
                                                                //responsedoorstatus[3]=byte(counter)

                                                                _, err := port.Write(responsedoorstatus)

                                                                if err != nil {
                                                                        log.Fatalf("port.Write: %v", err)
                                                                } else {
                                                                        log.Printf("TX:\n%v", hex.Dump(responsedoorstatus))
                                                                }




                                                                //counter = counter + 1

                                                                /*
                                                                passengercounterdataFront[3]=byte(counter)
                                                                _, err  = port.Write(passengercounterdataFront)

                                                                if err != nil {
                                                                        log.Fatalf("port.Write: %v", err)
                                                                } else {
                                                                        log.Printf("TX:\n%v", hex.Dump(passengercounterdataFront))
                                                                }
                                                                counter = counter + 1

                                                                passengercounterdataBack[3]=byte(counter)
                                                                _, err  = port.Write(passengercounterdataBack)

                                                                if err != nil {
                                                                        log.Fatalf("port.Write: %v", err)
                                                                } else {
                                                                        log.Printf("TX:\n%v", hex.Dump(passengercounterdataBack))
                                                                }
                                                                counter = counter + 1
                                                                */
                                                                break
                                                        }
                                                }
                                        }
                                        parse = make([]byte, 0)
                                }

                        }
                }

        }
}
