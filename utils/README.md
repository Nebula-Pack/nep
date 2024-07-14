# utils Package

This package provides utility functions to manage project directories and configurations.
Easy to use prompts and error handling are provided.

## Functions

### `CreateProject`

Creates a new project directory with a `nebpack` folder and `nebula-config.json` file.

**Parameters:**

- `projectName` (string): The name of the project directory.
- `useCurrentDir` (bool): Use the current directory if `true`, otherwise create a new directory.

**Returns:**

- `projectDir` (string): The path to the project directory.
- `error`: An error if the project creation fails.

**Example:**

```go
_, err := CreateProject("myProject", false)
if err != nil {
    fmt.Println(err)
}
```

### `UpdateConfig`

Updates the `nebula-config.json` file with given key-value pairs.

**Parameters:**

- `projectDir` (string): Path to the project directory containing the `nebula-config.json` file.
- `updates` ([]UpdatePath): A slice of `UpdatePath` structs defining the key-value pairs to update.

**Returns:**

- `error`: An error object if something went wrong during configuration update. Otherwise, nil.

**Example:**

```go
projectDir := "./your_project_directory"
// Example: Adding new key-value pairs and updating nested keys
updates := []utils.UpdatePath{
  {Path: []string{"key1"}, Value: "newValue1"},
  {Path: []string{"key2", "subKey"}, Value: "newValue2"},
  {Path: []string{"key3"}, Value: []string{"item1", "item2"}},
  {Path: []string{"key4", "nestedMap"}, Value: map[string]interface{}{"innerKey": "innerValue"}},
}
err := utils.UpdateConfig(projectDir, updates)
if err != nil {
  log.Fatalf("Error updating config: %v", err)
else {
  fmt.Println("Config updated successfully")
}

```

### `GroupedTextInput`

Collects multiple text inputs in a single prompt.
The user inputs are returned as a slice of strings.
Styled prompts are used to guide the user.

**Parameters:**

- `prompts` ([]string): The prompts for each text input.

**Returns:**

- `[]string`: The user inputs.
- `error`: An error if the input collection fails.

**Example:**

```go
prompts := []string{"Name", "Age", "Email"}
inputs, err := GroupedTextInput(prompts)
if err != nil {
    fmt.Println(err)
}
fmt.Println(inputs)
```

### `SelectFromList`

Creates an interactive list for selection.
The user can navigate through the list and select an item.
The selected item is returned as a string.

**Parameters:**

- `title` (string): The title of the list.
- `items` ([]Item): The items to display in the list.
- `itemsToShow` (int): The number of items to show at once.

**Returns:**

- `string`: The selected item.
- `error`: An error if the selection fails.

**Example:**

```go
items := []Item{
    {TitleText: "Option 1", Desc: "Description 1"},
    {TitleText: "Option 2", Desc: "Description 2"},
}
selected, err := SelectFromList("Choose an option", items, 3)
if err != nil {
    fmt.Println(err)
}
fmt.Println(selected)
```
