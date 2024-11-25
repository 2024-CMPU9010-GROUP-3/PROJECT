//go:build public

package dtos

type GetLocationHistoryListDto struct {
	HistoryEntries []LocationHistoryEntryDto `json:"history"`
	NextOffset     int32                     `json:"nextoffset"`
}
