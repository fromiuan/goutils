package snowflake

const defaultNodeServer = 1

var (
	defaultNode *Node
)

func init() {
	var err error
	defaultNode, err = NewNode(defaultNodeServer)
	if err != nil {
		panic("new node error")
	}
}

type SnowFlake struct {
	nodeServer int64
	n          *Node
}

func NewSnowFlake(nodeServer int64) *SnowFlake {
	nodeGenerate, err := NewNode(nodeServer)
	if err != nil {
		panic("new node error")
	}
	return &SnowFlake{
		nodeServer: nodeServer,
		n:          nodeGenerate,
	}
}

func (s *SnowFlake) SetNodeServer(nodeServer int64) error {
	nodeGenerate, err := NewNode(nodeServer)
	if err != nil {
		return err
	}
	s.n = nodeGenerate
	s.nodeServer = nodeServer
	return nil
}

func (s *SnowFlake) GetNode() int64 {
	return s.nodeServer
}

func (s *SnowFlake) GetID() ID {
	return s.n.Generate()
}

func (s *SnowFlake) String() string {
	return s.n.Generate().String()
}

func (s *SnowFlake) Int64() int64 {
	return s.n.Generate().Int64()
}

func GetID() ID {
	return defaultNode.Generate()
}

func GetString() string {
	return defaultNode.Generate().String()
}

func GetInt64() int64 {
	return defaultNode.Generate().Int64()
}
