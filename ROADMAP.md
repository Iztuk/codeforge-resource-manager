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

## Resource Binding Management

### Goal

Enable users to define and manage **resource bindings** between API operations and underlying data resources using the `x-resource` extension.

This allows API endpoints to be explicitly linked to resource contracts, enabling validation, auditing, and enforcement of data access rules.

### Includes

#### Add Resource Binding

- Select an API endpoint (path + method)
- Launch binding workflow
- Choose resource to bind (`resourceName.tableName`)
- Attach binding to the operation via `x-resource`
- Persist changes to the API contract document

#### Remove Resource Binding

- Select an API endpoint (path + method)
- Remove existing `x-resource` binding
- Persist changes to the API contract document

### Design Notes

- `x-resource` acts as a bridge between API contracts and resource definitions
- Designed to support CodeForge Observer's runtime auditing against **field level permissions**
- UI will reuse existing navigation and selection patterns for consistency

## Help View

### Goal

Enable users to view guides and controls within CodeForge Resource Manager.

This allows new users to quickly understand navigation patterns and available actions across the application.

### Includes

#### General Commands

- Display global keybindings available across all views:
  - Navigation (`h/j/k/l`)
  - Selection (`enter`)
  - Back (`backspace`)
  - Exit (`ctrl+c`)
- Provide a consistent reference for commonly used controls

#### Page Specific Commands

- Display contextual commands based on the current page:
  - Resource management and Resource Binding actions
    - Add (`ctrl+a`)
    - Remove/Delete (`ctrl+b`)
    - Toggle Field Permissions (`space`)
  - Resource Binding
    - Add (`ctrl+a`)
    - Remove/Delete (`ctrl+b`)
- Dynamically update based on `currentPage` and view state

#### Command Descriptions

- Provide short explanations for each command
- Clarify what action does (not just the keybinding)

#### Layout & Organization

- Ensure readability and alignment within the TUI
- Highlight important or frequently used actions

#### Accessibility & Discoverability

- Accessible from the Home Page
- Ensure help content is easy to scan quickly

### Design Notes

- Designed to improve onboarding and usability
- Reuse existing Lipgloss styling for consistency
