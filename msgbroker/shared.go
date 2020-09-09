package msgbroker

//Configuration for configure amqp
type Configuration struct {
	AMQPConnectionURL string
}

//AddTask for add some task
type AddTask struct {
	Number1 int
	Number2 int
}

//Config for configuration amqp
var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@localhost:5672/",
}
