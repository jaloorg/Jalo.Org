package repository

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoRepository Repositorio Singleton
type MongoRepository struct {
	db *mongo.Client
}

/* mongoCN objeto de conexión de la BD */
var instancia *MongoRepository
var once sync.Once

// Nombre de la base de datos
//TODO: Implementar extraccion desde archivo de configuracion
const dataBase string = "AyudaapDb"

//GetInstance Obtiene acceso a una instancia de conexion hacia MongoDb
func GetInstance() *MongoRepository {
	once.Do(conectarBD)
	return instancia
}

// ConectarDB inicia una conexión de hacia la BD
//TODO: Implementar consulta desde archivo de configuracion
func conectarBD() {
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USER")
	// pass := os.Getenv("DB_PASSWORD")

	//var cadenaConexion = fmt.Sprintf("mongodb+srv://%s:%s@%s:%s", user, pass, host, port)

	//cadenaConexion := fmt.Sprintf("mongodb://%s:%s@%s:%s/", user, pass, host, port)
	cadenaConexion := fmt.Sprintf("mongodb://localhost:27017")
	clienteOpts := options.Client().ApplyURI(cadenaConexion)
	cliente, err := mongo.Connect(context.TODO(), clienteOpts)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	err = cliente.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	instancia = &MongoRepository{db: cliente}
}

//ChequeoConnection Revisa si la conexion sigue activa
func (m MongoRepository) ChequeoConnection() int {

	if instancia != nil {
		instancia = GetInstance()
	}

	err := instancia.db.Ping(context.TODO(), nil)
	if err != nil {
		return 0
	}
	return 1
}

// GetCollection Obtienene la colecion y el contexto de trabaj
// `collection` Nombre de la conexion
func (m MongoRepository) GetCollection(collection string) (*mongo.Collection, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	if instancia != nil {
		instancia = GetInstance()
	}

	db := instancia.db.Database(dataBase)
	col := db.Collection(collection)

	return col, ctx, cancel
}
