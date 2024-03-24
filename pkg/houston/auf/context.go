package auf

import "context"

type userIDKey struct{}

func AddUserIDToContext[T any](ctx context.Context, userID T) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func ExtractUserIDFromContext[T any](ctx context.Context) T {
	return ctx.Value(userIDKey{}).(T)
}
