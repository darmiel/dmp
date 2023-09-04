package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

// User represents a user which can log in (obviously)
type User struct {
	// ID is generated by firebase (in the best case)
	ID string `gorm:"primarykey" json:"id"`
	// UserName is the username which is displayed in the frontend
	UserName string `gorm:"unique" json:"name"`
	// CreatedAt is the time when the entry was first created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the entry was modified
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is the time when the entry was (soft-) removed
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// OwnerProjects contains a list of all projects the user is owner of
	OwnerProjects []Project `gorm:"foreignKey:OwnerID" json:"owner_projects,omitempty"`
	// UserProjects contains a list of all projects the user has access to
	UserProjects []Project `gorm:"many2many:user_projects;" json:"user_projects,omitempty"`
	// AssignedTopics contains all topics the user was assigned to
	AssignedTopics []Topic `gorm:"many2many:user_topic_assignments" json:"assigned_users"`
	// AssignedActions contains all actions the user was assigned to
	AssignedActions []Action `gorm:"many2many:action_user_assignments" json:"assigned_actions"`
	// Actions contains all actions the user created
	CreatedActions []Action `gorm:"foreignKey:CreatorID" json:"actions"`
	// AssignedMeetings contains all meetings the user was assigned to
	AssignedMeetings []Meeting `gorm:"many2many:meeting_user_assignments" json:"assigned_meetings"`
}

// Comment represents a comment in a topic
type Comment struct {
	gorm.Model
	// AuthorID is the id of the author of the comment
	AuthorID string `json:"author_id"`
	// Author is the author of the comment
	Author User `json:"author,omitempty"`
	// Content is the comment content as markdown
	Content string `json:"content"`
	// TopicID is the ID of the topic the comment belongs to
	TopicID *uint `json:"topic_id"`
	// MeetingID is the ID of the meeting the comment belongs to
	MeetingID *uint `json:"meeting_id"`
	// ProjectID is the ID of the project the comment belongs to
	ProjectID *uint `json:"project_id"`
	// ActionID is the ID of the action the comment belongs to
	ActionID *uint `json:"action_id"`
}

// CheckProjectOwnership checks if the comment belongs to the project
func (c Comment) CheckProjectOwnership(projectID uint) bool {
	return c.ProjectID != nil && *c.ProjectID == projectID
}

// Topic represents a TO|DO-Point for the meeting
type Topic struct {
	gorm.Model
	// Title of the topic
	Title string `json:"title"`
	// Description of the topic
	Description string `json:"description"`
	// ID of the creator of the topic
	CreatorID string `json:"creator_id"`
	// Creator of the topic
	Creator User `json:"creator,omitempty"`
	// Comments in the topic
	Comments []Comment `json:"comments,omitempty"`
	// SolutionID is the ID of the comment which represents the solution of this topic
	SolutionID uint `json:"solution_id"`
	// Solution is the comment which represents the solution of this topic
	Solution Comment `json:"solution,omitempty"`
	// ClosedAt represents the time when the topic was resolved (if valid)
	ClosedAt sql.NullTime `json:"closed_at"`
	// ForceSolution requires a solution to be able to close topic if true
	ForceSolution bool `json:"force_solution"`
	// MeetingID is the ID of the meeting the topic belongs to
	MeetingID uint `json:"meeting_id"`
	// Meeting is the meeting the topic belongs to
	Meeting Meeting
	// AssignedUsers contains a list of users assigned to a topic
	AssignedUsers []User `gorm:"many2many:user_topic_assignments" json:"assigned_users"`
	// Actions contains a list of actions related to the topic
	Actions []Action `gorm:"many2many:topic_action_assignments" json:"actions"`
	// PriorityID is the ID of the priority of the topic
	PriorityID *uint `json:"priority_id"`
	// Priority is the priority of the topic
	Priority Priority `json:"priority,omitempty"`
	// Tags contains all tags of the topic
	Tags []Tag `gorm:"many2many:topic_tag_assignments" json:"tags"`
}

func (t Topic) CheckProjectOwnership(projectID uint) bool {
	// checking the Meeting ID for 0 is necessary because the meeting may not preloaded
	return t.Meeting.ID != 0 && t.Meeting.ProjectID == projectID
}

// Meeting represents a meeting (who would've guessed)
type Meeting struct {
	gorm.Model
	// Name of the meeting
	Name string `json:"name"`
	// Description of the meeting
	Description string `json:"description"`
	// StartDate of the meeting
	StartDate time.Time `json:"start_date"`
	// EndDate of the meeting
	EndDate time.Time `json:"end_date"`
	// Topics of the meeting
	Topics []Topic `json:"topics,omitempty"`
	// ProjectID is the project the meeting belongs to
	ProjectID uint `json:"project_id"`
	// ID of the creator of the meeting
	CreatorID string `json:"creator_id"`
	// Creator of the meeting
	Creator User `json:"creator,omitempty"`
	// Comments for the meeting
	Comments []Comment `json:"comments,omitempty"`
	// AssignedUsers contains all users assigned to the meeting
	AssignedUsers []User `gorm:"many2many:meeting_user_assignments" json:"assigned_users"`
	// Tags contains all tags of the meeting
	Tags []Tag `gorm:"many2many:meeting_tag_assignments" json:"tags"`
}

func (m Meeting) CheckProjectOwnership(projectID uint) bool {
	return m.ProjectID == projectID
}

// Project is a custom "realm" where meetings are saved
type Project struct {
	gorm.Model
	// Name is the name of the project and displayed in the frontend
	Name string `json:"name"`
	// Description of the project
	Description string `json:"description"`
	// PreviewURL is the display image
	PreviewURL string `json:"preview_url"`
	// OwnerID is the id of the creator of the project
	OwnerID string `json:"owner_id"`
	// Owner is the creator of the project
	Owner User `json:"owner,omitempty"`
	// Users contains all users which have access to the project
	Users []User `gorm:"many2many:user_projects;" json:"users,omitempty"`
	// Meetings contains all meetings in the project
	Meetings []Meeting `json:"meetings,omitempty"`
	// Priorities contains all priorities in the project
	Priorities []Priority `json:"priorities,omitempty"`
	// Tags contains all tags in the project
	Tags []Tag `json:"tags,omitempty"`
	// Actions contains all actions in the project
	Actions []Action `json:"actions,omitempty"`
	// Comments for the project
	Comments []Comment `json:"comments,omitempty"`
	// AIEnabled is true if the OpenAI API is enabled for the project
	AIEnabled bool `json:"ai_enabled"`
}

func (p Project) CheckProjectOwnership(projectID uint) bool {
	return p.ID == projectID
}

type Priority struct {
	gorm.Model
	// Title of the priority
	Title string `json:"title"`
	// Weight of the priority
	Weight int `json:"weight"`
	// Color of the priority
	Color string `json:"color"`
	// ProjectID is the ID of the project the priority belongs to
	ProjectID uint `json:"project_id"`
	// Project is the project the priority belongs to
	Project Project `json:"project,omitempty"`
}

func (p Priority) CheckProjectOwnership(projectID uint) bool {
	return p.ProjectID == projectID
}

type Action struct {
	gorm.Model
	// Title of the action
	Title string `json:"title"`
	// Description of the action
	Description string `json:"description"`
	// DueDate of the action (optional)
	DueDate sql.NullTime `json:"due_date"`
	// ProjectID is the ID of the project the action belongs to
	ProjectID uint `json:"project_id"`
	// Project is the project the action belongs to
	Project Project `json:"project,omitempty"`
	// Topics contains all topics the action is related to
	Topics []Topic `gorm:"many2many:action_topic_assignments" json:"topics"`
	// AssignedUsers contains all users assigned to the action
	AssignedUsers []User `gorm:"many2many:action_user_assignments" json:"assigned_users"`
	// PriorityID is the ID of the priority of the action
	PriorityID *uint `json:"priority_id"`
	// Priority is the priority of the action
	Priority Priority `json:"priority,omitempty"`
	// Tags contains all tags of the action
	Tags []Tag `gorm:"many2many:action_tag_assignments" json:"tags"`
	// ClosedAt represents the time when the action was resolved (if valid)
	ClosedAt sql.NullTime `json:"closed_at"`
	// CreatorID is the ID of the creator of the action
	CreatorID string `json:"creator_id"`
	// Comments for the action
	Comments []Comment `json:"comments,omitempty"`
}

func (a Action) CheckProjectOwnership(projectID uint) bool {
	return a.ProjectID == projectID
}

type Tag struct {
	gorm.Model
	// Title of the tag
	Title string `json:"title"`
	// Color of the tag
	Color string `json:"color"`
	// ProjectID is the ID of the project the tag belongs to
	ProjectID uint `json:"project_id"`
	// Project is the project the tag belongs to
	Project Project `json:"project,omitempty"`
	// Actions contains all actions the tag is related to
	Actions []Action `gorm:"many2many:action_tag_assignments" json:"actions"`
}

func (t Tag) CheckProjectOwnership(projectID uint) bool {
	return t.ProjectID == projectID
}

type Notification struct {
	gorm.Model
	// Title of the notification
	Title string `json:"title"`
	// Suffix of the notification, shown after the title in the notification view
	Suffix string `json:"suffix"`
	// Description of the notification
	Description string `json:"description"`
	// UserID is the ID of the user the notification belongs to
	UserID string `json:"user_id"`
	// User is the user the notification belongs to
	User User `json:"user,omitempty"`
	// ReadAt represents the time when the notification was read (if valid)
	ReadAt sql.NullTime `json:"read_at"`
	// Link is the link the notification points to (optional)
	Link string `json:"link"`
	// LinkTitle is the title of the link (optional)
	LinkTitle string `json:"link_title"`
}
