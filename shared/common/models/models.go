package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common columns for all tables
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User represents a user in the system
type User struct {
	BaseModel
	Email           string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash    string     `gorm:"column:password_hash" json:"-"`
	FirstName       string     `gorm:"size:100" json:"first_name"`
	LastName        string     `gorm:"size:100" json:"last_name"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	IsEmailVerified bool       `gorm:"default:false" json:"is_email_verified"`
	LastLogin       *time.Time `json:"last_login"`

	// BPMS-specific fields
	Roles         []Role   `gorm:"many2many:user_roles;" json:"roles"`
	ProcessGroups []string `gorm:"type:text[]" json:"process_groups"`
	Department    string   `gorm:"size:100" json:"department"`
	Position      string   `gorm:"size:100" json:"position"`

	// Authentication settings
	MFAEnabled bool   `gorm:"default:false" json:"mfa_enabled"`
	MFASecret  string `gorm:"column:mfa_secret" json:"-"`

	// Audit fields
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// Role represents a role in the RBAC system
type Role struct {
	BaseModel
	Name        string       `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Description string       `gorm:"size:255" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	IsSystem    bool         `gorm:"default:false" json:"is_system"`
}

// Permission represents a permission in the RBAC system
type Permission struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Resource    string `gorm:"size:50" json:"resource"` // process, task, user, etc.
	Action      string `gorm:"size:50" json:"action"`   // create, read, update, delete
}

// ProcessDefinition represents a business process definition
type ProcessDefinition struct {
	BaseModel
	Name        string `gorm:"not null;size:255" json:"name"`
	Key         string `gorm:"uniqueIndex;not null;size:100" json:"key"`
	Version     int    `gorm:"not null;default:1" json:"version"`
	Description string `gorm:"type:text" json:"description"`

	// Process definition data
	BPMN       string `gorm:"type:text" json:"bpmn"`         // BPMN XML
	JSONSchema string `gorm:"type:jsonb" json:"json_schema"` // JSON Schema for forms
	Variables  string `gorm:"type:jsonb" json:"variables"`   // Process variables

	// Process settings
	IsActive bool     `gorm:"default:true" json:"is_active"`
	Category string   `gorm:"size:100" json:"category"`
	Tags     []string `gorm:"type:text[]" json:"tags"`

	// AI enhancement
	AIEnabled bool   `gorm:"default:false" json:"ai_enabled"`
	AIConfig  string `gorm:"type:jsonb" json:"ai_config"`

	// Relationships
	Instances []ProcessInstance `gorm:"foreignKey:ProcessDefinitionID" json:"-"`

	// Audit fields
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// ProcessInstance represents a running instance of a process
type ProcessInstance struct {
	BaseModel
	ProcessDefinitionID uuid.UUID         `gorm:"type:uuid;not null" json:"process_definition_id"`
	ProcessDefinition   ProcessDefinition `gorm:"foreignKey:ProcessDefinitionID" json:"process_definition"`

	BusinessKey string `gorm:"size:255" json:"business_key"`
	Status      string `gorm:"size:50;not null" json:"status"` // active, completed, suspended, terminated

	// Instance data
	Variables string `gorm:"type:jsonb" json:"variables"`
	Context   string `gorm:"type:jsonb" json:"context"`

	// Timing
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	Duration  *int64     `json:"duration"` // milliseconds

	// Relationships
	Tasks []TaskInstance `gorm:"foreignKey:ProcessInstanceID" json:"tasks"`

	// Audit fields
	StartedBy *uuid.UUID `gorm:"type:uuid" json:"started_by"`
	EndedBy   *uuid.UUID `gorm:"type:uuid" json:"ended_by"`
}

// TaskInstance represents a task within a process instance
type TaskInstance struct {
	BaseModel
	ProcessInstanceID uuid.UUID       `gorm:"type:uuid;not null" json:"process_instance_id"`
	ProcessInstance   ProcessInstance `gorm:"foreignKey:ProcessInstanceID" json:"-"`

	TaskDefinitionKey string `gorm:"size:100;not null" json:"task_definition_key"`
	Name              string `gorm:"size:255" json:"name"`
	Description       string `gorm:"type:text" json:"description"`

	// Task assignment
	AssigneeID     *uuid.UUID `gorm:"type:uuid" json:"assignee_id"`
	Assignee       *User      `gorm:"foreignKey:AssigneeID" json:"assignee"`
	CandidateGroup string     `gorm:"size:100" json:"candidate_group"`

	// Task state
	Status       string     `gorm:"size:50;not null" json:"status"` // created, assigned, completed, cancelled
	Priority     int        `gorm:"default:50" json:"priority"`
	DueDate      *time.Time `json:"due_date"`
	FollowUpDate *time.Time `json:"follow_up_date"`

	// Task data
	FormData  string `gorm:"type:jsonb" json:"form_data"`
	Variables string `gorm:"type:jsonb" json:"variables"`

	// Timing
	CreatedAt   time.Time  `json:"created_at"`
	AssignedAt  *time.Time `json:"assigned_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Duration    *int64     `json:"duration"` // milliseconds

	// Audit fields
	AssignedBy  *uuid.UUID `gorm:"type:uuid" json:"assigned_by"`
	CompletedBy *uuid.UUID `gorm:"type:uuid" json:"completed_by"`
}

// BusinessRule represents a business rule
type BusinessRule struct {
	BaseModel
	Name        string `gorm:"not null;size:255" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Rule definition
	Expression string `gorm:"type:text;not null" json:"expression"` // expr language
	Language   string `gorm:"size:50;default:'expr'" json:"language"`

	// Rule metadata
	Category string   `gorm:"size:100" json:"category"`
	Tags     []string `gorm:"type:text[]" json:"tags"`
	Version  int      `gorm:"not null;default:1" json:"version"`
	IsActive bool     `gorm:"default:true" json:"is_active"`

	// AI enhancement
	AIGenerated bool   `gorm:"default:false" json:"ai_generated"`
	AIMetadata  string `gorm:"type:jsonb" json:"ai_metadata"`

	// Audit fields
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// FormSchema represents a dynamic form schema
type FormSchema struct {
	BaseModel
	Name        string `gorm:"not null;size:255" json:"name"`
	Key         string `gorm:"uniqueIndex;not null;size:100" json:"key"`
	Version     int    `gorm:"not null;default:1" json:"version"`
	Description string `gorm:"type:text" json:"description"`

	// Schema definition
	JSONSchema string `gorm:"type:jsonb;not null" json:"json_schema"`
	UISchema   string `gorm:"type:jsonb" json:"ui_schema"`

	// Schema metadata
	Category string   `gorm:"size:100" json:"category"`
	Tags     []string `gorm:"type:text[]" json:"tags"`
	IsActive bool     `gorm:"default:true" json:"is_active"`

	// AI enhancement
	AIGenerated bool `gorm:"default:false" json:"ai_generated"`

	// Audit fields
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`

	// Who did what
	UserID     *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User       *User      `gorm:"foreignKey:UserID" json:"user"`
	Action     string     `gorm:"size:100;not null" json:"action"`
	Resource   string     `gorm:"size:100;not null" json:"resource"`
	ResourceID *uuid.UUID `gorm:"type:uuid" json:"resource_id"`

	// Details
	Details   string `gorm:"type:jsonb" json:"details"`
	IPAddress string `gorm:"size:45" json:"ip_address"`
	UserAgent string `gorm:"size:500" json:"user_agent"`

	// Result
	Success      bool   `gorm:"not null" json:"success"`
	ErrorMessage string `gorm:"type:text" json:"error_message"`
}

// RefreshToken represents a JWT refresh token
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Token     string     `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	IPAddress string     `gorm:"size:45" json:"ip_address"`
	UserAgent string     `gorm:"size:500" json:"user_agent"`
	IsRevoked bool       `gorm:"default:false" json:"is_revoked"`
	RevokedAt *time.Time `json:"revoked_at"`
}
