### Procolo MDVR para simular A91

* Conexion del DVR al puerto serie: 115200,8,N,1

![](/images/image005.jpg?raw=true)


PASO 1. Establecer handshake  
Para este efecto es necesario analizar la información que llegue al puerto serie
. Lo primero es acumular los datos hasta que se reciba el retorno de carro (0x0d) y el cambio de linea (0x0a).
Cuando ocurra eso se debera de buscar la cadena que contenga los siguientes codigos en hexadecimal: 7e  9c 00 00 00 01 67 34 7e

Se debera de reponder inmediatamente de la misma manera (sin retorno de carro o cambio de linea)
:
```
7e 9c 00 00 00 01 67 34 7e
```

Una vez que ocurra esto el sistema mandara el mensaje "not check auth" lo cual significa que el handshake fue exitoso.

El sistema podra preguntar varias veces durante el tiempo que estemos simulando el A91

Ver el siguiente log donde muestra la comunicación exitosa con el handshake.


```
2019/11/12 11:32:59 RX:
00000000  20 28 65 6e 61 62 6c 65  29 20 45 6e 67 69 6e 65  | (enable) Engine|
00000010  20 4f 4e 20 21 21 21 7e  9c 00 00 00 01 67 34 7e  | ON !!!~.....g4~|
00000020  0d 0a                                             |..|
2019/11/12 11:32:59 TX:
00000000  7e 9c 00 00 00 01 67 34  7e                       |~.....g4~|
2019/11/12 11:32:59 RX:
00000000  20 2d 2d 2d 2d 2d 2d 2d  2d 2d 20 49 4e 20 3c 45  | --------- IN <E|
00000010  76 65 6e 74 5f 43 68 65  63 6b 5f 49 4e 5f 48 4f  |vent_Check_IN_HO|
00000020  44 4c 4f 43 4b 5f 49 4e  3e 20 41 63 74 69 6f 6e  |DLOCK_IN> Action|
00000030  20 2d 2d 2d 2d 2d 2d 2d  2d 0d 0a                 | --------..|
2019/11/12 11:32:59 RX:
00000000  20 3d 3d 48 69 73 69 2d  2d 62 65 67 69 6e 20 6f  | ==Hisi--begin o|
00000010  70 65 6e 0d 0a                                    |pen..|
2019/11/12 11:32:59 RX:
00000000  20 2d 2d 2d 48 69 73 69  20 55 73 61 72 74 20 49  | ---Hisi Usart I|
00000010  6e 69 74 20 6f 6b 0d 0a                           |nit ok..|
2019/11/12 11:32:59 RX:
00000000  20 20 3d 3d 3d 3d 6e 6f  74 20 63 68 65 63 6b 20  |  ====not check |
00000010  61 75 74 68 21 0d 0a                              |auth!..|
```

PASO 2: Enviar los pasajeros que suben o que bajan:

Se debera de enviar (por ejemplo) la siguiente cadena:

7E A5 00 00 00 01 67 38  01 00 00 00 03 00 01 7e

El primer y ultimo digito corresponde a 7E (flag byte) que nos indica que es un comando del A91 el cual estamos simulando

El siguiente byte es el Checksum: El cual es la suma de los campos entre los pipes:

7E A5 00 00 | 00 01 67 38  01 00 00 00 03 00 01 | 7e

la suma de: 01+67+38+01+03+01 = A5

Los dos byte siguientes (00 00) luego del A5 corresponden al serial number (en teoria se debe incrementar por uno, pero si no se hace no afecta en nada.

Luego sigue el Manufacture Number que siempre debe ser: 00 01

Luego sigue el peripheral number que siempre debe ser 67

Luego sigue el el function code que debe ser siempre 38

Luego sigue cual es puerta es: si mandamos (00 es la frontal y con 01 es la trasera)

Luego sigue el estado de la conexion de la camara (00 si es normal y 01 si es abnormal)

Luego sigue el estado de la camara (mandar 00 que es normal)

Luego sigue las Personas que entran (palabra=dos bytes) en este ejemplo son 3.

Luego sigue las Personas que bajan (palabra=dos bytes) en este ejemplo son 1.

Cuando se envia la cadena el sistema respondera con un mensaje de la siguiente forma:

in 3,out 1

Especificando que se recibio con exito la información Y que ya la envio a la plataforma.

Ejemplo del log:

```
2019/11/12 11:34:36 TX:
00000000  7e a5 00 00 00 01 67 38  01 00 00 00 03 00 01 7e  |~.....g8.......~|
2019/11/12 11:34:36 RX:
00000000  20 2b 43 53 51 3a 20 31  31 3b 20 43 61 6e 20 4f  | +CSQ: 11; Can O|
00000010  66 66 20 4c 69 6e 65 3b  20 47 50 53 3a 20 46 69  |ff Line; GPS: Fi|
00000020  78 3b 20 47 50 53 4e 75  6d 62 65 72 3a 20 36 3b  |x; GPSNumber: 6;|
00000030  20 47 50 52 53 3a 20 53  6f 63 6b 65 74 5b 30 5d  | GPRS: Socket[0]|
00000040  20 43 6f 6e 6e 65 63 74  20 4f 6b 3b 20 20 53 6f  | Connect Ok;  So|
00000050  63 6b 65 74 5b 31 5d 20  43 6f 6e 6e 65 63 74 69  |cket[1] Connecti|
00000060  6e 67 3b 20 20 20 20 54  65 6d 70 a3 ba 30 2e 30  |ng;    Temp..0.0|
00000070  0d 0a                                             |..|
2019/11/12 11:34:36 RX:
00000000  20 3d 3d 3d 20 47 50 52  53 20 45 76 65 6e 74 20  | === GPRS Event |
00000010  43 6f 64 65 20 3a 20 33  35 20 3d 3d 3d 0d 0a     |Code : 35 ===..|
2019/11/12 11:34:36 RX:
00000000  0d 0a                                             |..|
2019/11/12 11:34:36 RX:
00000000  2d 2d 2d 2d 2d 2d 2d 2d  20 54 52 5b 30 5d 20 30  |-------- TR[0] 0|
00000010  3a 31 20 2d 2d 2d 2d 2d  2d 2d 2d 0d 0a           |:1 --------..|
2019/11/12 11:34:36 RX:
00000000  0d 0a                                             |..|
2019/11/12 11:34:36 RX:
00000000  20 47 50 52 53 5b 30 5d  20 54 58 20 52 41 4d 0d  | GPRS[0] TX RAM.|
00000010  0a                                                |.|
2019/11/12 11:34:36 RX:
00000000  0d 0a                                             |..|
2019/11/12 11:34:36 RX:
00000000  20 23 23 46 4c 41 53 48  20 52 45 41 44 20 45 4e  | ##FLASH READ EN|
00000010  44 0d 0a                                          |D..|
2019/11/12 11:34:36 RX:
00000000  20 46 4c 41 53 48 20 43  4e 54 3a 30 2c 20 50 41  | FLASH CNT:0, PA|
00000010  43 4b 20 4e 55 4d 42 45  52 3a 31 0d 0a           |CK NUMBER:1..|
2019/11/12 11:34:36 RX:
00000000  20 20 3d 3d 3d 74 65 73  74 20 72 61 6d 5b 30 5d  |  ===test ram[0]|
00000010  20 63 6e 74 3a 31 0d 0a                           | cnt:1..|
2019/11/12 11:34:36 RX:
00000000  20 3d 3d 3d 20 53 65 6e  64 20 47 50 52 53 20 44  | === Send GPRS D|
00000010  61 74 61 5b 31 33 38 5d  3a 20 0d 0a              |ata[138]: ..|
2019/11/12 11:34:36 RX:
00000000  24 24 49 31 33 32 2c 38  36 31 31 30 37 30 33 39  |$$I132,861107039|
00000010  36 31 31 39 37 36 2c 43  43 45 2c 00 00 00 00 01  |611976,CCE,.....|
00000020  00 62 00 16 00 05 05 01  06 06 07 0b 14 00 15 02  |.b..............|
00000030  09 08 01 00 09 29 00 0a  12 00 0b 3d 09 16 00 00  |.....).....=....|
00000040  17 00 00 19 a1 01 1a 6d  04 40 23 00 06 02 fe e6  |.......m.@#.....|
00000050  29 01 03 19 bb 15 fa 04  ab a9 5d 25 0c a1 00 00  |).........]%....|
00000060  00 0d e9 7d 00 00 1c 00  40 00 00 02 0e 0c 03 00  |...}....@.......|
00000070  ea 09 38 25 b9 79 df 05  00 00 49 09 04 00 00 00  |..8%.y....I.....|
00000080  00 00 00 00 00 2a 33 32  0d 0a                    |.....*32..|
2019/11/12 11:34:36 RX:
00000000  0d 0a                                             |..|
2019/11/12 11:34:37 RX:
00000000  20 ba f3 c3 c5 3a 69 6e  20 33 2c 6f 75 74 20 31  | ....:in 3,out 1|
00000010  0d 0a                                             |..|
```
