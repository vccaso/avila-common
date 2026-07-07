package event

// EventMetadata defines structured context for UI-friendly event rendering and filtering.
type EventMetadata struct {
	Type        string `json:"type"`        // Unique machine name (e.g. "contact.created")
	Label       string `json:"label"`       // UI friendly label (e.g. "Contact Created")
	Service     string `json:"service"`     // Microservice prefix context (e.g. "contact")
	Description string `json:"description"` // Tooltip context details
	Icon        string `json:"icon"`        // Icon descriptor
}

// SystemEventsRegistry holds the metadata for all valid transactional outbox events in the EazySales ecosystem.
var SystemEventsRegistry = []EventMetadata{
	{
		Type:        "contact.created",
		Label:       "Contact Created",
		Service:     "contact",
		Description: "Triggered when a new contact is registered in the CRM system.",
		Icon:        "user-plus",
	},
	{
		Type:        "contact.updated",
		Label:       "Contact Updated",
		Service:     "contact",
		Description: "Triggered when any details or fields of a contact are updated.",
		Icon:        "user-check",
	},
	{
		Type:        "company.created",
		Label:       "Company Created",
		Service:     "contact",
		Description: "Triggered when a new business entity or company profile is created.",
		Icon:        "building-office",
	},
	{
		Type:        "company.updated",
		Label:       "Company Updated",
		Service:     "contact",
		Description: "Triggered when the core details or custom attributes of a company are updated.",
		Icon:        "building-office-2",
	},
	{
		Type:        "campaign.launched",
		Label:       "Campaign Launched",
		Service:     "contact",
		Description: "Triggered when a campaign run is created and recipients are enrolled.",
		Icon:        "megaphone",
	},
	{
		Type:        "campaign.recipient.queued",
		Label:       "Campaign Recipient Queued",
		Service:     "contact",
		Description: "Triggered when a campaign recipient is ready to be dispatched to an AI Agent execution.",
		Icon:        "send",
	},
	{
		Type:        "campaign.recipient.reply_received",
		Label:       "Campaign Recipient Reply Received",
		Service:     "contact",
		Description: "Triggered when an inbound email reply maps back to a campaign recipient.",
		Icon:        "reply",
	},
	{
		Type:        "agent.execution.completed",
		Label:       "Agent Execution Completed",
		Service:     "agent",
		Description: "Triggered when an AI Agent finishes its reasoning loop and successfully outputs a payload.",
		Icon:        "cpu",
	},
	{
		Type:        "agent.execution.started",
		Label:       "Agent Execution Started",
		Service:     "agent",
		Description: "Triggered when an AI Agent begins its reasoning run and task loop.",
		Icon:        "play",
	},
	{
		Type:        "agent.triggered",
		Label:       "Agent Triggered",
		Service:     "agent",
		Description: "Triggered when an AI Agent execution is initially dispatched.",
		Icon:        "bolt",
	},
	{
		Type:        "email.received",
		Label:       "Email Received",
		Service:     "notification",
		Description: "Triggered when a support or ticket-creation email is received from the queue.",
		Icon:        "envelope-open",
	},
	{
		Type:        "email.sent",
		Label:       "Email Sent",
		Service:     "notification",
		Description: "Triggered when any system notification or automated email response is sent.",
		Icon:        "paper-airplane",
	},
	{
		Type:        "user.affiliate_admin.created",
		Label:       "Affiliate Admin Created",
		Service:     "user",
		Description: "Triggered when a user is registered with the Affiliate Admin role (ID: 13)",
		Icon:        "user-plus",
	},
	{
		Type:        "customer.created",
		Label:       "Customer Created",
		Service:     "user",
		Description: "Triggered when a new Customer/Tenant organization is successfully created",
		Icon:        "building",
	},
	{
		Type:        "chat.lead_tagged",
		Label:       "Chat Lead Tagged",
		Service:     "contact",
		Description: "Triggered when an EazyChat operator flags a chat visitor as a lead.",
		Icon:        "user-circle",
	},
}
