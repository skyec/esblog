package main

// Unlike the original Scala implementation, this does not use immutable data structures.
// Immutability may be added as a later exercise/experiment to see how much work it
// actually is and how it will affect actual performance.

import (
	"sort"
	"sync"
)

// A post is a blog post. It has an ID and content.
type post struct {
	id      postId
	content postContent
}

// posts manages a collection of blog posts.
type posts struct {
	byId        map[string]*post
	byTimeAdded postSlice
	mx          sync.Mutex
}

/// newPosts creates a new collection of blog posts.
func newPosts() *posts {
	return &posts{
		byId:        map[string]*post{},
		byTimeAdded: postSlice{},
	}
}

// get returns the post that matches the provided postId
func (p *posts) get(id postId) *post {
	return p.byId[id.String()]
}

// mostRecent returns the most recent n posts. An empty slice is returned
// if there are no posts. All the posts are returned if more posts are
// requested than there are posts available.
func (p *posts) mostRecent(n int) []*post {
	orig := []*post(p.byTimeAdded)

	if n > len(orig) {
		n = len(orig)
	}

	r := make([]*post, n)
	copy(r, orig[len(orig)-n:])
	sort.Sort(sort.Reverse(postSlice(r)))
	return r
}

// apply updates the state based on the event type
func (p *posts) apply(event interface{}) *posts {
	p.mx.Lock()
	defer p.mx.Unlock()

	switch e := event.(type) {
	case postAdded:
		np := &post{
			id:      e.id,
			content: e.content,
		}
		p.byId[np.id.String()] = np
		p.byTimeAdded = append(p.byTimeAdded, np)

	case postEdited:
		p.byId[e.id.String()].content = e.content

	case postDeleted:
		delete(p.byId, e.id.String())
		p.byTimeAdded = p.byTimeAdded.Filter(func(pp interface{}) bool {
			if pp.(*post).id.String() == e.id.String() {
				return false
			}
			return true
		})
	}
	return p
}

// postSlice implements sort.Interface to make our slice of posts sortable
type postSlice []*post

func (s postSlice) Len() int           { return len(s) }
func (s postSlice) Swap(i, j int)      { s[j], s[i] = s[i], s[j] }
func (s postSlice) Less(i, j int) bool { return i < j }

// Filter calls the filter callback for each item in the slice. If the callback
// retruns true then the item is coppied to the new slice and is skipped if the
// return is false.
func (s postSlice) Filter(filter func(interface{}) bool) postSlice {
	snew := postSlice{}
	for _, p := range s {
		if filter(p) {
			snew = append(snew, p)
		}
	}
	return snew
}

// postsFromHistory is a heper that initialzes a posts object from historical events.
func postsFromHistory(events []interface{}) *posts {
	posts := newPosts()
	for _, v := range events {
		posts.apply(v)
	}
	return posts
}
