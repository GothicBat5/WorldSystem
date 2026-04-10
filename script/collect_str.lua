local roles = {}

-- Utility functions:
--Find person = not working, show groups have a bug, should be = No group found
local function trim(s)
  return s:match("^%s*(.-)%s*$")
end

local function capitalize(str)
  return str:sub(1,1):upper() .. str:sub(2):lower()
end

local function input(prompt)
  io.write(prompt)
  return trim(io.read() or "")
end

-- Core features
local function addPerson()
  local name = input("Enter name: ")
  if name == "" then
    print("Name cannot be empty.\n")
    return
  end

  local role = capitalize(input("Enter role: "))
  if role == "" then
    print("Role cannot be empty.\n")
    return
  end

  roles[role] = roles[role] or {}
  table.insert(roles[role], name)

  print("Added " .. name .. " as " .. role .. "\n")
end

local function showGroups()
  print("\n------ People Groups ------")

  -- Sort roles
  local sortedRoles = {}
  for role in pairs(roles) do
    table.insert(sortedRoles, role)
  end
  table.sort(sortedRoles)

  for _, role in ipairs(sortedRoles) do
    table.sort(roles[role])
    print(role .. "s: " .. table.concat(roles[role], ", "))
  end
  print()
end

local function deletePerson()
  local name = input("Enter name to delete: ")

  for role, people in pairs(roles) do
    for i, person in ipairs(people) do
      if person:lower() == name:lower() then
        table.remove(people, i)
        print("Removed " .. person .. " from " .. role .. "\n")
        return
      end
    end
  end

  print("Person not found.\n")
end

local function searchPerson()
  local name = input("Search name: ")

  for role, people in pairs(roles) do
    for _, person in ipairs(people) do
      if person:lower():find(name:lower()) then
        print(person .. " is a " .. role)
      end
    end
  end
  print()
end

-- Save/Load (basic text format)
local function saveToFile()
  local file = io.open("roles.txt", "w")
  if not file then
    print("Error saving file.\n")
    return
  end

  for role, people in pairs(roles) do
    for _, person in ipairs(people) do
      file:write(role .. "|" .. person .. "\n")
    end
  end

  file:close()
  print("Data saved!\n")
end

local function loadFromFile()
  local file = io.open("roles.txt", "r")
  if not file then return end

  for line in file:lines() do
    local role, name = line:match("(.+)|(.+)")
    if role and name then
      roles[role] = roles[role] or {}
      table.insert(roles[role], name)
    end
  end

  file:close()
end

-- Menu
local function showMenu()
  print("===== MENU =====")
  print("1. Add Person")
  print("2. Show Groups")
  print("3. Delete Person")
  print("4. Search Person")
  print("5. Save")
  print("6. Quit")
end

-- Main loop
loadFromFile()

while true do
  showMenu()
  local choice = input("Choose: ")

  if choice == "1" then
    addPerson()
  elseif choice == "2" then
    showGroups()
  elseif choice == "3" then
    deletePerson()
  elseif choice == "4" then
    searchPerson()
  elseif choice == "5" then
    saveToFile()
  elseif choice == "6" then
    print("Goodbye 👋")
    break
  else
    print("Invalid choice.\n")
  end
end
