# Rules System   ![](https://img.shields.io/badge/lang-go-greeng?style=flat)

## Contenido

1. [Acerca del proyecto](#Acerca-del-proyecto)
    * [Tecnologías usadas](#Tecnologías-usadas)
    
2. [Correr el Proyecto](#Correr-el-Proyecto)
   * [Instalación en docker (recomendado)](#Instalación-en-docker-(recomendado))
   * [Instalación en local](#Instalación-en-local)
     
3. [Uso](#Uso)

4. [Pruebas](#Pruebas)
   * [Pruebas Unitarias y de integración](#Pruebas-Unitarias-y-de-integración)
   * [Pruebas de estilo](#Pruebas-de-estilo)   

5. [Más](#Más)
   * [Links de Interés](#Links-de-Interés)
   * [Versionado](#Versionado)

6. [Catalogo de Codigos de Error](#Catalogo-de-Codigos-de-Error)
   * [Familia de Mccs](#Family-Mccs)
     * [Create](#Family-Mccs-Create)
     * [Delete](#Family-Mccs-Delete)
   * [Familia de Compañias](#Family-Companies)
     * [Create](#Family-Companies-Create)
     * [Delete](#Family-Companies-Delete)

## Acerca del proyecto

Microservicio que evalua cargos, mediante la aplicación de reglas personalizadas.

```
Módulos de reglas:

* Reglas de riesgo ( Globales / Por company / Mcc )
  
* Reglas de negocio
* Reglas de company
* Reglas de Policy & Compliance
```

### Tecnologías usadas

* [Golang Rules Engine](https://github.com/nikunjy/rules) - Motor que evalua reglas escritas en [antlr](https://tomassetti.me/antlr-mega-tutorial/)
* [Validador](https://github.com/go-playground/validator) - Validador de campos
* [MongoDB](https://www.mongodb.com/) - Validador de campos
* [Echo](https://echo.labstack.com/) - Framework web
* [Go-modules](https://blog.golang.org/using-go-modles) - Manejador de dependencias
* [Golangci-lint](https://github.com/golangci/golangci-lint) - Linter
* Golang version is 1.17

## Correr el Proyecto

### Instalación en docker (recomendado)

**Prerrequisitos:**
* Docker. [instalación](https://docs.docker.com/get-docker/)
* Docker Compose. [instalación](https://docs.docker.com/compose/install/)
* Variables de entorno
  * ```CONEKTA_GP_USER``` usuario de github 
  * ```CONEKTA_GP_TOKEN``` [token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token) de github
    
**Ejecutar:**
```shell
make docker-compose
```

### Instalación en local

Instalar dependencias necesarias, usar uno de los dos comandos
```shell
make download
```

Correr base de datos en local, usando docker
```shell
docker run -d  --name mongo-on-docker -p 27017:27017 mongo
```

Correr el proyecto mediante
```shell
go run ./cmd/httpserver/main/main.go
```

## Uso

Importar la [documentación](doc/postman) en un cliente rest, de preferencia usar postman.

Usar el cliente rest y listo.

## Pruebas

### Pruebas Unitarias y de integración

Para correr las pruebas unitarias en nuestra máquina podemos ejecutar
el siguiente comando
```shell
make test
```

Para las pruebas de integración tienes varias alternativas.

* Ejecutar las pruebas dentro de un container de docker
```shell
cd config/ci
make test-run
```

* correr las pruebas en tu IDE favorito
Para esto deberás disponibilizar por cuenta propia un servicio de mongo.
la variable de ambiente a configurar es `MONGODB_DATABASE`

### Pruebas de estilo

Para verificación de estilo de codificación contamos con una serie de comandos.
```shell
make lint-install
make lint
```
```shell
make code-format-check
```

## Más

### Links de Interés

- [Lineamientos de desarrollo y Pull Request](https://docs.google.com/document/d/1m4H5XOY3wmSIG2YaqHv3l3Fog96ck7pEawGl7yxOpoc/edit?usp=sharing)

### Versionado

Usamos [SemVer](http://semver.org/) para el versionado. Para todas las versiones disponibles, mira
los [tags en este repositorio](https://github.com/conekta/risk-rules/tags).

### Migrations

We use Liquibase as control version for Mongo
Database. [liquibase-mongo](https://github.com/liquibase/liquibase-mongodb)

Run Migrations

`liquibase update --changeLogFile=migrations/changelog/changelog.xml --log-level debug --url=mongodb://localhost:27018/rules`

##Catalogo-de-Codigos-de-Error
###Family-Mccs
####Create
001 - El nombre de la Familia de Mccs está duplicado.

002 - Algún mcc esta duplicado dentro de otra Familia de Mccs existente.
####Delete
003 - La Familia de Mccs esta asocida a una Regla.

###Family-Companies
####Create
004 - El nombre de la Familia de Companies está duplicado.
####Delete
005 - La Familia de Companies esta asocida a una Regla.