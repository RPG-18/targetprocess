package targetprocess

import "fmt"

const (
	storyEndpoint = `/api/v1/UserStories`
)

type UserStoryService struct {
	client *TPClient
}

type StoryGetReply struct {
	Prev  string
	Next  string
	Items []UseStory
}

func (s *UserStoryService) Search(query *Query) (StoryGetReply, error) {
	var reply StoryGetReply
	err := s.ExecSearch(query, &reply)
	return reply, err
}

func (s *UserStoryService) ExecSearch(query *Query, receiver interface{}) error {
	if query == nil {
		query = defaultQuery
	}

	return s.client.get(storyEndpoint, query.values(), receiver)
}

func (s *UserStoryService) Get(id int64) (UseStory, error) {
	var story UseStory
	err := s.ExecGet(id, &story)
	return story, err
}

func (s *UserStoryService) ExecGet(id int64, receiver interface{}) error {
	values := defaultQuery.values()

	endpoint := fmt.Sprintf("%s/%d", storyEndpoint, id)
	return s.client.get(endpoint, values, receiver)
}
