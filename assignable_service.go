package targetprocess

const (
	assignableEndpoint = `/api/v1/Assignables`
)

type AssignableService struct {
	client *TPClient
}

type AssignableGetReply struct {
	Prev  string
	Next  string
	Items []Assignable
}

func (s *AssignableService) Search(query *Query) (AssignableGetReply, error) {
	var reply AssignableGetReply
	err := s.ExecSearch(query, &reply)
	return reply, err
}

func (s *AssignableService) ExecSearch(query *Query, receiver interface{}) error {
	if query == nil {
		query = defaultQuery
	}

	return s.client.get(assignableEndpoint, query.values(), receiver)
}
