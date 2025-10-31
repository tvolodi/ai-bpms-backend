# HR Module Architecture Recommendation

## Executive Summary

**Recommendation: HR as a DOMAIN SERVICE (not core)**

**Reasoning**: HR has independent data models, processes, and access patterns separate from process management. Placing it in core would violate single responsibility and create unnecessary coupling.

**Architecture Position**: 
```
├── core-services/           (Process, Users, Auth)
└── domain-services/
    ├── hr-service/          ✅ RECOMMENDED
    ├── analytics-service/
    ├── notification-service/
    └── ai-service/
```

---

## Analysis: Core vs. Domain Service

### Option 1: HR in CORE (❌ NOT RECOMMENDED)

**What would happen:**
```
core-services/
├── process-engine/
│   ├── process management
│   ├── workflow execution
│   ├── task management
│   ├── HR functions (employees, payroll, leave)  ← WRONG PLACE
│   └── ...
└── user-management/
```

**Problems:**
```
1. SINGLE RESPONSIBILITY VIOLATION
   - Process Engine focused on workflow
   - HR focused on employee lifecycle
   - These are different domains!

2. DATA MODEL BLOAT
   - Process tables: processes, instances, tasks
   - HR tables: employees, departments, salaries, leaves, training
   - No overlap = separate concerns

3. ACCESS PATTERN MISMATCH
   - Process Engine: Real-time workflow (high throughput, low latency)
   - HR Module: Periodic batching (payroll, analytics, reporting)
   - Different performance requirements

4. TEAM SCALING ISSUES
   - Process team works on workflow logic
   - HR team needs to work on HR logic independently
   - Merging into one service creates conflicts

5. DEPLOYMENT COUPLING
   - Bug in Leave Management shouldn't restart process engine
   - Bug in Workflow shouldn't affect HR operations
   - Independent scaling requirements

6. TESTING COMPLEXITY
   - Process tests mixed with HR tests
   - Harder to maintain, longer test suites
   - More points of failure

7. DATABASE CONCERNS
   - Process queries optimized for workflow (time-series)
   - HR queries optimized for employee records (OLTP)
   - Conflicting optimization strategies
```

### Option 2: HR as DOMAIN SERVICE (✅ RECOMMENDED)

**What would happen:**
```
core-services/                    domain-services/
├── process-engine/               ├── hr-service/
│   ├── workflow                   │   ├── employee management
│   ├── tasks                      │   ├── payroll
│   └── rules                      │   ├── leave management
│                                  │   ├── performance mgmt
└── user-management/              │   └── training & dev
    ├── auth                       │
    ├── rbac                       ├── analytics-service/
    └── users                      ├── notification-service/
                                   └── ai-service/
```

**Benefits:**

```
1. ✅ CLEAN SEPARATION OF CONCERNS
   - HR Service owns: Employee lifecycle, payroll, benefits, leave
   - Process Engine owns: Workflow execution, task management
   - Clear boundaries, independent evolution

2. ✅ INDEPENDENT SCALING
   - HR Service: Batch processing optimized
   - Process Engine: Real-time throughput optimized
   - Scale each based on actual load patterns

3. ✅ TEAM AUTONOMY
   - HR team can develop independently
   - Process team can innovate freely
   - Minimal coordination needed

4. ✅ INDEPENDENT DEPLOYMENT
   - Deploy HR updates without touching process engine
   - Faster release cycles for each domain
   - Lower risk, isolated rollbacks

5. ✅ FOCUSED DATA MODELS
   - HR Database: Employees, salaries, leave, benefits
   - Process Database: Processes, instances, tasks, rules
   - Optimize each database separately

6. ✅ TECHNOLOGY FLEXIBILITY
   - Process Engine might use PostgreSQL + Redis
   - HR Service might use PostgreSQL + specialized payroll DB
   - Pick best tools per domain

7. ✅ TESTING ISOLATION
   - HR unit tests independent of process tests
   - Faster test execution
   - Clearer test organization

8. ✅ REUSABILITY
   - HR Service APIs usable by other systems
   - Mobile apps can call HR APIs directly
   - Vendor integrations (Workday, BambooHR) easier
```

---

## HR Module Scope Definition

### HR Core Domains

#### 1. **Employee Management**
```
employee/
├── models/
│   ├── employee.go         # Employee profile
│   ├── department.go       # Department/team
│   ├── job_title.go        # Role definitions
│   └── contract.go         # Employment contract
└── handlers/
    ├── employee_handler.go
    ├── department_handler.go
    └── contract_handler.go

Key Operations:
- Create/update employee records
- Manage departments
- Track employment dates
- Store contact information
```

#### 2. **Leave Management**
```
leave/
├── models/
│   ├── leave_type.go       # Annual, sick, unpaid, etc.
│   ├── leave_request.go    # Request for time off
│   ├── leave_balance.go    # Remaining leave
│   └── leave_rule.go       # Policy rules
└── handlers/
    ├── leave_request_handler.go
    ├── leave_approval_handler.go
    └── balance_handler.go

Key Operations:
- Request leave (with approval workflow)
- Track leave balances
- Enforce leave policies
- Generate leave reports
```

#### 3. **Payroll Management**
```
payroll/
├── models/
│   ├── salary.go           # Base salary
│   ├── allowances.go       # Bonuses, HRA, etc.
│   ├── deductions.go       # Taxes, loans, etc.
│   ├── payslip.go          # Monthly payslip
│   └── payment.go          # Payment records
└── handlers/
    ├── salary_handler.go
    ├── payslip_handler.go
    ├── payment_handler.go
    └── report_handler.go

Key Operations:
- Manage salary components
- Calculate monthly payroll
- Generate payslips
- Track payments
- Tax compliance
```

#### 4. **Performance Management**
```
performance/
├── models/
│   ├── appraisal.go        # Annual appraisal
│   ├── goal.go             # OKRs/Goals
│   ├── rating.go           # Performance rating
│   └── feedback.go         # 360 feedback
└── handlers/
    ├── appraisal_handler.go
    ├── goal_handler.go
    ├── rating_handler.go
    └── feedback_handler.go

Key Operations:
- Create appraisals
- Set performance goals
- Collect feedback
- Generate performance reports
```

#### 5. **Training & Development**
```
training/
├── models/
│   ├── course.go           # Training courses
│   ├── enrollment.go       # Employee enrollment
│   ├── certification.go    # Certifications
│   └── skill.go            # Skills tracking
└── handlers/
    ├── course_handler.go
    ├── enrollment_handler.go
    ├── cert_handler.go
    └── skill_handler.go

Key Operations:
- Manage training programs
- Track certifications
- Skill development
- Training effectiveness
```

#### 6. **Attendance & Compliance**
```
attendance/
├── models/
│   ├── attendance.go       # Daily attendance
│   ├── shift.go            # Shift definitions
│   ├── holiday.go          # Public holidays
│   └── compliance.go       # Regulatory compliance
└── handlers/
    ├── attendance_handler.go
    ├── shift_handler.go
    ├── holiday_handler.go
    └── compliance_handler.go

Key Operations:
- Track attendance
- Manage shifts
- Calculate overtime
- Generate compliance reports
```

---

## Integration Points Between Core & HR

### How HR Service Integrates with Process Engine

#### 1. **User Management Sync**
```
Process Engine                   HR Service
     ↓ (watches)                    ↓
   user-management          hr-service/employees
     ↓                              ↓
   Event: user.created             ↓
     ←──────────────────────────────←
          Subscribe: UserCreated Event
          
   Action: Create employee record automatically
```

#### 2. **Leave Request Workflow**
```
HR Service (leave request)  →  Process Engine
       ↓                              ↓
   leave_request.create      →   Trigger Workflow
       ↓                              ↓
   "Leave Approval Process"  →   process_instance.create
       ↓                              ↓
   Assign to manager task   ←   Create task
       ↓                              ↓
   Manager approves          →   task.complete
       ↓                              ↓
   Update leave balance      ←   Webhook callback
```

#### 3. **Employee Lifecycle Events**
```
HR Service emits events:
- employee.created
- employee.promoted
- employee.transferred
- employee.resigned
- salary.updated

Process Engine listens to:
- Trigger onboarding workflows
- Update role-based permissions
- Manage access controls
- Track organizational changes
```

---

## Architecture Diagram

### Service Dependency Map

```
┌─────────────────────────────────────────────────────────────────┐
│                         API Gateway                              │
│                    (Single Entry Point)                          │
└────────┬─────────────────────────────────────────────────────────┘
         │
         ├──────────────────────────────────────────────┐
         │                                              │
    ┌────▼─────────────────────┐          ┌────────────▼──────────┐
    │   CORE SERVICES          │          │  DOMAIN SERVICES      │
    ├──────────────────────────┤          ├───────────────────────┤
    │                          │          │                       │
    │ Process Engine           │          │ HR Service            │
    │ ├─ workflow execution    │          │ ├─ employees          │
    │ ├─ task management       │          │ ├─ payroll            │
    │ └─ rule engine           │          │ ├─ leave              │
    │                          │          │ ├─ performance        │
    │ User Management          │          │ ├─ training           │
    │ ├─ authentication        │          │ └─ attendance         │
    │ ├─ rbac                  │          │                       │
    │ └─ users                 │          │ Analytics Service     │
    │                          │          │ ├─ dashboards         │
    └──┬───────────────┬───────┘          │ └─ reports            │
       │               │                   │                       │
       │               │                   │ Notification Service  │
       │               │                   │ ├─ emails             │
       │               │                   │ └─ webhooks           │
       │               │                   │                       │
       │               │                   │ AI Service            │
       │               │                   │ ├─ predictions        │
       │               │                   │ └─ recommendations    │
       │               │                   │                       │
       └───────┬───────┴─────────────────┬─┘
               │                         │
         ┌─────▼─────────┐      ┌────────▼──────────┐
         │  Shared Libs  │      │   Shared Libs     │
         ├───────────────┤      ├───────────────────┤
         │ • Models      │      │ • Database        │
         │ • Cache       │      │ • Cache           │
         │ • Config      │      │ • Validation      │
         │ • Security    │      │ • Error handling  │
         └─────────────┬─┘      └────────┬──────────┘
                       │                 │
         ┌─────────────▼─────────────────▼──────────┐
         │         Event Bus (NATS)                  │
         │  ├─ Async communication                  │
         │  ├─ Event streaming                      │
         │  └─ Decoupled services                   │
         └──────────────────────────────────────────┘
                       │
         ┌─────────────▼──────────────────────────┐
         │      Infrastructure & Data             │
         │  ├─ PostgreSQL (shared)                │
         │  ├─ Redis (caching)                    │
         │  ├─ Prometheus (monitoring)            │
         │  └─ Grafana (dashboards)               │
         └───────────────────────────────────────┘
```

### Communication Flow Example: Leave Request

```
1. User requests leave via HR UI
   ↓
   POST /api/v1/leave-requests
   │
   └─→ HR Service validates
       ├─ Check leave balance
       ├─ Validate dates
       └─ Create leave_request record
           ↓
           Event: leave_request.created
                 {employee_id, start_date, end_date, reason}

2. Event published to NATS Event Bus
   ↓
   Process Engine listens to "leave_request.created"
   ├─ Trigger "Leave Approval" process
   ├─ Create process_instance
   └─ Create approval task for manager

3. Manager receives task notification
   ├─ Notification Service (subscribed to task.created)
   ├─ Send email/push notification
   └─ Manager reviews and approves

4. Manager completes task
   ↓
   Process Engine: task.completed event
   ↓
   HR Service listens to "leave_request.approved" workflow event
   ├─ Update leave_request status
   ├─ Deduct from leave balance
   └─ Send confirmation to employee

5. Async operations (non-blocking)
   ├─ Analytics Service: record leave statistics
   ├─ Notification Service: send confirmation email
   └─ Audit Service: log transaction
```

---

## Directory Structure for HR Service

```
domain-services/
└── hr-service/
    ├── cmd/
    │   └── main.go                    # Service entrypoint
    │
    ├── internal/
    │   ├── handlers/                  # HTTP handlers
    │   │   ├── employee_handler.go
    │   │   ├── leave_handler.go
    │   │   ├── payroll_handler.go
    │   │   ├── performance_handler.go
    │   │   ├── training_handler.go
    │   │   └── attendance_handler.go
    │   │
    │   ├── services/                  # Business logic
    │   │   ├── employee_service.go
    │   │   ├── leave_service.go
    │   │   ├── payroll_service.go
    │   │   ├── performance_service.go
    │   │   ├── training_service.go
    │   │   └── attendance_service.go
    │   │
    │   ├── repositories/              # Data access
    │   │   ├── employee_repo.go
    │   │   ├── leave_repo.go
    │   │   ├── payroll_repo.go
    │   │   ├── performance_repo.go
    │   │   ├── training_repo.go
    │   │   └── attendance_repo.go
    │   │
    │   ├── models/                    # Domain models
    │   │   ├── employee.go
    │   │   ├── leave.go
    │   │   ├── payroll.go
    │   │   ├── performance.go
    │   │   ├── training.go
    │   │   ├── attendance.go
    │   │   └── event_models.go        # Events published
    │   │
    │   ├── events/                    # Event handling
    │   │   ├── publisher.go           # Publish events
    │   │   ├── subscribers.go         # Subscribe to events
    │   │   └── event_handlers.go      # Process events
    │   │
    │   ├── middleware/                # HR-specific middleware
    │   │   ├── validation.go
    │   │   ├── authorization.go
    │   │   └── audit_logging.go
    │   │
    │   ├── cache/                     # HR caching strategies
    │   │   ├── employee_cache.go
    │   │   └── leave_balance_cache.go
    │   │
    │   └── config/                    # Service configuration
    │       └── config.go
    │
    ├── pkg/                           # Public APIs
    │   ├── dto/                       # Data transfer objects
    │   │   ├── employee_dto.go
    │   │   ├── leave_dto.go
    │   │   └── ...
    │   └── errors/                    # HR-specific errors
    │       └── errors.go
    │
    ├── migrations/                    # Database migrations
    │   ├── 001_create_employees.sql
    │   ├── 002_create_leave_tables.sql
    │   ├── 003_create_payroll.sql
    │   └── ...
    │
    ├── tests/                         # Integration tests
    │   ├── employee_test.go
    │   ├── leave_test.go
    │   └── ...
    │
    ├── config/
    │   ├── config.yaml
    │   └── config.local.yaml
    │
    └── README.md
```

---

## Deployment Architecture

### Independent Scaling

```
Load Balancer
│
├── Process Engine (2-3 instances)   ← High throughput
├── HR Service (1-2 instances)       ← Periodic batches
├── User Management (1 instance)     ← Auth-critical
└── API Gateway (2-4 instances)      ← Frontend

Databases:
├── Shared: user (auth required)
├── Process DB: processes, tasks, instances, rules
└── HR DB: employees, payroll, leave, performance
```

### Independent Deployment

```
Pipeline 1: Process Engine
├── Test changes
├── Deploy to staging
└── Production (without affecting HR)

Pipeline 2: HR Service
├── Test changes
├── Deploy to staging
└── Production (without affecting Process Engine)

Coordination:
├── Event schema versioning (NATS contracts)
└── API contract versioning (OpenAPI)
```

---

## Data Ownership & Consistency

### Who Owns What?

```
User Management owns:
├── users table
├── roles table
├── permissions table
└── auth credentials

Process Engine owns:
├── processes table
├── process_instances table
├── tasks table
├── task_assignments table
└── process_rules table

HR Service owns:
├── employees table (references users table)
├── departments table
├── job_titles table
├── salaries table
├── leave_requests table
├── leave_balances table
├── payslips table
├── appraisals table
├── training_enrollments table
└── attendance table
```

### Event-Based Synchronization

```
When employee is created:

1. HR Service creates employee record
2. Publishes: employee.created event
3. Process Engine listens:
   ├─ Creates process_context
   ├─ Trigger onboarding workflow
   └─ Assign initial tasks

Benefits:
✅ Loose coupling
✅ Eventual consistency
✅ Can retry failed events
✅ Each service stays independent
```

---

## Integration Scenarios

### Scenario 1: Employee Onboarding
```
HR receives new hire notification
  ↓
HR Service: Create employee record
  ↓
Event: employee.onboarding_started
  ↓
Process Engine: Trigger Onboarding Workflow
  ├─ Task: Complete IT setup
  ├─ Task: Complete HR paperwork
  ├─ Task: Assign to team
  └─ Task: Schedule training
  ↓
Workflow completion triggers HR callback
  ↓
HR Service: Mark onboarding complete
```

### Scenario 2: Leave Approval Workflow
```
Employee requests leave (HR UI)
  ↓
HR Service: Validate & create leave_request
  ↓
Event: leave.approval_needed
  ↓
Process Engine: Trigger Leave Approval Workflow
  ├─ Task: Manager review
  ├─ Task: HR compliance check
  └─ Task: Final approval
  ↓
Workflow completes
  ↓
HR Service: Update leave status & balance
```

### Scenario 3: Salary Review Workflow
```
HR initiates salary review (bulk operation)
  ↓
HR Service: Create salary_review record
  ↓
Event: salary.review_started
  ↓
Process Engine: Trigger Review Workflow for each employee
  ├─ Task: Manager recommendation
  ├─ Task: Finance approval
  └─ Task: HR final decision
  ↓
Each workflow completion triggers salary update
  ↓
HR Service: Update salaries & generate new payslips
```

---

## Technology Choices for HR Service

### Database Strategy

```
Primary: PostgreSQL (same as core)
Rationale:
├─ Consistent technology stack
├─ Shared connection pool management
├─ Familiar team expertise
└─ ACID compliance for critical HR data

Specialized considerations:
├─ Separate HR schema from Process schema
├─ Dedicated indexes for payroll queries
├─ Time-series data for analytics
└─ Archive old records separately
```

### Caching Strategy

```
L1: Local Cache (BigCache)
├─ Employee directory (frequently accessed)
├─ Leave policies (rarely changes)
└─ Salary structures (stable)

L2: Redis Distributed Cache
├─ Leave balance cache (per employee)
├─ Payroll calculations (temporary)
└─ Employee search results

Invalidation:
├─ Time-based (TTL)
├─ Event-based (when HR data changes)
└─ Manual (admin triggers)
```

### Background Processing

```
Async Jobs (NATS):
├─ Monthly payroll calculation
├─ Leave balance accrual
├─ Compliance report generation
├─ Performance review reminders
└─ Training expiration notices

Scheduled Tasks (Cron/Scheduler):
├─ Daily attendance sync
├─ Monthly financial reconciliation
├─ Quarterly performance reviews
└─ Annual compliance audits
```

---

## API Contract Between Core & HR

### Events Published by HR Service

```go
// Leave request created
{
  "event_type": "leave.requested",
  "employee_id": "emp_123",
  "start_date": "2025-11-01",
  "end_date": "2025-11-05",
  "reason": "vacation",
  "created_at": "2025-10-31T10:00:00Z"
}

// Salary updated
{
  "event_type": "salary.updated",
  "employee_id": "emp_123",
  "effective_date": "2025-11-01",
  "base_salary": 50000,
  "updated_by": "user_456",
  "timestamp": "2025-10-31T09:00:00Z"
}

// Employee lifecycle
{
  "event_type": "employee.lifecycle",
  "event_subtype": "promotion",
  "employee_id": "emp_123",
  "old_role": "Senior Developer",
  "new_role": "Tech Lead",
  "effective_date": "2025-11-01"
}
```

### Events Consumed from Process Engine

```go
// Workflow completed (for leave approval)
{
  "event_type": "process.workflow_completed",
  "workflow_id": "wf_leave_approval",
  "instance_id": "inst_789",
  "result": "approved",
  "context": {
    "leave_request_id": "leave_123",
    "employee_id": "emp_123"
  }
}

// Task assigned (for manager approval)
{
  "event_type": "process.task_assigned",
  "task_id": "task_456",
  "workflow_id": "wf_leave_approval",
  "assigned_to": "user_manager",
  "context": {
    "leave_request_id": "leave_123"
  }
}
```

---

## Risk Mitigation

### Risks of Separate HR Service

```
Risk 1: Data Consistency
├─ Mitigation: Event-sourcing, saga patterns
└─ Fallback: Eventual consistency with retry logic

Risk 2: Performance Impact
├─ Mitigation: Caching strategies, async processing
└─ Fallback: Batch operations during off-peak hours

Risk 3: Complexity
├─ Mitigation: Clear API contracts, comprehensive testing
└─ Fallback: Well-documented integration points

Risk 4: Team Coordination
├─ Mitigation: Interface contracts (OpenAPI + Events)
└─ Fallback: Regular sync meetings, change logs
```

---

## Migration Path

### Phase 1: Foundation (Weeks 1-2)
- Define HR data models
- Set up HR database schema
- Create employee & leave endpoints
- Test with mock process engine

### Phase 2: Integration (Weeks 3-4)
- Implement event publisher
- Subscribe to process engine events
- Create leave approval workflow
- End-to-end testing

### Phase 3: Expansion (Weeks 5-8)
- Add payroll management
- Add performance management
- Add training module
- Add attendance tracking

### Phase 4: Optimization (Weeks 9+)
- Performance tuning
- Advanced reporting
- AI integration
- Mobile app support

---

## Conclusion

**Why HR Should Be a Domain Service:**

1. ✅ **Different Concerns**: HR lifecycle ≠ Process workflow
2. ✅ **Independent Evolution**: HR changes shouldn't affect process engine
3. ✅ **Team Autonomy**: Separate teams can work independently
4. ✅ **Scalability**: Different performance profiles
5. ✅ **Reusability**: HR APIs can be used by other systems
6. ✅ **Testing**: Cleaner, faster test isolation
7. ✅ **Deployment**: Independent release cycles

**Architecture Decision**: 
```
domain-services/hr-service/ ← RIGHT PLACE
```

This follows the documented performance-optimized architecture where HR is a specialized domain service that integrates with core services through well-defined event contracts.
