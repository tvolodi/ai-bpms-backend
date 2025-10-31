package migration

import (
	"fmt"
	"log"

	"github.com/tvolodi/ai-bpms-backend/shared/common/models"
	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	Version     string
	Description string
	Up          func(*gorm.DB) error
	Down        func(*gorm.DB) error
}

// MigrationRecord tracks applied migrations
type MigrationRecord struct {
	Version   string `gorm:"primaryKey"`
	AppliedAt int64  `gorm:"autoCreateTime"`
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: getMigrations(),
	}
}

// Run executes all pending migrations
func (m *Migrator) Run() error {
	// Create migration tracking table
	if err := m.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		if !applied[migration.Version] {
			log.Printf("Running migration %s: %s", migration.Version, migration.Description)

			if err := migration.Up(m.db); err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.Version, err)
			}

			// Record migration as applied
			record := MigrationRecord{Version: migration.Version}
			if err := m.db.Create(&record).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
			}

			log.Printf("Migration %s completed successfully", migration.Version)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// Rollback rolls back the last migration
func (m *Migrator) Rollback() error {
	// Get the last applied migration
	var lastRecord MigrationRecord
	if err := m.db.Order("applied_at DESC").First(&lastRecord).Error; err != nil {
		return fmt.Errorf("no migrations to rollback: %w", err)
	}

	// Find the migration
	var targetMigration *Migration
	for _, migration := range m.migrations {
		if migration.Version == lastRecord.Version {
			targetMigration = &migration
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %s not found", lastRecord.Version)
	}

	log.Printf("Rolling back migration %s: %s", targetMigration.Version, targetMigration.Description)

	// Run the down migration
	if err := targetMigration.Down(m.db); err != nil {
		return fmt.Errorf("rollback %s failed: %w", targetMigration.Version, err)
	}

	// Remove migration record
	if err := m.db.Delete(&lastRecord).Error; err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", targetMigration.Version, err)
	}

	log.Printf("Migration %s rolled back successfully", targetMigration.Version)
	return nil
}

// getAppliedMigrations returns a map of applied migration versions
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var records []MigrationRecord
	if err := m.db.Find(&records).Error; err != nil {
		return nil, err
	}

	applied := make(map[string]bool)
	for _, record := range records {
		applied[record.Version] = true
	}

	return applied, nil
}

// getMigrations returns all available migrations
func getMigrations() []Migration {
	return []Migration{
		{
			Version:     "001_initial_schema",
			Description: "Create initial database schema",
			Up:          migration001Up,
			Down:        migration001Down,
		},
		{
			Version:     "002_rbac_system",
			Description: "Create RBAC (Role-Based Access Control) system",
			Up:          migration002Up,
			Down:        migration002Down,
		},
		{
			Version:     "003_process_engine",
			Description: "Create process engine tables",
			Up:          migration003Up,
			Down:        migration003Down,
		},
		{
			Version:     "004_audit_system",
			Description: "Create audit logging system",
			Up:          migration004Up,
			Down:        migration004Down,
		},
		{
			Version:     "005_indexes_optimization",
			Description: "Add performance indexes",
			Up:          migration005Up,
			Down:        migration005Down,
		},
	}
}

// migration001Up - Initial schema
func migration001Up(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Create schemas
	schemas := []string{"process_engine", "user_management", "analytics", "audit"}
	for _, schema := range schemas {
		if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error; err != nil {
			return err
		}
	}

	// Create base tables
	return db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.FormSchema{},
	)
}

func migration001Down(db *gorm.DB) error {
	// Drop tables in reverse order
	return db.Migrator().DropTable(
		&models.RefreshToken{},
		&models.FormSchema{},
		&models.User{},
	)
}

// migration002Up - RBAC system
func migration002Up(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Permission{},
		&models.Role{},
	); err != nil {
		return err
	}

	// Create default permissions
	permissions := []models.Permission{
		{Name: "process:create", Description: "Create process definitions", Resource: "process", Action: "create"},
		{Name: "process:read", Description: "Read process definitions", Resource: "process", Action: "read"},
		{Name: "process:update", Description: "Update process definitions", Resource: "process", Action: "update"},
		{Name: "process:delete", Description: "Delete process definitions", Resource: "process", Action: "delete"},
		{Name: "process:start", Description: "Start process instances", Resource: "process", Action: "start"},
		{Name: "task:assign", Description: "Assign tasks to users", Resource: "task", Action: "assign"},
		{Name: "task:complete", Description: "Complete tasks", Resource: "task", Action: "complete"},
		{Name: "task:view", Description: "View tasks", Resource: "task", Action: "view"},
		{Name: "task:delegate", Description: "Delegate tasks", Resource: "task", Action: "delegate"},
		{Name: "analytics:view", Description: "View analytics", Resource: "analytics", Action: "view"},
		{Name: "analytics:export", Description: "Export analytics data", Resource: "analytics", Action: "export"},
		{Name: "user:manage", Description: "Manage users", Resource: "user", Action: "manage"},
		{Name: "system:config", Description: "Configure system", Resource: "system", Action: "config"},
		{Name: "role:manage", Description: "Manage roles", Resource: "role", Action: "manage"},
	}

	for _, perm := range permissions {
		if err := db.FirstOrCreate(&perm, models.Permission{Name: perm.Name}).Error; err != nil {
			return err
		}
	}

	// Create default roles
	roles := []models.Role{
		{Name: "admin", Description: "System administrator", IsSystem: true},
		{Name: "manager", Description: "Department manager", IsSystem: true},
		{Name: "user", Description: "Regular user", IsSystem: true},
		{Name: "process-designer", Description: "Process designer", IsSystem: true},
		{Name: "analytics-viewer", Description: "Analytics viewer", IsSystem: true},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

func migration002Down(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.Role{},
		&models.Permission{},
	)
}

// migration003Up - Process engine
func migration003Up(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.ProcessDefinition{},
		&models.ProcessInstance{},
		&models.TaskInstance{},
		&models.BusinessRule{},
	)
}

func migration003Down(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.BusinessRule{},
		&models.TaskInstance{},
		&models.ProcessInstance{},
		&models.ProcessDefinition{},
	)
}

// migration004Up - Audit system
func migration004Up(db *gorm.DB) error {
	return db.AutoMigrate(&models.AuditLog{})
}

func migration004Down(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.AuditLog{})
}

// migration005Up - Performance indexes
func migration005Up(db *gorm.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)",
		"CREATE INDEX IF NOT EXISTS idx_process_instances_status ON process_instances(status)",
		"CREATE INDEX IF NOT EXISTS idx_process_instances_started_at ON process_instances(started_at)",
		"CREATE INDEX IF NOT EXISTS idx_task_instances_status ON task_instances(status)",
		"CREATE INDEX IF NOT EXISTS idx_task_instances_assignee_id ON task_instances(assignee_id)",
		"CREATE INDEX IF NOT EXISTS idx_task_instances_due_date ON task_instances(due_date)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource)",
		"CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at)",
		"CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func migration005Down(db *gorm.DB) error {
	indexes := []string{
		"DROP INDEX IF EXISTS idx_users_email",
		"DROP INDEX IF EXISTS idx_users_is_active",
		"DROP INDEX IF EXISTS idx_process_instances_status",
		"DROP INDEX IF EXISTS idx_process_instances_started_at",
		"DROP INDEX IF EXISTS idx_task_instances_status",
		"DROP INDEX IF EXISTS idx_task_instances_assignee_id",
		"DROP INDEX IF EXISTS idx_task_instances_due_date",
		"DROP INDEX IF EXISTS idx_audit_logs_timestamp",
		"DROP INDEX IF EXISTS idx_audit_logs_user_id",
		"DROP INDEX IF EXISTS idx_audit_logs_resource",
		"DROP INDEX IF EXISTS idx_refresh_tokens_expires_at",
		"DROP INDEX IF EXISTS idx_refresh_tokens_user_id",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			return err
		}
	}

	return nil
}
