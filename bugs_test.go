package targetprocess

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var (
	basicAuth = ClientOpt{
		User:     "John",
		Password: "123",
		Url:      "https://restapi.tpondemand.com",
	}

	accessToken = ClientOpt{
		AccessToken: "NjplaXdQeTJDOHVITFBta0QyME85QlhEOWpwTGdPM2V6VjIyelZlZ0NKWG1RPQ==",
		Url:         "https://restapi.tpondemand.com",
	}
)

func TestBugsRequest_Get(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)

	bug, err := cli.Bugs().Get(22)
	assert.NoError(err)
	assert.Equal(int64(22), bug.Id)

	assert.Equal(`Safari hangs when I see About Us page`, bug.Name)
}

func TestBugsRequest_GetByAccessToken(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(accessToken)

	bug, err := cli.Bugs().Get(22)
	assert.NoError(err)
	assert.Equal(int64(22), bug.Id)

	assert.Equal(`Safari hangs when I see About Us page`, bug.Name)
}

/*
Request from example.

curl -u John:123 -g --url "https://restapi.tpondemand.com/api/v1/Bugs?where=(EntityState.IsFinal%20eq%20%27false%27)&include=[Id,Name,Effort,Project]&take=10&orderbydesc=Effort&format=json"
*/
func TestBugsRequest_GetFromExample(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)

	reply, err := cli.Bugs().Search(
		Where(`(EntityState.IsFinal eq 'false')`).
			Inclde(`Id`, `Name`, `Effort`, `Project`).
			Order(OrderByDesc(`Effort`)).
			Take(10))
	assert.NoError(err)
	assert.Equal(`https://restapi.tpondemand.com/api/v1/Bugs/?include=[Id,Name,Effort,Project]&orderByDesc=Effort&format=json&take=10&skip=10`, reply.Next)
	assert.Len(reply.Items, 10)
}

func TestBugsService_Next(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)

	const url = `https://restapi.tpondemand.com/api/v1/Bugs/?include=[Id,Name,Effort,Project]&orderByDesc=Effort&format=json&take=10&skip=10`
	const next = `https://restapi.tpondemand.com/api/v1/Bugs/?include=[Id,Name,Effort,Project]&orderByDesc=Effort&format=json&take=10&skip=20`
	const prev = `https://restapi.tpondemand.com/api/v1/Bugs/?include=[Id,Name,Effort,Project]&orderByDesc=Effort&format=json&take=10&skip=0`

	reply, err := cli.Bugs().Next(url)
	assert.NoError(err)
	assert.Equal(next, reply.Next)
	assert.Equal(prev, reply.Prev)
	assert.Len(reply.Items, 10)
}

func TestBugsService_Create(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)

	bug, err := cli.Bugs().Create(BugDescription{
		Name:    "New Bug",
		Project: ProjectDescription{Id: 2},
	})

	assert.NoError(err)
	assert.Equal(`New Bug`, bug.Name)
	assert.Equal(int64(2), bug.Project.Id)

	_, err = cli.Bugs().Create(BugDescription{
		Name:    "New Bug",
		Project: ProjectDescription{Id: math.MaxInt32},
	})
	assert.Error(err)
}

func TestBugsService_CreateAssigned(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(accessToken)

	bug, err := cli.Bugs().Create(BugDescription{
		Name:    "New Bug",
		Project: ProjectDescription{Id: 2},
		Teams:   Teams{193},
	})

	assert.NoError(err)
	assert.Equal(`New Bug`, bug.Name)
}

func TestBugsService_Delete(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)

	reply, err := cli.Bugs().Search(nil)
	assert.NoError(err)
	assert.NotEmpty(reply.Items)
	if len(reply.Items) != 0 {
		err = cli.Bugs().Delete(reply.Items[0].Id)
		assert.NoError(err)
	}
}
