# Resource Manager Feature Roadmap

## File Loading & Initialization

**Goal**
Load contracts into memory and initialize app state

**Includes**

- Accept resource file path
- Accept API contract file path
- Parse JSON into internal models
- Initialize AppState
- Handle errors (invalid path, invalid JSON)

## Resource List View

**Goal**
Browse and select resources

**Includes**

- Display list of resources
- Navigation
- Open resource detail view

## Resource Add/Delete

**Goal**
Add and delete resources

**Includes**

- Add a resource
- Delete a resource
- Write to the resource contract file

## Resource Field Permission Editing

**Goal**
Allow users to modify field-level access rules for a resource table.

**Includes**

- Select editable permission cells ('read', 'write', 'mutable')
- Toggle permission values between 'true' and 'false'
- Persist field permission changes back to the resource contract file

## Resource Binding List

**Goal**  
Display API endpoints from the API contracts file and allow users to browse and inspect bindings between resources and endpoints.

**Includes**

- Load API contract data into application state
- Parse endpoints from OpenAPI / contracts file
- Display endpoints in the menu pane
- Support navigation through endpoints list
- Show endpoint details in the content pane:
  - HTTP method (GET, POST, PATCH, DELETE)
  - Route/path
  - Associated resource (via `x-resource` or equivalent metadata)
  - Operation type (readMany, readOne, create, update, delete)
