json = require("json")
function main(plan)
  print("=== RUN Extra RUN ===")
  print("Extra:", json.encode(plan:GetExtra()))

  local extra = plan:GetExtra()
  extra["xxx"] = "aa"
  extra["ccc"] = "aaa"
  plan:SetExtra(extra)

  local job_extra = plan:GetJobExtra()
  job_extra["ttt"] = "xxx"
  print("Job Extra:", json.encode(job_extra))
  plan:SetJobExtra(job_extra)
  sleep("1s")
  print("=== END Extra END ===")
end
