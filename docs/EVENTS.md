# Event Catalog

## Event: user.registered

- version: 1
- source: auth-service
- consumer: audit-service
- purpose: record user registration

## Event: note.created

- version: 1
- source: notes-service
- consumer: audit-service
- purpose: record note creation

## Event: note.updated

- version: 1
- source: notes-service
- consumer: audit-service
- purpose: record note modification

## Event: note.deleted

- version: 1
- source: notes-service
- consumer: audit-service
- purpose: record note deletion
