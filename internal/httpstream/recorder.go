package httpstream

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

// Record represents a single JSONL record.
type Record struct {
	Timestamp   string `json:"ts"`
	SessionID   string `json:"session"`
	SessionSeq  int64  `json:"seq"`   // Global session sequence number
	RecordIndex int64  `json:"index"` // Record index within session
	Type        string `json:"type"`  // request, response, sse, body, grpc, error

	// Request fields
	Method string `json:"method,omitempty"`
	URL    string `json:"url,omitempty"`
	Host   string `json:"host,omitempty"`

	// Response fields
	Status     int    `json:"status,omitempty"`
	StatusText string `json:"status_text,omitempty"`

	// SSE fields
	EventType string `json:"event_type,omitempty"`
	EventID   string `json:"event_id,omitempty"`
	EventData string `json:"event_data,omitempty"`

	// Headers - always included for request/response
	Headers map[string][]string `json:"headers,omitempty"`

	// Body fields
	Direction    string `json:"direction,omitempty"`     // C2S or S2C
	Size         int    `json:"size,omitempty"`          // Body size in bytes
	Body         string `json:"body,omitempty"`          // Full body (text)
	BodyBase64   string `json:"body_base64,omitempty"`   // Full body (base64 for binary)
	BodyEncoding string `json:"body_encoding,omitempty"` // "text" or "base64"
	ContentType  string `json:"content_type,omitempty"`  // Content-Type header

	// gRPC fields
	GRPCService    string `json:"grpc_service,omitempty"`     // e.g., "aiserver.v1.RepositoryService"
	GRPCMethod     string `json:"grpc_method,omitempty"`      // e.g., "SyncMerkleSubtreeV2"
	GRPCData       string `json:"grpc_data,omitempty"`        // JSON representation of protobuf message
	GRPCStreaming  bool   `json:"grpc_streaming,omitempty"`   // Is this a streaming RPC
	GRPCFrameIndex int    `json:"grpc_frame_index,omitempty"` // Frame index in streaming (0-based)
	GRPCCompressed bool   `json:"grpc_compressed,omitempty"`  // Frame compressed flag
	GRPCRawData    string `json:"grpc_raw,omitempty"`         // Base64 raw frame data (on error)

	// Error
	Error string `json:"error,omitempty"`
}

// RecordCallback is called when a record is written.
type RecordCallback func(Record)

// Recorder writes HTTP traffic to JSONL file with session tracking.
type Recorder struct {
	mu       sync.Mutex
	file     *os.File
	encoder  *json.Encoder
	logLevel LogLevel

	// Stats
	records    atomic.Int64
	sessionSeq atomic.Int64 // Session sequence counter

	// Callbacks
	onRecord RecordCallback

	// Memory cache for recent records (for initial frontend load)
	cacheMu      sync.RWMutex
	recordCache  []Record
	maxCacheSize int
}

// RecorderOption configures a Recorder.
type RecorderOption func(*Recorder)

// WithRecorderLogLevel sets the log level for recording.
func WithRecorderLogLevel(level LogLevel) RecorderOption {
	return func(r *Recorder) { r.logLevel = level }
}

// WithOnRecord sets a callback for each record written.
func WithOnRecord(cb RecordCallback) RecorderOption {
	return func(r *Recorder) { r.onRecord = cb }
}

// WithCacheSize sets the maximum number of records to cache in memory.
func WithCacheSize(size int) RecorderOption {
	return func(r *Recorder) { r.maxCacheSize = size }
}

// NewRecorder creates a new JSONL recorder.
// If path is empty, no file is written (in-memory only, for WebSocket/cache use).
func NewRecorder(path string, opts ...RecorderOption) (*Recorder, error) {
	r := &Recorder{
		logLevel:     LogLevelBasic,
		recordCache:  make([]Record, 0, 1000),
		maxCacheSize: 1000, // Keep last 1000 records for initial load
	}

	if path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0644)
		if err != nil {
			return nil, fmt.Errorf("open recorder file: %w", err)
		}
		r.file = file
		r.encoder = json.NewEncoder(file)
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

// Close closes the recorder file.
func (r *Recorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// write writes a record to the file (thread-safe, sync write).
func (r *Recorder) write(rec Record) error {
	if r.encoder != nil {
		r.mu.Lock()
		if err := r.encoder.Encode(rec); err != nil {
			r.mu.Unlock()
			return err
		}
		r.mu.Unlock()
	}

	r.records.Add(1)

	// Add to cache
	r.addToCache(rec)

	// Call callback if set
	if r.onRecord != nil {
		r.onRecord(rec)
	}

	return nil
}

// addToCache adds a record to the memory cache.
func (r *Recorder) addToCache(rec Record) {
	r.cacheMu.Lock()
	defer r.cacheMu.Unlock()

	// Add record to cache
	r.recordCache = append(r.recordCache, rec)

	// Trim if over max size (remove oldest)
	if len(r.recordCache) > r.maxCacheSize {
		r.recordCache = r.recordCache[1:]
	}
}

// RecordCount returns the number of records written.
func (r *Recorder) RecordCount() int64 {
	return r.records.Load()
}

// Session represents a tracked HTTP session.
type Session struct {
	ID          string
	Seq         int64 // Global session sequence number
	Host        string
	recorder    *Recorder
	recordIndex int64 // Record index counter within session
}

// NewSession creates a new tracked session.
func (r *Recorder) NewSession(host string) *Session {
	seq := r.sessionSeq.Add(1)
	return &Session{
		ID:          generateSessionID(),
		Seq:         seq,
		Host:        host,
		recorder:    r,
		recordIndex: 0,
	}
}

// nextRecordIndex returns and increments the record index.
func (s *Session) nextRecordIndex() int64 {
	s.recordIndex++
	return s.recordIndex
}

// Note: generateSessionID is defined in parser.go

// timestamp returns current time in RFC3339 format.
func timestamp() string {
	return time.Now().Format(time.RFC3339Nano)
}

// LogRequest logs an HTTP request.
// Note: Always records to JSONL regardless of log level.
func (s *Session) LogRequest(msg *HTTPMessage) {
	req := msg.Request
	if req == nil {
		return
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "request",
		Method:      req.Method,
		URL:         req.URL.RequestURI(),
		Host:        s.Host,
		Headers:     cloneHeaders(req.Header),
		ContentType: req.Header.Get("Content-Type"),
	}

	s.recorder.write(rec)
}

// LogResponse logs an HTTP response.
// Note: Always records to JSONL regardless of log level.
func (s *Session) LogResponse(msg *HTTPMessage) {
	resp := msg.Response
	if resp == nil {
		return
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "response",
		Status:      resp.StatusCode,
		StatusText:  resp.Status,
		Host:        s.Host,
		Headers:     cloneHeaders(resp.Header),
		ContentType: resp.Header.Get("Content-Type"),
	}

	s.recorder.write(rec)
}

// LogSSE logs an SSE event.
// Note: Always records to JSONL regardless of log level.
func (s *Session) LogSSE(host string, event *SSEEvent) {
	eventType := event.Event
	if eventType == "" {
		eventType = "message"
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "sse",
		Host:        host,
		EventType:   eventType,
		EventID:     event.ID,
		EventData:   truncateString(event.Data, 1000),
	}

	s.recorder.write(rec)
}

// LogBody logs body data (full content).
// Note: Always records to JSONL regardless of log level.
func (s *Session) LogBody(dir Direction, host string, data []byte) {
	if len(data) == 0 {
		return
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "body",
		Direction:   dir.String(),
		Host:        host,
		Size:        len(data),
	}

	// Store full body with appropriate encoding
	if utf8.Valid(data) && isPrintableText(data) {
		rec.Body = string(data)
		rec.BodyEncoding = "text"
	} else {
		rec.BodyBase64 = base64.StdEncoding.EncodeToString(data)
		rec.BodyEncoding = "base64"
	}

	s.recorder.write(rec)
}

// LogGRPC logs a gRPC message.
func (s *Session) LogGRPC(msg *GRPCMessage) {
	rec := Record{
		Timestamp:      timestamp(),
		SessionID:      s.ID,
		SessionSeq:     s.Seq,
		RecordIndex:    s.nextRecordIndex(),
		Type:           "grpc",
		Direction:      msg.Direction.String(),
		Host:           s.Host,
		GRPCService:    msg.Service,
		GRPCMethod:     msg.Method,
		URL:            msg.FullMethod,
		GRPCStreaming:  msg.IsStreaming,
		GRPCFrameIndex: msg.FrameIndex,
		GRPCCompressed: msg.Compressed,
	}

	if msg.JSON != "" {
		rec.GRPCData = msg.JSON
	} else if msg.Frame != nil {
		rec.Size = len(msg.Frame.Data)
	}

	if msg.Error != "" {
		rec.Error = msg.Error
		// Include raw data on error for debugging
		if msg.Frame != nil && len(msg.Frame.Data) > 0 {
			rec.GRPCRawData = base64.StdEncoding.EncodeToString(msg.Frame.Data)
		}
	}

	s.recorder.write(rec)
}

// isPrintableText checks if data is printable text.
func isPrintableText(data []byte) bool {
	for _, b := range data {
		// Allow printable ASCII, newlines, tabs
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			return false
		}
		// Reject DEL and most control chars
		if b == 127 {
			return false
		}
	}
	return true
}

// Debug logs debug information.
func (s *Session) Debug(format string, args ...interface{}) {
	if s.recorder.logLevel < LogLevelDebug {
		return
	}

	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "debug",
		Host:        s.Host,
		Error:       fmt.Sprintf(format, args...),
	}

	s.recorder.write(rec)
}

// LogError logs an error.
func (s *Session) LogError(err error) {
	rec := Record{
		Timestamp:   timestamp(),
		SessionID:   s.ID,
		SessionSeq:  s.Seq,
		RecordIndex: s.nextRecordIndex(),
		Type:        "error",
		Host:        s.Host,
		Error:       err.Error(),
	}

	s.recorder.write(rec)
}

// cloneHeaders creates a copy of headers.
func cloneHeaders(h http.Header) map[string][]string {
	if h == nil {
		return nil
	}
	clone := make(map[string][]string, len(h))
	for k, v := range h {
		clone[k] = append([]string(nil), v...)
	}
	return clone
}

// truncateString truncates string to max length.
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// SessionLogger wraps Session to implement Logger interface.
type SessionLogger struct {
	*Session
}

// Ensure SessionLogger implements Logger.
var _ Logger = (*SessionLogger)(nil)

// NewSessionLogger creates a logger for a session.
func (s *Session) Logger() *SessionLogger {
	return &SessionLogger{Session: s}
}

// WriteTo implements io.WriterTo for streaming records.
func (r *Recorder) WriteTo(w io.Writer) (int64, error) {
	if r.file != nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		return 0, r.file.Sync()
	}
	return 0, nil
}

// GetRecentRecords returns the most recent records (for initial frontend load).
func (r *Recorder) GetRecentRecords(limit int) []interface{} {
	r.cacheMu.RLock()
	defer r.cacheMu.RUnlock()

	if limit <= 0 || limit > len(r.recordCache) {
		limit = len(r.recordCache)
	}

	// Return the most recent records (last N items)
	start := len(r.recordCache) - limit
	if start < 0 {
		start = 0
	}

	results := make([]interface{}, 0, limit)
	for i := start; i < len(r.recordCache); i++ {
		results = append(results, r.recordCache[i])
	}

	return results
}
