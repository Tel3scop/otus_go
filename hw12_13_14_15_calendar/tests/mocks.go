package mocks

//go:generate mockgen -source=../internal/service/event.go -destination=service/event.go -package=eventServiceMocks
//go:generate mockgen -source=../internal/storage/event.go -destination=storage/event.go -package=eventRepositoryMocks
//go:generate mockgen -source=../internal/client/db/db.go -destination=db/db.go -package=dbmocks
