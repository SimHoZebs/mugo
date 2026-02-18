package adk

import (
	"context"
	"fmt"

	"google.golang.org/adk/session"
)

type LazySessionService struct {
	service session.Service
	err     error
}

func NewLazySessionService(service session.Service, err error) *LazySessionService {
	return &LazySessionService{
		service: service,
		err:     err,
	}
}

func (l *LazySessionService) Create(ctx context.Context, req *session.CreateRequest) (*session.CreateResponse, error) {
	if l.err != nil {
		return nil, fmt.Errorf("session service unavailable: %w", l.err)
	}
	return l.service.Create(ctx, req)
}

func (l *LazySessionService) Get(ctx context.Context, req *session.GetRequest) (*session.GetResponse, error) {
	if l.err != nil {
		return nil, fmt.Errorf("session service unavailable: %w", l.err)
	}
	return l.service.Get(ctx, req)
}

func (l *LazySessionService) List(ctx context.Context, req *session.ListRequest) (*session.ListResponse, error) {
	if l.err != nil {
		return nil, fmt.Errorf("session service unavailable: %w", l.err)
	}
	return l.service.List(ctx, req)
}

func (l *LazySessionService) Delete(ctx context.Context, req *session.DeleteRequest) error {
	if l.err != nil {
		return fmt.Errorf("session service unavailable: %w", l.err)
	}
	return l.service.Delete(ctx, req)
}

func (l *LazySessionService) AppendEvent(ctx context.Context, sess session.Session, event *session.Event) error {
	if l.err != nil {
		return fmt.Errorf("session service unavailable: %w", l.err)
	}
	return l.service.AppendEvent(ctx, sess, event)
}
