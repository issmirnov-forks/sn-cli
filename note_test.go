package sncli

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jonhadfield/gosn"
	"github.com/stretchr/testify/assert"
)

var testSession gosn.Session

func TestMain(m *testing.M) {
	testSession, _ = CliSignIn(os.Getenv("SN_EMAIL"), os.Getenv("SN_PASSWORD"), os.Getenv("SN_SERVER"))
	os.Exit(m.Run())
}

func TestWipeWith50(t *testing.T) {
	cleanUp(&testSession)
	defer cleanUp(&testSession)

	numNotes := 50
	textParas := 10
	err := createNotes(testSession, numNotes, textParas)
	assert.NoError(t, err)

	// check notes created
	noteFilter := gosn.Filter{
		Type: "Note",
	}
	filters := gosn.ItemFilters{
		Filters: []gosn.Filter{noteFilter},
	}
	gni := gosn.GetItemsInput{
		Session: testSession,
		Filters: filters,
	}
	var gno gosn.GetItemsOutput
	gno, err = gosn.GetItems(gni)

	assert.Equal(t, 50, len(gno.Items))
	wipeConfig := WipeConfig{
		Session: testSession,
	}
	var deleted int
	deleted, err = wipeConfig.Run()
	assert.NoError(t, err)
	assert.True(t, deleted >= numNotes, fmt.Sprintf("notes created: %d items deleted: %d", numNotes, deleted))
	time.Sleep(1 * time.Second)
}

func TestAddDeleteNoteByUUID(t *testing.T) {
	defer cleanUp(&testSession)

	// create note
	addNoteConfig := AddNoteConfig{
		Session: testSession,
		Title:   "TestNoteOne",
		Text:    "TestNoteOneText",
	}
	err := addNoteConfig.Run()
	assert.NoError(t, err, err)
	time.Sleep(1 * time.Second)

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
		Session: testSession,
		Filters: iFilter,
	}
	var preRes, postRes gosn.GetItemsOutput
	preRes, err = gnc.Run()
	time.Sleep(1 * time.Second)

	assert.NoError(t, err, err)

	newItemUUID := preRes.Items[0].UUID
	deleteNoteConfig := DeleteNoteConfig{
		Session:   testSession,
		NoteUUIDs: []string{newItemUUID},
	}
	var noDeleted int
	noDeleted, err = deleteNoteConfig.Run()
	time.Sleep(1 * time.Second)
	assert.Equal(t, noDeleted, 1)
	assert.NoError(t, err, err)

	postRes, err = gnc.Run()
	assert.NoError(t, err, err)
	assert.EqualValues(t, len(postRes.Items), 0, "note was not deleted")
}

func TestAddDeleteNoteByTitle(t *testing.T) {
	defer cleanUp(&testSession)

	addNoteConfig := AddNoteConfig{
		Session: testSession,
		Title:   "TestNoteOne",
	}
	err := addNoteConfig.Run()
	assert.NoError(t, err, err)

	deleteNoteConfig := DeleteNoteConfig{
		Session:    testSession,
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
		Session: testSession,
		Filters: iFilter,
	}
	var postRes gosn.GetItemsOutput
	postRes, err = gnc.Run()
	assert.NoError(t, err, err)
	assert.EqualValues(t, len(postRes.Items), 0, "note was not deleted")
}

func TestAddDeleteNoteByTitleRegex(t *testing.T) {
	defer cleanUp(&testSession)
	// add note
	addNoteConfig := AddNoteConfig{
		Session: testSession,
		Title:   "TestNoteOne",
	}
	err := addNoteConfig.Run()
	assert.NoError(t, err, err)

	// delete note
	deleteNoteConfig := DeleteNoteConfig{
		Session:    testSession,
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
		Session: testSession,
		Filters: iFilter,
	}
	var postRes gosn.GetItemsOutput
	postRes, err = gnc.Run()

	assert.NoError(t, err, err)
	assert.EqualValues(t, len(postRes.Items), 0, "note was not deleted")
}

func TestGetNote(t *testing.T) {
	defer cleanUp(&testSession)

	// create one note
	addNoteConfig := AddNoteConfig{
		Session: testSession,
		Title:   "TestNoteOne",
	}
	err := addNoteConfig.Run()
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
		Session: testSession,
		Filters: itemFilters,
	}
	var output gosn.GetItemsOutput
	output, err = getNoteConfig.Run()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(output.Items))

}

func TestCreateOneHundredNotes(t *testing.T) {
	defer cleanUp(&testSession)
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
}

func cleanUp(session *gosn.Session) {
	wipeConfig := WipeConfig{
		Session: *session,
	}
	_, err := wipeConfig.Run()
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
}
