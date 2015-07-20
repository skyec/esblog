package main

import "testing"

func TestAddPost(t *testing.T) {

	id := newRandomPostId()

	p := newPosts().
		apply(postAdded{postEvent{id}, postContent{"me", "the subject", "the body"}}).
		get(id)

	if p.id.String() != id.String() {
		t.Errorf("Expected: %s, got: %s", id.String(), p.id.String())
	}
}

func TestGetMostRecent(t *testing.T) {
	ps := newPosts().
		apply(postAdded{postEvent{newRandomPostId()}, postContent{"me", "the subject", "the body"}}).
		apply(postAdded{postEvent{newRandomPostId()}, postContent{"me", "the subject 2", "the body"}}).
		apply(postAdded{postEvent{newRandomPostId()}, postContent{"me", "the subject 3", "the body"}})

	// get just the most recent one
	recent := ps.mostRecent(1)
	if len(recent) != 1 {
		t.Errorf("Expected one recent post. Got: %d", len(recent))
	}
	if recent[0].content.title != "the subject 3" {
		t.Errorf("Expected the 3rd subject, got: %s", recent[0].content.title)
	}

	// negative case: ask for zero most recents
	recent = ps.mostRecent(0)
	if len(recent) != 0 {
		t.Errorf("Expected an empty slice, got: %d", len(recent))
	}

	// ask only for the most recent two (not the full set)
	recent = ps.mostRecent(2)
	if len(recent) != 2 {
		t.Errorf("Expected two posts, got: %d", len(recent))
	}
	if recent[1].content.title != "the subject 2" {
		t.Errorf("Expected the 2nd subject, got: %s", recent[1].content.title)
	}

	// ask for the full set
	recent = ps.mostRecent(3)
	if len(recent) != 3 {
		t.Errorf("Expected three posts, got: %d", len(recent))
	}
	for i, v := range []string{"the subject 3", "the subject 2", "the subject"} {
		if recent[i].content.title != v {
			t.Errorf("Expected %s, got: %s", v, recent[i].content.title)
		}
	}

	// asking for more than the full set returns just the full set
	recent = ps.mostRecent(99)
	if len(recent) != 3 {
		t.Errorf("Expected three posts, got: %d", len(recent))
	}
	for i, v := range []string{"the subject 3", "the subject 2", "the subject"} {
		if recent[i].content.title != v {
			t.Errorf("Expected %s, got: %s", v, recent[i].content.title)
		}
	}
}

func TestPostsFromHistory(t *testing.T) {

	deletedId := newRandomPostId()
	editedId := newRandomPostId()

	posts := postsFromHistory([]interface{}{
		postAdded{postEvent{deletedId},
			postContent{"me", "First Post!", "This is the first blog post"}},
		postAdded{postEvent{editedId},
			postContent{"me", "The worlds best pulled pork", "Here is the recipie for the worlds best pulled pork sandwiches"}},
		postEdited{postEvent{editedId},
			postContent{"me", "Sometimes the worlds best pulled pork", "Well, after last night, maybe this isn't the world's best."}},
		postAdded{postEvent{newRandomPostId()},
			postContent{"me", "Gradma's chocolate chip cookies", "An old family favourite"}},
		postDeleted{postEvent{deletedId}},
	})

	if posts == nil {
		t.Fatal("Posts is nil")
	}

	p := posts.get(editedId)
	if p == nil {
		t.Fatal("Edited post is nil")
	}
	expected := "Sometimes the worlds best pulled pork"
	if p.content.title != expected {
		t.Errorf("Expected: %s, got: %s", expected, p.content.title)
	}

	p = posts.get(deletedId)
	if p != nil {
		t.Errorf("Expected nil for deleted post but got %-v", p)
	}

}
