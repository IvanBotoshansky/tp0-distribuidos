# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |

### Servidor

Se trata de un "echo server", en donde los mensajes recibidos por el cliente se responden inmediatamente y sin alterar. 

Se ejecutan en bucle las siguientes etapas:

1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor retorna al paso 1.


### Cliente
 se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma:
 
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
4. Servidor responde al mensaje.
5. Servidor desconecta al cliente.
6. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

### Ejemplo

Al ejecutar el comando `make docker-compose-up`  y luego  `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```


## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).


### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `


### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).


### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

## Condiciones de Entrega
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios y que aprovechen el esqueleto provisto tanto (o tan poco) como consideren necesario.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deberán existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.
 (hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Se espera que se redacte una sección del README en donde se indique cómo ejecutar cada ejercicio y se detallen los aspectos más importantes de la solución provista, como ser el protocolo de comunicación implementado (Parte 2) y los mecanismos de sincronización utilizados (Parte 3).

Se proveen [pruebas automáticas](https://github.com/7574-sistemas-distribuidos/tp0-tests) de caja negra. Se exige que la resolución de los ejercicios pase tales pruebas, o en su defecto que las discrepancias sean justificadas y discutidas con los docentes antes del día de la entrega. El incumplimiento de las pruebas es condición de desaprobación, pero su cumplimiento no es suficiente para la aprobación. Respetar las entradas de log planteadas en los ejercicios, pues son las que se chequean en cada uno de los tests.

La corrección personal tendrá en cuenta la calidad del código entregado y casos de error posibles, se manifiesten o no durante la ejecución del trabajo práctico. Se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección informados  [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).

## Resolución y ejecución de ejercicios

### Ejercicio 1

Para llevar a cabo el ejercicio, se implementaron los siguientes archivos:

- `generar-compose.sh`: Script de bash que maneja la ejecución del script de python `mi-generador.py`.

- `mi-generador.py`: Script de python que genera el archivo docker-compose con la cantidad de clientes especificada.

Para ejecutar el script `generar-compose.sh` se debe correr el siguiente comando:

```bash
./generar-compose.sh <nombre_archivo_salida> <cantidad_clientes>
```

Por ejemplo:

```bash
./generar-compose.sh docker-compose-dev.yaml 5
```

### Ejercicio 2

Para llevar a cabo el ejercicio, se modificaron los archivos `docker-compose-dev.yaml` y `mi-generador.py` para que los archivos de configuración sean inyectados en los containers y persistidos por fuera de la imagen, utilizando volúmenes.

### Ejercicio 3

Para llevar a cabo el ejercicio, se implementó el script de bash `validar-echo-server.sh` que verifica el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo.

Para ejecutar el script `validar-echo-server.sh` se debe correr el siguiente comando:

```bash
./validar-echo-server.sh
```

Antes de ejecutar el script, para poder probar el caso de éxito, se debe asegurar que el servidor se encuentre corriendo.

Podemos levantar el servidor con 0 clientes con los siguientes comandos:

```bash
./generar-compose.sh docker-compose-dev.yaml 0

make docker-compose-up
```

### Ejercicio 4

Para implementar el cierre graceful ante la señal SIGTERM, tanto en el cliente como en el servidor, se realizaron las siguientes modificaciones:

#### Cliente
- Se implementó un mecanismo de notificación basado en canales con `c.done` que permite interrumpir ciclos de espera.
- Se desarrolló el método `Shutdown()` que realiza el cierre ordenado:
    1. Establece un indicador para salir del bucle principal (cerrando el canal `c.done`).
    2. Cierra y libera la conexión.
- Se registró un manejador de señal que invoca este método cuando se recibe SIGTERM.
- Una gorutina se encarga de detectar la recepción de la señal  y llamar al método `Shutdown()`.

#### Servidor
- Se implementó un manejador de señal usando el módulo `signal`.
- El método `shutdown()` del servidor:
  1. Establece un indicador para salir del bucle principal (`self._running = False`)
  2. Cierra y libera el socket.


Para probar el cierre graceful, se puede enviar la señal SIGTERM a los containers del cliente y del servidor:

```bash
docker stop <nombre_del_container>
```

Por ejemplo, ejecutando:

```bash
./generar-compose.sh docker-compose-dev.yaml 1
make docker-compose-up
make docker-compose-logs
```

Y en otra terminal:

```bash
docker stop client1
```

Podemos ver que los logs del cliente indican que se recibió la señal SIGTERM y se realizó el cierre graceful:

```
client1  | ... INFO     action: graceful_shutdown | result: in_progress | client_id: 1
client1  | ... INFO     action: close_connection | result: in_progress | client_id: 1
client1  | ... INFO     action: close_connection | result: success | client_id: 1
client1  | ... INFO     action: graceful_shutdown | result: success | client_id: 1
client1 exited with code 0
```

Y ejecutando `docker ps -a` verificamos que el container se cerró correctamente.

```
CONTAINER ID    ...    STATUS                               NAMES
4d29bbbd28cb    ...    Exited (0) About a minute ago        client1
```

### Ejercicio 5

Se agregaron variables de entorno en el archivo `docker-compose-dev.yaml` y en el generador `mi-generador.py` para que el cliente pueda recibir los campos que representan la apuesta de una persona.

Las variables son:
- CLI_NOMBRE
- CLI_APELLIDO
- CLI_DOCUMENTO
- CLI_NACIMIENTO
- CLI_NUMERO

#### Protocolo de comunicación

#### Estructura de los mensajes

Todos los mensajes que forman parte del protocolo de comunicación entre el cliente y el servidor tienen la siguiente estructura:

```
[length (4 bytes)][data (N bytes)]
```

Donde:
- `length` es un entero de 4 bytes que indica la longitud del campo `data`. Este campo es muy importante ya que permite al receptor saber cuántos bytes debe leer para obtener la data completa.
- `data` es una secuencia de bytes que representa el string que contiene la información del mensaje.


#### Mensajes que envía el cliente

- `BetMessage`: Mensaje que representa la apuesta de una persona.
    - Campos:
        - `Agency`: Agencia.
        - `FirstName`: Nombre de la persona.
        - `LastName`: Apellido de la persona.
        - `Document`: Documento de la persona.
        - `BirthDate`: Fecha de nacimiento de la persona.
        - `Number`: Número apostado por la persona.
    - Serialización:
        ```
        [length][bytes("Agency,FirstName,LastName,Document,BirthDate,Number")]
        ```

        Por ejemplo, si el mensaje contiene los siguientes campos:
        - `Agency`: "1"
        - `FirstName`: "Santiago Lionel"
        - `LastName`: "Lorca"
        - `Document`: "30904465"
        - `BirthDate`: "1999-03-17"
        - `Number`: "7574"

        La serialización sería:
        ```
        [bytes(48)][bytes("1,Santiago Lionel,Lorca,30904465,1999-03-17,7574")]
        
        =

        [00 00 00 30][31 2c 53 61 6e 74 69 61 67 6f 20 4c 69 6f 6e 65 6c 2c 4c 6f 72 63 61 2c 33 30 39 30 34 34 36 35 2c 31 39 39 39 2d 30 33 2d 31 37 2c 37 35 37 34]

        ```

        Cabe destacar que los caracteres son caracteres ASCII, por lo que cada uno ocupa un byte.

#### Mensajes que envía el servidor
- `ConfirmationMessage`: Mensaje de confirmación.
    - Campos:
        - `Status`: Estado de la confirmación ("success" o "fail").
    - Serialización:
        ```
        [length][bytes("success"/"fail")]
        ```

#### Flujo de mensajes
- Cliente crea un mensaje `BetMessage` y lo serializa.
- Cliente envía el mensaje serializado al servidor.
- Servidor recibe el mensaje, deserializa la información de la apuesta y la almacena. Luego, imprime `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.
- Servidor crea un mensaje `ConfirmationMessage` y lo serializa.
- Servidor responde al cliente con el mensaje serializado.
- Cliente recibe el mensaje y lo deserializa.
- Si la confirmación es de éxito, el cliente imprime `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

Tanto el cliente como el servidor evitan los fenómenos conocidos como *short read* y *short write*.

### Ejercicio 6

#### Cambios en el protocolo de comunicación

#### Mensajes que envía el cliente

Ahora, el cliente envía batchs de apuestas en vez de enviar una única apuesta. Los mensajes que envía el cliente (batchs) son listas de `BetMessage`s. Estos batchs son obtenidos mediante la lectura del archivo .csv que corresponda (`.data/agency-N.csv`).

La serialización de un batch de apuestas es la siguiente:

```
[length (4 bytes)][data (N bytes)]
```

Donde:
- `length` es un entero de 4 bytes que indica la longitud del campo `data`.
- `data` es una secuencia de bytes que representa el string que contiene la información de la/s apuesta/s.

El string que representa la información de varias apuestas tiene el siguiente formato:

```
"<campos de la apuesta 1>;<campos de la apuesta 2>;...;<campos de la apuesta N>"
```

Los campos de cada apuesta tienen el mismo formato que en el ejercicio 5.

#### Flujo de mensajes

- Cliente lee el archivo `.data/agency-N.csv` y obtiene los batches de apuestas.
- Por cada batch:
    - Cliente serializa el batch.
    - Cliente envía el mensaje serializado al servidor.
    - Servidor recibe el mensaje y deserializa la información de las apuestas.
        - Si las apuestas son válidas, las almacena e imprime `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`.
        -  En caso de detectar un error con alguna de las apuestas, imprime `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`
    - Servidor crea un mensaje `ConfirmationMessage` (con estado "success" o "fail") y lo serializa.
    - Servidor responde al cliente con el mensaje serializado.
    - Cliente recibe el mensaje y lo deserializa.
    - Si la confirmación es de éxito, el cliente imprime `action: batch_apuestas_enviado | result: success`.

#### Cálculo de `batch: maxAmount` para no exceder 8KB

Para calcular la cantidad máxima de apuestas por batch para que los mensajes no excedan 8KB, se realizó el siguiente análisis:

Según el protocolo definido, el tamaño máximo de los datos de una apuesta es de 128 bytes:

| Componente | Tamaño máximo (bytes)  |
|------------|------------------------|
| Agencia    | 1                      |
| Nombre     | 50                     |
| Apellido   | 50                     |
| Documento  | 8                      |
| Nacimiento | 10                     |
| Número     | 4                      |
| Comas (5)  | 5                      |
| **Total**  | **128 bytes**          |

Cabe destacar que se realizan validaciones sobre los campos de Nombre y Apellido para asegurar que no excedan los 50 caracteres. Los demás campos son de longitud conocida (por ejemplo, Documento tiene 8 caracteres y Nacimiento tiene 10 por el formato YYYY-MM-DD). Además, se toma como supuesto que el número de agencia es de un dígito y el número apostado es de máximo 4 dígitos.

Entonces, si consideramos la serialización de un batch de N apuestas:

- Datos de apuestas: N × 128 bytes
- Separadores ";" entre apuestas: (N - 1) bytes
- Cabecera de longitud: 4 bytes

Podemos calcular el tamaño total de un batch de N apuestas (en bytes):

```
Tamaño total = 4 + (N × 128) + (N - 1)
             = 4 + 128N + N - 1
             = 3 + 129N
```

Por lo tanto, para que el batch no exceda 8KB (8192 bytes):
```
Tamaño total ≤ 8192

3 + 129N ≤ 8192
129N ≤ 8189
N ≤ 8189 ÷ 129
N ≤ 63.48
```

Considerando que este cálculo asume el peor caso (todas las apuestas con campos de tamaño máximo), un valor seguro para `batch: maxAmount` sería:

```yaml
batch:
  maxAmount: 60
```

Reemplazando N por 60 en la fórmula del tamaño total de un batch de N apuestas:
```
Tamaño total = 3 + 129×60 = 3 + 7740 = 7743 bytes
```

El tamaño no excede 8KB y se tiene un margen de seguridad de `8192 bytes - 7743 bytes = 449 bytes`.

### Ejercicio 7

#### Cambios en el protocolo de comunicación

#### Estructura de los mensajes

Dado que tanto el servidor como el cliente podrían recibir mensajes de diferentes tipos en cierto momento del programa, se agregó un campo a la estructura de los mensajes para indicar de qué tipo de mensaje se trata.

```
[length (4 bytes)][message_type (1 byte)][data (N bytes)]
```

- `length` es un entero de 4 bytes que indica la suma de las longitudes de los campos `message_type` (1 byte) y `data`.
- `message_type` es un byte que indica el tipo de mensaje.
- `data` es una secuencia de bytes que representa el string que contiene la información del mensaje.

Los diferentes tipos de mensaje que existen son los siguientes:
- `BET_BATCH` (1): Mensaje que representa un batch de apuestas.
- `END_NOTIFICATION` (2): Mensaje que indica que el cliente terminó de enviar todos los batches.
- `WINNERS_REQUEST` (3): Mensaje que solicita al servidor la lista de ganadores.
- `CONFIRMATION` (11): Mensaje de confirmación.
- `WINNERS_LIST` (12): Mensaje que contiene la lista de ganadores.

#### Nuevos mensajes

- `EndNotificationMessage`: Mensaje que indica que el cliente terminó de enviar todos los batches.
    - Campos:
        - `Agency`: indica el número de agencia.
    - Serialización:
        ```
        [length][bytes(END_NOTIFICATION)][bytes("Agency")]
        ```
- `WinnersRequestMessage`: Mensaje que solicita al servidor la lista de ganadores.
    - Campos:
        - `Agency`: indica el número de agencia.
    - Serialización:
        ```
        [length][bytes(WINNERS_REQUEST)][bytes("Agency")]
        ```
- `WinnersListMessage`: Mensaje que contiene la lista de ganadores.
    - Campos:
        - `Winners`: lista de ganadores con formato `["documento_ganador1", "documento_ganador2", ..., "documento_ganadorN"]`.
    - Serialización:
        ```
        [length][bytes(WINNERS_LIST)][bytes("documento_ganador1,documento_ganador2,...,documento_ganadorN")]
        ```

#### Flujo de mensajes

- Misma lógica que el ejercicio 6.
- Cuando el Cliente termina de enviar todos los batches, serializa y envía un mensaje `EndNotificationMessage` al servidor.
- Servidor recibe el mensaje, lo deserializa y actualiza el registro de agencias que terminaron de enviar apuestas. Realiza el sorteo solo cuando todas las agencias terminaron de enviar apuestas.
- Cliente serializa y envía un mensaje `WinnersRequestMessage` al servidor.
    - Si se realizó el sorteo, el Servidor recibe el mensaje, lo deserializa y responde con un mensaje `WinnersListMessage` que contiene la lista de ganadores para la agencia en cuestión.
    - Si no se realizó el sorteo, el Servidor responde con un mensaje `ConfirmationMessage` con estado "fail".
- Cliente recibe el mensaje y lo deserializa.
    - Si no recibe un mensaje de tipo `WinnersList`, espera cierto tiempo y vuelve a enviar el mensaje `WinnersRequestMessage`.
    - Si recibe un mensaje de tipo `WinnersList`, lo deserializa e imprime `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.