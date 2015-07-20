package main

import "github.com/pborman/uuid"

type postEvent struct {
	id postId
}

type postAdded struct {
	postEvent
	content postContent
}

type postEdited struct {
	postEvent
	content postContent
}

type postDeleted struct {
	postEvent
}

type postContent struct {
	author, title, content string
}

// postId represents a unique postId
type postId struct {
	id uuid.UUID

	// wraps a uuid but with a smaller interface footprint
}

func newRandomPostId() postId {
	return postId{uuid.NewRandom()}
}

func (pid postId) String() string {
	return pid.id.String()
}
