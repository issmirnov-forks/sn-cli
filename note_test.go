package sncli

import (
	"fmt"
	"github.com/jonhadfield/gosn"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestWipeWith50(t *testing.T) {
	fmt.Printf("TestWipeWith50 start time: %+v\n", time.Now())
	session, err := CliSignIn(os.Getenv("SN_EMAIL"), os.Getenv("SN_PASSWORD"), os.Getenv("SN_SERVER"))
	assert.NoError(t, err)
	cleanUp(&session)

	numNotes := 50
	textParas := 10
	err = createNotes(session, numNotes, textParas)
	assert.NoError(t, err)

	// check notes created
	noteFilter := gosn.Filter{
		Type: "Note",
	}
	filters := gosn.ItemFilters{
		Filters: []gosn.Filter{noteFilter},
	}
	gni := gosn.GetItemsInput{
		Session: session,
		Filters: filters,
	}
	var gno gosn.GetItemsOutput
	gno, err = gosn.GetItems(gni)

	assert.Equal(t, len(gno.Items), 50)
	wipeConfig := WipeConfig{
		Session: session,
	}
	var deleted int
	deleted, err = wipeConfig.Run()
	assert.NoError(t, err)
	assert.True(t, deleted >= numNotes, fmt.Sprintf("notes created: %d items deleted: %d", numNotes, deleted))
	fmt.Printf("TestWipeWith50 end time: %+v\n", time.Now())
	time.Sleep(1 * time.Second)
}

func TestAddDeleteNoteByUUID(t *testing.T) {
	fmt.Printf("TestAddDeleteNoteByUUID start time: %+v\n", time.Now())

	session, err := CliSignIn(os.Getenv("SN_EMAIL"), os.Getenv("SN_PASSWORD"), os.Getenv("SN_SERVER"))
	assert.NoError(t, err)
	cleanUp(&session)

	// create note
	addNoteConfig := AddNoteConfig{
		Session: session,
		Title:   "TestNoteOne",
		Text:    "TestNoteOneText",
	}
	err = addNoteConfig.Run()
	assert.NoError(t, err, err)

	// get new note
	filter := gosn.Filter{
		Type:       "Note",
		Key:        "Title",
		Comparison: "==",
		Value:      "TestNoteOne",
	}

	iFilter := gosn.ItemFilters{
		Filters: []gosn.Filter{filter},
	}
	gnc := GetNoteConfig{
		Session: session,
		Filters: iFilter,
	}
	var preRes, postRes gosn.GetItemsOutput
	preRes, err = gnc.Run()
	assert.NoError(t, err, err)

	newItemUUID := preRes.Items[0].UUID
	deleteNoteConfig := DeleteNoteConfig{
		Session:   session,
		NoteUUIDs: []string{newItemUUID},
	}
	var noDeleted int
	noDeleted, err = deleteNoteConfig.Run()
	assert.Equal(t, noDeleted, 1)
	assert.NoError(t, err, err)

	postRes, err = gnc.Run()
	assert.NoError(t, err, err)
	assert.EqualValues(t, len(postRes.Items), 0, "note was not deleted")
	cleanUp(&session)

	fmt.Printf("TestAddDeleteNoteByUUID end time: %+v\n", time.Now())
	time.Sleep(1 * time.Second)
}

func TestAddDeleteNoteByTitle(t *testing.T) {
	session, err := CliSignIn(os.Getenv("SN_EMAIL"), os.Getenv("SN_PASSWORD"), os.Getenv("SN_SERVER"))
	assert.NoError(t, err, err)

	wipeConfig := WipeConfig{
		Session: session,
	}
	_, err = wipeConfig.Run()
	assert.NoError(t, err)

	addNoteConfig := AddNoteConfig{
		Session: session,
		Title:   "TestNoteOne",
	}
	err = addNoteConfig.Run()
	assert.NoError(t, err, err)

	deleteNoteConfig := DeleteNoteConfig{
		Session:    session,
		NoteTitles: []string{"TestNoteOne"},
	}
	var noDeleted int
	noDeleted, err = deleteNoteConfig.Run()
	assert.Equal(t, noDeleted, 1)
	assert.NoError(t, err, err)

	filter := gosn.Filter{
		Type:       "Note",
		Key:        "Title",
		Comparison: "==",
		Value:      "TestNoteOne",
	}

	iFilter := gosn.ItemFilters{
		Filters: []gosn.Filter{filter},
	}
	gnc := GetNoteConfig{
		Session: session,
		Filters: iFilter,
	}
	var postRes gosn.GetItemsOutput
	postRes, err = gnc.Run()
	assert.NoError(t, err, err)
	assert.EqualValues(t, len(postRes.Items), 0, "note was not deleted")

	cleanUp(&session)
	time.Sleep(1 * time.Second)
}

func TestAddDeleteNoteByTitleRegex(t *testing.T) {
	session, err := CliSignIn(os.Getenv("SN_EMAIL"), os.Getenv("SN_PASSWORD"), os.Getenv("SN_SERVER"))
	assert.NoError(t, err, err)
	cleanUp(&session)

	assert.NoError(t, err)
	// add note
	addNoteConfig := AddNoteConfig{
		Session: session,
		Title:   "TestNoteOne",
	}
	err = addNoteConfig.Run()
	assert.NoError(t, err, err)

	// delete note
	deleteNoteConfig := DeleteNoteConfig{
		Session:    session,
		NoteTitles: []string{"^T.*ote..[def]"},
		Regex:      true,
	}
	var noDeleted int
	noDeleted, err = deleteNoteConfig.Run()
	assert.Equal(t, noDeleted, 1)
	assert.NoError(t, err, err)

	// get same note again
	filter := gosn.Filter{
		Type:       "Note",
		Key:        "Title",
		Comparison: "==",
		Value:      "TestNoteOne",
	}
	iFilter := gosn.ItemFilters{
		Filters: []gosn.Filter{filter},
	}
	gnc := GetNoteConfig{
		Session: session,
		Filters: iFilter,
	}
	var postRes gosn.GetItemsOutput
	postRes, err = gnc.Run()

	assert.NoError(t, err, err)
	assert.EqualValues(t, len(postRes.Items), 0, "note was not deleted")

	cleanUp(&session)
	time.Sleep(1 * time.Second)
}

func TestGetNote(t *testing.T) {
	session, err := CliSignIn(os.Getenv("SN_EMAIL"), os.Getenv("SN_PASSWORD"), os.Getenv("SN_SERVER"))
	assert.NoError(t, err)
	cleanUp(&session)

	// create one note
	addNoteConfig := AddNoteConfig{
		Session: session,
		Title:   "TestNoteOne",
	}
	err = addNoteConfig.Run()
	assert.NoError(t, err)

	noteFilter := gosn.Filter{
		Type:       "Note",
		Key:        "Title",
		Comparison: "==",
		Value:      "TestNoteOne",
	}
	// retrieve one note
	itemFilters := gosn.ItemFilters{
		MatchAny: false,
		Filters:  []gosn.Filter{noteFilter},
	}
	getNoteConfig := GetNoteConfig{
		Session: session,
		Filters: itemFilters,
	}
	var output gosn.GetItemsOutput
	output, err = getNoteConfig.Run()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(output.Items))

	cleanUp(&session)
	time.Sleep(1 * time.Second)
}

func TestCreateOneHundredNotes(t *testing.T) {
	numNotes := 100
	textParas := 10
	session, err := CliSignIn(os.Getenv("SN_EMAIL"), os.Getenv("SN_PASSWORD"), os.Getenv("SN_SERVER"))
	assert.NoError(t, err)
	cleanUp(&session)

	err = createNotes(session, numNotes, textParas)
	assert.NoError(t, err)

	noteFilter := gosn.Filter{
		Type: "Note",
	}
	filter := gosn.ItemFilters{
		Filters: []gosn.Filter{noteFilter},
	}

	gnc := GetNoteConfig{
		Session: session,
		Filters: filter,
	}
	var res gosn.GetItemsOutput
	res, err = gnc.Run()
	assert.NoError(t, err)

	assert.True(t, len(res.Items) >= numNotes)
	wipeConfig := WipeConfig{
		Session: session,
	}
	var deleted int
	deleted, err = wipeConfig.Run()
	assert.NoError(t, err)
	assert.True(t, deleted >= numNotes)
	cleanUp(&session)
	time.Sleep(1 * time.Second)
}

func cleanUp(session *gosn.Session) {
	wipeConfig := WipeConfig{
		Session: *session,
	}
	_, err := wipeConfig.Run()
	if err != nil {
		panic(err)
	}
}
