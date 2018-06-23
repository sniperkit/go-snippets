package test

type DatabaseClient struct {
	Name	string
}

type RedisClient struct {
	Name	string
}

type DependencyContainer struct {
	Db    *DatabaseClient `inject`
	Redis *RedisClient    `inject`
}

type AdminDatabaseClient DatabaseClient
type AnonDatabaseClient DatabaseClient

type ObjectFactory struct {
}
