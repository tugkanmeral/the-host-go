package model

import "time"

type NewNoteRequest struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type NoteListingItemModel struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type NoteModel struct {
	Id             string    `json:"id"`
	Title          string    `json:"title"`
	Text           string    `json:"text"`
	Tags           []string  `json:"tags"`
	CreationDate   time.Time `json:"creationDate"`
	LastUpdateDate time.Time `json:"lastUpdateDate"`
}

type NotePartialUpdateRequest struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}
