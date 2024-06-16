package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Needs for SQLite start
)

type DB interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Ping() error
}

type Storage struct {
	db DB
	h  SQLHelpersReader
}

func New(db DB, helpers SQLHelpersReader) (*Storage, error) {
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{
		db: db,
		h:  helpers,
	}, nil
}

// Init - initializes DB entities
func (s *Storage) Init(ctx context.Context) error {
	initDb, err := s.h.ReadSQLFile("storage/queries/init_db.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	_, err = s.db.ExecContext(ctx, initDb)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}

// Migrations - starts the migration process and saves results in DB
func (s *Storage) Migrations(ctx context.Context) error {
	// Map with names of already executed migrations
	migrations := make(map[string]string)

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM migrations`)
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		migrations[file] = file
	}

	files, err := os.ReadDir("storage/migrations")
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		_, ok := migrations[file.Name()]

		// Call migration and save it to the database in case if it has not been previously executed
		if !ok {
			// Read and execute migration code
			migration, err := s.h.ReadSQLFile("storage/migrations/" + file.Name())
			if err != nil {
				return fmt.Errorf("failed to read SQL file: %w", err)
			}

			_, err = s.db.ExecContext(ctx, migration)
			if err != nil {
				return fmt.Errorf("failed to run migration %s: %w", file.Name(), err)
			}

			q := `INSERT INTO migrations (timestamp) VALUES (?)`
			_, err = s.db.ExecContext(ctx, q, file.Name()) // Save migration in DB
			if err != nil {
				return fmt.Errorf("failed to find migration %s: %w", file.Name(), err)
			}
		}
	}

	return nil
}

// GetClients - receives a list of clients. Optional parameter `clientID`. When passed, will receive only selected client
func (s *Storage) GetClients(ctx context.Context, clientID *int) ([]Client, error) {
	baseQuery, err := s.h.ReadSQLFile("storage/queries/clients.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file: %w", err)
	}

	var query string
	var args []interface{}
	if clientID != nil {
		query = fmt.Sprintf("%s WHERE c.id = ?", baseQuery)
		args = append(args, *clientID)
	} else {
		query = baseQuery
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	defer rows.Close()

	clientMap := make(map[int]*Client)

	for rows.Next() {
		var clientID int
		var clientName string
		var startDate, endDate string
		var priority Priority
		var leadCapacity int
		var leadID, leadStart, leadEnd sql.NullString

		err := rows.Scan(
			&clientID,
			&clientName,
			&startDate,
			&endDate,
			&priority,
			&leadCapacity,
			&leadID,
			&leadStart,
			&leadEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		client, exists := clientMap[clientID]
		if !exists {
			client = &Client{
				ID:           clientID,
				Name:         clientName,
				StartDate:    startDate,
				EndDate:      endDate,
				Priority:     priority,
				LeadCapacity: leadCapacity,
				Leads:        []Lead{},
			}
			clientMap[clientID] = client
		}

		if leadID.Valid {
			client.Leads = append(client.Leads, Lead{
				ClientID:  clientID,
				LeadID:    leadID.String,
				LeadStart: leadStart.String,
				LeadEnd:   leadEnd.String,
			})
		}
	}

	defer rows.Close()

	var clients []Client
	for _, client := range clientMap {
		clients = append(clients, *client)
	}

	return clients, nil
}

func (s *Storage) CreateClient(ctx context.Context, c ClientRequest) error {
	createClientQuery, err := s.h.ReadSQLFile("storage/queries/new_client.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	userID, err := s.generateUserID(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate client ID: %w", err)
	}

	_, err = s.db.ExecContext(
		ctx,
		createClientQuery,
		userID,
		c.Name,
		c.StartDate,
		c.EndDate,
		c.Priority,
		c.LeadCapacity,
	)
	if err != nil {
		return fmt.Errorf("failed to get clients: %w", err)
	}

	return nil
}

// AssignLead - Selects a suitable client for assignment. Assigns a Lead to him and returns ID of this client.
func (s *Storage) AssignLead(ctx context.Context, l AssignLeadRequest) (*Lead, error) {
	clients, err := s.GetClients(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	var availableClients []Client

	for _, client := range clients {
		if noCapacity(client) {
			continue
		}
		if unsuitableTime(client, l) {
			continue
		}

		availableClients = append(availableClients, client)
	}

	sortedClients := make([]Client, len(availableClients))
	copy(sortedClients, availableClients)

	sort.Slice(sortedClients, sortByPriorityAndCapacity(sortedClients))

	if len(sortedClients) == 0 {
		return nil, fmt.Errorf("there are no clients available to assign")
	}

	priorityUser := sortedClients[0]

	if &priorityUser != nil {
		assignLeadQuery, err := s.h.ReadSQLFile("storage/queries/assign_lead.sql")
		if err != nil {
			return nil, fmt.Errorf("failed to read SQL file: %w", err)
		}

		leadID, _ := uuid.NewUUID()

		_, err = s.db.ExecContext(
			ctx,
			assignLeadQuery,
			leadID.String(),
			priorityUser.ID,
			l.LeadStart,
			l.LeadEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("can't create lead: %w", err)
		}

		return &Lead{
			LeadID:    leadID.String(),
			ClientID:  priorityUser.ID,
			LeadStart: l.LeadStart,
			LeadEnd:   l.LeadEnd,
		}, nil
	}

	return nil, nil
}

func (s *Storage) generateUserID(ctx context.Context) (*int, error) {
	var count int
	countQuery := "select count(*) from clients;"

	err := s.db.QueryRowContext(ctx, countQuery).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("can't count clients: %w", err)
	}

	clientID := count + 1

	return &clientID, nil
}

func noCapacity(client Client) bool {
	return len(client.Leads) >= client.LeadCapacity
}

func unsuitableTime(client Client, lead AssignLeadRequest) bool {
	const layout = "2006-01-01 00:00:00"

	clientStart, err1 := time.Parse(layout, client.StartDate)
	clientEnd, err2 := time.Parse(layout, client.StartDate)
	leadStart, err3 := time.Parse(layout, lead.LeadStart)
	leadEnd, err4 := time.Parse(layout, lead.LeadEnd)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return true
	}

	if clientStart.After(leadStart) || clientEnd.Before(leadEnd) {
		return true
	}

	return false
}

func sortByPriorityAndCapacity(clients []Client) func(i, j int) bool {
	return func(i, j int) bool {
		if clients[i].Priority != clients[j].Priority {
			return PriorityMap[clients[i].Priority] > PriorityMap[clients[j].Priority]
		}

		freeLeadsPercentageI := ((clients[i].LeadCapacity - len(clients[i].Leads)) * 100) / clients[i].LeadCapacity
		freeLeadsPercentageJ := ((clients[j].LeadCapacity - len(clients[j].Leads)) * 100) / clients[j].LeadCapacity

		return freeLeadsPercentageI > freeLeadsPercentageJ
	}
}
