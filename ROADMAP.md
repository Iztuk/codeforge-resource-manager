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
