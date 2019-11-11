### Log DVR serial

Conexion del DVR 
![](/images/image005.jpg?raw=true)


* Acceder al archivo log.txt

* Buscar la cadena que empieza con: 7e 9c

Esta es la parte de log donde se hace el handshake:

```
2019/11/11 10:19:51 RX:
00000000  20 28 65 6e 61 62 6c 65  29 20 45 6e 67 69 6e 65  | (enable) Engine|
00000010  20 4f 4e 20 21 21 21 7e  9c 00 00 00 01 67 34 7e  | ON !!!~.....g4~|
00000020  0d 0a                                             |..|


2019/11/11 10:19:51 TX:
00000000  7e 9c 00 00 00 01 67 34  7e 0d 0a                 |~.....g4~..|


2019/11/11 10:19:52 RX:
00000000  20 20 3d 3d 3d 3d 6e 6f  74 20 63 68 65 63 6b 20  |  ====not check |
00000010  61 75 74 68 21 0d 0a                              |auth!..|
```

* La palabra not check auth quiere decir que ya hay un handshake

* Si se vuelve a pedir el handshake hay que responder igual

* Pero en mi experiencia una vez realizado ya no vuelve pedirlo


