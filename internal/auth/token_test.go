package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	ID     uuid.UUID = uuid.New()
	SECRET string    = "test-secrent"
)

func Test_ValidJWT(t *testing.T) {
	token, err := MakeJWT(ID, SECRET, time.Duration(time.Second*5))
	require.NoError(t, err)
	require.NotEmpty(t, token)

	id, err := ValidateJWT(token, SECRET)
	require.NoError(t, err)
	require.Equal(t, id, ID)
}

func Test_ExpiredJWT(t *testing.T) {
	token, err := MakeJWT(ID, SECRET, time.Duration(time.Second*1))
	require.NoError(t, err)
	require.NotEmpty(t, token)

	time.Sleep(time.Second * 2)

	id, err := ValidateJWT(token, SECRET)
	require.Error(t, err)
	require.Equal(t, id, uuid.Nil)
}

func Test_GetBearerToken(t *testing.T) {
	token, err := MakeJWT(ID, SECRET, time.Duration(time.Second*5))
	require.NoError(t, err)
	require.NotEmpty(t, token)

	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	headerToken, err := GetBearerToken(headers)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, token, headerToken)
}
