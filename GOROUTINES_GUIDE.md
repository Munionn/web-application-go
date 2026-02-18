# Using Goroutines in HTTP Handlers

## Important Note

**HTTP handlers in Go already run in separate goroutines!** The `http.Server` automatically handles each request in its own goroutine. You don't need goroutines just to handle requests concurrently.

## When to Use Goroutines in Handlers

Use goroutines when you need to:
1. **Run background tasks** after sending the response (logging, analytics, notifications)
2. **Execute parallel database queries** and combine results
3. **Perform async operations** that don't block the response

## Common Patterns

### Pattern 1: Fire-and-Forget Background Task

**Use case:** Logging, analytics, notifications that don't need to block the response.

```go
func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
    // ... fetch users ...
    
    // Send response first
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
    
    // Background task (runs after response is sent)
    go func(users []models.User) {
        log.Printf("Fetched %d users", len(users))
        // Send analytics, update cache, etc.
    }(users)
}
```

**⚠️ Important:** Don't write to `w` (ResponseWriter) inside the goroutine! The response must be sent before the handler returns.

### Pattern 2: Parallel Database Queries

**Use case:** Fetch multiple things concurrently and combine results.

```go
func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
    db := s.db.GetDB()
    
    // Channels to collect results
    activeChan := make(chan []models.User, 1)
    errorChan := make(chan error, 1)
    
    // Fetch in goroutine
    go func() {
        var users []models.User
        result := db.Where("active = ?", true).Find(&users)
        if result.Error != nil {
            errorChan <- result.Error
            return
        }
        activeChan <- users
    }()
    
    // Wait for result
    select {
    case users := <-activeChan:
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    case err := <-errorChan:
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

### Pattern 3: Using sync.WaitGroup

**Use case:** Wait for multiple goroutines to complete before responding.

```go
func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var allUsers []models.User
    
    wg.Add(2) // Number of goroutines
    
    // Goroutine 1
    go func() {
        defer wg.Done()
        var users []models.User
        db.Where("active = ?", true).Find(&users)
        mu.Lock()
        allUsers = append(allUsers, users...)
        mu.Unlock()
    }()
    
    // Goroutine 2
    go func() {
        defer wg.Done()
        // Another query...
    }()
    
    wg.Wait() // Wait for all goroutines
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(allUsers)
}
```

### Pattern 4: Using Context for Timeout

**Use case:** Set timeout for database operations.

```go
func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    resultChan := make(chan []models.User, 1)
    errorChan := make(chan error, 1)
    
    go func() {
        var users []models.User
        result := db.WithContext(ctx).Where("active = ?", true).Find(&users)
        if result.Error != nil {
            errorChan <- result.Error
            return
        }
        resultChan <- users
    }()
    
    select {
    case users := <-resultChan:
        json.NewEncoder(w).Encode(users)
    case err := <-errorChan:
        http.Error(w, err.Error(), http.StatusInternalServerError)
    case <-ctx.Done():
        http.Error(w, "Request timeout", http.StatusRequestTimeout)
    }
}
```

## ⚠️ Common Mistakes

### ❌ DON'T: Write to ResponseWriter in Goroutine

```go
// WRONG - This will panic!
go func() {
    json.NewEncoder(w).Encode(users) // ❌ Can't write after handler returns
}()
```

### ❌ DON'T: Use Goroutines Unnecessarily

```go
// WRONG - No benefit, just adds complexity
go func() {
    result := db.Find(&users) // ❌ Unnecessary goroutine
}()
```

### ✅ DO: Send Response First, Then Background Work

```go
// CORRECT
json.NewEncoder(w).Encode(users) // Send response first
go func() {
    // Background work after response sent ✅
    log.Printf("Processed %d users", len(users))
}()
```

## Best Practices

1. **Send response before goroutine** - Always write to ResponseWriter before starting background work
2. **Use channels for results** - When you need data from goroutines, use channels
3. **Handle errors** - Always check for errors from goroutines
4. **Use context** - Pass request context to goroutines for cancellation/timeout
5. **Protect shared data** - Use mutexes when multiple goroutines access shared variables

## Example Files

See `user_handlers_examples.go` for complete working examples of all patterns.
