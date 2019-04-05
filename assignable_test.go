package targetprocess

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssignableServiceService_Search(t *testing.T) {
	assert := assert.New(t)
	cli := NewClient(basicAuth)

	response, err := cli.Assignables().Search(
		Where(`(EntityState.IsFinal eq 'false')`).
			Inclde(`Id`, `Name`, `Effort`, `Project`, `EntityType`).
			Order(OrderByDesc(`Effort`)).
			Take(10))

	assert.NoError(err)
	assert.Equal(`https://restapi.tpondemand.com/api/v1/Assignables/?include=[Id,Name,Effort,Project,EntityType]&where=(EntityState.IsFinal eq 'false')&orderByDesc=Effort&format=json&take=10&skip=10`,
		response.Next)
	assert.Len(response.Items, 10)
}
