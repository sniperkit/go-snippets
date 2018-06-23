## Go con Ejemplos

Contenido, herramientas y servidor web de [Go con Ejemplos][1].

### Información General

El sitio de [Go con Ejemplos][1] se construye analizando el código y los
comentarios de los archivos fuente en el folder `examples` y mostrando esta
información en el sitio usando `templates` (plantillas). Los programas que
realizan este proceso de publicación se encuentran en el directorio `tools`.

El proceso de publicación produce un directorio de archivos estáticos
(`public`) perfecto para ser servido por cualquier servidor HTTP moderno.
Además, se incluye un servidor web Go ligero en `server.go`.

### Compilación de ejemplos

Para compilar el sitio:

```console
$ go get github.com/extemporalgenome/slug
$ go get github.com/russross/blackfriday
$ ./tools/build
```

Para compilar constantemente en un ciclo:

```console
$ ./tools/build-loop
```

### Despliegue local

Para lanzar un servidor local que puedes consultar en
[127.0.0.1:8000](http://127.0.0.1:8000).

```console
$ ./tools/server
```

### Instrucciones para colaborar

Para colaborar revisa la [lista de pendientes][2] y escoge una traducción
pendiente.

Una vez que sepas cuál ejercicio te gustaría traducir, sigue estos pasos:

1. Revisa `examples.txt` y busca el nombre del ejercicio en idioma original,
   por ejemplo `Reading Files`.
2. Revisa el directorio `examples/` e identifica el directorio que corresponde
   al nombre del ejercicio. Generalmente es el mismo nombre convertido a
   minúsculas y reemplazando caracteres especiales y espacios con `-`. Por
   ejemplo `Reading Files` se convertiría en `reading-files`.
3. Cambia el nombre del ejercicio en `examples.txt` por el nombre en español,
   en nuestro ejemplo sería `Lectura de Archivos`.
4. Usa `git mv` para cambiar el nombre del directorio al que corresponda de
   acuerdo al nuevo nombre en español, por ejemplo `lectura-de-archivos`. En
   caso de tener un nombre con acentos el directorio deberá usar una letra
   minúscula sin acento. Por ejemplo `Análisis de Formato de Fecha` se
   convertiría en `analisis-de-formato-de-fecha`.
5. Usa `git mv` para cambiar el nombre de los archivos `.go` y `.sh` dentro del
   nuevo directorio a su nuevo nombre en español.
6. Verifica tu traducción corriendo `tools/build`.
7. Haz un [pull request][3] para que aceptemos tu traducción y cerremos el
   pendiente de la lista.

**Verifica** que hayas traducido también las **variables** en el código, así como los
**mensajes de salida** y **comentarios**. 

No olvides que además de traducir comentarios, el código también es importante! 
Mantener el idioma español como estándar en el repositorio ayudará a comprender 
mejor cada ejercicio. 

### License

This work is copyright Mark McGranaghan and licensed under a
[Creative Commons Attribution 3.0 Unported License](http://creativecommons.org/licenses/by/3.0/).

The Go Gopher is copyright [Renée French](http://reneefrench.blogspot.com/) and licensed under a
[Creative Commons Attribution 3.0 Unported License](http://creativecommons.org/licenses/by/3.0/).


### Traducciones

Algunas traducciones hechas por contribuidores:

* [Chino](http://everyx.github.io/gobyexample/) by [everyx](https://github.com/everyx)

### Gracias

Gracias a [Jeremy Ashkenas](https://github.com/jashkenas) por
[Docco](http://jashkenas.github.com/docco/), lo que inspiró este proyecto.

[1]: http://goconejemplos.com
[2]: https://github.com/dabit/gobyexample/issues
[3]: https://help.github.com/articles/creating-a-pull-request
