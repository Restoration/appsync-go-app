```mermaid
sequenceDiagram
    participant User as User
    participant AppSync as AppSync
    participant GoAPI as GoAPI
    participant PostgreSQL as PostgreSQL

    alt Fetch Messages
        User->>AppSync: Request to fetch messages
        AppSync->>GoAPI: Fetch messages request
        GoAPI->>PostgreSQL: SELECT query for messages
        PostgreSQL-->>GoAPI: Return Messages Data
        GoAPI-->>AppSync: Respond with messages
        AppSync-->>User: Display fetched messages
    else Send Message
        User->>AppSync: Request to send a new message
        AppSync->>GoAPI: Send message request
        GoAPI->>PostgreSQL: INSERT query with new message
        PostgreSQL-->>GoAPI: Confirm Inserted
        GoAPI-->>AppSync: Respond with confirmation
        AppSync-->>User: Display message sent confirmation
    end
