json = require("json")
function main(plan)
  print("=== RUN Files RUN ===")
  print("Files:", json.encode(plan:GetFiles()))

  xpcall(function()
    print(plan:GetFileContent("test_files"))
  end,
  function()
    plan:SetFileContent("test_files", "test.txt", "233")
  end)

  local filename, content = plan:GetFileContent("test_files")
  if content == "233" then
    plan:SetFileContent("test_files", "test2.txt", "456")
  else
    plan:SetFileContent("test_files", "test.txt", "233")
  end
  print(plan:GetFileContent("test_files"))


  -- Job
  print("Job Files:", json.encode(plan:GetJobFiles()))

  xpcall(function()
    print(plan:GetJobFileContent("test_files"))
  end,
  function()
    plan:SetJobFileContent("test_files", "test.txt", "233")
  end)

  local filename, content = plan:GetJobFileContent("test_files")
  if content == "233" then
    plan:SetJobFileContent("test_files", "test2.txt", "456")
  else
    plan:SetJobFileContent("test_files", "test.txt", "233")
  end
  print(plan:GetJobFileContent("test_files"))


  sleep("1s")
  print("=== END Files END ===")
end
