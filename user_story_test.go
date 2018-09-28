package targetprocess

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserStoryService_Search(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)

	reply, err := cli.Story().Search(
		Where(`(EntityState.IsFinal eq 'false')`).
			Inclde(`Id`, `Name`, `Effort`, `Project`).
			Order(OrderByDesc(`Effort`)).
			Append("Bugs-Count").
			Take(10))
	assert.NoError(err)
	assert.Equal(`https://restapi.tpondemand.com/api/v1/UserStories/?include=[Id,Name,Effort,Project]&append=[Bugs-Count]&where=(EntityState.IsFinal eq 'false')&orderByDesc=Effort&format=json&take=10&skip=10`,
		reply.Next)
	assert.Len(reply.Items, 10)
}

type storyWithBugsCntory struct {
	UseStory
	BugsCount int `json:"Bugs-Count"`
}

func TestUserStoryService_ExecSearch(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)

	var response struct {
		Prev  string
		Next  string
		Items []storyWithBugsCntory
	}

	err := cli.Story().ExecSearch(
		Where(`(EntityState.IsFinal eq 'false')`).
			Inclde(`Id`, `Name`, `Effort`, `Project`).
			Order(OrderByDesc(`Effort`)).
			Append("Bugs-Count").
			Take(10), &response)
	assert.NoError(err)
	assert.Equal(`https://restapi.tpondemand.com/api/v1/UserStories/?include=[Id,Name,Effort,Project]&append=[Bugs-Count]&where=(EntityState.IsFinal eq 'false')&orderByDesc=Effort&format=json&take=10&skip=10`,
		response.Next)
	assert.Len(response.Items, 10)
	assert.True(response.Items[6].BugsCount != 0)
}

func TestUserStoryService_Get(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)
	story, err := cli.story.Get(133)
	assert.NoError(err)
	assert.Equal(`Create templates`, story.Name)

}
